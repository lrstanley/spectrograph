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

    <ContainerPanel v-bind="$attrs">
      <template v-if="!!$slots.header" #header>
        <slot name="header" />
      </template>

      <template v-if="props.barePanel">
        <slot />
      </template>
      <template v-else>
        <div class="px-6 py-5 space-y-6">
          <div class="grid grid-cols-2 gap-6">
            <slot />
          </div>
        </div>
      </template>

      <template v-if="!!$slots.actions" #actions>
        <slot name="actions" />
      </template>
    </ContainerPanel>
  </div>
</template>

<script setup lang="ts">
const props = defineProps<{
  title: string
  description: string
  barePanel?: boolean
}>()
</script>
