<script setup lang="ts">
import { Radio, RefreshCw } from "@lucide/vue"
import type { SchemaGuildEventRead, SchemaGuildEventTypeEnum } from "#open-fetch-schemas/api"
import { cn } from "@/lib/utils"

defineProps<{
  events: SchemaGuildEventRead[]
  state?: StreamState
  loading?: boolean
  // Error loading the event history, if any.
  error?: unknown
  hasMore?: boolean
}>()

const emit = defineEmits<{ loadMore: []; retry: [] }>()

const typeClass: Record<SchemaGuildEventTypeEnum, string> = {
  DEBUG: "bg-muted-foreground",
  INFO: "bg-magma-2",
  WARNING: "bg-magma-5",
  ERROR: "bg-destructive",
}

const stateLabel: Record<StreamState, string> = {
  connecting: "Connecting",
  live: "Live",
  reconnecting: "Reconnecting",
  offline: "Offline",
}

function metadata(event: SchemaGuildEventRead) {
  return Object.entries(event.metadata ?? {}).map(([key, value]) => ({
    key,
    value: typeof value === "object" ? JSON.stringify(value) : String(value),
  }))
}
</script>

<template>
  <section aria-labelledby="event-log-title" class="flex min-h-0 flex-col border bg-card">
    <header class="flex items-center justify-between gap-3 border-b px-4 py-3">
      <div>
        <h2 id="event-log-title" class="font-heading text-base font-semibold">Activity</h2>
        <p class="text-xs text-muted-foreground">What Spectrograph did in this server, newest first.</p>
      </div>
      <span
        v-if="state"
        :class="
          cn(
            'inline-flex items-center gap-1.5 text-xs',
            state === 'live' ? 'text-success' : state === 'offline' ? 'text-destructive' : 'text-muted-foreground',
          )
        "
        role="status"
      >
        <Radio :class="cn('size-3.5', state === 'live' && 'motion-safe:animate-pulse')" aria-hidden="true" />
        {{ stateLabel[state] }}
      </span>
    </header>

    <Alert v-if="error && !loading" variant="destructive" class="m-4 w-auto">
      <AlertTitle>Couldn't load earlier activity</AlertTitle>
      <AlertDescription>
        <p>{{ apiErrorMessage(error) }}</p>
        <Button variant="outline" size="sm" class="mt-2" @click="emit('retry')">
          <RefreshCw data-icon="inline-start" />
          Try again
        </Button>
      </AlertDescription>
    </Alert>
    <div v-if="loading && events.length === 0" class="flex flex-col gap-3 p-4">
      <Skeleton v-for="n in 5" :key="n" class="h-5 w-full" />
    </div>
    <Empty v-else-if="events.length === 0 && !error" class="border-0">
      <EmptyHeader>
        <EmptyTitle>No activity yet</EmptyTitle>
        <EmptyDescription>
          Events show up here as people join and leave matching voice channels.
        </EmptyDescription>
      </EmptyHeader>
    </Empty>
    <!-- Relative, so the absolutely positioned screen reader labels stay inside the scroll area. -->
    <ol v-else-if="events.length" class="relative max-h-[28rem] divide-y overflow-y-auto" aria-label="Server events">
      <li
        v-for="event in events"
        :key="event.id"
        class="grid grid-cols-[auto_1fr] gap-x-3 px-4 py-2.5 text-sm"
        :data-event-id="event.id"
      >
        <span
          :class="cn('mt-1.5 size-2 shrink-0 rounded-full', typeClass[event.type])"
          :title="event.type"
          aria-hidden="true"
        />
        <div class="flex min-w-0 flex-col gap-1">
          <div class="flex flex-wrap items-baseline gap-x-2">
            <span class="sr-only">{{ event.type }}:</span>
            <span class="min-w-0 break-words">{{ event.message }}</span>
            <time
              :datetime="event.create_time"
              :title="formatDate(event.create_time)"
              class="ml-auto font-mono text-[11px] tabular-nums text-muted-foreground"
            >
              {{ formatClock(event.create_time) }}
            </time>
          </div>
          <div v-if="metadata(event).length || event.type !== 'INFO'" class="flex flex-wrap gap-1">
            <Badge v-if="event.type !== 'INFO'" :variant="event.type === 'ERROR' ? 'destructive' : 'outline'">
              {{ event.type.toLowerCase() }}
            </Badge>
            <Badge v-for="item in metadata(event)" :key="item.key" variant="secondary" class="font-mono font-normal">
              {{ item.key }}={{ item.value }}
            </Badge>
          </div>
        </div>
      </li>
    </ol>
    <footer v-if="hasMore" class="border-t p-2">
      <Button variant="ghost" size="sm" class="w-full" :disabled="loading" @click="emit('loadMore')">
        <Spinner v-if="loading" data-icon="inline-start" />
        Show older events
      </Button>
    </footer>
  </section>
</template>
