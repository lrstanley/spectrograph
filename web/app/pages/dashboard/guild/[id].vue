<script setup lang="ts">
import { ExternalLink, RefreshCw, ShieldAlert } from "@lucide/vue"
import { toast } from "vue-sonner"
import type { SchemaGuildAdminConfig } from "#open-fetch-schemas/api"

definePageMeta({ title: "Server", key: (route) => route.fullPath })

const route = useRoute("dashboard-guild-id")
const { $api } = useNuxtApp()
const { isAdmin } = useSelf()
const guildID = Number(route.params.id)

// Deep, so saved configs can be written back into the fetched guild.
const guild = useApi("/guilds/{guildID}", { path: { guildID }, deep: true })
// Only system admins can see other users, so this is an admin-only view.
const admins = useLazyApi("/guilds/{guildID}/admins", {
  path: { guildID },
  immediate: isAdmin.value,
  query: { per_page: 25, sort: "username", order: "asc" },
})
await guild

if (guild.error.value || !guild.data.value) {
  const status = apiErrorStatus(guild.error.value) ?? 500
  throw createError({
    statusCode: status,
    statusMessage: status === 404 ? "Server not found" : "Couldn't load this server",
    message:
      status === 404
        ? "This server doesn't exist, or you don't manage it. If you were just made an admin, log out and back in."
        : apiErrorMessage(guild.error.value),
    data: { requestID: apiRequestID(guild.error.value) },
    fatal: true,
  })
}

useHead({ title: () => guild.data.value?.name ?? "Server" })
usePageTitle(() => guild.data.value?.name)

const status = computed(() => (guild.data.value ? guildStatus(guild.data.value) : "not-joined"))
const config = computed(() => guild.data.value?.edges.guild_config)
const adminConfig = computed(() => guild.data.value?.edges.guild_admin_config)
const toggling = ref(false)

function onAdminConfigSaved(config: SchemaGuildAdminConfig) {
  if (!guild.data.value) return
  guild.data.value.edges.guild_admin_config = config
  refreshNuxtData("my-guilds")
}

async function setEnabled(enabled: boolean) {
  if (!config.value || !guild.data.value) return
  toggling.value = true
  try {
    guild.data.value.edges.guild_config = await $api("/guild-configs/{guildconfigID}", {
      method: "PATCH",
      path: { guildconfigID: config.value.id },
      body: { enabled },
    })
    refreshNuxtData("my-guilds")
    toast.success(enabled ? "Spectrograph resumed" : "Spectrograph paused", {
      description: enabled
        ? `Managing voice channels in ${guild.data.value.name} again.`
        : `Channels in ${guild.data.value.name} are left as they are.`,
    })
  } catch (err) {
    toast.error("Couldn't change the server's status", { description: apiErrorMessage(err) })
  } finally {
    toggling.value = false
  }
}
</script>

<template>
  <template v-if="guild.data.value">
    <header class="flex flex-col gap-4 sm:flex-row sm:items-center">
      <GuildIcon :guild="guild.data.value" :status="status" class="size-14 text-lg" />
      <div class="flex min-w-0 flex-1 flex-col gap-1">
        <h1 class="truncate font-heading text-2xl font-semibold tracking-tight">{{ guild.data.value.name }}</h1>
        <p class="flex flex-wrap items-center gap-x-3 gap-y-1 text-sm text-muted-foreground">
          <Badge
            :variant="status === 'active' ? 'default' : status === 'not-joined' ? 'outline' : 'destructive'"
          >
            {{ guildStatusLabel[status] }}
          </Badge>
          <span v-if="guild.data.value.member_count">{{ guild.data.value.member_count.toLocaleString() }} members</span>
          <span v-if="guild.data.value.joined_at" :title="formatDate(guild.data.value.joined_at)">
            Added {{ timeAgo(guild.data.value.joined_at) }}
          </span>
        </p>
      </div>
      <Field v-if="config" orientation="horizontal" class="w-auto items-center">
        <FieldLabel for="guild-enabled" class="text-sm">Manage channels</FieldLabel>
        <Switch
          id="guild-enabled"
          :model-value="config.enabled !== false"
          :disabled="toggling || adminConfig?.enabled === false"
          @update:model-value="setEnabled"
        />
      </Field>
    </header>

    <Alert v-if="adminConfig?.enabled === false" variant="destructive">
      <ShieldAlert />
      <AlertTitle>An administrator disabled Spectrograph in this server</AlertTitle>
      <AlertDescription>
        <p>Your settings are kept, but the bot won't manage channels here until it's re-enabled.</p>
        <p v-if="adminConfig.comment">Reason: {{ adminConfig.comment }}</p>
      </AlertDescription>
    </Alert>

    <Empty v-if="status === 'not-joined' || !config" class="border">
      <EmptyHeader>
        <EmptyMedia>
          <GuildIcon :guild="guild.data.value" class="size-16 text-xl" />
        </EmptyMedia>
        <EmptyTitle>Spectrograph isn't in {{ guild.data.value.name }} yet</EmptyTitle>
        <EmptyDescription>
          Add the bot to start managing voice channels. Discord asks you to confirm the server and its permissions:
          view and manage channels, plus the voice permissions it copies onto cloned channels.
        </EmptyDescription>
      </EmptyHeader>
      <EmptyContent class="flex-row justify-center gap-2">
        <Button as-child>
          <a :href="inviteURL(guild.data.value.guild_id)" target="_blank" rel="noopener">
            Add to {{ guild.data.value.name }}
            <ExternalLink data-icon="inline-end" />
          </a>
        </Button>
        <Button variant="outline" :disabled="guild.status.value === 'pending'" @click="guild.refresh()">
          <RefreshCw data-icon="inline-start" />
          Check again
        </Button>
      </EmptyContent>
    </Empty>

    <div v-else class="grid items-start gap-6 lg:grid-cols-[minmax(0,1fr)_minmax(0,1fr)]">
      <div class="flex min-w-0 flex-col gap-6">
        <GuildConfigForm
          :config="config"
          :disabled="adminConfig?.enabled === false"
          @saved="guild.data.value.edges.guild_config = $event"
        />
        <GuildAdminConfigForm
          v-if="isAdmin && adminConfig"
          :config="adminConfig"
          @saved="onAdminConfigSaved"
        />
        <Card v-if="isAdmin && admins.data.value?.content.length">
          <CardHeader>
            <CardTitle>Who can manage this server</CardTitle>
            <CardDescription>Discord users with Manage Server here who have logged in to Spectrograph.</CardDescription>
          </CardHeader>
          <CardContent>
            <ul class="flex flex-col gap-3">
              <li v-for="admin in admins.data.value.content" :key="admin.id" class="flex items-center gap-3">
                <UserAvatar :user="admin" class="size-8" />
                <div class="min-w-0 text-sm leading-tight">
                  <p class="truncate font-medium">{{ displayName(admin) }}</p>
                  <p class="truncate text-xs text-muted-foreground">{{ userTag(admin) }}</p>
                </div>
              </li>
            </ul>
          </CardContent>
        </Card>
      </div>
      <GuildActivity :guild-id="guild.data.value.id" class="lg:sticky lg:top-20" />
    </div>
  </template>
</template>
