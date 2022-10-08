/**
 * Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
 * of this source code is governed by the MIT license that can be found in
 * the LICENSE file.
 */

import { createClient as createWSClient } from "graphql-ws"
import { loadingBar } from "@/lib/core/status"
import { retryExchange } from "@urql/exchange-retry"
import {
  cacheExchange, createClient, dedupExchange, fetchExchange, subscriptionExchange
} from "@urql/vue"

export * from "@/lib/api/graphql"

function fetchWithTimeout(url: RequestInfo, opts: RequestInit): Promise<Response> {
  loadingBar.value = true
  const controller = new AbortController()
  const id = setTimeout(() => controller.abort(), 5000)

  const promise = new Promise<Response>((resolve, reject) => {
    fetch(url, {
      ...opts,
      signal: controller.signal,
    })
      .then((resp) => {
        clearTimeout(id)
        resolve(resp)
      })
      .catch((err) => {
        resolve(err)
      })
      .finally(() => {
        loadingBar.value = false
      })
  })
  return promise
}

export const wsClient = createWSClient({
  url: `${window.location.protocol == "https" ? "wss" : "ws"}://${window.location.host}/-/graphql`,
  keepAlive: 10000,
})

export const client = createClient({
  url: "/-/graphql",
  requestPolicy: "cache-and-network",
  fetch: fetchWithTimeout,
  exchanges: [
    dedupExchange,
    cacheExchange,
    retryExchange({
      initialDelayMs: 1000,
      maxDelayMs: 15000,
      randomDelay: true,
      maxNumberAttempts: 3,
      retryIf: (err) => (err && err.networkError ? true : false),
    }),
    fetchExchange,
    subscriptionExchange({
      forwardSubscription(operation) {
        return {
          subscribe: (sink) => {
            const dispose = wsClient.subscribe(operation, sink)
            return {
              unsubscribe: dispose,
            }
          },
        }
      },
    }),
  ],
})
