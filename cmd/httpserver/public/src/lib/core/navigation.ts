/**
 * Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
 * of this source code is governed by the MIT license that can be found in
 * the LICENSE file.
 */
import { h } from "vue"
import GuildIcon from "@/components/guild/icon.vue"
import { LinkStatus } from "@/lib/core/types"

import type { Link, DashboardLink } from "@/lib/core/types"
import type { ComputedRef } from "vue"
import type { RouteLocationNormalized } from "vue-router/auto"
import type { Guild } from "@/lib/api"

const state = useState()

export const headerLinks: Link[] = [
  {
    name: "Features",
    description: "Features provided by Spectrograph",
    href: "/#more-info",
    icon: IEmojiSparkles,
  },
  {
    name: "Documentation",
    description: "Our documentation includes guides on proper setup of Spectrograph",
    href: "#",
    icon: IFasBook,
  },
  {
    name: "Service Health",
    description: "Health of all components of the platform",
    to: "/info/service-health",
    icon: IFasHeartPulse,
  },
]

export const mainLinks: Link[] = [
  {
    name: "Documentation",
    description: "Our documentation includes guides on proper setup of Spectrograph",
    href: "#",
    icon: IFasBook,
  },
  {
    name: "Service Health",
    description: "Health of all components of the platform",
    to: "/info/service-health",
    icon: IFasHeartPulse,
  },
  {
    name: "Privacy",
    description: "Spectrograph's Privacy Policy",
    to: "/info/privacy",
    icon: IFasFingerprint,
  },
  {
    name: "Terms",
    description: "Spectrograph's Terms of Service",
    to: "/info/terms",
    icon: IFasBuildingColumns,
  },
]

export function guildLink(guild: Guild): DashboardLink {
  const status = guild.joinedAt
    ? (guild.guildConfig?.enabled && guild.guildAdminConfig?.enabled && LinkStatus.Healthy) ||
      LinkStatus.Unhealthy
    : LinkStatus.NotJoined

  let to: RouteLocationNormalized

  if (status == LinkStatus.NotJoined) {
    to = {
      name: "/dashboard/guild/[id]/invite",
      params: { id: guild.id },
    } as RouteLocationNormalized<"/dashboard/guild/[id]/invite">
  } else {
    to = {
      name: "/dashboard/guild/[id]/",
      params: { id: guild.id },
    } as RouteLocationNormalized<"/dashboard/guild/[id]/">
  }

  return {
    guild: guild,
    name: guild.name,
    description: `View ${guild.name}'s dashboard`,
    to: to,
    icon: h(GuildIcon, { guild: guild }, {}),
    isChild: true,
    hasStatus: true,
    status: status,
  }
}

export const dashboardLinks: ComputedRef<DashboardLink[]> = computed(() => {
  const guilds = state.base.self?.userGuilds.edges?.map(({ node }) => node) ?? []
  return [
    {
      name: "Dashboard",
      description: "View your dashboard",
      to: {
        name: "/dashboard/",
      } as RouteLocationNormalized<"/dashboard/">,
      icon: IFasHouse,
    },
    ...guilds.map(guildLink),
    {
      name: "Documentation",
      description: "Our documentation includes guides on proper setup of Spectrograph",
      href: "#",
      icon: IFasBook,
    },
    {
      name: "Service Health",
      description: "Health of all components of the platform",
      to: "/info/service-health",
      icon: IFasHeartPulse,
    },
  ]
})

export const adminDashboardLinks: DashboardLink[] = [
  {
    name: "Guilds",
    description: "List all guilds",
    to: "/admin/guilds",
    icon: IFasServer,
  },
  {
    name: "Guild Events",
    description: "List all guild events",
    to: "/admin/events",
    icon: IFasClock,
  },
  {
    name: "Users",
    description: "List all users",
    to: "/admin/users",
    icon: IFasUsers,
  },
]

export const socialLinks: Link[] = [
  {
    name: "Discord",
    description: "Join our Discord server!",
    href: "https://liam.sh/chat",
    icon: IFabDiscord,
  },
  {
    name: "GitHub",
    description: "Authors Github profile",
    href: "https://github.com/lrstanley",
    icon: IFabGithub,
  },
]
