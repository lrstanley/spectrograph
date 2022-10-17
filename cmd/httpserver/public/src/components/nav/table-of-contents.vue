<template>
  <div v-if="props.element && links.length">
    <div class="text-lg uppercase text-balance-500">Table of Contents</div>
    <NavTocWrapper :items="links" />
  </div>
</template>

<script setup lang="ts">
import Link from "@/components/nav/toc/link.vue"
import { createTOC } from "@/lib/util"
const props = defineProps<{
  element: Node
}>()

const links = computed(() => (props.element ? createTOC(props.element, Link) : []))
</script>
