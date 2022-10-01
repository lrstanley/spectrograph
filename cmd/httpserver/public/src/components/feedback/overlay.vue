<template>
  <div class="relative">
    <div
      v-if="props.loading"
      class="absolute z-10 flex items-center justify-center w-full h-full transition-all duration-200 rounded bg-chat-500/60"
    >
      <div class="flex items-center">
        <i-fas-circle-notch class="h-12 text-4xl align-middle animate-spin text-discord-500" />
      </div>
    </div>

    <slot />
  </div>
</template>

<script setup lang="ts">
const props = defineProps<{
  loading: boolean
}>()
</script>
