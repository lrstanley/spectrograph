import { execSync } from "node:child_process"
import tailwindcss from "@tailwindcss/vite"
import { markdown } from "./build/markdown"

// Version shown in the UI and prefilled into bug reports. Falls back to the
// current commit when not provided (e.g. local builds).
function appVersion() {
  if (process.env.NUXT_PUBLIC_VERSION) {
    return process.env.NUXT_PUBLIC_VERSION
  }
  try {
    return execSync("git describe --tags --always --dirty", { stdio: ["ignore", "pipe", "ignore"] })
      .toString()
      .trim()
  } catch {
    return "unknown"
  }
}

export default defineNuxtConfig({
  app: {
    head: {
      htmlAttrs: { lang: "en" },
      script: [
        {
          // Apply the saved theme (dark by default) before the app mounts, to
          // avoid a flash. useColorMode() takes over from here.
          innerHTML:
            "try{var m=localStorage.getItem('spectrograph-color-mode');document.documentElement.classList.add(m==='light'?'light':'dark')}catch(e){document.documentElement.classList.add('dark')}",
        },
      ],
      title: "Spectrograph",
      meta: [
        {
          name: "description",
          content:
            "Spectrograph is a Discord bot that keeps one empty voice channel open in each group, cloning and removing channels as people join and leave.",
        },
        { name: "theme-color", content: "#16131f" },
      ],
      link: [{ rel: "icon", type: "image/png", href: "/favicon.png" }],
    },
  },
  modules: ["shadcn-nuxt", "nuxt-open-fetch", "@vueuse/nuxt"],
  compatibilityDate: "2026-09-05",
  experimental: {
    asyncContext: true,
    typedPages: true,
    asyncEntry: true,
  },
  devtools: {
    enabled: false,
    viteInspect: false,
    componentInspector: false,
  },
  ssr: false,
  runtimeConfig: {
    public: {
      API_URL: "/-",
      version: appVersion(),
    },
  },
  css: ["~/assets/main.css"],
  nitro: {
    // Served as an SPA by the Go binary (which falls back to index.html), so
    // only the entrypoint is generated, not a directory per route. Nuxt already
    // renders index.html (as "/index.html") for client-only apps, and also
    // rendering "/" writes the same file concurrently, which can corrupt it.
    prerender: { crawlLinks: false, ignore: [/^\/$/] },
  },
  devServer: { port: 8081 },
  telemetry: false,
  vite: {
    plugins: [tailwindcss(), markdown()],
    optimizeDeps: {
      include: ["@lucide/vue", "clsx", "tailwind-merge", "vue-sonner", "@tanstack/vue-table", "@vueuse/router"],
    },
    server: {
      open: false,
      strictPort: true,
    },
  },
  shadcn: {
    prefix: "",
    componentDir: "@/components/ui",
  },
  openFetch: {
    disableNitroPlugin: true,
    disableNuxtPlugin: true,
    openAPITS: {
      rootTypes: true,
      // entrest adds defaults on *Update component schemas; openapi-typescript
      // treats those as required when defaultNonNullable is true (default), even
      // for PATCH bodies referenced via $ref.
      defaultNonNullable: false,
    },
    clients: {
      api: {
        schema: "../internal/database/ent/rest/openapi.json",
      },
    },
  },
})
