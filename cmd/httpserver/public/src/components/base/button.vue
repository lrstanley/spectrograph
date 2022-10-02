<template>
  <component
    :is="props.el"
    :class="{
      'inline-flex grow-0 shrink': !props.block,
      'flex justify-center': props.block,
      ...(!props.noPadding
        ? {
            'px-2.5 py-1.5': props.size === 'xs',
            'px-3 py-2': props.size === 'sm',
            'px-4 py-2': props.size === 'md' || props.size === 'lg',
            'px-6 py-3': props.size === 'xl',
          }
        : {}),
      'text-xs': props.size === 'xs',
      'text-sm leading-4': props.size === 'sm',
      'text-sm': props.size === 'md',
      'text-base': props.size === 'lg' || props.size === 'xl',

      'bg-online-700 hover:bg-online-800 disabled:bg-online-900/70': props.type === 'success',
      'bg-discord-500 hover:bg-discord-600 disabled:bg-discord-900/70': props.type === 'info',
      'bg-idle-600 hover:bg-idle-700 disabled:bg-dnd-900/70': props.type === 'warning',
      'bg-dnd-600 hover:bg-dnd-700 disabled:bg-dnd-900/70': props.type === 'error',
      'bg-opacity-70': props.transparent,
    }"
    class="items-center font-medium text-white transition-all duration-100 ease-in-out border border-transparent rounded shadow-sm focus:ring-discord-500 disabled:text-white/70 focus:outline-none focus:ring-2 focus:ring-offset-2 disabled:cursor-not-allowed"
  >
    <div
      v-if="props.icon || props.loading"
      :class="{
        '-ml-0.5 mr-2 h-4 w-4': props.size === 'xs',
        '-ml-1 mr-2 h-5 w-5': props.size === 'sm',
        '-ml-1 mr-3 h-5 w-5': props.size === 'md' || props.size === 'lg' || props.size === 'xl',
      }"
    >
      <i-fas-circle-notch v-if="props.loading" class="animate-spin" />
      <slot v-else name="icon" />
    </div>
    <slot />
  </component>
</template>

<script setup lang="ts">
const props = withDefaults(
  defineProps<{
    icon?: boolean
    el?: string
    loading?: boolean
    noPadding?: boolean
    block?: boolean
    transparent?: boolean
    active?: boolean // TODO: do something with this.
    size?: "xs" | "sm" | "md" | "lg" | "xl"
    type?: "success" | "error" | "warning" | "info" | "default"
  }>(),
  {
    el: "a",
    size: "md",
    type: "default",
  }
)
</script>
