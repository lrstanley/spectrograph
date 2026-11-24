<script setup lang="ts">
import { toast } from "vue-sonner"
import type { SchemaUser } from "#open-fetch-schemas/api"

const props = defineProps<{ user: SchemaUser | null }>()
const emit = defineEmits<{ close: []; changed: [user: SchemaUser] }>()

const { $api } = useNuxtApp()
const reason = ref("")
const pending = ref(false)
const error = ref("")

watch(
  () => props.user,
  () => {
    reason.value = ""
    error.value = ""
  },
)

const reasonError = computed(() =>
  reason.value.trim().length === 0
    ? "Give a reason, it's shown to other admins."
    : reason.value.length > 500
      ? "Use 500 characters or fewer."
      : "",
)

async function submit() {
  const user = props.user
  if (!user) return
  if (!user.banned && reasonError.value) return
  pending.value = true
  error.value = ""
  try {
    const updated = user.banned
      ? await $api("/users/{userID}/unban", { method: "POST", path: { userID: user.id } })
      : await $api("/users/{userID}/ban", {
          method: "POST",
          path: { userID: user.id },
          body: { reason: reason.value.trim() },
        })
    toast.success(user.banned ? `Unbanned ${displayName(user)}` : `Banned ${displayName(user)}`)
    emit("changed", updated)
    emit("close")
  } catch (err) {
    error.value = apiErrorMessage(err)
  } finally {
    pending.value = false
  }
}
</script>

<template>
  <Dialog :open="user !== null" @update:open="!$event && emit('close')">
    <DialogContent v-if="user">
      <DialogHeader>
        <DialogTitle>{{ user.banned ? "Unban" : "Ban" }} {{ displayName(user) }}?</DialogTitle>
        <DialogDescription v-if="user.banned">
          They'll be able to log in and manage their servers again.
          <template v-if="user.ban_reason"> They were banned for: "{{ user.ban_reason }}".</template>
        </DialogDescription>
        <DialogDescription v-else>
          They'll be signed out and can't log in again. Their servers keep working with Spectrograph.
        </DialogDescription>
      </DialogHeader>
      <form id="ban-user-form" novalidate @submit.prevent="submit">
        <FieldGroup v-if="!user.banned">
          <Field :data-invalid="(reason.length > 0 && !!reasonError) || undefined">
            <FieldLabel for="ban-reason">Reason</FieldLabel>
            <Textarea
              id="ban-reason"
              v-model="reason"
              rows="3"
              maxlength="500"
              placeholder="Spamming channel creation"
              :aria-invalid="(reason.length > 0 && !!reasonError) || undefined"
            />
            <FieldDescription>{{ reason.length }}/500</FieldDescription>
            <FieldError v-if="reason.length > 0 && reasonError" :errors="[reasonError]" />
          </Field>
        </FieldGroup>
      </form>
      <Alert v-if="error" variant="destructive">
        <AlertDescription>{{ error }}</AlertDescription>
      </Alert>
      <DialogFooter>
        <Button variant="ghost" :disabled="pending" @click="emit('close')">Cancel</Button>
        <Button
          type="submit"
          form="ban-user-form"
          :variant="user.banned ? 'default' : 'destructive'"
          :disabled="pending || (!user.banned && !!reasonError)"
        >
          <Spinner v-if="pending" data-icon="inline-start" />
          {{ user.banned ? "Unban user" : "Ban user" }}
        </Button>
      </DialogFooter>
    </DialogContent>
  </Dialog>
</template>
