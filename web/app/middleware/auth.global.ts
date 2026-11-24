export default defineNuxtRouteMiddleware(async (to) => {
  if (import.meta.server) {
    return
  }

  const { load, banned } = useSelf()
  let self: Awaited<ReturnType<typeof load>>
  try {
    self = await load()
  } catch (err) {
    // The API is unreachable or failing; public pages still render.
    if (isProtectedPath(to.path)) {
      return abortNavigation(
        createError({
          statusCode: apiErrorStatus(err) ?? 503,
          statusMessage: "Couldn't reach the Spectrograph API",
          message: apiErrorMessage(err),
          data: { requestID: apiRequestID(err) },
        }),
      )
    }
    return
  }

  if (to.path === "/login") {
    if (self) {
      return navigateTo(safeNext(to.query.next), { replace: true })
    }
    return
  }

  if (isProtectedPath(to.path) && !self) {
    // Banned users are signed out, but logging in again won't help.
    return navigateTo(
      banned.value ? { path: "/login", query: { reason: "banned" } } : { path: "/login", query: { next: to.fullPath } },
      { replace: true },
    )
  }

  if (pathWithin(to.path, "/admin") && !self?.admin) {
    return abortNavigation(
      createError({
        statusCode: 403,
        statusMessage: "Admins only",
        message: "This page is only available to Spectrograph administrators.",
      }),
    )
  }
})
