// Overrides the dashboard breadcrumb title (the route's meta.title by default)
// for the current page, e.g. with the name of the server being viewed. Keyed by
// path, so a stale title is never shown on another page.
export function usePageTitle(title: MaybeRefOrGetter<string | undefined>) {
  const route = useRoute()
  const state = usePageTitleState()
  const path = route.path

  watchEffect(() => {
    state.value = { path, title: toValue(title) }
  })
  onScopeDispose(() => {
    if (state.value?.path === path) state.value = undefined
  })
}

export function usePageTitleState() {
  return useState<{ path: string; title?: string } | undefined>("page-title", () => undefined)
}
