/**
 * Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
 * of this source code is governed by the MIT license that can be found in
 * the LICENSE file.
 */

import type { Component } from "vue"
import type { RouteNamedMap } from "vue-router/auto/routes"
import type { RouteLocationNormalized } from "vue-router/auto"
import type { Guild } from "@/lib/api"

export interface Link {
  name: string
  description: string
  href?: string
  to?: keyof RouteNamedMap | RouteLocationNormalized
  icon: Component
}

export interface DashboardLink extends Link {
  guild?: Guild
  isChild?: boolean
  hasStatus?: boolean
  status?: LinkStatus
  hover?: boolean
}

export enum LinkStatus {
  Healthy = "healthy",
  Unhealthy = "unhealthy",
  NotJoined = "not-joined",
}
