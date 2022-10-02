/**
 * Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
 * of this source code is governed by the MIT license that can be found in
 * the LICENSE file.
 */

import { defineStore } from "pinia"
import { BaseDocument } from "@/lib/api"
import { client } from "@/lib/api/client"
import { useStorage } from "@vueuse/core"

import type { CombinedError } from "@urql/vue"
import type { BaseQuery } from "@/lib/api"

export interface History {
  title: string
  path: string
  timestamp: string
}

export interface State {
  base: BaseQuery | null
  history: History[]
  sidebarCollapsed: boolean
}

export const useState = defineStore("state", {
  state: () => {
    return useStorage("state", {
      base: null,
      history: [],
      sidebarCollapsed: false,
    } as State)
  },
  actions: {
    async fetchBase(): Promise<CombinedError | null> {
      let error: CombinedError

      await client
        .query<BaseQuery>(BaseDocument, {}, { requestPolicy: "network-only" })
        .toPromise()
        .then((resp) => {
          Object.assign(this.base, resp.data)

          if (resp.error !== null) {
            error = resp.error
          }
        })

      return error ?? null
    },
    addToHistory(item: History) {
      // Truncate to a max size.
      if (this.history.length > 4) {
        this.history.shift()
      }

      // Remove any previous duplicates with the exact same path.
      for (let i = this.history.length - 1; i >= 0; i--) {
        if (this.history[i].path === item.path) {
          this.history.splice(i, 1)
        }
      }

      this.history.push(item)
    },
  },
})
