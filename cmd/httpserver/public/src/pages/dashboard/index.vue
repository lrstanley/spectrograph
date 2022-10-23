<route lang="yaml">
meta:
  layout: dashboard
  title: Dashboard
</route>

<template>
  <div>
    <h1 class="pb-4 text-2xl font-semibold text-white">Dashboard</h1>

    <GuildList :guilds="guilds" show-empty />
  </div>
</template>

<script setup lang="ts">
import type { Guild } from "@/lib/api"

const state = useState()
const guilds = state.base.self?.userGuilds.edges?.map(({ node }) => node as Guild) ?? []
</script>
