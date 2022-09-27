<template>
  <div
    class="p-4 rounded-md"
    :class="{
      'bg-online-600': props.status === 'success',
      'bg-discord-500': props.status === 'info',
      'bg-idle-600': props.status === 'warning',
      'bg-dnd-500': props.status === 'error',
    }"
  >
    <div class="flex">
      <div class="shrink-0">
        <i-fas-circle-check v-if="props.status === 'success'" class="w-5 h-5 text-white" />
        <i-fas-circle-info v-if="props.status === 'info'" class="w-5 h-5 text-white" />
        <i-fas-triangle-exclamation v-if="props.status === 'warning'" class="w-5 h-5 text-white" />
        <i-fas-circle-xmark v-if="props.status === 'error'" class="w-5 h-5 text-white" />
      </div>
      <div class="ml-3">
        <h3 v-if="props.title" class="text-sm font-medium">{{ props.title }}</h3>
        <slot />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const props = defineProps<{
  status: "success" | "error" | "warning" | "info"
  title?: string
}>()
</script>

<style scoped></style>
