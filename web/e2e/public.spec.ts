import { expect, test } from "./fixtures"

test("landing page links to the bot invite", async ({ page }) => {
  await page.goto("/")
  await expect(page.getByRole("heading", { level: 1 })).toHaveText("Voice channels that make room for everyone.")
  await expect(page.getByRole("link", { name: "Add to Discord" }).first()).toHaveAttribute("href", "/-/invite")
})

test("legal pages render", async ({ page }) => {
  await page.goto("/info/privacy")
  await expect(page.getByRole("heading", { name: "Privacy Policy", level: 1 })).toBeVisible()
  await page.goto("/info/terms")
  await expect(page.getByRole("heading", { name: "Terms of Service", level: 1 })).toBeVisible()
})

test("dashboard asks signed-out users to log in with Discord", async ({ page, baseURL }) => {
  await page.goto("/dashboard")
  await expect(page).toHaveURL(/\/login\?next=/)
  await expect(page.getByRole("link", { name: "Continue with Discord" })).toHaveAttribute(
    "href",
    `/-/auth/providers/discord?next=${encodeURIComponent(`${baseURL}/dashboard`)}`,
  )
})

test("admin pages require a login", async ({ page }) => {
  await page.goto("/admin/users")
  await expect(page).toHaveURL(/\/login\?next=/)
})

test("unknown pages offer a bug report without identifiers", async ({ page }) => {
  await page.goto("/nope/42")
  await expect(page.getByRole("heading", { name: "Page not found" })).toBeVisible()
  const href = await page.getByRole("link", { name: "Report this problem" }).getAttribute("href")
  const url = new URL(href ?? "")
  expect(url.origin + url.pathname).toBe("https://github.com/lrstanley/spectrograph/issues/new")
  expect(url.searchParams.get("template")).toBe("bug_report.yml")
  expect(url.searchParams.get("environment")).toContain("Page: /nope/:id")
  expect(href).not.toContain("42")
})

test("login ignores off-site next URLs", async ({ page, baseURL }) => {
  await page.goto("/login?next=//example.org/dashboard")
  await expect(page.getByRole("link", { name: "Continue with Discord" })).toHaveAttribute(
    "href",
    `/-/auth/providers/discord?next=${encodeURIComponent(`${baseURL}/dashboard`)}`,
  )
})
