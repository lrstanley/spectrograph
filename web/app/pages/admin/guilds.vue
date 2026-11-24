<script setup lang="ts">
import { Server } from "@lucide/vue"
import { createColumnHelper } from "@tanstack/vue-table"
import { useRouteQuery } from "@vueuse/router"
import type { SchemaGuildRead, SchemaGuildSortableFields } from "#open-fetch-schemas/api"
import type { DataTableFeatures } from "@/utils/data-table"

definePageMeta({ layout: "wide", title: "All servers" })

const filters = {
  name: useRouteQuery<string>("name", ""),
  added: useRouteQuery<string>("added", ""),
}
const defaultSort: SchemaGuildSortableFields = "create_time"
const pagination = usePagination<SchemaGuildSortableFields>({
  sort: defaultSort,
  order: "desc",
  resetChanged: Object.values(filters),
})

const guilds = useApi("/guilds", {
  query: computed(() => ({
    page: pagination.page.value,
    per_page: pagination.perPage.value,
    sort: pagination.sort.value,
    order: pagination.order.value,
    "name.ihas": filters.name.value || undefined,
    "joinedAt.null": filters.added.value === "true" ? false : filters.added.value === "false" ? true : undefined,
  })),
})
await guilds

const columnHelper = createColumnHelper<DataTableFeatures, SchemaGuildRead>()
const columns = columnHelper.columns([
  columnHelper.accessor("name", {
    header: "Server",
    meta: {
      headerClass: "w-[35%]",
      filter: {
        type: "text",
        placeholder: "name contains",
        debounce: 300,
        value: filters.name,
        onUpdate: (value) => {
          filters.name.value = String(value ?? "")
        },
      },
    },
  }),
  columnHelper.display({
    id: "status",
    header: "Status",
    meta: {
      filter: {
        type: "select",
        placeholder: "Any",
        options: [
          { label: "Bot added", value: "true" },
          { label: "Bot not added", value: "false" },
        ],
        value: filters.added,
        onUpdate: (value) => {
          filters.added.value = String(value ?? "")
        },
      },
    },
  }),
  columnHelper.accessor("guild_id", { header: "Discord ID", enableSorting: false }),
  columnHelper.accessor("member_count", { header: "Members", meta: { headerClass: "w-28" } }),
  columnHelper.accessor("joined_at", { header: "Added", meta: { headerClass: "w-32" } }),
])
</script>

<template>
  <div class="flex flex-col gap-1">
    <h1 class="font-heading text-2xl font-semibold tracking-tight">All servers</h1>
    <p class="text-sm text-muted-foreground">Every server a Spectrograph user manages, whether or not the bot was added.</p>
  </div>

  <Alert v-if="guilds.error.value" variant="destructive">
    <AlertTitle>Couldn't load servers</AlertTitle>
    <AlertDescription>{{ apiErrorMessage(guilds.error.value) }}</AlertDescription>
  </Alert>

  <DataTable
    v-model:sort="pagination.sort.value"
    v-model:order="pagination.order.value"
    :columns="columns"
    :data="guilds.data.value?.content ?? []"
    :pagination="{
      page: guilds.data.value?.page ?? pagination.page.value,
      pageSize: pagination.perPage.value,
      pageCount: guilds.data.value?.last_page ?? 1,
      total: guilds.data.value?.total_count ?? 0,
    }"
    :default-sort="defaultSort"
    :page-size-options="pagination.sizes"
    :loading="guilds.status.value === 'pending'"
    @update:pagination="
      ({ page, pageSize }) => {
        pagination.page.value = page
        pagination.perPage.value = pageSize
      }
    "
  >
    <template #cell-name="{ row }">
      <NuxtLink :to="`/dashboard/guild/${row.id}`" class="flex min-w-0 items-center gap-2 hover:underline">
        <GuildIcon :guild="row" class="size-6 text-[10px] ring-0 ring-offset-0" />
        <span class="truncate font-medium">{{ row.name }}</span>
      </NuxtLink>
    </template>
    <template #cell-status="{ row }">
      <Badge
        :variant="guildStatus(row) === 'active' ? 'default' : guildStatus(row) === 'not-joined' ? 'outline' : 'destructive'"
      >
        {{ guildStatusLabel[guildStatus(row)] }}
      </Badge>
    </template>
    <template #cell-guild_id="{ row }">
      <span class="font-mono text-xs text-muted-foreground">{{ row.guild_id }}</span>
    </template>
    <template #cell-member_count="{ row }">
      <span class="tabular-nums">{{ (row.member_count ?? 0).toLocaleString() }}</span>
    </template>
    <template #cell-joined_at="{ row }">
      <span :title="formatDate(row.joined_at)">{{ row.joined_at ? timeAgo(row.joined_at) : "—" }}</span>
    </template>
    <template #empty>
      <Empty class="border-0 p-2">
        <EmptyHeader>
          <EmptyMedia variant="icon"><Server /></EmptyMedia>
          <EmptyTitle>No servers match</EmptyTitle>
          <EmptyDescription>Clear a filter to see more.</EmptyDescription>
        </EmptyHeader>
      </Empty>
    </template>
  </DataTable>
</template>
