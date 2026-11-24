<script setup lang="ts">
import { ExternalLink } from "@lucide/vue"
import { toast } from "vue-sonner"
import type { SchemaGuildConfig, SchemaGuildConfigUpdate } from "#open-fetch-schemas/api"

const props = defineProps<{
  config: SchemaGuildConfig
  disabled?: boolean
}>()

const emit = defineEmits<{ saved: [config: SchemaGuildConfig] }>()

const { $api } = useNuxtApp()

type Draft = Required<Pick<SchemaGuildConfigUpdate, "regex_match" | "default_max_clones" | "contact_email">>

function fromConfig(config: SchemaGuildConfig): Draft {
  return {
    regex_match: config.regex_match ?? "",
    default_max_clones: config.default_max_clones ?? 0,
    contact_email: config.contact_email ?? "",
  }
}

// Editable draft of the settings, reset when the saved settings change. Other
// updates to the config (e.g. pausing the bot) keep unsaved edits.
const draft = ref<Draft>(fromConfig(props.config))
watch([() => props.config.regex_match, () => props.config.default_max_clones, () => props.config.contact_email], () => {
  draft.value = fromConfig(props.config)
})

const pending = ref(false)
const serverError = ref("")
// A validation error from the API for one field, shown until that field changes.
const fieldError = ref<{ field: keyof Draft; value: unknown; message: string }>()

function apiError(field: keyof Draft) {
  const err = fieldError.value
  return err?.field === field && err.value === draft.value[field] ? err.message : ""
}

const errors = computed(() => ({
  // Patterns are Go (RE2) regular expressions, so only the backend can check them.
  regex_match: byteLength(draft.value.regex_match) > 100 ? "Use 100 characters or fewer." : apiError("regex_match"),
  default_max_clones:
    !Number.isInteger(draft.value.default_max_clones) ||
    draft.value.default_max_clones < 0 ||
    draft.value.default_max_clones > 20
      ? "Enter a whole number from 0 to 20."
      : apiError("default_max_clones"),
  contact_email:
    byteLength(draft.value.contact_email) > 255
      ? "Use 255 characters or fewer."
      : draft.value.contact_email && !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(draft.value.contact_email)
        ? "Enter an email address like name@example.com."
        : apiError("contact_email"),
}))

const changes = computed(() => {
  const body: SchemaGuildConfigUpdate = {}
  if (draft.value.regex_match !== (props.config.regex_match ?? "")) body.regex_match = draft.value.regex_match
  if (draft.value.default_max_clones !== (props.config.default_max_clones ?? 0))
    body.default_max_clones = draft.value.default_max_clones
  if (draft.value.contact_email !== (props.config.contact_email ?? "")) body.contact_email = draft.value.contact_email
  return body
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
    const config = await $api("/guild-configs/{guildconfigID}", {
      method: "PATCH",
      path: { guildconfigID: props.config.id },
      body,
    })
    emit("saved", config)
    toast.success("Settings saved")
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

function reset() {
  draft.value = fromConfig(props.config)
  serverError.value = ""
  fieldError.value = undefined
}
</script>

<template>
  <Card>
    <CardHeader>
      <CardTitle>Channel settings</CardTitle>
      <CardDescription>
        Choose which voice channels Spectrograph manages, and how far it can scale them. You can also change these from
        Discord with <code class="font-mono text-foreground">/spectrograph settings</code>.
      </CardDescription>
    </CardHeader>
    <CardContent>
      <form id="guild-config-form" novalidate @submit.prevent="save">
        <FieldGroup>
          <Field :data-invalid="!!errors.regex_match || undefined" :data-disabled="disabled || undefined">
            <FieldLabel for="regex-match">Channel name pattern</FieldLabel>
            <Input
              id="regex-match"
              v-model="draft.regex_match"
              class="font-mono"
              placeholder="^.* \+$"
              maxlength="100"
              autocomplete="off"
              spellcheck="false"
              :disabled="disabled"
              :aria-invalid="!!errors.regex_match || undefined"
              aria-describedby="regex-match-help"
            />
            <FieldDescription id="regex-match-help">
              A regular expression. Voice channels with matching names are managed. Leave it empty to use the default,
              which matches names ending in a space and a plus, like "Squad +".
              <a href="https://github.com/google/re2/wiki/Syntax" target="_blank" rel="noopener" class="inline-flex items-center gap-0.5">
                Syntax reference
                <ExternalLink class="size-3" aria-hidden="true" />
              </a>
            </FieldDescription>
            <FieldError v-if="errors.regex_match" :errors="[errors.regex_match]" />
          </Field>

          <Field :data-invalid="!!errors.default_max_clones || undefined" :data-disabled="disabled || undefined">
            <FieldLabel for="max-clones">Maximum clones per group</FieldLabel>
            <Input
              id="max-clones"
              v-model.number="draft.default_max_clones"
              type="number"
              inputmode="numeric"
              min="0"
              max="20"
              step="1"
              class="w-32"
              :disabled="disabled"
              :aria-invalid="!!errors.default_max_clones || undefined"
            />
            <FieldDescription>How many extra channels Spectrograph may add to one group. 0 uses the default.</FieldDescription>
            <FieldError v-if="errors.default_max_clones" :errors="[errors.default_max_clones]" />
          </Field>

          <Field :data-invalid="!!errors.contact_email || undefined" :data-disabled="disabled || undefined">
            <FieldLabel for="contact-email">Contact email</FieldLabel>
            <Input
              id="contact-email"
              v-model.trim="draft.contact_email"
              type="email"
              autocomplete="email"
              placeholder="you@example.com"
              :disabled="disabled"
              :aria-invalid="!!errors.contact_email || undefined"
            />
            <FieldDescription>Optional. Where we can reach this server's admins about problems with the bot.</FieldDescription>
            <FieldError v-if="errors.contact_email" :errors="[errors.contact_email]" />
          </Field>
        </FieldGroup>
      </form>
      <Alert v-if="serverError" variant="destructive" class="mt-6">
        <AlertTitle>Settings weren't saved</AlertTitle>
        <AlertDescription>{{ serverError }}</AlertDescription>
      </Alert>
    </CardContent>
    <CardFooter class="flex flex-wrap items-center gap-3 border-t">
      <span class="text-xs text-muted-foreground" :title="formatDate(config.update_time)">
        Last saved {{ timeAgo(config.update_time) }}
      </span>
      <div class="ml-auto flex gap-2">
        <Button variant="ghost" type="button" :disabled="!dirty || pending" @click="reset">Discard changes</Button>
        <Button type="submit" form="guild-config-form" :disabled="disabled || !dirty || invalid || pending">
          <Spinner v-if="pending" data-icon="inline-start" />
          Save settings
        </Button>
      </div>
    </CardFooter>
  </Card>
</template>
