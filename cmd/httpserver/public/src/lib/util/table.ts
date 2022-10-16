/**
 * Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
 * of this source code is governed by the MIT license that can be found in
 * the LICENSE file.
 */

import merge from "ts-deepmerge"
import { OrderDirection } from "@/lib/api"
import { resetCursor, usePagination } from "@/lib/util/pagination"
import { refDebounced } from "@vueuse/core"

import type { Filter } from "@/lib/util/pagination"
import type { Ref, ComputedRef, Component } from "vue"

export interface IDType {
  id: string | number
}

export interface ColumnDef<T extends IDType, OrderField, WhereInput extends Record<string, any>> {
  id: string
  header: string
  sortField?: OrderField
  type?: "text" | "number" | "boolean"
  field?: keyof T
  accessorFn?: (data: T) => any
  filterFn?: (value: any) => WhereInput
  renderFn?: (value: any, data: T) => Component
  clickFn?: (data: T, event: MouseEvent) => void
}

export interface ColumnState<T extends IDType, OrderField, WhereInput extends Record<string, any>>
  extends ColumnDef<T, OrderField, WhereInput> {
  filterValue: Ref<any>
  filterValueDebounced: Readonly<Ref<any>>

  canClick(): boolean
  canFilter(): boolean
  canSort(): boolean
  isSorted(): OrderDirection | boolean
  toggleSort(): void
}

export interface Row<T extends IDType, OrderField, WhereInput extends Record<string, any>> {
  id: string | number
  data: T
  columns(): ColumnState<T, OrderField, WhereInput>[]
  cells: ComputedRef<Cell<T, OrderField, WhereInput>[]>
}

export interface Cell<T extends IDType, OrderField, WhereInput extends Record<string, any>> {
  column: ColumnState<T, OrderField, WhereInput>
  data: T
  value?: any
  render(): Component
}

export interface SortingState<OrderField> {
  field: OrderField
  direction: OrderDirection
}

export interface QueryState<OrderField, WhereInput extends Record<string, any>> extends Filter {
  where?: WhereInput
  orderBy?: OrderField
  order?: OrderDirection
}

export interface Table<T extends IDType, OrderField, WhereInput extends Record<string, any>> {
  columnDefs: ColumnDef<T, OrderField, WhereInput>[]
  columns: ColumnState<T, OrderField, WhereInput>[]

  sortState: Ref<SortingState<OrderField>>
  filterState: Ref<WhereInput>
  queryState: Ref<Record<string, any>>

  data: Ref<T[]>
  rows: ComputedRef<Row<T, OrderField, WhereInput>[]>
  cursor: Ref<string>

  canSort(): boolean
  isSorted(): boolean

  canFilter(): boolean
  isFiltered(): boolean
}

export class CoreTable<T extends IDType, OrderField, WhereInput extends Record<string, any>>
  implements Table<T, OrderField, WhereInput>
{
  columnDefs: ColumnDef<T, OrderField, WhereInput>[]
  columns: ColumnState<T, OrderField, WhereInput>[]
  sortState: Ref<SortingState<OrderField>>
  filterState: ComputedRef<WhereInput>
  queryState: ComputedRef<QueryState<OrderField, WhereInput>>
  data: Ref<T[]>
  cursor: Ref<string>

  rows: ComputedRef<Row<T, OrderField, WhereInput>[]>
  _pageSize: number
  _pagination: Filter

  constructor(columnDefs: ColumnDef<T, OrderField, WhereInput>[], pageSize: number) {
    this.columnDefs = columnDefs

    this.columns = columnDefs.map((col) => {
      const newRef = ref<any>(null)
      return {
        ...col,
        filterValue: newRef,
        filterValueDebounced: refDebounced(newRef, 250),
        canClick: () => !!col.clickFn,
        canFilter: () => !!col.filterFn,
        canSort: () => !!col.sortField,
        isSorted: () =>
          this.sortState.value.field === col.sortField ? this.sortState.value.direction : false,
        toggleSort: () => this._toggleSortColumn(col),
      }
    })

    this.sortState = ref({
      field: null,
      direction: null,
    }) as Ref<SortingState<OrderField>>

    this.filterState = computed(() => {
      let filters = {} as WhereInput

      this.columns.forEach((col) => {
        if (!col.canFilter()) return

        const value = col.filterFn(col.filterValueDebounced.value)
        if (value === null) return

        filters = merge(filters, value) as WhereInput
      })

      return filters
    })

    this.data = ref([]) as Ref<T[]>
    this.cursor = ref("") as Ref<string>

    this.rows = computed(() => {
      return this.data.value.map((data) => {
        return {
          id: data.id,
          data,
          columns: () => this.columns,
          cells: computed(() => this.columns.map((col) => this._generateCell(data, col))),
        }
      })
    })

    this._pageSize = pageSize
    this._pagination = usePagination(this.cursor, this._pageSize)

    resetCursor(this.cursor, [this.sortState, this.filterState])

    this.queryState = computed(() => {
      const q = {} as QueryState<OrderField, WhereInput>

      if (this.isSorted()) {
        q.orderBy = this.sortState.value.field
        q.order = this.sortState.value.direction
      }

      if (this.isFiltered()) {
        q.where = this.filterState.value
      }

      return {
        ...q,
        ...this._pagination,
      }
    })
  }

  _toggleSortColumn(col: ColumnDef<T, OrderField, WhereInput>) {
    if (!col.sortField) return

    if (this.sortState.value.field !== col.sortField) {
      this.sortState.value.field = col.sortField
      this.sortState.value.direction = OrderDirection.Asc
      return
    }

    if (this.sortState.value.direction === OrderDirection.Asc) {
      this.sortState.value.direction = OrderDirection.Desc
      return
    }

    this.sortState.value.field = null
    this.sortState.value.direction = null
  }

  _generateCell(data: T, cell: ColumnState<T, OrderField, WhereInput>): Cell<T, OrderField, WhereInput> {
    let value: any

    if (cell.field) {
      value = data[cell.field]
    } else if (cell.accessorFn) {
      value = cell.accessorFn(data)
    } else {
      value = null
    }

    return {
      column: cell,
      data,
      value,
      render: () => (cell.renderFn ? cell.renderFn(value, data) : h("span", value)),
    }
  }

  canSort() {
    return this.columns.some((col) => col.canSort())
  }

  isSorted() {
    return this.sortState.value.field !== null
  }

  canFilter() {
    return this.columns.some((col) => col.canFilter())
  }

  isFiltered() {
    return this.filterState.value ? Object.keys(this.filterState.value).length > 0 : false
  }
}
