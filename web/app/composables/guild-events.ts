import type { SchemaGuildEventRead } from "#open-fetch-schemas/api"

export type StreamState = "connecting" | "live" | "reconnecting" | "offline"

const pageSize = 50
const maxPageSize = 500

// Guild events: paginated history (newest first) merged with a live SSE stream.
// On connect, the stream replays recent events (or only newer ones, when
// reconnecting with Last-Event-ID), so events are de-duplicated by ID.
export function useGuildEvents(guildID: MaybeRefOrGetter<number>) {
  const perPage = ref(pageSize)
  const history = useApi("/guilds/{guildID}/guild-events", {
    path: computed(() => ({ guildID: toValue(guildID) })),
    query: computed(() => ({ per_page: perPage.value, sort: "create_time" as const, order: "desc" as const })),
  })

  const live = ref<SchemaGuildEventRead[]>([])
  const state = ref<StreamState>("connecting")
  const config = useRuntimeConfig()

  watch(
    () => toValue(guildID),
    (id) => {
      live.value = []
      let source: EventSource | undefined
      let retry: ReturnType<typeof setTimeout> | undefined
      let attempts = 0
      let stopped = false

      const connect = () => {
        state.value = attempts === 0 ? "connecting" : "reconnecting"
        const base = String(config.public.API_URL).replace(/\/$/, "")
        source = new EventSource(`${base}/guilds/${id}/guild-events/stream`, { withCredentials: true })
        source.onopen = () => {
          attempts = 0
          state.value = "live"
        }
        // Events have no "event:" name, so they arrive as plain messages.
        source.onmessage = (message) => {
          try {
            const event = JSON.parse(message.data) as SchemaGuildEventRead
            if (!live.value.some((e) => e.id === event.id)) {
              live.value = [event, ...live.value].slice(0, maxPageSize)
            }
          } catch {
            // Ignore malformed payloads.
          }
        }
        source.onerror = () => {
          // The browser retries on its own (sending Last-Event-ID) unless the
          // server rejected the stream, which closes it.
          if (source?.readyState === EventSource.CONNECTING) {
            state.value = "reconnecting"
            return
          }
          source?.close()
          if (stopped) return
          attempts++
          state.value = attempts > 6 ? "offline" : "reconnecting"
          retry = setTimeout(connect, Math.min(30_000, 1000 * 2 ** attempts))
        }
      }

      connect()

      onWatcherCleanup(() => {
        stopped = true
        clearTimeout(retry)
        source?.close()
      })
    },
    { immediate: true },
  )

  const events = computed(() => {
    const seen = new Set<number>()
    return [...live.value, ...(history.data.value?.content ?? [])]
      .filter((e) => !seen.has(e.id) && seen.add(e.id))
      .sort((a, b) => Date.parse(b.create_time) - Date.parse(a.create_time) || b.id - a.id)
  })

  const hasMore = computed(
    () => !!history.data.value && !history.data.value.is_last_page && perPage.value < maxPageSize,
  )

  function loadMore() {
    perPage.value = Math.min(maxPageSize, perPage.value + pageSize)
  }

  return { events, history, state, hasMore, loadMore }
}
