/**
 * Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
 * of this source code is governed by the MIT license that can be found in
 * the LICENSE file.
 */
import { setupLayouts } from "virtual:generated-layouts"
import { createRouter, createWebHistory } from "vue-router/auto"
import { BaseDocument } from "@/lib/api"
import { client } from "@/lib/api/client"
import { useState } from "@/lib/core/state"
import { loadingBar } from "@/lib/core/status"
import { titleCase } from "@/lib/util"

import type { RouteRecordRaw } from "vue-router/auto"
import type { CombinedError } from "@urql/vue"

function recursiveLayouts(route: RouteRecordRaw): RouteRecordRaw {
  if (route.children) {
    for (let i = 0; i < route.children.length; i++) {
      route.children[i] = recursiveLayouts(route.children[i])
    }

    return route
  }

  return setupLayouts([route])[0]
}

const router = createRouter({
  history: createWebHistory("/"),
  extendRoutes(routes) {
    console.log(routes)
    return routes.map((route) => {
      if (route.path.startsWith("/admin") || route.path.startsWith("/dashboard")) {
        route = {
          ...route,
          meta: {
            auth: true,
            ...route.meta,
          },
        }
      }

      return recursiveLayouts(route)
    })
  },
})

router.beforeEach(async (to, from, next) => {
  const state = useState()

  if (from.name != to.name || JSON.stringify(from.params) != JSON.stringify(to.params)) {
    loadingBar.value = true
  }

  let error: CombinedError

  if (state.base == null || (from.path == "/" && from.name == undefined)) {
    await client
      .query(BaseDocument, {}, { requestPolicy: "network-only" })
      .toPromise()
      .then((resp) => {
        state.base = resp.data

        if (resp.error !== null) {
          error = resp.error
        }
      })
  }

  if (to.meta.auth == true && state.base.self == null) {
    window.location.href = `/-/auth/providers/discord?next=${window.location.origin + to.path}`
    return
  }

  if (
    error !== undefined &&
    !error.graphQLErrors?.some((e) => e.path?.includes("self")) &&
    to.name !== "catchall"
  ) {
    console.log(error)
    next({ name: "catchall", params: { catchall: error.name } })
    return
  }

  next()
})

router.afterEach((to) => {
  const state = useState()

  let title: string

  if (to.meta?.title) {
    title = to.meta.title as string
  } else {
    let args = to.path
      .split("/")
      .reverse()
      .filter((item) => item != "")

    if (args.length > 2) {
      args = args.slice(0, 2)
    }

    // for (let i = 0; i < args.length; i++) {
    //   if (args[i] == "p") {
    //     args[i] = "Posts"
    //   }
    // }

    title = titleCase(args.reverse().join(" · ").replace(/-/g, " "))

    if (title.length < 2) {
      title = "Home"
    }
  }

  document.title = `${title} · Spectrograph`
  state.addToHistory({ title, path: to.path, timestamp: new Date().toISOString() })

  // Scroll to anchor, just in case the page happens to not render fast enough.
  nextTick(() => {
    loadingBar.value = false

    if (location.hash && !to.meta.disableAnchor) {
      const el = document.getElementById(location.hash.slice(1))
      if (el) {
        el.scrollIntoView()
      }
    }
  })
  return
})

export default router
