import path from "path"
import AutoImport from "unplugin-auto-import/vite"
import IconsResolver from "unplugin-icons/resolver"
import Icons from "unplugin-icons/vite"
import { HeadlessUiResolver } from "unplugin-vue-components/resolvers"
import Components from "unplugin-vue-components/vite"
import { defineConfig } from "vite"
import Pages from "vite-plugin-pages"
import Vue from "@vitejs/plugin-vue"

export default defineConfig({
  resolve: {
    alias: {
      "@/": `${path.resolve(__dirname, "src")}/`,
    },
  },
  publicDir: `${path.resolve(__dirname, "src")}/assets`,
  plugins: [
    Pages({
      dirs: "src/pages",
      routeBlockLang: "yaml",
      extendRoute(route) {
        // route, parent
        if (route.path.startsWith("/admin")) {
          route = {
            ...route,
            meta: {
              auth: true,
            },
          }
        }

        return route
      },
    }),
    Vue({}),
    Components({
      dts: true,
      directives: true,
      directoryAsNamespace: true,
      resolvers: [
        HeadlessUiResolver(),
        IconsResolver({
          componentPrefix: "i",
          enabledCollections: ["fa6-solid", "fa6-regular", "fa6-brands", "twemoji"],
          alias: {
            fas: "fa6-solid",
            far: "fa6-regular",
            fab: "fa6-brands",
            emoji: "twemoji",
          },
        }),
      ],
    }),
    AutoImport({
      dts: true,
      imports: [
        "vue",
        "vue-router",
        "@vueuse/core",
        {
          "@/lib/core/state": ["useState"],
        },
      ],
      resolvers: [
        IconsResolver({
          componentPrefix: "icon",
          enabledCollections: ["fa6-solid", "fa6-regular", "fa6-brands", "twemoji"],
          alias: {
            fas: "fa6-solid",
            far: "fa6-regular",
            fab: "fa6-brands",
            emoji: "twemoji",
          },
        }),
      ],
      eslintrc: {
        enabled: true,
      },
    }),
    Icons({
      autoInstall: true,
      defaultClass: "icon",
    }),
  ],
  base: "/",
  build: {
    sourcemap: false,
    emptyOutDir: true,
  },
  preview: {
    port: 8081,
  },
  server: {
    base: "/",
    port: 8081,
    strictPort: true,
    proxy: {
      "^/(-|security\\.txt|robots\\.txt)(/.*|$)": {
        target: "http://localhost:8080",
        xfwd: true,
      },
    },
  },
})
