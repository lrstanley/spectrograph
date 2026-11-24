<script setup lang="ts">
import { Toaster } from "@/components/ui/sonner"
import "vue-sonner/style.css"

const route = useRoute()
const loading = useLoadingIndicator()
const nuxtApp = useNuxtApp()
const lastError = useLastApiError()
const { markBanned } = useSelf()
const title = computed(() => String(route.meta.title ?? "Spectrograph"))

// Applies the "dark"/"light" class on <html>. Dark is the default.
useColorMode({ storageKey: "spectrograph-color-mode", initialValue: "dark" })

useHead({
  title,
  titleTemplate: (value) => (value === "Spectrograph" ? value : `${value} | Spectrograph`),
})

nuxtApp.hook("openFetch:onRequest", () => {
  if (loading.isLoading.value) return
  loading.start()
})
// Network failures (e.g. the API is down) get no response.
nuxtApp.hook("openFetch:onRequestError", () => {
  loading.finish({ error: true })
})
nuxtApp.hook("openFetch:onResponse", (ctx) => {
  if (ctx.response.status < 300) {
    loading.finish()
    return
  }
  loading.finish({ error: true })
  // Signed-out checks of /self aren't errors worth reporting.
  if (ctx.response.status === 401) return
  const body = (ctx.response._data ?? {}) as ApiErrorBody
  const url = typeof ctx.request === "string" ? ctx.request : ctx.request.url
  // Banned mid-session: sign out and explain, from any page that needs a
  // login. The auth middleware handles the initial /self check itself.
  if (isBannedError({ status: ctx.response.status, data: body }) && !/\/self$/.test(url) && markBanned()) {
    if (isProtectedPath(route.path)) {
      navigateTo({ path: "/login", query: { reason: "banned" } }, { replace: true })
    }
    return
  }
  lastError.value = {
    status: ctx.response.status,
    requestID: body.request_id,
    at: body.timestamp ?? new Date().toISOString(),
  }
})
</script>

<template>
  <NuxtRouteAnnouncer />
  <NuxtLoadingIndicator color="var(--primary)" />
  <Toaster position="bottom-right" />
  <NuxtLayout>
    <NuxtPage />
  </NuxtLayout>
</template>
