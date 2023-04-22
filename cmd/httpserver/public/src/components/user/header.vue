<template>
  <div class="md:flex md:items-center md:justify-between md:space-x-5">
    <div class="flex items-start space-x-5">
      <div class="shrink-0">
        <div class="relative">
          <div class="avatar">
            <div
              class="w-16 h-16 rounded-full ring"
              :class="user.banned ? 'ring-error' : 'ring-success'"
            >
              <img :src="user.avatarURL" />
            </div>
          </div>
        </div>
      </div>
      <div class="pt-0 md:pt-1.5">
        <h1 class="text-2xl font-bold text-nitro-500">{{ user.username }}</h1>
        <span class="flex flex-col gap-2 md:flex-row">
          <span class="badge badge-primary">#{{ user.discriminator }}</span>
          <span class="badge badge-success">joined: {{ joinedAt }}</span>
        </span>
      </div>
    </div>
    <div
      class="flex flex-col-reverse mt-6 space-y-4 space-y-reverse justify-stretch sm:flex-row-reverse sm:justify-end sm:space-y-0 sm:space-x-3 sm:space-x-reverse md:mt-0 md:flex-row md:space-x-3"
    >
      <a
        v-if="user.id === state.base?.self?.id"
        class="h-10 min-h-0 btn btn-error"
        href="/-/auth/logout"
      >
        <i-fas-circle-xmark class="h-4 mr-2" />
        Logout
      </a>
    </div>
  </div>

  <FeedbackAlert v-if="user.banned" type="error" class="mt-10">
    User has been banned by an administrator
    <span v-if="user.banReason" class="rounded badge">reason: {{ user.banReason }}</span>
  </FeedbackAlert>
</template>

<script setup lang="ts">
import type { User } from "@/lib/api"

const props = defineProps<{
  user: User
}>()

const state = useState()
const joinedAt = useTimeAgo(computed(() => props.user.createTime))
</script>
