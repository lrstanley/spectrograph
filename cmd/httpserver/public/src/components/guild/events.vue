<template>
  <div
    v-show="guildEvents.length > 0"
    v-bind="$attrs"
    class="overflow-x-auto overflow-y-scroll border rounded border-chat-900"
  >
    <table class="table w-full table-compact" aria-label="guild events">
      <thead>
        <tr>
          <th class="sticky top-0 z-10 pl-4 bg-channel-500/80 backdrop-blur sm:pl-6">Event timestamp</th>
          <th class="sticky top-0 z-10 bg-channel-500/80 backdrop-blur">Status</th>
          <th class="sticky top-0 z-10 bg-channel-500/80 backdrop-blur">Event message</th>
          <th class="sticky top-0 z-10 bg-channel-500/80 backdrop-blur">Metadata</th>
        </tr>
      </thead>
      <tbody class="max-h-32">
        <tr v-for="event in guildEvents" :key="event.id">
          <td class="pl-4 sm:pl-6 text-chat-300">
            {{ useTimeAgo(event.createTime).value }}
          </td>
          <td
            :class="{
              'text-bravery-500': event.type === GuildEventType.Info,
              'text-idle-500': event.type === GuildEventType.Warning,
              'text-dnd-500': event.type === GuildEventType.Error,
              'text-nitro-500': event.type === GuildEventType.Debug,
            }"
            class="font-semibold"
          >
            {{ event.type }}
          </td>
          <td>{{ event.message }}</td>
          <td>
            <template v-if="event.metadata">
              <span v-for="(value, key) in event.metadata" :key="key" class="badge badge-secondary">
                {{ key }}: {{ value }}
              </span>
            </template>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script setup lang="ts">
import { useTimeAgo } from "@vueuse/core"
import { useGuildEventStreamSubscription, GuildEventType } from "@/lib/api"
import type { Guild, GuildEventStreamSubscription } from "@/lib/api"

const state = useState()

const props = defineProps<{
  guild: Guild
}>()

const eventStream = useGuildEventStreamSubscription(
  {
    variables: {
      guildID: props.guild.id,
      types: [
        ...(state.base.self?.admin ? [GuildEventType.Debug] : []),
        GuildEventType.Error,
        GuildEventType.Warning,
        GuildEventType.Info,
      ],
    },
  },
  (
    messages: GuildEventStreamSubscription[] = [],
    response: GuildEventStreamSubscription
  ): GuildEventStreamSubscription[] => {
    return [...messages, response]
  }
)

const guildEvents = computed(() => {
  const ids = eventStream.data?.value?.map((e) => e.guildEventAdded.id)
  if (!ids) return []

  const events = eventStream.data?.value
    ?.map((e) => e.guildEventAdded)
    .filter(({ id }, index) => !ids.includes(id, index + 1))

  if (!events) return []

  return events.sort((a, b) => Number(b.id) - Number(a.id))
})
</script>
