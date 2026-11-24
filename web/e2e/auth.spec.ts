import { expect, test } from "./fixtures"

test("failed Discord sign-ins land on the login page with a retry", async ({ page, baseURL }) => {
  // Without a pending OAuth2 flow (no state cookie), the callback fails.
  await page.goto("/-/auth/providers/discord/callback?error=access_denied&state=nope")
  await expect(page).toHaveURL(/\/login\?reason=error&request_id=/)

  const alert = page.getByRole("alert").filter({ hasText: "Sign-in failed" })
  await expect(alert).toBeVisible()
  await expect(page.getByRole("link", { name: "Try again with Discord" })).toHaveAttribute(
    "href",
    `/-/auth/providers/discord?next=${encodeURIComponent(`${baseURL}/dashboard`)}`,
  )

  const requestID = new URL(page.url()).searchParams.get("request_id") ?? ""
  expect(requestID).not.toBe("")
  await expect(alert).toContainText(requestID)

  const href = await alert.getByRole("link", { name: "Report this problem" }).getAttribute("href")
  const url = new URL(href ?? "")
  expect(url.origin + url.pathname).toBe("https://github.com/lrstanley/spectrograph/issues/new")
  expect(url.searchParams.get("context")).toContain(`Request ID: ${requestID}`)
  expect(url.searchParams.get("environment")).toContain("Page: /login")
})

test("login ignores request IDs that aren't IDs", async ({ page }) => {
  await page.goto(`/login?reason=error&request_id=${encodeURIComponent("Jane Doe <jane@example.com>")}`)
  const alert = page.getByRole("alert").filter({ hasText: "Sign-in failed" })
  await expect(alert).toBeVisible()
  await expect(alert).not.toContainText("jane")
  const href = await alert.getByRole("link", { name: "Report this problem" }).getAttribute("href")
  expect(decodeURIComponent(href ?? "")).not.toContain("jane")
})

test("banned users see why they can't log in", async ({ page }) => {
  await page.goto("/login?reason=banned")
  await expect(page.getByText("This account is banned")).toBeVisible()
  await expect(page.getByText("Sign-in failed")).toHaveCount(0)
})

test("static assets are cached, and HTML is revalidated", async ({ page, request }) => {
  const assets: string[] = []
  page.on("response", (res) => {
    if (new URL(res.url()).pathname.startsWith("/_nuxt/")) assets.push(res.url())
  })
  await page.goto("/")
  await expect(page.getByRole("heading", { level: 1 })).toBeVisible()
  expect(assets.length).toBeGreaterThan(0)

  const asset = await request.get(assets[0] ?? "")
  expect(asset.headers()["cache-control"]).toBe("public, max-age=31536000, immutable")

  for (const path of ["/", "/dashboard/guild/1"]) {
    const html = await request.get(path)
    expect(html.status()).toBe(200)
    expect(html.headers()["cache-control"]).toBe("no-cache")
    expect(html.headers().etag).toBeTruthy()
  }

  const head = await request.head("/")
  expect(head.status()).toBe(200)
})
