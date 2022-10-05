<template>
  <div class="grid grid-cols-1 gap-3 md:grid-cols-3 md:gap-6">
    <div class="md:col-span-1">
      <div class="px-4 sm:px-0">
        <h3 class="text-lg font-medium leading-6 text-white">{{ props.title }}</h3>
        <p class="mt-1 text-sm text-channel-200">
          {{ props.description }}
        </p>
      </div>
    </div>

    <ContainerPanel>
      <slot />

      <template #actions>
        <slot name="actions" />
      </template>
    </ContainerPanel>
  </div>
</template>

<script setup lang="ts">
const props = defineProps<{
  title: string
  description: string
}>()
</script>
