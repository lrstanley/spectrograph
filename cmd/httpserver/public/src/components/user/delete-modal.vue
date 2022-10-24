<template>
  <FeedbackModal :open="open" :error="error">
    <template #header>
      <div
        class="items-center justify-center hidden w-10 h-10 mx-auto rounded-full bg-dnd-300 shrink-0 md:flex"
      >
        <i-fas-triangle-exclamation class="w-6 h-6 -mt-1 text-dnd-600" aria-hidden="true" />
      </div>

      <h3 class="text-lg font-bold">
        Delete account
        <span class="text-bravery-500">{{ props.user.username }}</span>
        <span class="text-balance-400">#{{ props.user.discriminator }}?</span>
      </h3>
    </template>

    <p class="text-center">
      This will delete your account and all of your data. However, it will not remove Spectrograph from
      your guilds, as there can be multiple admins within Spectrograph for a given guild. Please either
      disable Spectrograph for your guilds, or Remove the bot from your guild via the Discord client.
    </p>

    <template #actions>
      <button class="btn btn-sm" :disabled="loading" @click="close">Cancel</button>
      <button
        class="btn btn-sm btn-error"
        :class="loading && 'loading'"
        :disabled="loading"
        @click="deleteAccount"
      >
        Delete account
      </button>
    </template>
  </FeedbackModal>
</template>

<script setup lang="ts">
import { useDeleteAccountMutation } from "@/lib/api"
import type { User } from "@/lib/api"
const props = defineProps<{
  modelValue: boolean
  user: User
}>()

const emit = defineEmits<{
  (e: "update:modelValue", value: boolean): void
}>()

const open = useVModel(props, "modelValue", emit)
const loading = ref(false)
const error = ref<Error>(null)

const deleteMutation = useDeleteAccountMutation()

function close() {
  open.value = false
  error.value = null
  loading.value = false
}

function deleteAccount() {
  loading.value = true
  error.value = null

  deleteMutation
    .executeMutation({})
    .then(({ error: e }) => {
      if (e) {
        error.value = e
        return
      }

      emit("update:modelValue", false)
      window.location.replace("/-/auth/logout")
    })
    .finally(() => (loading.value = false))
}
</script>
