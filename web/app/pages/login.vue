<script setup lang="ts">
import { Bug, RotateCcw } from "@lucide/vue"
import logo from "@/assets/img/mic.png"

definePageMeta({ title: "Log in", layout: "landing" })

const route = useRoute()
const { banned: selfBanned } = useSelf()
const next = computed(() => safeNext(route.query.next))
// "reason" is set by the backend when the Discord login fails, and by the auth
// middleware for banned users.
const banned = computed(() => route.query.reason === "banned" || selfBanned.value)
const failed = computed(() => !banned.value && route.query.reason === "error")
// Only a request ID is accepted (it's put into bug reports), never free text.
const requestID = computed(() => {
  const value = route.query.request_id
  return typeof value === "string" && /^[\w.:/-]{1,128}$/.test(value) ? value : undefined
})

const bugURL = useBugReportURL(() => ({
  title: "Discord sign-in failed",
  status: undefined,
  errorType: undefined,
  timestamp: undefined,
  requestID: requestID.value,
}))
</script>

<template>
  <div class="mx-auto flex max-w-md flex-col items-center gap-6 px-4 py-24 text-center">
    <img :src="logo" alt="" class="size-16" />
    <div class="flex flex-col gap-2">
      <h1 class="font-heading text-3xl font-semibold tracking-tight">Log in to Spectrograph</h1>
      <p class="text-muted-foreground">
        Sign in with Discord to configure Spectrograph in the servers you manage. We read your profile, email, and server
        list; we never see your password.
      </p>
    </div>
    <Alert v-if="banned" variant="destructive" class="text-left">
      <AlertTitle>This account is banned</AlertTitle>
      <AlertDescription>
        Your account can't use Spectrograph. If you think this is a mistake, email questions@spectrograph.app.
      </AlertDescription>
    </Alert>
    <Alert v-else-if="failed" variant="destructive" class="text-left" role="alert">
      <AlertTitle>Sign-in failed</AlertTitle>
      <AlertDescription>
        <p>
          Discord didn't complete the sign-in. This happens if you cancel on Discord's page, or if the sign-in took
          too long. Try again, and if it keeps failing, report it.
        </p>
        <p v-if="requestID" class="text-xs">
          Request ID <code class="font-mono text-foreground">{{ requestID }}</code>
        </p>
        <a :href="bugURL" target="_blank" rel="noopener" class="inline-flex items-center gap-1 underline underline-offset-4">
          <Bug class="size-3.5" aria-hidden="true" />
          Report this problem
        </a>
        <p class="text-xs">Reports are public, so leave out anything personal.</p>
      </AlertDescription>
    </Alert>
    <Button as-child size="lg" class="w-full">
      <a :href="loginURL(next)">
        <RotateCcw v-if="failed" data-icon="inline-start" />
        {{ failed ? "Try again with Discord" : "Continue with Discord" }}
      </a>
    </Button>
    <p class="text-xs text-muted-foreground">
      By continuing you agree to the
      <NuxtLink to="/info/terms" class="underline underline-offset-4">terms of service</NuxtLink>
      and
      <NuxtLink to="/info/privacy" class="underline underline-offset-4">privacy policy</NuxtLink>.
    </p>
  </div>
</template>
