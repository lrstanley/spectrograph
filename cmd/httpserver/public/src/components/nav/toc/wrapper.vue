<template>
  <ul class="pl-0 my-3 list-none">
    <component :is="item" v-for="(item, index) in props.items" :key="index" :depth="1" />
  </ul>
</template>

<script setup lang="ts">
import type { VNode } from "vue"

const props = defineProps<{
  items: VNode[]
}>()
</script>
