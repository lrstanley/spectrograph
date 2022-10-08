<template>
  <img
    v-if="props.guild.iconURL"
    class="rounded-full"
    :class="{
      'border-2 border-solid border-balance-500': props.status === LinkStatus.Healthy,
      'border-2 border-solid border-dnd-400': props.status === LinkStatus.Unhealthy,
    }"
    :alt="guild.name + `'s Guild Icon`"
    :src="props.guild.iconURL"
    v-bind="$attrs"
  />
  <span
    v-else
    class="inline-flex items-center justify-center rounded-full bg-channel-400 shrink-0"
    :class="{
      'border-2 border-solid border-balance-500': props.status === LinkStatus.Healthy,
      'border-2 border-solid border-dnd-400': props.status === LinkStatus.Unhealthy,
    }"
    v-bind="$attrs"
  >
    <span
      class="font-medium leading-none text-white capitalize"
      :class="{
        'text-sm': props.size === 'sm',
        'text-md': props.size === 'md',
        'text-lg': props.size === 'lg',
        'text-xl': props.size === 'xl',
        'text-2xl': props.size === '2xl',
      }"
    >
      {{ props.guild.name[0] }}
    </span>
  </span>
</template>

<script setup lang="ts">
import { LinkStatus } from "@/lib/core/navigation"
import type { Guild } from "@/lib/api"

const props = withDefaults(
  defineProps<{
    guild: Guild
    size?: "sm" | "md" | "lg" | "xl" | "2xl"
    status?: LinkStatus | undefined
  }>(),
  { size: "sm", status: undefined }
)
</script>
