import type { CellData, RowData, TableFeatures } from "@tanstack/table-core"
import {
  columnFilteringFeature,
  columnVisibilityFeature,
  createFilteredRowModel,
  createPaginatedRowModel,
  createSortedRowModel,
  filterFns,
  rowPaginationFeature,
  rowSelectionFeature,
  rowSortingFeature,
  sortFns,
  tableFeatures,
} from "@tanstack/vue-table"
import type { MaybeRefOrGetter } from "vue"

export type DataTableFilterType = "text" | "number" | "boolean" | "select"

export interface DataTableFilterOption {
  label: string
  value: string
}

export interface DataTableFilterConfig {
  type: DataTableFilterType
  placeholder?: string
  options?: readonly DataTableFilterOption[]
  debounce?: number
  value?: MaybeRefOrGetter<unknown>
  onUpdate?: (value: unknown) => void
}

export interface DataTableColumnMeta {
  label?: string
  filter?: DataTableFilterConfig
  cellClass?: string
  headerClass?: string
}

export interface DataTablePagination {
  page: number
  pageSize: number
  pageCount?: number
  total?: number
}

declare module "@tanstack/table-core" {
  interface ColumnMeta<TFeatures extends TableFeatures, TData extends RowData, TValue extends CellData>
    extends DataTableColumnMeta {}
}

export const dataTableFeatures = tableFeatures({
  columnFilteringFeature,
  columnVisibilityFeature,
  rowPaginationFeature,
  rowSelectionFeature,
  rowSortingFeature,
  filteredRowModel: createFilteredRowModel(),
  paginatedRowModel: createPaginatedRowModel(),
  sortedRowModel: createSortedRowModel(),
  filterFns,
  sortFns,
})

export type DataTableFeatures = typeof dataTableFeatures
