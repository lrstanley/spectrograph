<template>
  <table class="table w-full table-compact">
    <thead>
      <tr v-for="headerGroup in table.getHeaderGroups()" :key="headerGroup.id">
        <th
          v-for="header in headerGroup.headers"
          :key="header.id"
          :colSpan="header.colSpan"
          class="first:pl-5"
        >
          <div
            :class="header.column.getCanSort() ? 'cursor-pointer select-none flex' : ''"
            @click="header.column.getToggleSortingHandler()?.($event)"
          >
            <FlexRender
              v-if="!header.isPlaceholder"
              :render="header.column.columnDef.header"
              :props="header.getContext()"
            />

            <span v-if="header.column.getCanSort()" class="inline-flex items-center ml-auto">
              <i-fas-sort-up
                v-if="header.column.getIsSorted() === 'asc'"
                class="w-3 h-3 text-nitro-500"
              />
              <i-fas-sort-down
                v-else-if="header.column.getIsSorted() === 'desc'"
                class="w-3 h-3 text-nitro-500"
              />
              <i-fas-sort v-else class="w-3 h-3" />
            </span>
          </div>

          <div
            v-if="table.getAllColumns().some((c) => c.getCanFilter())"
            :class="header.column.getCanFilter() ? '' : 'invisible'"
          >
            <!-- TODO: change boolean to select with n/a, true, and false (for indeterminate)? -->
            <select
              v-if="(header.column.columnDef as Column<any, any, any>).meta?.input === 'boolean'"
              class="w-full max-w-xs py-0 select select-bordered select-xs"
              @change="
                header.column.setFilterValue(
                  ($event.target as HTMLSelectElement).value === 'true'
                    ? true
                    : ($event.target as HTMLSelectElement).value === 'false'
                    ? false
                    : null
                )
              "
            >
              <option>any</option>
              <option>true</option>
              <option>false</option>
            </select>
            <input
              v-else
              type="text"
              placeholder="Search..."
              class="w-full max-w-xs input input-bordered input-xs"
              @input="header.column.setFilterValue(($event.target as HTMLInputElement).value)"
            />
          </div>
        </th>
      </tr>
    </thead>
    <tbody>
      <tr v-for="row in table.getRowModel().rows" :key="row.id">
        <td
          v-for="cell in row.getVisibleCells()"
          :key="cell.id"
          class="first:pl-5"
          :class="{
            'text-dnd-400 font-bold': typeof cell.getValue() === 'boolean' && !cell.getValue(),
            'text-online-400 font-bold': typeof cell.getValue() === 'boolean' && cell.getValue(),
          }"
        >
          <FlexRender :render="cell.column.columnDef.cell" :props="cell.getContext()" />
        </td>
      </tr>
    </tbody>
    <tfoot>
      <tr v-for="footerGroup in table.getFooterGroups()" :key="footerGroup.id">
        <th v-for="header in footerGroup.headers" :key="header.id" :colSpan="header.colSpan">
          <FlexRender
            v-if="!header.isPlaceholder"
            :render="header.column.columnDef.footer"
            :props="header.getContext()"
          />
        </th>
      </tr>

      <tr v-if="props.pageInfo">
        <CorePagination v-model="cursor" :page-info="props.pageInfo" />
      </tr>
    </tfoot>
  </table>
</template>

<script setup lang="ts">
import { useVModel } from "@vueuse/core"
import { FlexRender } from "@/lib/util/table"
import type { Table, Column } from "@/lib/util/table"
import type { PageInfo } from "@/lib/api"

const props = defineProps<{
  modelValue: Table<any, any, any>
  cursor?: string
  pageInfo?: PageInfo
}>()

const emit = defineEmits<{
  (e: "update:modelValue", value: Table<any, any, any>): void
  (e: "update:cursor", value: string): void
}>()

const table = useVModel(props, "modelValue", emit)
const cursor = useVModel(props, "cursor", emit)
// const table = computed({
//   get: () => props.modelValue,
//   set: (val) => emit("update:modelValue", val),
// })
</script>
