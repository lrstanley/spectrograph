<script setup lang="ts">
import type { HTMLAttributes } from "vue"
import type { SchemaUser } from "#open-fetch-schemas/api"

defineProps<{
  user: Pick<SchemaUser, "username" | "global_name" | "avatar_url">
  class?: HTMLAttributes["class"]
}>()
</script>

<template>
  <Avatar :class="$props.class">
    <AvatarImage :src="user.avatar_url" alt="" />
    <AvatarFallback>{{ initials(displayName(user)) }}</AvatarFallback>
  </Avatar>
</template>
