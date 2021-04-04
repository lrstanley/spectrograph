import Vue from 'vue'
import VueRouter from 'vue-router'

import state from '~/lib/core/state'

import LayoutDefault from '~/components/core/layout-default'
import LayoutUser from '~/components/core/layout-user'
import Index from '~/pages/index'
import NotFound from '~/pages/not-found'
import Auth from '~/pages/auth'

Vue.use(VueRouter)

const routes = [
    {
        path: '/',
        component: LayoutDefault,
        children: [
            {
                path: '/',
                name: Index.name,
                component: Index,
                meta: { auth: true }
            },
            // {
            //     path: '/admin/settings',
            //     name: 'admin-settings',
            //     component: AdminSettings,
            // }
        ]
    },
    { path: '/auth/:method', name: Auth.name, component: Auth, props: true },
    { path: '/404', name: NotFound.name, component: NotFound },
    { path: '*', name: 'catchall', redirect: '/404' }
]

const router = new VueRouter({ routes, mode: 'history' })

router.beforeEach((to, from, next) => {
    router.app.$Progress.start()

    // ensure we fetch auth information before we load any pages.
    state.dispatch('get_auth', true).then(() => {
    }).catch((e) => {
        console.log("unable to fetch user details:", e) // TODO: error page?
    }).finally(() => {
        beforeEach(to, from, next)
    })
})

function beforeEach(to, from, next) {
    const isAuthed = state.getters.authed

    let _next = function (...vars) {
        if (!to.meta.noDisableLoading) {
            state.commit('disableLoading')
        }
        next(...vars)
    }

    // target route has meta fields.
    if (to.meta && Object.keys(to.meta).length > 0) {
        if (to.meta.admin && isAuthed && !state.getters.admin) {
            // page requires admin, but user is not admin.
            // TODO: error page of some kind.
            console.log("[route:perm] denied, not admin")
            router.app.$Progress.fail()
            return _next({ name: Index.name })
        }

        // handle page titles, go through each component, from most specific
        // to most generic in the routing table, until we find one that returns
        // a title.
        let title = null
        if (to.matched.length > 0) {
            for (let i = to.matched.length - 1; i >= 0; i--) {
                title = to.matched[i].components.default.title
                if (typeof title === 'function') {
                    title = title(to.params)
                }
                if (title) { break }
            }
        }
        state.commit('setTitle', title)
    }

    // check if the user is authed, and redirect based on state.
    if (isAuthed) {
        // If going to auth routes, and we're already logged in... redirect
        // to index (or cancel if already @ index).
        if (to.name == Auth.name && to.params.method != "logout") {
            router.app.$Progress.fail()
            return from.name != Index.name ? _next({ name: Index.name }) : _next(false)
        }
        router.app.$Progress.finish()
        return _next()
    } else if (to.meta.authed) {
        // target route requires authentication, but we're not authed.
        router.app.$Progress.fail()
        return _next({
            name: Auth.name,
            params: {
                method: "redirect",
                next: { path: to.path, query: to.query }
            }
        })
    } else {
        router.app.$Progress.finish()
        return _next()
    }
}

export default router
