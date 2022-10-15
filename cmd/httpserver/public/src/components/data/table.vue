<template>
  <div class="rounded bg-base-200">
    <table class="table w-full table-compact">
      <thead>
        <tr>
          <th v-for="col in table.columns" :key="col.id" :colSpan="1" class="first:pl-5">
            <div :class="col.canSort() ? 'cursor-pointer select-none flex' : ''" @click="col.toggleSort">
              <span>{{ col.header }}</span>

              <span v-if="col.canSort()" class="inline-flex items-center ml-auto">
                <i-fas-sort-up
                  v-if="col.isSorted() === OrderDirection.Asc"
                  class="w-3 h-3 text-nitro-500"
                />
                <i-fas-sort-down
                  v-else-if="col.isSorted() === OrderDirection.Desc"
                  class="w-3 h-3 text-nitro-500"
                />
                <i-fas-sort v-else class="w-3 h-3" />
              </span>
            </div>

            <div v-if="table.isFiltered()" :class="table.canFilter() ? '' : 'invisible'">
              <select
                v-if="col.type === 'boolean'"
                v-model="col.filterValue.value"
                class="w-full max-w-xs py-0 select select-bordered select-xs"
              >
                <option :value="null">any</option>
                <option :value="true">true</option>
                <option :value="false">false</option>
              </select>
              <input
                v-else-if="col.type === 'text'"
                v-model="col.filterValue.value"
                type="text"
                placeholder="Search..."
                class="w-full max-w-xs input input-bordered input-xs"
              />
              <input v-else class="invisible w-full max-w-xs input input-bordered input-xs" />
            </div>
          </th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="row in table.rows.value" :key="row.id">
          <td
            v-for="cell in row.cells.value"
            :key="cell.column.field"
            class="first:pl-5"
            :class="{
              'text-dnd-400 font-bold': cell.column.type === 'boolean' && !cell.value,
              'text-online-400 font-bold': cell.column.type === 'boolean' && cell.value,
            }"
          >
            <component :is="cell.render()" />
          </td>
        </tr>
      </tbody>
    </table>
    <!-- no data or no results found -->
    <div
      v-if="table.rows.value?.length < 1"
      class="flex flex-row items-center justify-center gap-3 px-10 py-4"
    >
      <div class="badge badge-secondary">0</div>
      <div>results found with provided search query</div>
    </div>
    <div v-if="$slots.footer" class="flex flex-row items-center gap-3 p-2">
      <slot name="footer" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { OrderDirection } from "@/lib/api"
import type { Table } from "@/lib/util/table"

const props = defineProps<{
  modelValue: Table<any, any, any>
}>()

const table = computed(() => props.modelValue)
</script>
