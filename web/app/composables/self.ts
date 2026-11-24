const selfKey = "self"

// The authenticated user, or null when signed out. Loaded once by the global
// auth middleware, and shared through the Nuxt data cache. Banned users (the
// API responds with 403 "user is banned") are treated as signed out, with
// banned set.
export function useSelf() {
  const request = useApi("/self", {
    key: selfKey,
    dedupe: "defer",
    immediate: false,
    watch: false,
  })
  const loaded = useState("self-loaded", () => false)
  const banned = useState("self-banned", () => false)

  async function load() {
    if (loaded.value) {
      return request.data.value ?? null
    }
    await request.execute()
    const err = request.error.value
    if (err && apiErrorStatus(err) !== 401) {
      if (!isBannedError(err)) {
        throw err
      }
      request.clear()
      banned.value = true
    }
    loaded.value = true
    return request.data.value ?? null
  }

  // Marks the user as signed out (e.g. after deleting their account), clearing
  // the cached response as well as the data.
  function clear() {
    request.clear()
    loaded.value = true
  }

  // Marks the user as banned and signed out, e.g. when another request is
  // rejected because they were banned mid-session. Returns false if they were
  // already marked as banned.
  function markBanned() {
    clear()
    if (banned.value) return false
    banned.value = true
    return true
  }

  return {
    self: request.data,
    banned: readonly(banned),
    load,
    clear,
    markBanned,
    isAdmin: computed(() => request.data.value?.admin === true),
  }
}
