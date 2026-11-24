<script setup lang="ts">
import {
  Bug,
  ChevronsUpDown,
  CircleDot,
  LayoutGrid,
  ListTree,
  LogOut,
  MessagesSquare,
  Server,
  ShieldCheck,
  UserRound,
  Users,
} from "@lucide/vue"

const route = useRoute()
const { self, isAdmin } = useSelf()
const guilds = useMyGuilds()
const bugURL = useBugReportURL()

const adminItems = [
  { title: "All servers", to: "/admin/guilds", icon: Server },
  { title: "Event log", to: "/admin/events", icon: ListTree },
  { title: "Users", to: "/admin/users", icon: Users },
]

function isActive(to: string, exact = false) {
  const current = route.path.replace(/\/+$/, "") || "/"
  return exact ? current === to : current === to || current.startsWith(`${to}/`)
}
</script>

<template>
  <Sidebar collapsible="icon">
    <SidebarHeader>
      <SidebarMenu>
        <SidebarMenuItem>
          <SidebarMenuButton size="lg" as-child tooltip="Spectrograph">
            <BrandMark to="/dashboard" />
          </SidebarMenuButton>
        </SidebarMenuItem>
      </SidebarMenu>
    </SidebarHeader>
    <SidebarContent>
      <SidebarGroup>
        <SidebarGroupContent>
          <SidebarMenu>
            <SidebarMenuItem>
              <SidebarMenuButton as-child :is-active="isActive('/dashboard', true)" tooltip="Your servers">
                <NuxtLink to="/dashboard">
                  <LayoutGrid />
                  <span>Your servers</span>
                </NuxtLink>
              </SidebarMenuButton>
            </SidebarMenuItem>
            <SidebarMenuItem v-if="self">
              <SidebarMenuButton as-child :is-active="isActive(`/dashboard/user/${self.id}`)" tooltip="Account">
                <NuxtLink :to="`/dashboard/user/${self.id}`">
                  <UserRound />
                  <span>Account</span>
                </NuxtLink>
              </SidebarMenuButton>
            </SidebarMenuItem>
          </SidebarMenu>
        </SidebarGroupContent>
      </SidebarGroup>

      <SidebarGroup v-if="guilds.status.value === 'pending' || guilds.data.value?.content.length">
        <SidebarGroupLabel>Servers</SidebarGroupLabel>
        <SidebarGroupContent>
          <SidebarMenu>
            <template v-if="guilds.status.value === 'pending' && !guilds.data.value">
              <SidebarMenuItem v-for="n in 3" :key="n">
                <SidebarMenuSkeleton show-icon />
              </SidebarMenuItem>
            </template>
            <SidebarMenuItem v-for="guild in guilds.data.value?.content ?? []" :key="guild.id">
              <SidebarMenuButton
                as-child
                :is-active="isActive(`/dashboard/guild/${guild.id}`)"
                :tooltip="`${guild.name} (${guildStatusLabel[guildStatus(guild)].toLowerCase()})`"
              >
                <NuxtLink :to="`/dashboard/guild/${guild.id}`">
                  <GuildIcon :guild="guild" compact class="size-4 text-[10px] ring-0 ring-offset-0" />
                  <span class="truncate">{{ guild.name }}</span>
                  <CircleDot
                    :class="[
                      'ml-auto size-3!',
                      guildStatus(guild) === 'active' && 'text-success',
                      guildStatus(guild) === 'paused' && 'text-warning',
                      guildStatus(guild) === 'admin-disabled' && 'text-destructive',
                      guildStatus(guild) === 'not-joined' && 'text-muted-foreground',
                    ]"
                    :aria-label="guildStatusLabel[guildStatus(guild)]"
                  />
                </NuxtLink>
              </SidebarMenuButton>
            </SidebarMenuItem>
          </SidebarMenu>
        </SidebarGroupContent>
      </SidebarGroup>

      <SidebarGroup v-if="isAdmin">
        <SidebarGroupLabel>Administration</SidebarGroupLabel>
        <SidebarGroupContent>
          <SidebarMenu>
            <SidebarMenuItem v-for="item in adminItems" :key="item.to">
              <SidebarMenuButton as-child :is-active="isActive(item.to)" :tooltip="item.title">
                <NuxtLink :to="item.to">
                  <component :is="item.icon" />
                  <span>{{ item.title }}</span>
                </NuxtLink>
              </SidebarMenuButton>
            </SidebarMenuItem>
          </SidebarMenu>
        </SidebarGroupContent>
      </SidebarGroup>

      <SidebarGroup class="mt-auto">
        <SidebarGroupLabel>Help</SidebarGroupLabel>
        <SidebarGroupContent>
          <SidebarMenu>
            <SidebarMenuItem>
              <SidebarMenuButton as-child tooltip="Report an issue">
                <a :href="bugURL" target="_blank" rel="noopener">
                  <Bug />
                  <span>Report an issue</span>
                </a>
              </SidebarMenuButton>
            </SidebarMenuItem>
            <SidebarMenuItem>
              <SidebarMenuButton as-child tooltip="Known issues">
                <a :href="issuesURL" target="_blank" rel="noopener">
                  <MessagesSquare />
                  <span>Known issues</span>
                </a>
              </SidebarMenuButton>
            </SidebarMenuItem>
          </SidebarMenu>
        </SidebarGroupContent>
      </SidebarGroup>
    </SidebarContent>

    <SidebarFooter v-if="self">
      <SidebarMenu>
        <SidebarMenuItem>
          <DropdownMenu>
            <DropdownMenuTrigger as-child>
              <SidebarMenuButton size="lg" class="data-[state=open]:bg-sidebar-accent">
                <UserAvatar :user="self" class="size-8" />
                <div class="grid min-w-0 flex-1 text-left text-sm leading-tight">
                  <span class="flex items-center gap-1 truncate font-semibold">
                    {{ displayName(self) }}
                    <ShieldCheck v-if="self.admin" class="size-3.5! text-primary" aria-label="Administrator" />
                  </span>
                  <span class="truncate text-xs text-muted-foreground">{{ userTag(self) }}</span>
                </div>
                <ChevronsUpDown class="ml-auto" />
              </SidebarMenuButton>
            </DropdownMenuTrigger>
            <DropdownMenuContent side="top" align="start" class="w-(--reka-dropdown-menu-trigger-width) min-w-56">
              <DropdownMenuGroup>
                <DropdownMenuItem as-child>
                  <NuxtLink :to="`/dashboard/user/${self.id}`">
                    <UserRound />
                    Account
                  </NuxtLink>
                </DropdownMenuItem>
                <DropdownMenuItem as-child>
                  <a :href="bugURL" target="_blank" rel="noopener">
                    <Bug />
                    Report an issue
                  </a>
                </DropdownMenuItem>
              </DropdownMenuGroup>
              <DropdownMenuSeparator />
              <DropdownMenuGroup>
                <DropdownMenuItem as-child>
                  <a :href="logoutURL('/')">
                    <LogOut />
                    Log out
                  </a>
                </DropdownMenuItem>
              </DropdownMenuGroup>
            </DropdownMenuContent>
          </DropdownMenu>
        </SidebarMenuItem>
      </SidebarMenu>
    </SidebarFooter>
    <SidebarRail />
  </Sidebar>
</template>
