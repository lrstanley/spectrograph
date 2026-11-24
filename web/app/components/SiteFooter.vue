<script setup lang="ts">
const config = useRuntimeConfig()
const bugURL = useBugReportURL()
const year = new Date().getFullYear()
</script>

<template>
  <footer class="border-t">
    <div class="mx-auto flex max-w-6xl flex-col gap-8 px-4 py-10 sm:px-6 md:flex-row md:items-start md:justify-between">
      <div class="flex max-w-sm flex-col gap-3">
        <BrandMark />
        <p class="text-muted-foreground text-sm">
          Keeps one empty voice channel open in every group, so nobody has to ask a moderator for another.
        </p>
      </div>
      <nav aria-label="Footer" class="grid grid-cols-2 gap-x-12 gap-y-2 text-sm">
        <NuxtLink to="/info/privacy" class="text-muted-foreground hover:text-foreground">Privacy policy</NuxtLink>
        <a :href="repoURL" target="_blank" rel="noopener" class="text-muted-foreground hover:text-foreground">Source code</a>
        <NuxtLink to="/info/terms" class="text-muted-foreground hover:text-foreground">Terms of service</NuxtLink>
        <a :href="issuesURL" target="_blank" rel="noopener" class="text-muted-foreground hover:text-foreground">Known issues</a>
        <a href="https://liam.sh/chat" target="_blank" rel="noopener" class="text-muted-foreground hover:text-foreground">Discord chat</a>
        <a :href="bugURL" target="_blank" rel="noopener" class="text-muted-foreground hover:text-foreground">Report an issue</a>
      </nav>
    </div>
    <div class="mx-auto flex max-w-6xl flex-wrap items-center justify-between gap-2 px-4 pb-8 text-xs text-muted-foreground sm:px-6">
      <a href="https://liam.sh" target="_blank" rel="noopener" class="hover:text-foreground">&copy; {{ year }} Liam Stanley</a>
      <span class="font-mono">{{ config.public.version }}</span>
    </div>
  </footer>
</template>
