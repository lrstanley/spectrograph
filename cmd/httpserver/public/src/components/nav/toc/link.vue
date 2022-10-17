<template>
  <li class="pl-0 truncate">
    <a
      :href="props.href"
      class="text-sm tracking-tighter no-underline link"
      :class="{
        'text-bravery-500/80 hover:text-bravery-600 !text-base': props.depth <= 1,
        'text-nitro-500/80 hover:text-nitro-600': props.depth === 2,
        'text-bravery-400/80 hover:text-bravery-500': props.depth === 3,
        'text-nitro-400/80 hover:text-nitro-500': props.depth === 4,
        'before:border-l-[3px] before:border-l-online-500 before:mr-[calc(1ch-3px)]': active,
        'ml-[1ch]': !active,
      }"
      :title="props.title"
    >
      {{ props.title }}
    </a>

    <ul v-if="$slots.default" class="my-1.5 list-none" :class="props.depth > 1 ? 'pl-3' : ''">
      <slot />
    </ul>
  </li>
</template>

<script setup lang="ts">
import { useIntersectionObserver } from "@vueuse/core"

const props = defineProps<{
  href: string
  title: string
  depth?: number
}>()

const active = ref(false)

console.log(props.href.substring(1))
const { stop } = useIntersectionObserver(
  document.getElementById(props.href.substring(1)),
  ([{ isIntersecting }], observerElement) => {
    active.value = isIntersecting
  }
)
</script>
