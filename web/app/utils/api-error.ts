// API errors are JSON objects: {error, type, code, request_id, timestamp}.
// Validation errors also have "field", and their error is "<field>: <reason>".
export type ApiErrorBody = {
  error?: string
  type?: string
  code?: number
  request_id?: string
  timestamp?: string
  field?: string
}

function errorBody(err: unknown): ApiErrorBody | undefined {
  if (typeof err === "object" && err && "data" in err) {
    const data = (err as { data?: unknown }).data
    if (typeof data === "object" && data) {
      return data as ApiErrorBody
    }
  }
  return undefined
}

export function apiErrorMessage(err: unknown): string {
  const data = errorBody(err)
  if (typeof data?.error === "string" && data.error) {
    return data.error
  }
  if (err instanceof Error && err.message) {
    return err.message
  }
  return "request failed"
}

export function apiErrorStatus(err: unknown): number | undefined {
  if (typeof err === "object" && err) {
    const value = (err as { statusCode?: unknown; status?: unknown }).statusCode ?? (err as { status?: unknown }).status
    if (typeof value === "number") {
      return value
    }
  }
  return errorBody(err)?.code
}

export function apiRequestID(err: unknown): string | undefined {
  return errorBody(err)?.request_id || undefined
}

// The backend rejects banned users with a "user is banned" error.
export function isBannedError(err: unknown) {
  return apiErrorStatus(err) === 403 && /\bbanned\b/i.test(apiErrorMessage(err))
}

// Returns the field and reason of a validation error (e.g. an invalid
// "regex_match"), so forms can show it on the matching input.
export function apiFieldError(err: unknown): { field: string; message: string } | undefined {
  const data = errorBody(err)
  if (!data?.field || typeof data.error !== "string") return undefined
  const reason = data.error.startsWith(`${data.field}: `) ? data.error.slice(data.field.length + 2) : data.error
  return { field: data.field, message: reason.charAt(0).toUpperCase() + reason.slice(1) }
}
