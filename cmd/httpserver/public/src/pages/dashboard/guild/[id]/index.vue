<route lang="yaml">
meta:
  layout: dashboard
  title: Dashboard
</route>

<template>
  <FeedbackAlert v-if="error" type="error">{{ error }}</FeedbackAlert>
  <div v-else>
    <GuildHeader :guild="guild" @toggle-enabled="toggleEnabled" />
    <GuildEvents :guild="guild" class="mt-10 max-h-64" />

    <ContainerSettings
      title="Guild Configuration"
      description="This configuration will change how Spectrograph will interact with your guild."
      class="mt-12"
    >
      <!-- contact email -->
      <div class="col-span-2 lg:col-span-1 form-control">
        <label id="contact-email" class="label">
          <span class="label-text text-secondary">Contact email</span>
        </label>
        <input
          v-model="guild.guildConfig.contactEmail"
          type="email"
          placeholder="Contact email"
          class="w-full rounded input input-bordered"
          aria-labelledby="contact-email"
        />
      </div>

      <!-- channel match rule -->
      <div class="col-span-2 lg:col-span-1 form-control">
        <label id="channel-match-rule" class="label">
          <span class="label-text text-secondary">Channel match rule</span>
          <a
            class="label-text-alt link link-primary"
            href="https://github.com/google/re2/wiki/Syntax"
            target="_blank"
          >
            how-to
          </a>
        </label>
        <input
          v-model="guild.guildConfig.regexMatch"
          type="text"
          placeholder="default: '^.* \+$'"
          class="w-full rounded input input-bordered"
          aria-labelledby="channel-match-rule"
        />
      </div>

      <!-- maximum clones -->
      <div class="col-span-2 lg:col-span-1 form-control">
        <label id="maximum-clones" class="label">
          <span class="label-text text-secondary">Maximum allowed clones</span>
        </label>
        <input
          v-model="guild.guildConfig.defaultMaxClones"
          type="number"
          min="0"
          placeholder=""
          class="w-full rounded input input-bordered"
          aria-labelledby="maximum-clones"
        />
        <label class="label">
          <span class="text-xs label-text">Set to 0 for max allowed</span>
        </label>
      </div>

      <template #actions>
        <div class="hidden italic text-chat-300 lg:block" aria-hidden="true">
          last updated: {{ configUpdatedAt }}
        </div>

        <button
          class="h-10 min-h-0 ml-auto btn btn-accent"
          :disabled="!guild.guildAdminConfig.enabled"
          @click="updateConfig"
        >
          <i-fas-circle-check class="h-4 mr-2" />
          Save
        </button>
      </template>
    </ContainerSettings>

    <ContainerSettings
      v-if="state.base.self?.admin"
      title="Admin Configuration"
      description="Administrative configuration which has the potential to override user-level configs."
      class="mt-12"
    >
      <!-- admin: maximum clones -->
      <div class="col-span-2 lg:col-span-1 form-control">
        <label id="admin-maximum-clones" class="label">
          <span class="label-text text-secondary">Maximum allowed clones</span>
        </label>
        <input
          v-model="guild.guildAdminConfig.defaultMaxClones"
          type="number"
          min="0"
          placeholder=""
          class="w-full rounded input input-bordered"
          aria-labelledby="admin-maximum-clones"
        />
        <label class="label">
          <span class="text-xs label-text">Set to 0 for max allowed</span>
        </label>
      </div>

      <!-- admin: maximum channels -->
      <div class="col-span-2 lg:col-span-1 form-control">
        <label id="admin-maximum-channels" class="label">
          <span class="label-text text-secondary">Maximum allowed channels</span>
        </label>
        <input
          v-model="guild.guildAdminConfig.defaultMaxChannels"
          type="number"
          min="0"
          placeholder=""
          class="w-full rounded input input-bordered"
          aria-labelledby="admin-maximum-channels"
        />
        <label class="label">
          <span class="text-xs label-text">Set to 0 for max allowed</span>
        </label>
      </div>

      <!-- admin: comment -->
      <div class="col-span-2 form-control">
        <label id="admin-comment" class="label">
          <span class="label-text text-secondary">Administrative comment</span>
        </label>
        <textarea
          v-model="guild.guildAdminConfig.comment"
          class="w-full h-24 rounded input input-bordered"
          aria-labelledby="admin-comment"
        />
        <label class="label">
          <span class="text-xs label-text">May be visible to the user</span>
        </label>
      </div>

      <template #actions>
        <div class="hidden italic text-chat-300 lg:block" aria-hidden="true">
          last updated: {{ adminConfigUpdatedAt }}
        </div>

        <div class="ml-auto space-x-3">
          <button
            class="h-10 min-h-0 btn"
            :class="guild.guildAdminConfig?.enabled ? 'btn-error' : 'btn-success'"
            @click="toggleAdminEnabled"
          >
            <i-fas-circle-xmark v-if="guild.guildAdminConfig?.enabled" class="h-4 mr-2" />
            <i-fas-circle-check v-else class="h-4 mr-2" />
            {{ guild.guildAdminConfig?.enabled ? "Disable (admin)" : "Enable (admin)" }}
          </button>
          <button class="h-10 min-h-0 btn btn-accent" @click="updateAdminConfig">
            <i-fas-circle-check class="h-4 mr-2" />
            Save
          </button>
        </div>
      </template>
    </ContainerSettings>
  </div>
</template>

<script setup lang="ts" async>
import { useTimeAgo } from "@vueuse/core"
import {
  useGetGuildQuery,
  useUpdateGuildConfigMutation,
  useUpdateGuildAdminConfigMutation,
} from "@/lib/api"
import type { Guild } from "@/lib/api"

const state = useState()
const route = useRoute("/dashboard/guild/[id]/")

const {
  data,
  error,
  executeQuery: refetch,
} = await useGetGuildQuery({ variables: { id: route.params.id } })
const guild = computed(() => data?.value.node as Guild)

const config = useUpdateGuildConfigMutation()
const configUpdatedAt = useTimeAgo(computed(() => guild.value.guildConfig.updateTime))

function toggleEnabled() {
  config
    .executeMutation({
      id: guild.value.guildConfig.id,
      input: { enabled: !guild.value.guildConfig.enabled },
    })
    .then(() => {
      refetch()
      state.fetchBase()
    })
}

function updateConfig() {
  config
    .executeMutation({
      id: guild.value.guildConfig.id,
      input: {
        contactEmail: guild.value.guildConfig.contactEmail,
        regexMatch: guild.value.guildConfig.regexMatch,
        defaultMaxClones: guild.value.guildConfig.defaultMaxClones,
      },
    })
    .then(() => {
      refetch()
      state.fetchBase()
    })
}

const adminconfig = useUpdateGuildAdminConfigMutation()
const adminConfigUpdatedAt = useTimeAgo(computed(() => guild.value.guildAdminConfig.updateTime))

function toggleAdminEnabled() {
  adminconfig
    .executeMutation({
      id: guild.value.guildAdminConfig.id,
      input: { enabled: !guild.value.guildAdminConfig.enabled },
    })
    .then(() => {
      refetch()
      state.fetchBase()
    })
}

function updateAdminConfig() {
  adminconfig
    .executeMutation({
      id: guild.value.guildAdminConfig.id,
      input: {
        defaultMaxChannels: guild.value.guildAdminConfig.defaultMaxChannels,
        defaultMaxClones: guild.value.guildAdminConfig.defaultMaxClones,
        comment: guild.value.guildAdminConfig.comment,
      },
    })
    .then(() => {
      refetch()
      state.fetchBase()
    })
}
</script>
