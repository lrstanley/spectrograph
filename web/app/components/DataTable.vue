<script setup lang="ts" generic="TData extends object, TSort extends string = string">
import { ArrowDown, ArrowUp, ArrowUpDown, ChevronLeft, ChevronRight, Columns3 } from "@lucide/vue"
import type {
  Column,
  ColumnDef,
  ColumnFiltersState,
  ColumnVisibilityState,
  Header,
  PaginationState,
  RowSelectionState,
  SortingState,
  Updater,
} from "@tanstack/vue-table"
import {
  FlexRender,
  useTable,
} from "@tanstack/vue-table"
import type { HTMLAttributes } from "vue"
import { computed, h, onBeforeUnmount, reactive, ref, toValue, watch } from "vue"
import Checkbox from "@/components/ui/checkbox/Checkbox.vue"
import {
  Table,
  TableBody,
  TableCell,
  TableEmpty,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
import type {
  DataTableFeatures,
  DataTableFilterConfig,
  DataTablePagination,
} from "@/utils/data-table"
import { dataTableFeatures } from "@/utils/data-table"

const props = withDefaults(defineProps<{
  columns: ColumnDef<DataTableFeatures, TData>[]
  data: TData[]
  pagination?: DataTablePagination
  pageSizeOptions?: readonly number[]
  initialPageSize?: number
  initialSorting?: SortingState
  initialColumnFilters?: ColumnFiltersState
  initialColumnVisibility?: ColumnVisibilityState
  sorting?: SortingState
  sort?: TSort
  order?: "asc" | "desc"
  defaultSort?: TSort
  defaultOrder?: "asc" | "desc"
  columnFilters?: ColumnFiltersState
  columnVisibility?: ColumnVisibilityState
  rowSelection?: RowSelectionState
  selectable?: boolean
  rowId?: (row: TData, index: number) => string
  showColumnVisibility?: boolean
  showPagination?: boolean
  tableClass?: HTMLAttributes["class"]
  loading?: boolean
  emptyText?: string
}>(), {
  pageSizeOptions: () => [10, 20, 50, 100],
  initialPageSize: 10,
  defaultOrder: "desc",
  selectable: false,
  showColumnVisibility: true,
  showPagination: true,
  loading: false,
  emptyText: "No results.",
})

const emit = defineEmits<{
  "update:pagination": [value: DataTablePagination]
  "update:sorting": [value: SortingState]
  "update:sort": [value: TSort | undefined]
  "update:order": [value: "asc" | "desc"]
  "update:columnFilters": [value: ColumnFiltersState]
  "update:columnVisibility": [value: ColumnVisibilityState]
  "update:rowSelection": [value: RowSelectionState]
  "selection-change": [rows: TData[]]
}>()

function resolveUpdater<T>(updater: Updater<T>, current: T): T {
  return typeof updater === "function"
    ? (updater as (value: T) => T)(current)
    : updater
}

const internalPagination = ref<PaginationState>({
  pageIndex: 0,
  pageSize: props.pagination?.pageSize ?? props.initialPageSize,
})
const internalSorting = ref<SortingState>(props.initialSorting ?? [])
const internalColumnFilters = ref<ColumnFiltersState>(props.initialColumnFilters ?? [])
const internalColumnVisibility = ref<ColumnVisibilityState>(props.initialColumnVisibility ?? {})
const internalRowSelection = ref<RowSelectionState>(props.rowSelection ?? {})
const filterDrafts = reactive<Record<string, unknown>>({})
const filterTimers = new Map<string, ReturnType<typeof setTimeout>>()

const selectionColumn: ColumnDef<DataTableFeatures, TData> = {
  id: "__select",
  header: ({ table }) => h(Checkbox, {
    "modelValue": table.getIsAllPageRowsSelected()
      || (table.getIsSomePageRowsSelected() && "indeterminate"),
    "onUpdate:modelValue": (value: boolean | "indeterminate") => {
      table.toggleAllPageRowsSelected(value === true)
    },
    "aria-label": "Select all rows",
  }),
  cell: ({ row }) => h(Checkbox, {
    "modelValue": row.getIsSelected(),
    "onUpdate:modelValue": (value: boolean | "indeterminate") => {
      row.toggleSelected(value === true)
    },
    "aria-label": "Select row",
  }),
  enableSorting: false,
  enableHiding: false,
  enableColumnFilter: false,
  meta: {
    label: "Select",
    cellClass: "w-10",
  },
}

const tableColumns = computed(() => props.selectable
  ? [selectionColumn, ...props.columns]
  : props.columns)

function getRowId(row: TData, index: number) {
  if (props.rowId) {
    return props.rowId(row, index)
  }
  if ("id" in row && (typeof row.id === "string" || typeof row.id === "number")) {
    return String(row.id)
  }
  return String(index)
}

const controlledPagination = computed<PaginationState>(() => props.pagination
  ? {
      pageIndex: Math.max(0, props.pagination.page - 1),
      pageSize: props.pagination.pageSize,
    }
  : internalPagination.value)

const controlledSorting = computed<SortingState>(() => {
  if (props.sorting !== undefined) {
    return props.sorting
  }
  if (props.sort !== undefined) {
    const hasColumn = tableColumns.value.some(column =>
      column.id === props.sort
      || ("accessorKey" in column && column.accessorKey === props.sort),
    )
    return hasColumn
      ? [{ id: props.sort, desc: (props.order ?? props.defaultOrder) === "desc" }]
      : []
  }
  return internalSorting.value
})
const controlledColumnFilters = computed(() => props.columnFilters ?? internalColumnFilters.value)
const controlledColumnVisibility = computed(() => props.columnVisibility ?? internalColumnVisibility.value)
const controlledRowSelection = computed(() => props.rowSelection ?? internalRowSelection.value)

const table = useTable({
  features: dataTableFeatures,
  get data() {
    return props.data
  },
  get columns() {
    return tableColumns.value
  },
  state: {
    get pagination() {
      return controlledPagination.value
    },
    get sorting() {
      return controlledSorting.value
    },
    get columnFilters() {
      return controlledColumnFilters.value
    },
    get columnVisibility() {
      return controlledColumnVisibility.value
    },
    get rowSelection() {
      return controlledRowSelection.value
    },
  },
  initialState: {
    pagination: {
      pageIndex: 0,
      pageSize: props.pagination?.pageSize ?? props.initialPageSize,
    },
    sorting: props.initialSorting,
    columnFilters: props.initialColumnFilters,
    columnVisibility: props.initialColumnVisibility,
    rowSelection: props.rowSelection,
  },
  getRowId,
  get pageCount() {
    if (props.pagination?.pageCount != null) {
      return Math.max(1, props.pagination.pageCount)
    }
    if (props.pagination?.total != null) {
      return Math.max(1, Math.ceil(props.pagination.total / props.pagination.pageSize))
    }
    return undefined
  },
  manualPagination: true,
  manualSorting: true,
  manualFiltering: true,
  get enableRowSelection() {
    return props.selectable
  },
  onPaginationChange: (updater) => {
    const next = resolveUpdater(updater, controlledPagination.value)
    if (props.pagination) {
      emit("update:pagination", {
        ...props.pagination,
        page: next.pageIndex + 1,
        pageSize: next.pageSize,
      })
    }
    else {
      internalPagination.value = next
    }
  },
  onSortingChange: (updater) => {
    const next = resolveUpdater(updater, controlledSorting.value)
    if (props.sort !== undefined && props.sorting === undefined) {
      const first = next[0]
      emit("update:sort", (first?.id ?? props.defaultSort) as TSort | undefined)
      emit("update:order", first
        ? (first.desc ? "desc" : "asc")
        : props.defaultOrder)
      return
    }
    internalSorting.value = next
    emit("update:sorting", next)
  },
  onColumnFiltersChange: (updater) => {
    const next = resolveUpdater(updater, controlledColumnFilters.value)
    internalColumnFilters.value = next
    emit("update:columnFilters", next)
  },
  onColumnVisibilityChange: (updater) => {
    const next = resolveUpdater(updater, controlledColumnVisibility.value)
    internalColumnVisibility.value = next
    emit("update:columnVisibility", next)
  },
  onRowSelectionChange: (updater) => {
    const next = resolveUpdater(updater, controlledRowSelection.value)
    internalRowSelection.value = next
    emit("update:rowSelection", next)
    emit(
      "selection-change",
      table.getCoreRowModel().rows
        .filter(row => next[row.id])
        .map(row => row.original),
    )
  },
})

const tablePagination = computed(() => table.atoms.pagination.get())
const currentPage = computed(() => tablePagination.value.pageIndex + 1)
const pageCount = computed(() => {
  if (props.pagination?.pageCount != null) {
    return Math.max(1, props.pagination.pageCount)
  }
  if (props.pagination?.total != null) {
    return Math.max(1, Math.ceil(props.pagination.total / tablePagination.value.pageSize))
  }
  return Math.max(1, table.getPageCount())
})
const visibleColumns = computed(() => table.getVisibleLeafColumns())
const selectableCount = computed(() => Object.values(controlledRowSelection.value).filter(Boolean).length)
const selectedRows = computed(() => table.getSelectedRowModel().rows.map(row => row.original))

function columnMeta(column: Column<DataTableFeatures, TData>) {
  return column.columnDef.meta
}

function columnLabel(column: Column<DataTableFeatures, TData>) {
  const meta = columnMeta(column)
  const header = column.columnDef.header
  return meta?.label
    ?? (typeof header === "string" ? header : column.id.replaceAll("_", " "))
}

function headerLabel(header: Header<DataTableFeatures, TData>) {
  const value = header.column.columnDef.header
  return typeof value === "string" ? value : columnLabel(header.column)
}

function hasCustomHeader(header: Header<DataTableFeatures, TData>) {
  return typeof header.column.columnDef.header === "function"
}

function hasCustomCell(column: Column<DataTableFeatures, TData>) {
  return column.columnDef.cell !== undefined
}

function formatCellValue(value: unknown) {
  if (value == null || value === "") {
    return "—"
  }
  if (typeof value === "object") {
    try {
      return JSON.stringify(value)
    }
    catch {
      return String(value)
    }
  }
  return String(value)
}

function toggleSorting(column: Column<DataTableFeatures, TData>) {
  const nextOrder = column.getNextSortingOrder()
  if (nextOrder === false) {
    column.clearSorting()
    return
  }
  column.toggleSorting(nextOrder === "desc")
}

function ariaSort(column: Column<DataTableFeatures, TData>) {
  if (!column.getCanSort()) {
    return undefined
  }
  const sorted = column.getIsSorted()
  return sorted === "asc" ? "ascending" : sorted === "desc" ? "descending" : "none"
}

function filterLabel(column: Column<DataTableFeatures, TData>) {
  return `Filter by ${columnLabel(column).toLowerCase()}`
}

function filterConfig(column: Column<DataTableFeatures, TData>) {
  return columnMeta(column)?.filter
}

function currentFilterValue(column: Column<DataTableFeatures, TData>) {
  if (Object.hasOwn(filterDrafts, column.id)) {
    return filterDrafts[column.id]
  }
  const value = filterConfig(column)?.value
  if (value !== undefined) {
    return toValue(value)
  }
  return column.getFilterValue()
}

function booleanFilterValue(column: Column<DataTableFeatures, TData>) {
  const value = currentFilterValue(column)
  return value === true ? "true" : value === false ? "false" : "any"
}

function selectFilterValue(column: Column<DataTableFeatures, TData>) {
  const value = currentFilterValue(column)
  return value == null || value === "" ? "all" : String(value)
}

function normalizeFilterValue(value: unknown, config: DataTableFilterConfig) {
  if (value === "" || value === "any" || value === "all" || value == null) {
    return undefined
  }
  if (config.type === "number") {
    const number = Number(value)
    return Number.isFinite(number) ? number : undefined
  }
  if (config.type === "boolean") {
    return value === true || value === "true"
  }
  return value
}

function commitFilter(column: Column<DataTableFeatures, TData>, value: unknown) {
  const config = filterConfig(column)
  if (config) {
    const next = normalizeFilterValue(value, config)
    column.setFilterValue(next)
    config.onUpdate?.(next)
    delete filterDrafts[column.id]
  }
}

function filterOptions(column: Column<DataTableFeatures, TData>) {
  return filterConfig(column)?.options ?? []
}

function setColumnFilter(
  column: Column<DataTableFeatures, TData>,
  value: unknown,
  immediate = false,
) {
  const config = filterConfig(column)
  if (!config) {
    return
  }
  filterDrafts[column.id] = value
  const existingTimer = filterTimers.get(column.id)
  if (existingTimer) {
    clearTimeout(existingTimer)
    filterTimers.delete(column.id)
  }
  const debounce = immediate ? 0 : Math.max(0, config.debounce ?? 0)
  if (!debounce) {
    commitFilter(column, value)
    return
  }
  filterTimers.set(column.id, setTimeout(() => {
    filterTimers.delete(column.id)
    commitFilter(column, value)
  }, debounce))
}

function setBooleanFilter(column: Column<DataTableFeatures, TData>, value: unknown) {
  setColumnFilter(column, value, true)
}

function goToPage(page: number) {
  const next = Math.min(Math.max(1, page), pageCount.value)
  if (props.pagination) {
    emit("update:pagination", {
      ...props.pagination,
      page: next,
    })
  }
  else {
    table.setPageIndex(next - 1)
  }
}

function setPageSize(pageSize: number) {
  if (!props.pageSizeOptions.includes(pageSize)) {
    return
  }
  if (props.pagination) {
    emit("update:pagination", {
      ...props.pagination,
      page: 1,
      pageSize,
    })
  }
  else {
    table.setPageSize(pageSize)
    table.setPageIndex(0)
  }
}

function isPageActive(page: number) {
  return page === currentPage.value
}

function clearSelection() {
  const next = {}
  internalRowSelection.value = next
  emit("update:rowSelection", next)
  emit("selection-change", [])
}

watch(() => props.columnFilters, () => {
  for (const key of Object.keys(filterDrafts)) {
    delete filterDrafts[key]
  }
}, { deep: true })

onBeforeUnmount(() => {
  for (const timer of filterTimers.values()) {
    clearTimeout(timer)
  }
})

defineExpose({
  table,
  selectedRows,
  clearSelection,
})
</script>

<template>
  <div class="flex min-w-0 flex-col gap-3" :aria-busy="loading">
    <div
      v-if="$slots.toolbar"
      class="flex flex-wrap items-center gap-2"
    >
      <slot
        name="toolbar"
        :table="table"
        :selected-rows="selectedRows"
        :selected-count="selectableCount"
      />
    </div>

    <div class="min-w-0 overflow-hidden rounded-none border">
      <Table
        :class="[tableClass, { 'opacity-60': loading && data.length > 0 }]"
      >
        <TableHeader>
          <TableRow
            v-for="headerGroup in table.getHeaderGroups()"
            :key="headerGroup.id"
          >
            <TableHead
              v-for="header in headerGroup.headers"
              :key="header.id"
              :colspan="header.colSpan"
              :class="columnMeta(header.column)?.headerClass"
              :aria-sort="ariaSort(header.column)"
            >
              <slot
                :name="`header-${header.column.id}`"
                :column="header.column"
                :header="header"
              >
                <template v-if="header.isPlaceholder" />
                <template v-else-if="header.column.getCanSort() && !hasCustomHeader(header)">
                  <Button
                    variant="ghost"
                    size="sm"
                    class="-ml-2 max-w-full truncate px-2"
                    :title="headerLabel(header)"
                    @click="toggleSorting(header.column)"
                  >
                    <span class="truncate">{{ headerLabel(header) }}</span>
                    <ArrowUp v-if="header.column.getIsSorted() === 'asc'" data-icon="inline-end" />
                    <ArrowDown v-else-if="header.column.getIsSorted() === 'desc'" data-icon="inline-end" />
                    <ArrowUpDown v-else data-icon="inline-end" />
                  </Button>
                </template>
                <FlexRender v-else-if="!header.isPlaceholder" :header="header" />
              </slot>
            </TableHead>
          </TableRow>
          <TableRow v-if="visibleColumns.some(column => filterConfig(column))">
            <TableHead
              v-for="column in visibleColumns"
              :key="`filter-${column.id}`"
              :class="['h-8 px-1 py-0.5', columnMeta(column)?.headerClass]"
            >
              <slot
                v-if="filterConfig(column)"
                :name="`filter-${column.id}`"
                :column="column"
                :value="currentFilterValue(column)"
                :set-value="(value: unknown) => setColumnFilter(column, value, true)"
              >
                <template v-if="filterConfig(column)?.type === 'select'">
                  <Select
                    :model-value="selectFilterValue(column)"
                    @update:model-value="setColumnFilter(column, $event, true)"
                  >
                    <SelectTrigger class="h-7 w-full px-2 py-0.5" size="sm" :aria-label="filterLabel(column)">
                      <SelectValue :placeholder="filterConfig(column)?.placeholder ?? 'Any'" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectGroup>
                        <SelectItem value="all">Any</SelectItem>
                        <SelectItem
                          v-for="option in filterOptions(column)"
                          :key="option.value"
                          :value="option.value"
                        >
                          {{ option.label }}
                        </SelectItem>
                      </SelectGroup>
                    </SelectContent>
                  </Select>
                </template>
                <Select
                  v-else-if="filterConfig(column)?.type === 'boolean'"
                  :model-value="booleanFilterValue(column)"
                  @update:model-value="setBooleanFilter(column, $event)"
                >
                  <SelectTrigger class="h-7 w-full px-2 py-0.5" size="sm" :aria-label="filterLabel(column)">
                    <SelectValue :placeholder="filterConfig(column)?.placeholder ?? 'Any'" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectGroup>
                      <SelectItem value="any">Any</SelectItem>
                      <SelectItem value="true">True</SelectItem>
                      <SelectItem value="false">False</SelectItem>
                    </SelectGroup>
                  </SelectContent>
                </Select>
                <Input
                  v-else
                  class="h-7 px-2 py-0.5"
                  :model-value="String(currentFilterValue(column) ?? '')"
                  :type="filterConfig(column)?.type === 'number' ? 'number' : 'search'"
                  :placeholder="filterConfig(column)?.placeholder"
                  :aria-label="filterLabel(column)"
                  :inputmode="filterConfig(column)?.type === 'number' ? 'decimal' : undefined"
                  @update:model-value="setColumnFilter(column, $event)"
                />
              </slot>
            </TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <template v-if="loading && data.length === 0">
            <TableRow v-for="index in Math.min(tablePagination.pageSize, 5)" :key="`loading-${index}`">
              <TableCell v-for="column in visibleColumns" :key="column.id">
                <Skeleton class="h-4 w-full" />
              </TableCell>
            </TableRow>
          </template>
          <template v-else-if="table.getRowModel().rows.length">
            <TableRow
              v-for="row in table.getRowModel().rows"
              :key="row.id"
              :data-state="row.getIsSelected() ? 'selected' : undefined"
            >
              <TableCell
                v-for="cell in row.getVisibleCells()"
                :key="cell.id"
                :class="columnMeta(cell.column)?.cellClass"
              >
                <slot
                  :name="`cell-${cell.column.id}`"
                  :row="row.original"
                  :value="cell.getValue()"
                  :cell="cell"
                >
                  <FlexRender v-if="hasCustomCell(cell.column)" :cell="cell" />
                  <span v-else class="block max-w-full truncate">
                    {{ formatCellValue(cell.getValue()) }}
                  </span>
                </slot>
              </TableCell>
            </TableRow>
          </template>
          <TableEmpty v-else :colspan="visibleColumns.length">
            <slot name="empty">
              {{ emptyText }}
            </slot>
          </TableEmpty>
        </TableBody>
      </Table>
    </div>

    <div
      v-if="showPagination && (pagination || data.length > 0)"
      class="flex flex-wrap items-center justify-between gap-3"
    >
      <div class="text-muted-foreground text-xs">
        <span v-if="pagination?.total != null">
          {{ pagination.total }} result{{ pagination.total === 1 ? "" : "s" }}
        </span>
        <span v-else>
          {{ selectableCount }} selected
        </span>
      </div>
      <div class="flex flex-wrap items-center gap-3">
        <DropdownMenu v-if="showColumnVisibility">
          <DropdownMenuTrigger as-child>
            <Button variant="outline" size="sm">
              <Columns3 data-icon="inline-start" />
              Columns
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="start" class="w-44">
            <DropdownMenuLabel>Show columns</DropdownMenuLabel>
            <DropdownMenuSeparator />
            <DropdownMenuCheckboxItem
              v-for="column in table.getAllLeafColumns().filter(item => item.getCanHide())"
              :key="column.id"
              :model-value="column.getIsVisible()"
              @update:model-value="column.toggleVisibility(!!$event)"
            >
              {{ columnLabel(column) }}
            </DropdownMenuCheckboxItem>
          </DropdownMenuContent>
        </DropdownMenu>
        <div class="text-muted-foreground flex items-center gap-2 text-xs">
          <span>Rows</span>
          <Select
            :model-value="String(tablePagination.pageSize)"
            @update:model-value="setPageSize(Number($event))"
          >
            <SelectTrigger class="w-18" size="sm">
              <SelectValue />
            </SelectTrigger>
            <SelectContent position="popper" side="top">
              <SelectGroup>
                <SelectItem
                  v-for="size in pageSizeOptions"
                  :key="size"
                  :value="String(size)"
                >
                  {{ size }}
                </SelectItem>
              </SelectGroup>
            </SelectContent>
          </Select>
        </div>
        <span class="text-muted-foreground min-w-24 text-center text-xs">
          Page {{ currentPage }} of {{ pageCount }}
        </span>
        <div class="flex items-center gap-1">
          <Button
            variant="outline"
            size="icon-sm"
            aria-label="Previous page"
            :disabled="currentPage <= 1"
            @click="goToPage(currentPage - 1)"
          >
            <ChevronLeft data-icon="inline-start" />
          </Button>
          <Button
            v-for="page in pageCount <= 7 ? pageCount : []"
            :key="page"
            variant="outline"
            size="icon-sm"
            :aria-label="`Page ${page}`"
            :aria-current="isPageActive(page) ? 'page' : undefined"
            :class="{ 'bg-accent': isPageActive(page) }"
            @click="goToPage(page)"
          >
            {{ page }}
          </Button>
          <Button
            variant="outline"
            size="icon-sm"
            aria-label="Next page"
            :disabled="currentPage >= pageCount"
            @click="goToPage(currentPage + 1)"
          >
            <ChevronRight data-icon="inline-end" />
          </Button>
        </div>
      </div>
    </div>
  </div>
</template>
