export const repoURL = "https://github.com/lrstanley/spectrograph"
export const issuesURL = `${repoURL}/issues`

// Template and field IDs come from the bug report form inherited from
// https://github.com/lrstanley/.github/blob/master/.github/ISSUE_TEMPLATE/bug_report.yml
const bugTemplate = "bug_report.yml"

export type BugReportContext = {
  version: string
  path: string
  userAgent?: string
  title?: string
  status?: number
  errorType?: string
  requestID?: string
  timestamp?: string
}

// Numeric path segments are user or guild IDs. Replace them so links never carry
// identifiers.
export function redactPath(path: string) {
  return (path.split(/[?#]/)[0] ?? "/").replace(/\/\d+(?=\/|$)/g, "/:id")
}

export function bugReportURL(ctx: BugReportContext) {
  const environment = [
    `- Version: ${ctx.version}`,
    `- Page: ${redactPath(ctx.path)}`,
    ctx.userAgent ? `- Browser: ${ctx.userAgent}` : undefined,
  ].filter(Boolean)

  const context = [
    ctx.status ? `- Status: ${ctx.status}${ctx.errorType ? ` (${ctx.errorType})` : ""}` : undefined,
    ctx.requestID ? `- Request ID: ${ctx.requestID}` : undefined,
    ctx.timestamp ? `- Time: ${ctx.timestamp}` : undefined,
  ].filter(Boolean)

  const params = new URLSearchParams({ template: bugTemplate })
  if (ctx.title) {
    params.set("title", `bug: ${ctx.title}`)
  }
  params.set("environment", environment.join("\n"))
  if (context.length) {
    params.set("context", context.join("\n"))
  }
  return `${repoURL}/issues/new?${params.toString()}`
}
