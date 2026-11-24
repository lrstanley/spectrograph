function apiURL(path: string) {
  return `${String(useRuntimeConfig().public.API_URL).replace(/\/$/, "")}${path}`
}

function absolute(path: string) {
  return new URL(path, window.location.origin).href
}

// True if path is prefix, or a page under it.
export function pathWithin(path: string, prefix: string) {
  return path === prefix || path.startsWith(`${prefix}/`)
}

// Pages that need a login.
export function isProtectedPath(path: string) {
  return pathWithin(path, "/dashboard") || pathWithin(path, "/admin")
}

// Returns the "next" query value when it's a same-site path, or the fallback.
// Protocol-relative ("//host") and backslash paths would leave the site.
export function safeNext(value: unknown, fallback = "/dashboard") {
  return typeof value === "string" && /^\/(?![/\\])/.test(value) ? value : fallback
}

// Starts the Discord OAuth2 flow (a full page navigation). The backend only
// redirects back to URLs on its own host.
export function loginURL(next = "/dashboard") {
  return `${apiURL("/auth/providers/discord")}?next=${encodeURIComponent(absolute(next))}`
}

export function logoutURL(next = "/") {
  return `${apiURL("/auth/logout")}?next=${encodeURIComponent(absolute(next))}`
}

// Redirects to Discord to add the bot, optionally to a specific guild
// (snowflake).
export function inviteURL(guildSnowflake?: string) {
  return apiURL(guildSnowflake ? `/invite/${guildSnowflake}` : "/invite")
}
