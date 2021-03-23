import Vue from 'vue'
import VueRouter from 'vue-router'
import store from './store.js'

import Index from './pages/Index'
import AdminSettings from './pages/AdminSettings'
import NotFound from './pages/NotFound'
import Auth from './pages/Auth'

Vue.use(VueRouter)

const routes = [
    {
        name: 'index',
        path: '/',
        component: Index,
        meta: { title: 'Index' },
    },
    {
        name: 'auth-github-callback',
        path: '/ui/auth/github/callback',
        component: Auth,
    },
    {
        name: 'admin-settings',
        path: '/admin/settings',
        component: AdminSettings,
    },
    {
        name: 'catchall',
        path: '*',
        redirect: '/404'
    },
    {
        name: 'not-found',
        path: '/404',
        component: NotFound,
        meta: { title: 'Page not found' },
    }
]

const router = new VueRouter({ routes, mode: 'history' })
router.beforeEach((to, from, next) => {
    if (to.meta.title !== undefined) {
        document.title = `${to.meta.title} · spectrograph`
    } else {
        document.title = "spectrograph"
    }

    if (from.name == null && !to.name.startsWith("auth-")) {
        store.dispatch('get_auth').catch((err) => {
            console.log(err)
            if (to.name.startsWith('admin-') || !!to.meta.auth_required) {
                window.location.replace("/api/v1/auth/github/redirect")
            }
        }).finally(() => {
            next()
        })
    } else {
        if (!store.state.auth) {
            // Assume we already fetched auth earlier.
            if (to.name.startsWith('admin-') || !!to.meta.auth_required) {
                window.location.replace("/api/v1/auth/github/redirect")
            }
        }
        next()
    }
})

export default router
