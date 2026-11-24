import { useRouteQuery } from "@vueuse/router"

export type PaginationOptions<T extends string = string> = {
  page?: number
  perPage?: number
  sort?: T
  order?: "asc" | "desc"
  resetChanged?: Ref<unknown>[]
}

export function usePagination<T extends string = string>({
  page = 1,
  perPage = 15,
  sort,
  order = "desc",
  resetChanged = [],
}: PaginationOptions<T>) {
  const results = {
    page: useRouteQuery<number>("page", page, { transform: Number }),
    perPage: useRouteQuery<number>("per_page", perPage, { transform: Number }),
    sort: useRouteQuery<T>("sort", sort),
    order: useRouteQuery<"asc" | "desc">("order", order),
    sizes: [5, 15, 25, 50, 100, 250],
  }

  if (resetChanged.length) {
    watch([...resetChanged, results.perPage, results.sort, results.order], () => {
      results.page.value = 1
    })
  }

  return results
}
