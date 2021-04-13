import { defineConfig } from 'vite'
import { createVuePlugin } from 'vite-plugin-vue2'
import ViteComponents, { VuetifyResolver } from 'vite-plugin-components'
import path from 'path'

export default defineConfig({
    plugins: [
        createVuePlugin(),
        ViteComponents({
            // https://github.com/antfu/vite-plugin-components#configuration
            dirs: ['./src/components'],
            extensions: ['vue'],
            deep: true,
            customComponentResolvers: [VuetifyResolver()]
        })
    ],
    publicDir: "src/static",
    resolve: {
        alias: {
            // see also: jsconfig.json
            '@': path.resolve(__dirname, './src')
        },
    },
    build: {
        target: 'es2015',
        sourcemap: true
    },
    server: {
        port: 8081,
        strictPort: true,
        proxy: {
            // '^/dist/.*': {
            //     target: 'http://127.0.0.1:8081',
            //     toProxy: true,
            //     xfwd: true,
            //     rewrite: (path) => path.replace(/\/dist/, '')
            // },
            '^/api/.*': {
                target: 'http://http-server:8080',
                xfwd: true,
            }
        },
        force: true
    },
    sourcemap: true
});
