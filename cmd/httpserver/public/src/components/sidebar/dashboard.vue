<template>
  <div>
    <!-- Sidebar component, swap this element with another sidebar if you like -->
    <div class="flex flex-col flex-auto min-h-0 bg-channel-500">
      <div class="flex flex-col flex-auto pt-5 pb-4 overflow-y-auto">
        <a href="#" class="inline-flex items-center px-4 shrink-0">
          <img class="w-auto h-14 md:h-12" src="/img/mic.png" />
          <h2 class="text-3xl md:text-xl text-gradient bg-gradient-to-r from-nitro-500 to-dnd-400">
            Spectrograph
          </h2>
        </a>

        <nav class="flex flex-col flex-auto px-4 mt-5 space-y-2 md:space-y-1">
          <template v-for="item in dashboardLinks" :key="item.name">
            <router-link
              :to="item.to || item.href"
              active-class="text-white bg-channel-900/50"
              class="flex items-center p-3 text-lg font-medium transition-all duration-100 ease-in rounded text-channel-300 hover:bg-channel-900/50 hover:text-white group md:py-2 md:text-sm"
              :class="item.isChild ? 'ml-3' : ''"
            >
              <component
                :is="item.icon"
                class="w-5 h-5 mr-3 text-gray-400 shrink-0 group-hover:text-gray-300"
                aria-hidden="true"
              />
              {{ item.name }}

              <span v-if="item.hasStatus" class="pl-1 ml-auto">
                <i-fas-check v-if="item.status == LinkStatus.Healthy" class="text-online-500" />
                <i-fas-xmark v-else-if="item.status == LinkStatus.Unhealthy" class="text-dnd-500" />
                <i-fas-circle-plus
                  v-else-if="item.status == LinkStatus.NotJoined"
                  class="text-channel-400"
                />
              </span>
            </router-link>
          </template>
        </nav>
      </div>
      <div class="flex p-4 grow-0 shrink-0 bg-channel-700">
        <a href="#" class="block w-full group shrink-0">
          <div class="flex items-center">
            <div>
              <img
                class="inline-block w-12 h-12 rounded-full md:h-9 md:w-9"
                :src="state.base.self.avatarURL"
              />
            </div>
            <div class="ml-3">
              <span class="inline-flex flex-auto text-lg font-medium text-balance-500 md:text-sm">
                {{ state.base.self.username }}
                <span class="text-channel-300">#{{ state.base.self.discriminator }}</span>
              </span>
              <p class="text-lg font-medium text-gray-300 md:text-xs group-hover:text-gray-200">
                View profile
              </p>
            </div>
          </div>
        </a>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { dashboardLinks, LinkStatus } from "@/lib/core/navigation"

const state = useState()

const guilds = state.base.self?.userGuilds.edges?.map(({ node }) => node)
</script>
