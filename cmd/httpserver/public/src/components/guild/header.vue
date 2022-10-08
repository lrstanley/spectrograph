<template>
  <div class="md:flex md:items-center md:justify-between md:space-x-5">
    <div class="flex items-start space-x-5">
      <div class="shrink-0">
        <div class="relative">
          <GuildIcon
            :guild="guild"
            class="w-16 h-16 text-4xl border-2 border-solid"
            :class="enabled ? 'border-balance-500' : 'border-dnd-400'"
            size="2xl"
          />
        </div>
      </div>

      <div class="pt-1.5">
        <h1 class="text-2xl font-bold text-white">Guild: {{ guild.name }}</h1>
        <span class="badge badge-info">joined: {{ joinedAt }}</span>
      </div>
    </div>
    <div
      class="flex flex-col-reverse mt-6 space-y-4 space-y-reverse justify-stretch sm:flex-row-reverse sm:justify-end sm:space-y-0 sm:space-x-3 sm:space-x-reverse md:mt-0 md:flex-row md:space-x-3"
    >
      <button
        class="h-10 min-h-0 btn"
        :class="guild.guildConfig?.enabled ? 'btn-error' : 'btn-success'"
        :disabled="!guild.guildAdminConfig.enabled"
        @click="emit('toggleEnabled')"
      >
        <i-fas-circle-xmark v-if="guild.guildConfig?.enabled" class="h-4 mr-2" />
        <i-fas-circle-check v-else class="h-4 mr-2" />
        {{ guild.guildConfig?.enabled ? "Disable" : "Enable" }}
      </button>
    </div>
  </div>

  <FeedbackAlert v-if="!guild.guildAdminConfig.enabled" type="error" class="mt-10">
    Guild has been disabled by an administrator
    <span v-if="guild.guildAdminConfig.comment" class="rounded badge">
      reason: {{ guild.guildAdminConfig.comment }}
    </span>
  </FeedbackAlert>
</template>

<script setup lang="ts">
import type { Guild } from "@/lib/api"

const props = defineProps<{
  guild: Guild
}>()

const emit = defineEmits<{
  (e: "toggleEnabled"): void
}>()

const joinedAt = useTimeAgo(computed(() => props.guild.joinedAt))
const enabled = computed(() => props.guild.guildConfig?.enabled && props.guild.guildAdminConfig?.enabled)
</script>
