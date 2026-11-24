import type { SchemaGuild, SchemaGuildEdges, SchemaUser } from "#open-fetch-schemas/api"

export type GuildStatus = "active" | "paused" | "admin-disabled" | "not-joined"

type GuildWithEdges = SchemaGuild & { edges?: SchemaGuildEdges }

// Mirrors how the worker decides whether to act on a guild.
export function guildStatus(guild: GuildWithEdges): GuildStatus {
  if (!guild.joined_at) return "not-joined"
  if (guild.edges?.guild_admin_config?.enabled === false) return "admin-disabled"
  if (guild.edges?.guild_config?.enabled === false) return "paused"
  return "active"
}

export const guildStatusLabel: Record<GuildStatus, string> = {
  active: "Active",
  paused: "Paused",
  "admin-disabled": "Disabled by an admin",
  "not-joined": "Bot not added",
}

export function initials(name: string) {
  return (
    name
      .split(/\s+/)
      .filter(Boolean)
      .slice(0, 2)
      .map((part) => part[0])
      .join("")
      .toUpperCase() || "?"
  )
}

export function displayName(user: Pick<SchemaUser, "username" | "global_name">) {
  return user.global_name || user.username
}

// Legacy Discord tags; migrated accounts have "0" or no discriminator.
export function userTag(user: Pick<SchemaUser, "username" | "discriminator">) {
  return user.discriminator && user.discriminator !== "0"
    ? `${user.username}#${user.discriminator}`
    : `@${user.username}`
}
