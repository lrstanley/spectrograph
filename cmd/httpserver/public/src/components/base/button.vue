<template>
  <component
    :is="props.el"
    :class="{
      'inline-flex grow-0 shrink': !props.block,
      'flex justify-center': props.block,
      'px-2.5 py-1.5': !props.noPadding && props.size === 'xs',
      'px-3 py-2': !props.noPadding && props.size === 'sm',
      'px-4 py-2': !props.noPadding && (props.size === 'md' || props.size === 'lg'),
      'px-6 py-3': !props.noPadding && props.size === 'xl',

      'text-xs': props.size === 'xs',
      'text-sm leading-4': props.size === 'sm',
      'text-sm': props.size === 'md',
      'text-base': props.size === 'lg' || props.size === 'xl',

      'bg-online-600 hover:bg-online-700 focus:ring-online-500': props.type === 'success',
      'bg-discord-500 hover:bg-discord-600 focus:ring-discord-400': props.type === 'info',
      'bg-idle-600 hover:bg-idle-700 focus:ring-idle-500': props.type === 'warning',
      'bg-dnd-500 hover:bg-dnd-600 focus:ring-dnd-400': props.type === 'error',
      'bg-opacity-70': props.transparent,
    }"
    class="items-center font-medium transition-all duration-100 ease-in-out border border-transparent rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-offset-2"
  >
    <div
      v-if="props.icon"
      :class="{
        '-ml-0.5 mr-2 h-4 w-4': props.size === 'xs',
        '-ml-1 mr-2 h-5 w-5': props.size === 'sm',
        '-ml-1 mr-3 h-5 w-5': props.size === 'md' || props.size === 'lg' || props.size === 'xl',
      }"
    >
      <slot name="icon" />
    </div>
    <slot />
  </component>
</template>

<script setup lang="ts">
const props = withDefaults(
  defineProps<{
    icon?: boolean
    el?: string
    noPadding?: boolean
    block?: boolean
    transparent?: boolean
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
