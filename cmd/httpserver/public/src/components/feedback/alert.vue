<template>
  <div
    class="flex flex-row items-center justify-start rounded alert"
    :class="{
      'alert-success': props.type === 'success',
      'alert-info': props.type === 'info',
      'alert-warning': props.type === 'warning',
      'alert-error': props.type === 'error',
    }"
  >
    <i-fas-circle-check v-if="props.type === 'success'" class="w-5 h-5 text-white" />
    <i-fas-circle-info v-if="props.type === 'info'" class="w-5 h-5 text-white" />
    <i-fas-triangle-exclamation v-if="props.type === 'warning'" class="w-5 h-5 text-white" />
    <i-fas-circle-xmark v-if="props.type === 'error'" class="w-5 h-5 text-white" />

    <span class="!mt-0"><slot /></span>
  </div>
</template>

<script setup lang="ts">
const props = defineProps<{
  type: "success" | "error" | "warning" | "info"
}>()
</script>
