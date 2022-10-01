/**
 * Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
 * of this source code is governed by the MIT license that can be found in
 * the LICENSE file.
 */
import { h } from "vue"

import type { Component, ComputedRef } from "vue"
import type { RouteNamedMap } from "vue-router/auto/routes"
import type { RouteLocationNormalized } from "vue-router/auto"
import type { Guild } from "@/lib/api"

const state = useState()

export const headerLinks = [
  { name: "Features", href: "/#more-info", current: true },
  { name: "Documentation", href: "#", current: false },
  { name: "Service Health", href: "#", current: false },
]

export interface Link {
  name: string
  description: string
  href?: string
  to?: keyof RouteNamedMap | RouteLocationNormalized
  icon: Component
}

export interface DashboardLink extends Link {
  isChild?: boolean
  hasStatus?: boolean
  status?: LinkStatus
}

export enum LinkStatus {
  Healthy = "healthy",
  Unhealthy = "unhealthy",
  NotJoined = "not-joined",
}

export const mainLinks: Link[] = [
  {
    name: "Documentation",
    description: "Our documentation includes guides on proper setup of Spectrograph",
    href: "#",
    icon: IconFasBook,
  },
  {
    name: "Service Health",
    description: "Health of all components of the platform",
    to: "/info/service-health",
    icon: IconFasHeartPulse,
  },
  {
    name: "Privacy",
    description: "Spetrograph's Privacy Policy",
    to: "/info/privacy",
    icon: IconFasFingerprint,
  },
  {
    name: "Terms",
    description: "Spectrograph's Terms of Service",
    to: "/info/terms",
    icon: IconFasBuildingColumns,
  },
]

function guildIcon(name: string, url: string): Component {
  if (url) {
    return h("img", { src: url, alt: `${name}'s Guild Icon`, class: "rounded-full" })
  }

  return h(
    "span",
    { class: "inline-flex items-center justify-center rounded-full bg-channel-400 shrink-0" },
    [h("span", { class: "text-sm font-medium leading-none text-white capitalize" }, name[0])]
  )
}

function guildLink(guild: Guild): DashboardLink {
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
    name: guild.name,
    description: `View ${guild.name}'s dashboard`,
    to: to,
    icon: guildIcon(guild.name, guild.iconURL),
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
      icon: IconFasHouse,
    },
    ...guilds.map(guildLink),
    {
      name: "Service Health",
      description: "Health of all components of the platform",
      to: "/info/service-health",
      icon: IconFasHeartPulse,
    },
  ]
})

export const socialLinks: Link[] = [
  {
    name: "Discord",
    description: "Join our Discord server!",
    href: "https://liam.sh/chat",
    icon: IconFabDiscord,
  },
  {
    name: "GitHub",
    description: "Authors Github profile",
    href: "https://github.com/lrstanley",
    icon: IconFabGithub,
  },
]
