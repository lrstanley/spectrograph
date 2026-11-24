import type { FetchOptions, ResponseType } from "ofetch"
import type { OpenFetchClient } from "#imports"
import type { paths as ApiPaths } from "#open-fetch-schemas/api"

export default defineNuxtPlugin({
  enforce: "pre",
  setup(nuxtApp): { provide: { api: OpenFetchClient<ApiPaths> } } {
    const options = useRuntimeConfig().public.openFetch.api as FetchOptions<ResponseType>
    const baseURL = String(useRuntimeConfig().public.API_URL).replace(/\/$/, "")

    return {
      provide: {
        api: createOpenFetch<ApiPaths>(
          (localOptions) => ({
            ...options,
            ...localOptions,
            baseURL,
            credentials: "include",
          }),
          useRequestFetch() as typeof globalThis.$fetch,
          "api",
          // Nuxt 4.6 types hooks as NuxtHookRegistry, which omits hookable's private fields.
          nuxtApp.hooks as unknown as Parameters<typeof createOpenFetch>[3],
        ),
      },
    }
  },
})
