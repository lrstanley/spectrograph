<script setup lang="ts">
import { Moon, Sun } from "@lucide/vue"

const mode = useColorMode({ storageKey: "spectrograph-color-mode", initialValue: "dark" })
const dark = computed(() => mode.value === "dark")
</script>

<template>
  <Button
    variant="ghost"
    size="icon-sm"
    :aria-label="dark ? 'Switch to light theme' : 'Switch to dark theme'"
    :title="dark ? 'Switch to light theme' : 'Switch to dark theme'"
    @click="mode = dark ? 'light' : 'dark'"
  >
    <Sun v-if="dark" />
    <Moon v-else />
  </Button>
</template>
