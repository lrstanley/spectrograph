<route lang="yaml">
meta:
  layout: dashboard
  title: Dashboard
</route>

<template>
  <FeedbackAlert v-if="error" type="error">{{ error }}</FeedbackAlert>
  <div v-else class="space-y-12">
    <UserBanModal v-model="openBanModal" :user="user" @banned="refetch" />
    <UserHeader :user="user" />

    <ContainerDescription
      title="Profile Information"
      subtitle="Information about your profile (mostly pulled from Discord)"
    >
      <ContainerDescriptionItem name="Username">
        <span class="text-nitro-500">{{ user.username }}#{{ user.discriminator }}</span>
      </ContainerDescriptionItem>
      <ContainerDescriptionItem name="Email" :value="user.email">
        {{ showEmail ? user.email : userEmail }}

        <button class="btn btn-xs" @click="showEmail = !showEmail">
          {{ showEmail ? "hide" : "show" }}
        </button>
      </ContainerDescriptionItem>
      <ContainerDescriptionItem name="Locale" :value="user.locale" />
      <ContainerDescriptionItem name="Verification Status">
        <span :class="user.verified ? 'text-online-500' : 'text-dnd-400'">
          {{ user.verified ? "Verified" : "Not verified" }}
        </span>
      </ContainerDescriptionItem>
      <ContainerDescriptionItem name="MFA Enabled">
        <span :class="user.mfaEnabled ? 'text-online-500' : 'text-dnd-400'">
          {{ user.mfaEnabled ? "Enabled" : "Disabled" }}
        </span>
      </ContainerDescriptionItem>
      <ContainerDescriptionItem name="First Login">
        {{ useTimeAgo(user.createTime).value }}
      </ContainerDescriptionItem>
      <ContainerDescriptionItem name="User Last Updated">
        {{ useTimeAgo(user.updateTime).value }}
      </ContainerDescriptionItem>
    </ContainerDescription>

    <GuildList :guilds="guilds" :show-empty="state.base.self?.id === user.id" />

    <ContainerSettings
      title="Destructive Actions"
      description="Irreversible actions that may delete your account and all associated data."
      bare-panel
    >
      <div class="flex flex-col flex-auto divide-y divide-channel-800 border-chat-900">
        <div class="flex flex-col items-start px-6 py-5 md:items-center md:flex-row">
          <div>
            <div class="text-base uppercase">Delete Account</div>
            <div class="text-sm text-gray-500">
              This will delete your account and all of your data (however will not remove Spectrograph
              from your guilds).
            </div>
          </div>
          <button class="w-full mt-2 ml-auto btn btn-sm btn-error md:w-auto md:mt-0">
            Delete account
          </button>
        </div>
        <div
          v-if="state.base.self?.admin"
          class="flex flex-col items-start px-6 py-5 md:items-center md:flex-row"
        >
          <div>
            <div class="text-base uppercase">Ban user (admin only)</div>
            <div class="text-sm text-gray-500">This will ban the user from using Spectrograph.</div>
          </div>
          <button
            class="w-full mt-2 ml-auto btn btn-sm btn-error md:w-auto md:mt-0"
            :disabled="state.base.self?.id === user.id"
            @click="openBanModal = true"
          >
            {{ user.banned ? "Unban" : "Ban" }} user
          </button>
        </div>
      </div>
    </ContainerSettings>
  </div>
</template>

<script setup lang="ts" async>
import { useTimeAgo } from "@vueuse/core"
import { useGetUserQuery } from "@/lib/api"
import type { User, Guild } from "@/lib/api"

const state = useState()
const route = useRoute("/dashboard/user/[id]")

const {
  data,
  error,
  executeQuery: refetch,
} = await useGetUserQuery({ variables: { id: route.params.id } })
const user = computed(() => data?.value.node as User)
const guilds = computed(() => user.value?.userGuilds.edges?.map(({ node }) => node as Guild) ?? [])

const userEmail = computed(() => {
  if (!user.value) return ""

  const [username, domain] = user.value.email.split("@", 2)

  return username.replaceAll(/[^@]/g, "*") + "@" + domain
})

const showEmail = ref(false)
const openBanModal = ref(false)
</script>
