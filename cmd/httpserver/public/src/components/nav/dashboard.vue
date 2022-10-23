<template>
  <div class="flex flex-col flex-auto min-h-0 bg-channel-500">
    <div class="flex flex-col flex-auto pt-5 pb-4 overflow-y-auto">
      <router-link to="/" class="inline-flex items-center px-4 shrink-0" aria-label="Spectrograph logo">
        <img loading="eager" class="h-14 w-14 md:h-12 md:w-12" :src="imgLogo" aria-hidden="true" />
        <h2 class="text-3xl md:text-xl text-gradient bg-gradient-to-r from-nitro-500 to-dnd-400">
          Spectrograph
        </h2>
      </router-link>

      <nav
        class="flex flex-col flex-auto px-4 mt-5 space-y-2 md:space-y-1"
        aria-label="sidebar navigation"
      >
        <template v-for="group in groups" :key="group.id">
          <div v-if="group.title" class="pt-5">
            <span class="px-3 text-sm font-bold text-idle-500">{{ group.title }}</span>
          </div>
          <template v-for="link in group.links" :key="link.name">
            <component
              :is="link.to ? 'router-link' : 'a'"
              v-bind="
                link.to
                  ? { to: link.to, 'active-class': 'text-white bg-channel-900/50' }
                  : { href: link.href }
              "
              class="flex items-center p-3 text-lg font-medium transition-all duration-100 ease-in rounded text-channel-300 hover:bg-channel-900/50 hover:text-white group md:py-2 md:text-sm"
              :class="link.isChild ? 'ml-3' : ''"
            >
              <component
                :is="link.icon"
                class="w-5 h-5 mr-3 text-gray-400 shrink-0 group-hover:text-gray-300"
                aria-hidden="true"
              />
              {{ link.name }}

              <GuildStatus
                v-if="link.hasStatus"
                class="pl-1 ml-auto"
                :status="link.status"
                aria-hidden="true"
              />
            </component>
          </template>
        </template>
      </nav>
    </div>
    <div class="flex flex-col grow-0 shrink-0 bg-channel-700">
      <div class="flex items-center p-4">
        <img
          aria-hidden="true"
          class="rounded-full w-9 h-9 max-w-none"
          :src="state.base.self.avatarURL"
        />
        <div class="flex flex-col ml-3 truncate">
          <span class="inline-flex flex-auto text-lg font-medium text-balance-500 md:text-sm">
            {{ state.base.self.username }}
            <span class="text-channel-300">#{{ state.base.self.discriminator }}</span>
          </span>
          <span class="inline-flex text-lg font-medium md:text-sm text-channel-300">
            {{ userEmail }}
          </span>
        </div>
      </div>
      <div class="grid grid-flow-col gap-2 px-4 py-2 bg-channel-800">
        <router-link
          role="button"
          :to="{ name: '/dashboard/user/[id]', params: { id: state.base.self.id } }"
          class="btn btn-block btn-xs btn-primary"
        >
          Profile
        </router-link>
        <a role="button" href="/-/auth/logout" class="btn btn-block btn-xs">Logout</a>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import imgLogo from "@/assets/img/mic.png?format=webp&imagetools"
import { dashboardLinks, adminDashboardLinks } from "@/lib/core/navigation"
import type { DashboardLink } from "@/lib/core/types"

const state = useState()

interface Group {
  id: string
  title?: string
  links: DashboardLink[]
}

const groups = computed(() => {
  return [
    { id: "user-links", links: dashboardLinks.value },
    ...(state.base.self?.admin
      ? [{ id: "admin-links", title: "ADMIN", links: adminDashboardLinks }]
      : []),
  ] as Group[]
})

const userEmail = computed(() => {
  if (!state.base.self?.email) return ""

  const [username, domain] = state.base.self.email.split("@", 2)

  return username.replaceAll(/[^@]/g, "*") + "@" + domain
})
</script>
