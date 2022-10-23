<template>
  <FeedbackModal :open="open" :error="error">
    <template #header>
      <div
        class="items-center justify-center hidden w-10 h-10 mx-auto rounded-full bg-dnd-300 shrink-0 md:flex"
      >
        <i-fas-triangle-exclamation class="w-6 h-6 -mt-1 text-dnd-600" aria-hidden="true" />
      </div>

      <h3 class="text-lg font-bold">
        {{ props.user.banned ? "Unban" : "Ban" }} user
        <span class="text-bravery-500">{{ props.user.username }}</span>
        <span class="text-balance-400">#{{ props.user.discriminator }}?</span>
      </h3>
    </template>

    <div v-if="!props.user.banned" class="w-full mt-3 form-control">
      <input
        v-model="banReason"
        type="text"
        placeholder="Provide ban reason"
        class="w-full input input-bordered"
      />
    </div>

    <template #actions>
      <button class="btn btn-sm" :disabled="loading" @click="close">Cancel</button>
      <button
        v-if="props.user.banned"
        class="btn btn-sm btn-error"
        :class="loading && 'loading'"
        :disabled="loading"
        @click="unbanUser"
      >
        Unban user
      </button>
      <button
        v-else
        class="btn btn-sm btn-error"
        :class="loading && 'loading'"
        :disabled="loading"
        @click="banUser"
      >
        Ban user
      </button>
    </template>
  </FeedbackModal>
</template>

<script setup lang="ts">
import { useBanUserMutation, useUnbanUserMutation } from "@/lib/api"
import type { User } from "@/lib/api"
const props = defineProps<{
  modelValue: boolean
  user: User
}>()

const emit = defineEmits<{
  (e: "banned", value: boolean): void
  (e: "update:modelValue", value: boolean): void
}>()

const open = useVModel(props, "modelValue", emit)
const loading = ref(false)
const error = ref<Error>(null)

const banReason = ref(props.user.banReason)
const banMutation = useBanUserMutation()
const unbanMutation = useUnbanUserMutation()

function close() {
  open.value = false
  error.value = null
  loading.value = false
}

function banUser() {
  loading.value = true
  error.value = null

  banMutation
    .executeMutation({
      id: props.user.id,
      reason: banReason.value,
    })
    .then(({ error: e }) => {
      if (e) {
        error.value = e
        return
      }

      emit("banned", true)
      emit("update:modelValue", false)
    })
    .finally(() => (loading.value = false))
}

function unbanUser() {
  loading.value = true
  error.value = null

  unbanMutation
    .executeMutation({
      id: props.user.id,
    })
    .then(({ error: e }) => {
      if (e) {
        error.value = e
        return
      }

      emit("banned", false)
      emit("update:modelValue", false)
    })
    .finally(() => (loading.value = false))
}
</script>
