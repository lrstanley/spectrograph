<script setup lang="ts">
import type { HTMLAttributes } from "vue"
import type { SchemaGuild } from "#open-fetch-schemas/api"
import { cn } from "@/lib/utils"

const props = defineProps<{
  guild: Pick<SchemaGuild, "name" | "icon_url">
  status?: GuildStatus
  // Show only the first initial (for tiny icons).
  compact?: boolean
  class?: HTMLAttributes["class"]
}>()
</script>

<template>
  <Avatar
    :class="
      cn(
        'ring-2 ring-offset-2 ring-offset-background',
        status === 'active' && 'ring-success',
        (status === 'paused' || status === 'admin-disabled') && 'ring-destructive',
        (!status || status === 'not-joined') && 'ring-transparent',
        props.class,
      )
    "
  >
    <AvatarImage v-if="guild.icon_url" :src="guild.icon_url" :alt="''" />
    <AvatarFallback class="font-semibold">{{ compact ? initials(guild.name).slice(0, 1) : initials(guild.name) }}</AvatarFallback>
  </Avatar>
</template>
