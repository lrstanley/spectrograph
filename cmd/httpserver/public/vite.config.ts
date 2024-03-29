import anchor from "markdown-it-anchor"
import path from "path"
import { visualizer } from "rollup-plugin-visualizer"
import AutoImport from "unplugin-auto-import/vite"
import IconsResolver from "unplugin-icons/resolver"
import Icons from "unplugin-icons/vite"
import { HeadlessUiResolver } from "unplugin-vue-components/resolvers"
import Components from "unplugin-vue-components/vite"
import { VueRouterAutoImports } from "unplugin-vue-router"
import VueRouter from "unplugin-vue-router/vite"
import { defineConfig } from "vite"
import { imagetools } from "vite-imagetools"
import codegen from "vite-plugin-graphql-codegen"
import Markdown from "vite-plugin-md"
import Layouts from "vite-plugin-vue-layouts"
import Vue from "@vitejs/plugin-vue"
import link from "@yankeeinlondon/link-builder"

const icons = IconsResolver({
  componentPrefix: "i",
  enabledCollections: ["fa6-solid", "fa6-regular", "fa6-brands", "twemoji"],
  alias: {
    fas: "fa6-solid",
    far: "fa6-regular",
    fab: "fa6-brands",
    emoji: "twemoji",
  },
})

export default defineConfig({
  resolve: {
    alias: {
      "@/": `${path.resolve(__dirname, "src")}/`,
    },
  },
  plugins: [
    visualizer({
      filename: "./dist/stats.html",
    }),
    codegen({
      enableWatcher: true,
      config: {
        errorsOnly: true,
        schema: "./../../../internal/database/graphql/schema/*.gql",
        documents: "./src/lib/api/*.gql",
        generates: {
          "./src/lib/api/graphql.ts": {
            plugins: ["typescript", "typescript-operations", "typescript-vue-urql"],
            config: {
              preResolveTypes: true,
              nonOptionalTypename: true,
              skipTypeNameForRoot: true,
              useTypeImports: true,
              inputMaybeValue: "T | Ref<T> | ComputedRef<T>",
            },
            // hooks: {
            //   afterOneFileWrite: ["pnpm exec prettier --write"],
            // },
          },
        },
      },
    }),
    VueRouter({
      routesFolder: "src/pages",
      routeBlockLang: "yaml",
      extensions: ["vue", "md"],
    }),
    Vue({
      include: [/\.vue$/, /\.md$/],
    }),
    Markdown({
      markdownItOptions: {
        html: true,
        linkify: true,
        typographer: false,
      },
      markdownItSetup(md) {
        md.use(anchor)
      },
      builders: [link()],
      wrapperClasses: "prose max-w-none",
    }),
    Layouts({
      layoutsDirs: "src/layouts",
      defaultLayout: "default",
    }),
    Components({
      extensions: ["vue", "md"],
      directives: true,
      directoryAsNamespace: true,
      importPathTransform: (path) => path.replace(/.*\/src\//, "@/"),
      resolvers: [HeadlessUiResolver(), icons],
    }),
    AutoImport({
      dts: true,
      imports: [
        "vue",
        "@vueuse/core",
        {
          "@/lib/core/state": ["useState"],
        },
        VueRouterAutoImports,
      ],
      resolvers: [icons],
      eslintrc: {
        enabled: true,
      },
    }),
    Icons({
      autoInstall: true,
      defaultClass: "icon",
    }),
    imagetools(),
  ],
  base: "/",
  build: {
    sourcemap: "hidden",
    emptyOutDir: true,
  },
  preview: {
    port: 8081,
    open: false,
  },
  server: {
    base: "/",
    port: 8081,
    open: false,
    strictPort: true,
    proxy: {
      "^/(-|security\\.txt|robots\\.txt)(/.*|$)": {
        target: "http://localhost:8080",
        xfwd: true,
        ws: true,
      },
    },
  },
})
