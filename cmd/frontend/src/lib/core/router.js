import Vue from 'vue'
import VueRouter from 'vue-router'

import state from '~/lib/core/state'

import LayoutDefault from '~/components/core/layout-default.vue'
import LayoutUser from '~/components/core/layout-user.vue'
import Index from '~/views/index.vue'
import NotFound from '~/views/not-found.vue'
import Auth from '~/views/auth.vue'

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
                // meta: { auth: true }
            },
            {
                path: '/user',
                redirect: '/user/details',
                component: LayoutUser,
                meta: { auth: true },
                children: [
                    {
                        path: '/user/details',
                        name: Index.name,
                        component: Index
                    }
                ]
            }
        ]
    },
    { path: '/auth/:method', name: Auth.name, component: Auth, props: true },
    { path: '/404', name: NotFound.name, component: NotFound },
    { path: '*', name: 'catchall', redirect: '/404' }
]

const router = new VueRouter({
    routes,
    mode: 'history',
    scrollBehavior(to, from, savedPosition) {
        if (savedPosition) { return savedPosition }
        return { x: 0, y: 0 }
    }
})

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

    // see if any routes we want to go to require admin or authentication.
    let wantsAdmin = false
    let wantsAuth = false
    for (let route of to.matched) {
        if (route.meta.admin) { wantsAdmin = true }
        if (route.meta.auth) { wantsAuth = true }
    }

    // target route has meta fields.
    if (to.meta && Object.keys(to.meta).length > 0) {
        if (wantsAdmin && isAuthed && !state.getters.admin) {
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
    } else if (wantsAuth) {
        // target route requires authentication, but we're not authed.
        router.app.$Progress.fail()
        return _next({
            name: Auth.name,
            params: {
                method: "redirect",
                // after we authenticate, go back to the page they were trying
                // to go to.
                next: { path: to.path, query: to.query }
            }
        })
    } else {
        router.app.$Progress.finish()
        // if logging in, tell auth to redirect back to the page they were on.
        if (to.name == Auth.name && to.params.method == "redirect" && !to.params.next) {
            to.params.next = { path: from.path, query: from.query }
        }
        return _next()
    }
}

Vue.prototype.updateQuery = function (update) {
    // TODO: router.app.$route
    // router.app.$router
    let query = Object.assign({}, this.$route.query, update)
    for (let key in query) {
        // if a GET param is false, just remote it.
        if (!query[key]) { delete query[key] }
    }
    this.$router.push({ query: query }).catch(() => { })
}

Vue.prototype.anchor = function (refName) {
    const element = this.$refs[refName]
    if (!element) { return }
    element.scrollIntoView()
}

Vue.prototype.resetComponentData = function () {
    Object.assign(this.$data, this.$options.data.apply(this))
}

export default router
