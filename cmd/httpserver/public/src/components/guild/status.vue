<template>
  <span v-bind="$attrs">
    <i-fas-check v-if="props.status == LinkStatus.Healthy" class="text-balance-500" />
    <i-fas-xmark v-else-if="props.status == LinkStatus.Unhealthy" class="text-dnd-500" />
    <i-fas-circle-plus v-else-if="props.status == LinkStatus.NotJoined" class="text-channel-400" />
  </span>
</template>

<script setup lang="ts">
import { LinkStatus } from "@/lib/core/types"

const props = withDefaults(
  defineProps<{
    status?: LinkStatus | undefined
  }>(),
  { status: undefined }
)
</script>
