<route lang="yaml">
meta:
  layout: dashboard
  title: Guilds
</route>

<template>
  <div class="p-2">
    <DataTable v-model="table">
      <template
        v-if="data?.guilds.pageInfo.hasNextPage || data?.guilds.pageInfo.hasPreviousPage"
        #footer
      >
        <CorePagination
          v-model="table.cursor.value"
          :page-info="data?.guilds.pageInfo"
          class="ml-auto"
        />
      </template>
    </DataTable>
  </div>
</template>

<script setup lang="ts">
import GuildIcon from "@/components/guild/icon.vue"
import { CoreTable } from "@/lib/util/table"
import { GuildOrderField, useGetAllGuildsQuery } from "@/lib/api"
import type { Guild, GuildWhereInput } from "@/lib/api"

const table = new CoreTable<Guild, GuildOrderField, GuildWhereInput>(
  [
    {
      id: "name",
      header: "Name",
      field: "name",
      sortField: GuildOrderField.Name,
      type: "text",
      filterFn: (val: string) => ({ nameContainsFold: val }),
      renderFn: (value, data) =>
        h("div", { class: "flex items-center gap-2" }, [
          h(GuildIcon, { guild: data, style: "height: 20px; width: 20px" }),
          h("div", {}, value),
        ]),
    },
    {
      id: "joinedAt",
      header: "Joined date",
      field: "joinedAt",
      sortField: GuildOrderField.JoinedAt,
    },
    {
      id: "memberCount",
      header: "Members",
      field: "memberCount",
    },
    {
      id: "large",
      header: "Large",
      field: "large",
      type: "boolean",
      filterFn: (val: boolean) => (val ? { large: val } : null),
    },
    {
      id: "enabled",
      header: "Enabled (user)",
      accessorFn: (data) => data.guildConfig?.enabled ?? false,
      type: "boolean",
      filterFn: (val: boolean) => (val !== null ? { hasGuildConfigWith: { enabled: val } } : null),
    },
    {
      id: "adminenabled",
      header: "Enabled (admin)",
      accessorFn: (data) => data.guildAdminConfig?.enabled ?? false,
      type: "boolean",
      filterFn: (val: boolean) => (val !== null ? { hasGuildAdminConfigWith: { enabled: val } } : null),
    },
  ],
  2
)

const { data } = useGetAllGuildsQuery({
  variables: computed(() => table.queryState.value),
})
table.data = computed(() => (data.value?.guilds.edges?.map(({ node }) => node) ?? []) as Guild[])
</script>
