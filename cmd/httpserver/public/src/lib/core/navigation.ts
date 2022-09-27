/**
 * Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
 * of this source code is governed by the MIT license that can be found in
 * the LICENSE file.
 */

import type { Component } from "vue"

export const headerLinks = [
  { name: "Features", href: "#", current: true },
  { name: "Documentation", href: "#", current: false },
  { name: "Service Health", href: "#", current: false },
]

export interface Link {
  name: string
  description: string
  href?: string
  to?: Record<string, any>
  icon: Component
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
    href: "#",
    icon: IconFasHeartPulse,
  },
  { name: "Privacy", description: "Spetrograph's Privacy Policy", href: "#", icon: IconFasFingerprint },
  {
    name: "Terms",
    description: "Spectrograph's Terms of Service",
    href: "#",
    icon: IconFasBuildingColumns,
  },
]

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
