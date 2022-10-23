<template>
  <div class="sm:col-span-1">
    <dt class="text-sm font-medium text-gray-500">
      <slot name="name">
        {{ props.name }}
      </slot>
    </dt>
    <dd class="mt-1 text-sm text-white">
      <slot>
        {{ props.value }}
      </slot>
    </dd>
  </div>
</template>

<script setup lang="ts">
const props = defineProps<{
  name?: string
  value?: any
}>()
</script>
