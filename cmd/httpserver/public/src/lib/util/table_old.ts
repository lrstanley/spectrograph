/**
 * Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
 * of this source code is governed by the MIT license that can be found in
 * the LICENSE file.
 */

import { defineComponent } from "vue"
import { OrderDirection } from "@/lib/api"
import { resetCursor, usePagination } from "@/lib/util/pagination"
import { getCoreRowModel, useVueTable } from "@tanstack/vue-table"
import { refDebounced } from "@vueuse/core"

import type { Ref, ComputedRef } from "vue"
import type {
  ColumnDef,
  SortingState,
  ColumnFiltersState,
  Table as TSTable,
  OnChangeFn,
} from "@tanstack/vue-table"

export const FlexRender = defineComponent({
  props: ["render", "props"],
  setup: (props: { render: any; props: any }) => {
    return () => {
      if (typeof props.render === "function" || typeof props.render === "object") {
        return h(props.render, props.props)
      }

      return props.render
    }
  },
})

function updateValue<T>(val: Ref<T>): OnChangeFn<T> {
  return (updaterOrValue) => {
    val.value = updaterOrValue instanceof Function ? updaterOrValue(val.value) : updaterOrValue
  }
}

export type Column<T, OrderField, WhereInput> = ColumnDef<T> & {
  meta?: {
    sort?: OrderField
    input?: string
    filter?: (val: any) => WhereInput
  }
}

export interface TableConfig<T, OrderField, WhereInput> {
  columns: Column<T, OrderField, WhereInput>[]
  pageSize: number
}

export interface Table<T, OrderField, WhereInput> extends TSTable<T> {
  config: TableConfig<T, OrderField, WhereInput>

  inputData?: Ref<T[]>
  sorting: Ref<SortingState>
  filters: Ref<ColumnFiltersState>
  cursor: Ref<string>
  queryVariables: ComputedRef<Record<string, any>>

  setInputData: (val: Ref<T[]> | ComputedRef<T[]>) => () => void
}

export function makeTable<T, OrderField, WhereInput>(
  config: TableConfig<T, OrderField, WhereInput>
): Table<T, OrderField, WhereInput> {
  const sorting = ref<SortingState>([])
  const filters = ref<ColumnFiltersState>([])
  const filtersDebounced = refDebounced(filters, 250)

  const cursor = ref<string>("")
  const pagination = usePagination(cursor, config.pageSize)
  resetCursor(cursor, [sorting, filtersDebounced])

  let table = {} as any as Table<T, OrderField, WhereInput>

  const queryVariables = computed(() => {
    const q: Record<string, any> = {}

    if (sorting.value?.length > 0) {
      const def = table.getColumn(sorting.value[0].id).columnDef as Column<T, OrderField, WhereInput>
      const field = def.meta?.sort

      if (field) {
        q.orderBy = field
        q.order = sorting.value[0].desc ? OrderDirection.Desc : OrderDirection.Asc
      }
    }

    for (const f of filtersDebounced.value) {
      const def = table.getColumn(f.id).columnDef as Column<T, OrderField, WhereInput>
      const filterFn = def.meta?.filter

      if (!filterFn) {
        continue
      }

      const filters = filterFn(f.value)

      if (f.value !== undefined && f.value !== null && filters) {
        q.where = {
          ...(q.where || {}),
          ...filters,
        }
      }
    }

    return {
      ...q,
      ...pagination,
    }
  })

  table = {
    ...table,
    config,
    sorting,
    filters,
    cursor,
    queryVariables,
    ...useVueTable<T>({
      get data() {
        return table.inputData?.value || []
      },
      get columns() {
        return config.columns.map((c: Column<T, OrderField, WhereInput>) => ({
          ...c,
          enableSorting: !!c.meta?.sort,
          enableColumnFilter: !!c.meta?.filter,
        }))
      },
      state: {
        get sorting() {
          return sorting.value
        },
        get columnFilters() {
          return filters.value
        },
      },

      onSortingChange: updateValue(sorting),
      onColumnFiltersChange: updateValue(filters),
      getCoreRowModel: getCoreRowModel(),
      enableSorting: true,
      manualSorting: true,
      enableColumnFilters: true,
      manualFiltering: true,
      enableGrouping: false,
    }),
  }

  return table
}
