// The most recent failed API request, used to prefill bug reports with the
// request ID (which maps to server logs) and nothing user-identifying.
export type LastApiError = {
  status: number
  requestID?: string
  at: string
}

export function useLastApiError() {
  return useState<LastApiError | null>("last-api-error", () => null)
}
