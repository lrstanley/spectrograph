<script setup lang="ts">
import logo from "@/assets/img/mic.png"

withDefaults(defineProps<{ to?: string; compact?: boolean }>(), { to: "/", compact: false })
</script>

<template>
  <NuxtLink :to="to" class="group inline-flex items-center gap-2 rounded-sm outline-none focus-visible:ring-2 focus-visible:ring-ring">
    <img :src="logo" alt="" aria-hidden="true" class="size-8 shrink-0 transition-transform duration-300 group-hover:-rotate-12" />
    <span v-if="!compact" class="font-heading text-lg font-semibold tracking-tight">Spectrograph</span>
    <span v-else class="sr-only">Spectrograph</span>
  </NuxtLink>
</template>
