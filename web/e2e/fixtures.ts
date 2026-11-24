import { test as base, expect } from "@playwright/test"
import { type SpectrographServer, startSpectrograph } from "./server"

// Same variable the Go database tests use. Tests are skipped without it.
const databaseURL = process.env.TEST_DATABASE_URL ?? process.env.DATABASE_URL ?? ""

type Worker = { server: SpectrographServer }

export const test = base.extend<object, Worker>({
  server: [
    // biome-ignore lint/correctness/noEmptyPattern: Playwright requires an object pattern for fixture arguments.
    async ({}, use) => {
      const server = await startSpectrograph(databaseURL)
      await use(server)
      await server.stop()
    },
    { scope: "worker", timeout: 200_000 },
  ],
  baseURL: async ({ server }, use) => {
    await use(server.url)
  },
})

test.skip(!databaseURL, "TEST_DATABASE_URL (or DATABASE_URL) is not set")

export { expect }
