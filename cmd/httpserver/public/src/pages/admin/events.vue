<route lang="yaml">
meta:
  layout: dashboard
  title: Events
</route>

<template>
  <div class="p-2">
    <DataTable v-model="table">
      <template
        v-if="data?.guildevents.pageInfo.hasNextPage || data?.guildevents.pageInfo.hasPreviousPage"
        #footer
      >
        <CorePagination
          v-model="table.cursor.value"
          :page-info="data?.guildevents.pageInfo"
          class="ml-auto"
        />
      </template>
    </DataTable>
  </div>
</template>

<script setup lang="ts">
import GuildIcon from "@/components/guild/icon.vue"
import { useTimeAgo } from "@vueuse/core"
import { CoreTable } from "@/lib/util"
import { GuildEventType, GuildEventOrderField, useGetAllGuildEventsQuery } from "@/lib/api"
import type { GuildEvent, GuildEventWhereInput } from "@/lib/api"

const router = useRouter()

const table = new CoreTable<GuildEvent, GuildEventOrderField, GuildEventWhereInput>(
  [
    {
      id: "guild",
      header: "Guild",
      accessorFn: (e) => e.guild.name,
      type: "text",
      clickFn: (data, event) => {
        router.push({ name: "/dashboard/guild/[id]/", params: { id: data.guild.id } })
        event.preventDefault()
      },
      filterFn: (val: string) =>
        val === null
          ? null
          : {
              hasGuildWith: {
                or: [{ nameContainsFold: val }, { guildIDContainsFold: val }],
              },
            },
      renderFn: (value, e) =>
        h("div", { class: "flex items-center gap-2" }, [
          h(GuildIcon, { guild: e.guild, style: "height: 20px; width: 20px" }),
          h("div", {}, value),
        ]),
    },
    {
      id: "createTime",
      header: "Timestamp",
      field: "createTime",
      sortField: GuildEventOrderField.CreatedAt,
      renderFn: (value) => h("div", {}, useTimeAgo(value).value),
    },
    {
      id: "type",
      header: "Type",
      field: "type",
      renderFn: (value, e) =>
        h(
          "div",
          {
            class: {
              "text-bravery-500": e.type === GuildEventType.Info,
              "text-idle-500": e.type === GuildEventType.Warning,
              "text-dnd-500": e.type === GuildEventType.Error,
              "text-nitro-500": e.type === GuildEventType.Debug,
            },
          },
          e.type
        ),
    },
    {
      id: "message",
      header: "Message",
      field: "message",
      type: "text",
      filterFn: (val: string) => (val !== null ? { messageContainsFold: val } : null),
      renderFn: (value) => h("div", { class: "truncate", style: "max-width: 400px;" }, value),
    },
    {
      id: "metadata",
      header: "Metadata",
      field: "metadata",
      renderFn: (value, e) => {
        return e.metadata
          ? h(
              "div",
              { class: "flex flex-col gap-2" },
              Object.entries(e.metadata).map(([key, value]) =>
                h("span", { class: "badge badge-secondary" }, `${key}: ${value}`)
              )
            )
          : null
      },
    },
  ],
  15
)

const { data } = useGetAllGuildEventsQuery({
  variables: computed(() => table.queryState.value),
})
table.data = computed(
  () => (data.value?.guildevents.edges?.map(({ node }) => node) ?? []) as GuildEvent[]
)
</script>
