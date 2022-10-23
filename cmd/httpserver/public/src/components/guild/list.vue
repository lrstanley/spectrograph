<template>
  <div
    v-if="guilds?.length || props.showEmpty"
    class="overflow-hidden border rounded shadow bg-chat-700 border-chat-900"
  >
    <ul v-if="guilds?.length" role="list" class="divide-y divide-chat-800">
      <li v-for="guild in guilds" :key="guild.guild.id">
        <router-link :to="guild.to" class="block hover:bg-chat-800">
          <div class="flex items-center p-4 sm:px-6">
            <div class="flex items-center flex-1 min-w-0">
              <component :is="guild.icon" :status="guild.status" class="w-12 h-12 text-4xl shrink-0" />
              <div class="grid flex-1 min-w-0 grid-cols-1 gap-4 px-4 md:grid-cols-2">
                <div>
                  <p class="text-sm font-medium text-indigo-600 truncate">
                    {{ guild.name }}
                  </p>
                  <p class="flex items-center mt-2 text-sm text-gray-500">
                    <span class="truncate">{{ guild.description }}</span>
                  </p>
                </div>
                <div class="hidden md:block">
                  <div>
                    <p class="text-sm text-bravery-500">STATUS</p>
                    <p class="flex items-center mt-2 text-sm text-gray-500">
                      {{ guild.status }}

                      <GuildStatus
                        class="ml-1.5 h-4 w-4 shrink-0"
                        :status="guild.status"
                        aria-hidden="true"
                      />
                    </p>
                  </div>
                </div>
              </div>
            </div>
            <div>
              <i-fas-chevron-right class="w-5 h-5 text-gray-400" aria-hidden="true" />
            </div>
          </div>
        </router-link>
      </li>
    </ul>
    <div v-else-if="props.showEmpty" class="p-4 text-center text-gray-400">
      <p class="text-lg">You are not an owner or administrator of any Discord guilds.</p>
      <p class="mt-2">
        Create a guild or become a guild administrator, and
        <a class="link link-primary" href="/-/auth/logout">logout</a> and back in to add Spectrograph to
        them here.
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { guildLink } from "@/lib/core/navigation"
import type { Guild } from "@/lib/api"

const props = defineProps<{
  guilds: Guild[]
  showEmpty?: boolean
}>()

const guilds = computed(() => props.guilds.map((guild) => guildLink(guild)))
</script>
