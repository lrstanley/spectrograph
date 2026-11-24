const relative = new Intl.RelativeTimeFormat("en", { numeric: "auto" })
const absolute = new Intl.DateTimeFormat("en", { dateStyle: "medium", timeStyle: "short" })
const clock = new Intl.DateTimeFormat("en", { hour: "2-digit", minute: "2-digit", second: "2-digit", hour12: false })

const units: [Intl.RelativeTimeFormatUnit, number][] = [
  ["year", 31_536_000],
  ["month", 2_592_000],
  ["week", 604_800],
  ["day", 86_400],
  ["hour", 3_600],
  ["minute", 60],
  ["second", 1],
]

export function timeAgo(value?: string | null, now = Date.now()) {
  if (!value) return "never"
  const seconds = (new Date(value).getTime() - now) / 1000
  for (const [unit, size] of units) {
    if (Math.abs(seconds) >= size || unit === "second") {
      return relative.format(Math.round(seconds / size), unit)
    }
  }
  return ""
}

export function formatDate(value?: string | null) {
  return value ? absolute.format(new Date(value)) : "never"
}

export function formatClock(value: string) {
  return clock.format(new Date(value))
}
