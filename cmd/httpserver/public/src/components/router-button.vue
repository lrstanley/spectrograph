<template>
  <router-link v-slot="{ isActive, href, navigate }" custom :to="props.to">
    <BaseButton v-bind="$attrs" :active="isActive" :href="href" @click="navigate">
      <slot />
    </BaseButton>
  </router-link>
</template>

<script setup lang="ts">
import type { RouterLinkProps } from "vue-router/auto"

export interface Props extends RouterLinkProps {
  init?: boolean
}

const props = defineProps<Props>()
</script>
