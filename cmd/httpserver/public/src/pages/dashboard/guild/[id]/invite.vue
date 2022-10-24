<route lang="yaml">
meta:
  layout: dashboard
  title: Invite Guild
</route>

<template>
  <div class="flex flex-col h-full gap-4 mx-auto place-content-center">
    <FeedbackOverlay :loading="fetching">
      <a
        class="relative block w-full px-16 py-6 text-center transition-all duration-200 ease-in-out border border-gray-900 rounded-lg focus:outline-none focus:ring-2 focus:ring-bravery-500 bg-chat-700 hover:bg-chat-800/70"
        :href="'/-/invite/' + guild?.guildID"
        target="_blank"
        rel="noopener noreferrer"
      >
        <img alt="Guild icon" :src="guild?.iconURL" class="w-16 h-16 mx-auto rounded-full" />
        <span class="block mt-2 text-lg font-medium text-white">
          Invite Spectrograph to {{ guild?.name }}
        </span>
      </a>
    </FeedbackOverlay>
  </div>
</template>

<script setup lang="ts">
import { useGetGuildIdQuery } from "@/lib/api"

const route = useRoute("/dashboard/guild/[id]/invite")

const { data, fetching } = useGetGuildIdQuery({
  variables: { id: route.params.id },
})

interface Guild {
  __typename: "Guild"
  id: string
  guildID: string
  iconURL: string
  name: string
}

const guild = computed(() => data.value?.node as Guild)
</script>
