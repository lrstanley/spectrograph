<script setup lang="ts">
import { Users } from "@lucide/vue"
import { createColumnHelper } from "@tanstack/vue-table"
import { useRouteQuery } from "@vueuse/router"
import type { SchemaUser, SchemaUserSortableFields } from "#open-fetch-schemas/api"
import type { DataTableFeatures } from "@/utils/data-table"

definePageMeta({ layout: "wide", title: "Users" })

const { self } = useSelf()

const filters = {
  username: useRouteQuery<string>("username", ""),
  email: useRouteQuery<string>("email", ""),
  banned: useRouteQuery<string>("banned", ""),
  admin: useRouteQuery<string>("admin", ""),
}
const defaultSort: SchemaUserSortableFields = "create_time"
const pagination = usePagination<SchemaUserSortableFields>({
  sort: defaultSort,
  order: "desc",
  resetChanged: Object.values(filters),
})

function bool(value: string) {
  return value === "true" ? true : value === "false" ? false : undefined
}

const users = useApi("/users", {
  // Deep, so ban/unban results can replace rows.
  deep: true,
  query: computed(() => ({
    page: pagination.page.value,
    per_page: pagination.perPage.value,
    sort: pagination.sort.value,
    order: pagination.order.value,
    "username.ihas": filters.username.value || undefined,
    "email.ihas": filters.email.value || undefined,
    "banned.eq": bool(filters.banned.value),
    "admin.eq": bool(filters.admin.value),
  })),
})
await users

const banTarget = ref<SchemaUser | null>(null)
const columnHelper = createColumnHelper<DataTableFeatures, SchemaUser>()

function textFilter(ref: Ref<string>, placeholder: string) {
  return {
    type: "text" as const,
    placeholder,
    debounce: 300,
    value: ref,
    onUpdate: (value: unknown) => {
      ref.value = String(value ?? "")
    },
  }
}

function boolFilter(ref: Ref<string>) {
  return {
    type: "boolean" as const,
    value: () => bool(ref.value),
    onUpdate: (value: unknown) => {
      ref.value = value === undefined ? "" : String(value)
    },
  }
}

const columns = columnHelper.columns([
  columnHelper.accessor("username", {
    header: "User",
    meta: { filter: textFilter(filters.username, "username contains"), headerClass: "w-[30%]" },
  }),
  columnHelper.accessor("email", {
    header: "Email",
    meta: { filter: textFilter(filters.email, "email contains"), cellClass: "max-w-56" },
  }),
  columnHelper.accessor("admin", {
    header: "Admin",
    enableSorting: false,
    meta: { filter: boolFilter(filters.admin), headerClass: "w-28" },
  }),
  columnHelper.accessor("banned", {
    header: "Banned",
    enableSorting: false,
    meta: { filter: boolFilter(filters.banned), headerClass: "w-28" },
  }),
  columnHelper.accessor("create_time", { header: "First login", meta: { headerClass: "w-36" } }),
  columnHelper.display({ id: "actions", header: "", enableHiding: false, meta: { headerClass: "w-24" } }),
])

// The table memoizes rows by array identity, so swap in a new array.
function replaceUser(updated: SchemaUser) {
  if (!users.data.value) return
  users.data.value.content = users.data.value.content.map((row) => (row.id === updated.id ? updated : row))
}
</script>

<template>
  <div class="flex flex-col gap-1">
    <h1 class="font-heading text-2xl font-semibold tracking-tight">Users</h1>
    <p class="text-sm text-muted-foreground">Everyone who has logged in to Spectrograph.</p>
  </div>

  <Alert v-if="users.error.value" variant="destructive">
    <AlertTitle>Couldn't load users</AlertTitle>
    <AlertDescription>{{ apiErrorMessage(users.error.value) }}</AlertDescription>
  </Alert>

  <DataTable
    v-model:sort="pagination.sort.value"
    v-model:order="pagination.order.value"
    :columns="columns"
    :data="users.data.value?.content ?? []"
    :pagination="{
      page: users.data.value?.page ?? pagination.page.value,
      pageSize: pagination.perPage.value,
      pageCount: users.data.value?.last_page ?? 1,
      total: users.data.value?.total_count ?? 0,
    }"
    :default-sort="defaultSort"
    :page-size-options="pagination.sizes"
    :loading="users.status.value === 'pending'"
    @update:pagination="
      ({ page, pageSize }) => {
        pagination.page.value = page
        pagination.perPage.value = pageSize
      }
    "
  >
    <template #cell-username="{ row }">
      <NuxtLink :to="`/dashboard/user/${row.id}`" class="flex min-w-0 items-center gap-2 hover:underline">
        <UserAvatar :user="row" class="size-6" />
        <span class="truncate font-medium">{{ displayName(row) }}</span>
        <span class="truncate text-muted-foreground">{{ userTag(row) }}</span>
      </NuxtLink>
    </template>
    <template #cell-email="{ row }">
      <span class="block truncate">{{ row.email || "—" }}</span>
    </template>
    <template #cell-admin="{ row }">
      <Badge v-if="row.admin" variant="outline">Admin</Badge>
      <span v-else class="text-muted-foreground">No</span>
    </template>
    <template #cell-banned="{ row }">
      <Badge v-if="row.banned" variant="destructive" :title="row.ban_reason">Banned</Badge>
      <span v-else class="text-muted-foreground">No</span>
    </template>
    <template #cell-create_time="{ row }">
      <span :title="formatDate(row.create_time)">{{ timeAgo(row.create_time) }}</span>
    </template>
    <template #cell-actions="{ row }">
      <Button
        v-if="!row.admin && row.id !== self?.id"
        size="xs"
        :variant="row.banned ? 'outline' : 'destructive'"
        @click="banTarget = row"
      >
        {{ row.banned ? "Unban" : "Ban" }}
      </Button>
    </template>
    <template #empty>
      <Empty class="border-0 p-2">
        <EmptyHeader>
          <EmptyMedia variant="icon"><Users /></EmptyMedia>
          <EmptyTitle>No users match</EmptyTitle>
          <EmptyDescription>Clear a filter to see more.</EmptyDescription>
        </EmptyHeader>
      </Empty>
    </template>
  </DataTable>

  <BanUserDialog :user="banTarget" @close="banTarget = null" @changed="replaceUser" />
</template>
