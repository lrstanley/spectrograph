<script setup lang="ts">
import { ListTree } from "@lucide/vue"
import { createColumnHelper } from "@tanstack/vue-table"
import { useRouteQuery } from "@vueuse/router"
import type {
  SchemaGuildEventRead,
  SchemaGuildEventSortableFields,
  SchemaGuildEventTypeEnum,
} from "#open-fetch-schemas/api"
import type { DataTableFeatures } from "@/utils/data-table"

definePageMeta({ layout: "wide", title: "Event log" })

const types: SchemaGuildEventTypeEnum[] = ["DEBUG", "INFO", "WARNING", "ERROR"]

const filters = {
  type: useRouteQuery<string>("type", ""),
  message: useRouteQuery<string>("message", ""),
}
const defaultSort: SchemaGuildEventSortableFields = "create_time"
const pagination = usePagination<SchemaGuildEventSortableFields>({
  perPage: 25,
  sort: defaultSort,
  order: "desc",
  resetChanged: Object.values(filters),
})

const events = useApi("/guild-events", {
  query: computed(() => ({
    page: pagination.page.value,
    per_page: pagination.perPage.value,
    sort: pagination.sort.value,
    order: pagination.order.value,
    "type.eq": (filters.type.value || undefined) as SchemaGuildEventTypeEnum | undefined,
    "message.ihas": filters.message.value || undefined,
  })),
})
await events

// Slot names can't contain a "." in templates.
const serverSlot = "cell-guild.name"

const columnHelper = createColumnHelper<DataTableFeatures, SchemaGuildEventRead>()
const columns = columnHelper.columns([
  columnHelper.accessor("create_time", { header: "Time", meta: { headerClass: "w-36" } }),
  // The ID matches the API's sort field.
  columnHelper.accessor((row) => row.edges.guild.name, { id: "guild.name", header: "Server", meta: { headerClass: "w-52" } }),
  columnHelper.accessor("type", {
    header: "Type",
    enableSorting: false,
    meta: {
      headerClass: "w-32",
      filter: {
        type: "select",
        placeholder: "Any",
        options: types.map((value) => ({ label: value.toLowerCase(), value })),
        value: filters.type,
        onUpdate: (value) => {
          filters.type.value = String(value ?? "")
        },
      },
    },
  }),
  columnHelper.accessor("message", {
    header: "Message",
    enableSorting: false,
    meta: {
      cellClass: "max-w-0 w-[40%]",
      filter: {
        type: "text",
        placeholder: "message contains",
        debounce: 300,
        value: filters.message,
        onUpdate: (value) => {
          filters.message.value = String(value ?? "")
        },
      },
    },
  }),
  columnHelper.accessor("metadata", { header: "Details", enableSorting: false, meta: { cellClass: "max-w-0" } }),
])
</script>

<template>
  <div class="flex flex-col gap-1">
    <h1 class="font-heading text-2xl font-semibold tracking-tight">Event log</h1>
    <p class="text-sm text-muted-foreground">Events from every server, newest first. Events are kept for 7 days.</p>
  </div>

  <Alert v-if="events.error.value" variant="destructive">
    <AlertTitle>Couldn't load events</AlertTitle>
    <AlertDescription>{{ apiErrorMessage(events.error.value) }}</AlertDescription>
  </Alert>

  <DataTable
    v-model:sort="pagination.sort.value"
    v-model:order="pagination.order.value"
    :columns="columns"
    :data="events.data.value?.content ?? []"
    :pagination="{
      page: events.data.value?.page ?? pagination.page.value,
      pageSize: pagination.perPage.value,
      pageCount: events.data.value?.last_page ?? 1,
      total: events.data.value?.total_count ?? 0,
    }"
    :default-sort="defaultSort"
    :page-size-options="pagination.sizes"
    :loading="events.status.value === 'pending'"
    table-class="min-w-4xl table-fixed"
    @update:pagination="
      ({ page, pageSize }) => {
        pagination.page.value = page
        pagination.perPage.value = pageSize
      }
    "
  >
    <template #cell-create_time="{ row }">
      <time :datetime="row.create_time" :title="formatDate(row.create_time)" class="text-muted-foreground">
        {{ timeAgo(row.create_time) }}
      </time>
    </template>
    <template #[serverSlot]="{ row }">
      <NuxtLink
        :to="`/dashboard/guild/${row.edges.guild.id}`"
        class="flex min-w-0 items-center gap-2 font-medium hover:underline"
        :title="row.edges.guild.name"
      >
        <GuildIcon :guild="row.edges.guild" compact class="size-6 text-[10px]" aria-hidden="true" />
        <span class="truncate">{{ row.edges.guild.name }}</span>
      </NuxtLink>
    </template>
    <template #cell-type="{ row }">
      <Badge :variant="row.type === 'ERROR' ? 'destructive' : row.type === 'WARNING' ? 'default' : 'outline'">
        {{ row.type.toLowerCase() }}
      </Badge>
    </template>
    <template #cell-message="{ row }">
      <span class="block truncate" :title="row.message">{{ row.message }}</span>
    </template>
    <template #cell-metadata="{ row }">
      <span class="block truncate font-mono text-xs text-muted-foreground">
        {{
          Object.entries(row.metadata ?? {})
            .map(([key, value]) => `${key}=${typeof value === "object" ? JSON.stringify(value) : value}`)
            .join(" ") || "—"
        }}
      </span>
    </template>
    <template #empty>
      <Empty class="border-0 p-2">
        <EmptyHeader>
          <EmptyMedia variant="icon"><ListTree /></EmptyMedia>
          <EmptyTitle>No events match</EmptyTitle>
          <EmptyDescription>Clear a filter, or wait for the bot to do something.</EmptyDescription>
        </EmptyHeader>
      </Empty>
    </template>
  </DataTable>
</template>
