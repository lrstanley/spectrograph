<script setup lang="ts">
import { Ban, LogOut, ShieldCheck } from "@lucide/vue"
import type { SchemaUser } from "#open-fetch-schemas/api"

definePageMeta({ title: "Account", key: (route) => route.fullPath })

const route = useRoute("dashboard-user-id")
const { self, isAdmin } = useSelf()
const userID = Number(route.params.id)

const user = useApi("/users/{userID}", { path: { userID } })
const guilds = useLazyApi("/users/{userID}/user-guilds", {
  path: { userID },
  query: { per_page: 100, sort: "name", order: "asc" },
})
await user

if (user.error.value || !user.data.value) {
  const status = apiErrorStatus(user.error.value) ?? 500
  throw createError({
    statusCode: status,
    statusMessage: status === 404 ? "User not found" : "Couldn't load this user",
    message: status === 404 ? "This user doesn't exist, or you can't view them." : apiErrorMessage(user.error.value),
    data: { requestID: apiRequestID(user.error.value) },
    fatal: true,
  })
}

const isSelf = computed(() => self.value?.id === user.data.value?.id)
const pageTitle = computed(() => (isSelf.value ? "Account" : (user.data.value?.username ?? "User")))
useHead({ title: pageTitle })
usePageTitle(pageTitle)

const showEmail = ref(false)
const deleteOpen = ref(false)
const banTarget = ref<SchemaUser | null>(null)

const maskedEmail = computed(() => {
  const [name = "", domain = ""] = (user.data.value?.email ?? "").split("@", 2)
  return name ? `${name[0]}${"•".repeat(Math.max(2, name.length - 1))}@${domain}` : ""
})

const details = computed(() => {
  const u = user.data.value
  if (!u) return []
  return [
    { label: "Discord ID", value: u.user_id, mono: true },
    { label: "Locale", value: u.locale || "Unknown" },
    { label: "Email verified", value: u.verified ? "Yes" : "No" },
    { label: "Two-factor auth", value: u.mfa_enabled ? "On" : "Off" },
    { label: "First login", value: timeAgo(u.create_time), title: formatDate(u.create_time) },
    { label: "Last updated", value: timeAgo(u.update_time), title: formatDate(u.update_time) },
  ]
})
</script>

<template>
  <template v-if="user.data.value">
    <header class="flex flex-col gap-4 sm:flex-row sm:items-center">
      <UserAvatar
        :user="user.data.value"
        :class="['size-16 ring-2 ring-offset-2 ring-offset-background', user.data.value.banned ? 'ring-destructive' : 'ring-primary']"
      />
      <div class="flex min-w-0 flex-1 flex-col gap-1">
        <h1 class="flex items-center gap-2 truncate font-heading text-2xl font-semibold tracking-tight">
          {{ displayName(user.data.value) }}
          <Badge v-if="user.data.value.admin" variant="outline"><ShieldCheck data-icon="inline-start" />Admin</Badge>
          <Badge v-if="user.data.value.banned" variant="destructive"><Ban data-icon="inline-start" />Banned</Badge>
        </h1>
        <p class="text-sm text-muted-foreground">{{ userTag(user.data.value) }}</p>
      </div>
      <Button v-if="isSelf" as-child variant="outline">
        <a :href="logoutURL('/')"><LogOut data-icon="inline-start" />Log out</a>
      </Button>
    </header>

    <Alert v-if="user.data.value.banned" variant="destructive">
      <Ban />
      <AlertTitle>Banned from Spectrograph</AlertTitle>
      <AlertDescription v-if="user.data.value.ban_reason">Reason: {{ user.data.value.ban_reason }}</AlertDescription>
    </Alert>

    <Card>
      <CardHeader>
        <CardTitle>Profile</CardTitle>
        <CardDescription>Synced from Discord each time {{ isSelf ? "you log" : "they log" }} in.</CardDescription>
      </CardHeader>
      <CardContent>
        <dl class="grid gap-x-8 gap-y-4 sm:grid-cols-2">
          <div class="flex flex-col gap-0.5">
            <dt class="text-xs text-muted-foreground">Email</dt>
            <dd class="flex items-center gap-2 text-sm">
              <span class="truncate">{{ showEmail ? user.data.value.email : maskedEmail || "Not shared" }}</span>
              <Button
                v-if="user.data.value.email"
                variant="ghost"
                size="xs"
                :aria-label="showEmail ? 'Hide email address' : 'Show email address'"
                @click="showEmail = !showEmail"
              >
                {{ showEmail ? "Hide" : "Show" }}
              </Button>
            </dd>
          </div>
          <div v-for="item in details" :key="item.label" class="flex flex-col gap-0.5">
            <dt class="text-xs text-muted-foreground">{{ item.label }}</dt>
            <dd :class="['text-sm', item.mono && 'font-mono']" :title="item.title">{{ item.value }}</dd>
          </div>
        </dl>
      </CardContent>
    </Card>

    <Card>
      <CardHeader>
        <CardTitle>Servers</CardTitle>
        <CardDescription>Servers {{ isSelf ? "you manage" : "this user manages" }} with Spectrograph.</CardDescription>
      </CardHeader>
      <CardContent>
        <div v-if="guilds.status.value === 'pending'" class="flex flex-col gap-2">
          <Skeleton v-for="n in 2" :key="n" class="h-10 w-full" />
        </div>
        <ItemGroup v-else-if="guilds.data.value?.content.length" class="gap-1">
          <Item v-for="guild in guilds.data.value.content" :key="guild.id" as-child size="sm">
            <NuxtLink :to="`/dashboard/guild/${guild.id}`">
              <ItemMedia>
                <GuildIcon :guild="guild" :status="guildStatus(guild)" class="size-8" />
              </ItemMedia>
              <ItemContent>
                <ItemTitle>{{ guild.name }}</ItemTitle>
                <ItemDescription>{{ guildStatusLabel[guildStatus(guild)] }}</ItemDescription>
              </ItemContent>
            </NuxtLink>
          </Item>
        </ItemGroup>
        <p v-else class="text-sm text-muted-foreground">No servers.</p>
      </CardContent>
    </Card>

    <Card v-if="isSelf || (isAdmin && !user.data.value.admin)" class="border-destructive/40">
      <CardHeader>
        <CardTitle>{{ isSelf ? "Delete account" : "Ban user" }}</CardTitle>
        <CardDescription v-if="isSelf">
          Removes your account and signs you out. Servers keep their Spectrograph settings.
        </CardDescription>
        <CardDescription v-else>
          {{ user.data.value.banned ? "Let this user log in again." : "Stop this user from logging in to Spectrograph." }}
        </CardDescription>
      </CardHeader>
      <CardFooter>
        <Button v-if="isSelf" variant="destructive" @click="deleteOpen = true">Delete account</Button>
        <Button v-else :variant="user.data.value.banned ? 'outline' : 'destructive'" @click="banTarget = user.data.value">
          {{ user.data.value.banned ? "Unban user" : "Ban user" }}
        </Button>
      </CardFooter>
    </Card>

    <DeleteAccountDialog v-if="isSelf" v-model:open="deleteOpen" />
    <BanUserDialog :user="banTarget" @close="banTarget = null" @changed="user.data.value = $event" />
  </template>
</template>
