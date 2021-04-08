import { defineConfig } from 'vite';
import { createVuePlugin } from 'vite-plugin-vue2';
import ViteComponents, { VuetifyResolver } from 'vite-plugin-components';

export default defineConfig({
    plugins: [
        createVuePlugin(),
        ViteComponents({
            customComponentResolvers: [VuetifyResolver()]
        })
    ],
    publicDir: "src/static",
    resolve: {
        alias: {
            '~': __dirname + '/src/',
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
            '^/api/.*': {
                target: 'http://http-server:8080',
                xfwd: true,
            }
        }
    },
    sourcemap: true
});
