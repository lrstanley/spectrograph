// Guilds the signed-in user can manage (shared by the sidebar and dashboard).
export function useMyGuilds() {
  const { self } = useSelf()
  return useApi("/users/{userID}/user-guilds", {
    key: "my-guilds",
    path: computed(() => ({ userID: self.value?.id ?? 0 })),
    query: { per_page: 100, sort: "name", order: "asc" },
    // Also stops a refetch for a signed-out user (e.g. after deleting the account).
    enabled: () => !!self.value,
  })
}
