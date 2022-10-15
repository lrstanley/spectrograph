<route lang="yaml">
meta:
  layout: dashboard
  title: Users
</route>

<template>
  <div class="p-2">
    <DataTable v-model="table" v-motion-fade>
      <template v-if="data?.users.pageInfo.hasNextPage || data?.users.pageInfo.hasPreviousPage" #footer>
        <CorePagination v-model="table.cursor.value" :page-info="data?.users.pageInfo" class="ml-auto" />
      </template>
    </DataTable>
  </div>
</template>

<script setup lang="ts">
import { CoreTable } from "@/lib/util/table"
import { UserOrderField, useGetAllUsersQuery } from "@/lib/api"
import type { User, UserWhereInput } from "@/lib/api"

const table = new CoreTable<User, UserOrderField, UserWhereInput>(
  [
    {
      id: "username",
      header: "Username",
      accessorFn: (user) => `${user.username}#${user.discriminator}`,
      sortField: UserOrderField.Username,
      type: "text",
      filterFn: (val: string) => ({ usernameContainsFold: val }),
      renderFn: (value, user) =>
        h("div", { class: "flex items-center gap-2" }, [
          h("img", { src: user.avatarURL, style: "height: 20px; width: 20px" }),
          h("div", {}, value),
        ]),
    },
    {
      id: "email",
      header: "Email",
      field: "email",
      sortField: UserOrderField.Email,
      type: "text",
      filterFn: (val: string) => ({ emailContainsFold: val }),
    },
    {
      id: "locale",
      header: "Locale",
      field: "locale",
      type: "text",
      filterFn: (val: string) => ({ localeContainsFold: val }),
    },
    {
      id: "verified",
      header: "Verified",
      field: "verified",
      type: "boolean",
      filterFn: (val: boolean) => (val !== null ? { verified: val } : null),
    },
    {
      id: "mfaEnabled",
      header: "MFA",
      field: "mfaEnabled",
      type: "boolean",
      filterFn: (val: boolean) => (val !== null ? { mfaEnabled: val } : null),
    },
    {
      id: "admin",
      header: "Admin",
      field: "admin",
      type: "boolean",
      filterFn: (val: boolean) => (val !== null ? { admin: val } : null),
    },
  ],
  20
)

const { data } = useGetAllUsersQuery({
  variables: computed(() => table.queryState.value),
})
table.data = computed(() => (data.value?.users.edges?.map(({ node }) => node) ?? []) as User[])
</script>
