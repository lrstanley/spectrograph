<script setup lang="ts">
const props = defineProps<{ guildId: number }>()
const { events, history, state, hasMore, loadMore } = useGuildEvents(() => props.guildId)
</script>

<template>
  <GuildEventLog
    :events="events"
    :state="state"
    :loading="history.status.value === 'pending'"
    :error="history.error.value"
    :has-more="hasMore"
    @load-more="loadMore"
    @retry="history.refresh()"
  />
</template>
