<script setup lang="ts">
import { ChevronRight, ExternalLink, Server } from "@lucide/vue"

definePageMeta({ title: "Your servers" })

const guilds = useMyGuilds()
await guilds
</script>

<template>
  <div class="flex flex-wrap items-end justify-between gap-4">
    <div class="flex flex-col gap-1">
      <h1 class="font-heading text-2xl font-semibold tracking-tight">Your servers</h1>
      <p class="text-sm text-muted-foreground">Servers where you have Manage Server permission in Discord.</p>
    </div>
    <Button as-child variant="outline">
      <a :href="inviteURL()" target="_blank" rel="noopener">
        Add to another server
        <ExternalLink data-icon="inline-end" />
      </a>
    </Button>
  </div>

  <Alert v-if="guilds.error.value" variant="destructive">
    <AlertTitle>Couldn't load your servers</AlertTitle>
    <AlertDescription>{{ apiErrorMessage(guilds.error.value) }}</AlertDescription>
  </Alert>

  <ItemGroup v-else-if="guilds.data.value?.content.length" class="gap-2">
    <Item v-for="guild in guilds.data.value.content" :key="guild.id" as-child variant="outline">
      <NuxtLink :to="`/dashboard/guild/${guild.id}`">
        <ItemMedia>
          <GuildIcon :guild="guild" :status="guildStatus(guild)" class="size-10" />
        </ItemMedia>
        <ItemContent>
          <ItemTitle>{{ guild.name }}</ItemTitle>
          <ItemDescription>
            {{ guildStatusLabel[guildStatus(guild)] }}
            <template v-if="guild.member_count"> with {{ guild.member_count.toLocaleString() }} members</template>
          </ItemDescription>
        </ItemContent>
        <ItemActions>
          <Badge v-if="guildStatus(guild) === 'not-joined'" variant="outline">Add bot</Badge>
          <ChevronRight class="size-4 text-muted-foreground" aria-hidden="true" />
        </ItemActions>
      </NuxtLink>
    </Item>
  </ItemGroup>

  <Empty v-else class="border">
    <EmptyHeader>
      <EmptyMedia variant="icon">
        <Server />
      </EmptyMedia>
      <EmptyTitle>No servers to manage yet</EmptyTitle>
      <EmptyDescription>
        Spectrograph lists servers you own or have Manage Server permission in. If you were just given that
        permission, log out and back in to refresh the list.
      </EmptyDescription>
    </EmptyHeader>
    <EmptyContent>
      <Button as-child variant="outline">
        <a :href="logoutURL('/login')">Log out</a>
      </Button>
    </EmptyContent>
  </Empty>
</template>
