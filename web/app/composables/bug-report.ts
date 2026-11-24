// Builds a prefilled bug report link for the current page. Only the app
// version, a redacted path, the browser user agent, and (when present) the
// request ID of the last failed API call are included.
export function useBugReportURL(extra?: MaybeRefOrGetter<Partial<BugReportContext> | undefined>) {
  const route = useRoute()
  const config = useRuntimeConfig()
  const lastError = useLastApiError()

  return computed(() => {
    const overrides = toValue(extra) ?? {}
    const err = lastError.value
    return bugReportURL({
      version: String(config.public.version),
      path: route.path,
      userAgent: navigator.userAgent,
      status: err?.status,
      requestID: err?.requestID,
      timestamp: err?.at,
      ...overrides,
    })
  })
}
