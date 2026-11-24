<script setup lang="ts">
import { Bug, House, LogIn } from "@lucide/vue"
import type { NuxtError } from "#app"

const props = defineProps<{ error: NuxtError }>()

const route = useRoute()
const lastError = useLastApiError()

const status = computed(() => props.error.statusCode || 500)
const requestID = computed(() => {
  const data = props.error.data as { requestID?: string } | string | undefined
  return (typeof data === "object" ? data?.requestID : undefined) ?? lastError.value?.requestID
})

const copy = computed(() => {
  switch (status.value) {
    case 401:
      return { title: "Log in to continue", body: "This page needs you to be logged in with Discord." }
    case 403:
      return { title: props.error.statusMessage || "You can't view this page", body: props.error.message }
    case 404:
      return {
        title: props.error.statusMessage?.startsWith("Page not found")
          ? "Page not found"
          : props.error.statusMessage || "Page not found",
        body: props.error.statusMessage?.startsWith("Page not found")
          ? "There's nothing at this address. Check the link, or start again from the home page."
          : props.error.message,
      }
    default:
      return {
        title: props.error.statusMessage || "Something went wrong",
        body:
          props.error.message && props.error.message !== props.error.statusMessage
            ? props.error.message
            : "Try again in a moment. If it keeps happening, report it so it can be fixed.",
      }
  }
})

const bugURL = useBugReportURL(() => ({
  title: `${status.value} on ${redactPath(route.path)}`,
  status: status.value,
  errorType: copy.value.title,
  requestID: requestID.value,
}))

useHead({ title: () => copy.value.title })
</script>

<template>
  <div class="flex min-h-svh flex-col">
    <header class="mx-auto flex h-14 w-full max-w-6xl items-center px-4 sm:px-6">
      <BrandMark />
    </header>
    <main class="mx-auto flex w-full max-w-2xl flex-1 flex-col justify-center gap-8 px-4 py-16 sm:px-6">
      <div class="flex flex-col gap-3">
        <p class="font-heading text-7xl font-bold tracking-tighter text-primary tabular-nums sm:text-8xl" aria-hidden="true">
          {{ status }}
        </p>
        <h1 class="font-heading text-3xl font-semibold tracking-tight">{{ copy.title }}</h1>
        <p v-if="copy.body" class="max-w-prose text-muted-foreground">{{ copy.body }}</p>
        <p v-if="requestID" class="text-xs text-muted-foreground">
          Request ID <code class="font-mono text-foreground">{{ requestID }}</code>
        </p>
      </div>
      <div class="flex flex-wrap gap-2">
        <Button v-if="status === 401" @click="clearError({ redirect: `/login?next=${encodeURIComponent(route.fullPath)}` })">
          <LogIn data-icon="inline-start" />
          Log in
        </Button>
        <Button :variant="status === 401 ? 'outline' : 'default'" @click="clearError({ redirect: '/' })">
          <House data-icon="inline-start" />
          Go to the home page
        </Button>
        <Button as-child variant="outline">
          <a :href="bugURL" target="_blank" rel="noopener">
            <Bug data-icon="inline-start" />
            Report this problem
          </a>
        </Button>
      </div>
      <p class="text-sm text-muted-foreground">
        Before reporting, check the
        <a :href="issuesURL" target="_blank" rel="noopener" class="text-foreground underline underline-offset-4">known issues</a>.
        Reports are public, so leave out anything personal.
      </p>
    </main>
  </div>
</template>
