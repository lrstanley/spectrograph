<script setup lang="ts">
import { toast } from "vue-sonner"
import type { SchemaGuildAdminConfig, SchemaGuildAdminConfigUpdate } from "#open-fetch-schemas/api"

const props = defineProps<{ config: SchemaGuildAdminConfig }>()
const emit = defineEmits<{ saved: [config: SchemaGuildAdminConfig] }>()

const { $api } = useNuxtApp()

type Draft = Required<SchemaGuildAdminConfigUpdate>

function fromConfig(config: SchemaGuildAdminConfig): Draft {
  return {
    enabled: config.enabled ?? true,
    default_max_channels: config.default_max_channels ?? 0,
    default_max_clones: config.default_max_clones ?? 0,
    comment: config.comment ?? "",
  }
}

const draft = ref<Draft>(fromConfig(props.config))
watch(
  () => props.config,
  (config) => {
    draft.value = fromConfig(config)
  },
)

const pending = ref(false)
const serverError = ref("")
// A validation error from the API for one field, shown until that field changes.
const fieldError = ref<{ field: keyof Draft; value: unknown; message: string }>()

function apiError(field: keyof Draft) {
  const err = fieldError.value
  return err?.field === field && err.value === draft.value[field] ? err.message : ""
}

function countError(field: "default_max_channels" | "default_max_clones") {
  const value = draft.value[field]
  return !Number.isInteger(value) || value < 0 ? "Enter a whole number, 0 or more." : apiError(field)
}

const errors = computed(() => ({
  default_max_channels: countError("default_max_channels"),
  default_max_clones: countError("default_max_clones"),
  comment: apiError("comment"),
}))

const changes = computed(() => {
  const saved = fromConfig(props.config)
  return Object.fromEntries(
    (Object.keys(saved) as (keyof Draft)[])
      .filter((key) => draft.value[key] !== saved[key])
      .map((key) => [key, draft.value[key]]),
  ) as SchemaGuildAdminConfigUpdate
})
const dirty = computed(() => Object.keys(changes.value).length > 0)
const invalid = computed(() => Object.values(errors.value).some(Boolean))

async function save() {
  if (invalid.value || !dirty.value) return
  pending.value = true
  serverError.value = ""
  fieldError.value = undefined
  // Errors refer to the submitted values, which may be edited while saving.
  const body = changes.value
  try {
    emit(
      "saved",
      await $api("/guild-admin-configs/{guildadminconfigID}", {
        method: "PATCH",
        path: { guildadminconfigID: props.config.id },
        body,
      }),
    )
    toast.success("Admin overrides saved")
  } catch (err) {
    const field = apiFieldError(err)
    if (field && field.field in errors.value) {
      const key = field.field as keyof Draft
      fieldError.value = { field: key, value: body[key], message: field.message }
    } else {
      serverError.value = apiErrorMessage(err)
    }
  } finally {
    pending.value = false
  }
}
</script>

<template>
  <Card class="border-primary/40">
    <CardHeader>
      <CardTitle>Admin overrides</CardTitle>
      <CardDescription>
        Only Spectrograph administrators see this. These values take precedence over the server's own settings.
      </CardDescription>
    </CardHeader>
    <CardContent>
      <form id="guild-admin-config-form" novalidate @submit.prevent="save">
        <FieldGroup>
          <Field orientation="horizontal">
            <Switch id="admin-enabled" v-model="draft.enabled" />
            <FieldContent>
              <FieldLabel for="admin-enabled">Allow Spectrograph in this server</FieldLabel>
              <FieldDescription>Turning this off stops the bot here, whatever the server's admins choose.</FieldDescription>
            </FieldContent>
          </Field>
          <div class="grid gap-6 sm:grid-cols-2">
            <Field :data-invalid="!!errors.default_max_channels || undefined">
              <FieldLabel for="admin-max-channels">Maximum channels</FieldLabel>
              <Input
                id="admin-max-channels"
                v-model.number="draft.default_max_channels"
                type="number"
                min="0"
                step="1"
                :aria-invalid="!!errors.default_max_channels || undefined"
              />
              <FieldDescription>0 uses the system default.</FieldDescription>
              <FieldError v-if="errors.default_max_channels" :errors="[errors.default_max_channels]" />
            </Field>
            <Field :data-invalid="!!errors.default_max_clones || undefined">
              <FieldLabel for="admin-max-clones">Maximum clones</FieldLabel>
              <Input
                id="admin-max-clones"
                v-model.number="draft.default_max_clones"
                type="number"
                min="0"
                step="1"
                :aria-invalid="!!errors.default_max_clones || undefined"
              />
              <FieldDescription>0 uses the system default.</FieldDescription>
              <FieldError v-if="errors.default_max_clones" :errors="[errors.default_max_clones]" />
            </Field>
          </div>
          <Field :data-invalid="!!errors.comment || undefined">
            <FieldLabel for="admin-comment">Note</FieldLabel>
            <Textarea id="admin-comment" v-model="draft.comment" rows="3" :aria-invalid="!!errors.comment || undefined" />
            <FieldDescription>Shown to the server's admins when the bot is disabled here.</FieldDescription>
            <FieldError v-if="errors.comment" :errors="[errors.comment]" />
          </Field>
        </FieldGroup>
      </form>
      <Alert v-if="serverError" variant="destructive" class="mt-6">
        <AlertTitle>Overrides weren't saved</AlertTitle>
        <AlertDescription>{{ serverError }}</AlertDescription>
      </Alert>
    </CardContent>
    <CardFooter class="flex flex-wrap items-center gap-3 border-t">
      <span class="text-xs text-muted-foreground" :title="formatDate(config.update_time)">
        Last saved {{ timeAgo(config.update_time) }}
      </span>
      <Button
        class="ml-auto"
        type="submit"
        form="guild-admin-config-form"
        :disabled="!dirty || invalid || pending"
      >
        <Spinner v-if="pending" data-icon="inline-start" />
        Save overrides
      </Button>
    </CardFooter>
  </Card>
</template>
