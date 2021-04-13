import Vue from 'vue'
import Vuex from 'vuex'

import config from '@/lib/core/config'
import api from '@/lib/http/api'

// import createPersistedState from 'vuex-persistedstate'
// import createLogger from 'vuex/dist/logger

Vue.use(Vuex)

export default new Vuex.Store({
    strict: config.debug ? true : false,
    state: {
        loading: true,
        title: "",
        auth: null,
        auth_last_checked: null,
    },
    mutations: {
        disableLoading: (state) => { state.loading = false },
        enableLoading: (state) => { state.loading = true }, // TODO: handle this better..
        setTitle: (state, val) => {
            state.title = val ? val.replace(" - ", " · ") : ""
            document.title = state.title != "" ? `${state.title} · ${config.application}` : config.application
        },
        set_auth(state, auth) {
            state.auth = auth
            state.auth_last_checked = new Date()
        }
    },
    actions: {
        get_auth: async function ({ commit, state }, routeGuard) {
            if (routeGuard && state.auth !== null) {
                // determine if auth-check is older than 5 minutes, and
                // we should re-check for auth just in case.
                if ((new Date() - state.auth_last_checked) < 5 * 60 * 1000) {
                    return
                }
            }

            try {
                let resp = await api.auth.self()
                commit('set_auth', resp.data)
                return resp.data
            } catch (err) {
                commit('set_auth', false)
                throw err
            }
        },
        logout: async function ({ commit }) {
            let resp = api.auth.logout()
            commit('set_auth', false)
            return resp.data
        }
    },
    getters: {
        user: (state) => { return state.auth?.user },
        authed: (state) => { return state.auth?.authenticated },
        admin: (state) => { return state.auth?.admin },
        loading: (state) => { return state.loading }
    }
})
