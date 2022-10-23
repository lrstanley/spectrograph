<template>
  <ContainerPanel bare-panel>
    <template v-if="props.title || props.subtitle" #header>
      <h3 class="text-lg font-medium leading-6 text-white">
        {{ props.title }}
      </h3>
      <p v-if="props.subtitle" class="max-w-2xl mt-1 text-sm text-gray-500">
        {{ props.subtitle }}
      </p>
    </template>

    <div class="px-6 py-5 border-chat-900" :class="(props.title || props.subtitle) && 'border-t'">
      <dl class="grid grid-cols-1 gap-x-4 gap-y-8 md:grid-cols-2 lg:grid-cols-3">
        <slot />
      </dl>
    </div>
  </ContainerPanel>
</template>

<script setup lang="ts">
const props = defineProps<{
  title?: string
  subtitle?: string
}>()
</script>
