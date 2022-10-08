<route lang="yaml">
meta:
  layout: dashboard
  title: Dashboard
</route>

<template>
  <div>
    <h1 class="pb-4 text-2xl font-semibold text-white">Dashboard</h1>

    <div class="overflow-hidden border rounded shadow bg-chat-600 border-chat-900">
      <ul role="list" class="divide-y divide-chat-800">
        <li v-for="guild in guilds" :key="guild.guild.id">
          <router-link :to="guild.to" class="block hover:bg-chat-700">
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
    </div>
  </div>
</template>

<script setup lang="ts">
import { guildLink } from "@/lib/core/navigation"
import type { Guild } from "@/lib/api"

const state = useState()
const guilds = state.base.self?.userGuilds.edges?.map(({ node }) => guildLink(node as Guild)) ?? []
</script>
