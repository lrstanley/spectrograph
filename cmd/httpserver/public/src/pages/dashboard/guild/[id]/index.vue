<route lang="yaml">
meta:
  layout: dashboard
  title: Dashboard
</route>

<template>
  <FeedbackAlert v-if="error" status="error">{{ error }}</FeedbackAlert>
  <div v-else>
    <div class="md:flex md:items-center md:justify-between md:space-x-5">
      <div class="flex items-start space-x-5">
        <div class="shrink-0">
          <div class="relative">
            <GuildIcon
              :guild="guild"
              class="w-16 h-16 text-4xl border-2 border-solid border-balance-500"
              size="2xl"
            />
          </div>
        </div>
        <!--
        Use vertical padding to simulate center alignment when both lines of text are one line,
        but preserve the same layout if the text wraps without making the image jump around.
      -->
        <div class="pt-1.5">
          <h1 class="text-2xl font-bold text-white">Guild: {{ guild.name }}</h1>
          <!-- <p class="text-sm font-medium text-white">joined: {{ joinedAt }}</p> -->
          <span
            class="inline-flex items-center rounded bg-discord-600 px-2 py-0.5 text-xs font-medium text-white"
          >
            joined: {{ joinedAt }}
          </span>
        </div>
      </div>
      <div
        class="flex flex-col-reverse mt-6 space-y-4 space-y-reverse justify-stretch sm:flex-row-reverse sm:justify-end sm:space-y-0 sm:space-x-3 sm:space-x-reverse md:mt-0 md:flex-row md:space-x-3"
      >
        <BaseButton
          :type="guild.guildConfig?.enabled ? 'error' : 'success'"
          class="cursor-pointer"
          @click="toggleEnabled"
        >
          {{ guild.guildConfig?.enabled ? "Disable" : "Enable" }}
        </BaseButton>
        <!--
        <button
          type="button"
          class="inline-flex items-center justify-center px-4 py-2 text-sm font-medium text-white bg-indigo-600 border border-transparent rounded-md shadow-sm hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 focus:ring-offset-gray-100"
        >
          Advance to offer
        </button> -->
      </div>
    </div>

    <div class="mt-5 border-4 border-gray-200 border-dashed rounded-lg h-96" />
  </div>
</template>

<script setup lang="ts" async>
import { useTimeAgo } from "@vueuse/core"
import { useGetGuildQuery, useUpdateGuildConfigMutation } from "@/lib/api"
import type { Guild } from "@/lib/api"

const route = useRoute("/dashboard/guild/[id]/")

const {
  data,
  error,
  executeQuery: refetch,
} = await useGetGuildQuery({ variables: { id: route.params.id } })
const guild = computed(() => data?.value.node as Guild)
const joinedAt = useTimeAgo(guild.value.joinedAt)

const config = useUpdateGuildConfigMutation()

function toggleEnabled() {
  config
    .executeMutation({
      id: guild.value.guildConfig.id,
      input: { enabled: !guild.value.guildConfig.enabled },
    })
    .then(refetch)
}
</script>
