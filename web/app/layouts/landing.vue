<script setup lang="ts">
const { self } = useSelf()
</script>

<template>
  <div class="flex min-h-svh flex-col">
    <header class="sticky top-0 z-20 border-b bg-background/85 backdrop-blur">
      <div class="mx-auto flex h-14 max-w-6xl items-center gap-6 px-4 sm:px-6">
        <BrandMark />
        <nav aria-label="Main" class="hidden items-center gap-5 text-sm text-muted-foreground sm:flex">
          <NuxtLink to="/#how-it-works" class="hover:text-foreground">How it works</NuxtLink>
          <NuxtLink to="/#setup" class="hover:text-foreground">Setup</NuxtLink>
        </nav>
        <div class="ml-auto flex items-center gap-2">
          <ColorModeToggle />
          <Button v-if="self" as-child size="sm">
            <NuxtLink to="/dashboard">Open dashboard</NuxtLink>
          </Button>
          <Button v-else as-child size="sm" variant="outline">
            <NuxtLink to="/login">Log in</NuxtLink>
          </Button>
        </div>
      </div>
    </header>
    <main class="flex-1">
      <slot />
    </main>
    <SiteFooter />
  </div>
</template>
