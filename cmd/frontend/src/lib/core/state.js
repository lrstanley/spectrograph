import Vue from 'vue'
import Vuex from 'vuex'

import config from '~/lib/core/config'
import api from '~/lib/http/api'

// import createPersistedState from 'vuex-persistedstate'
// import createLogger from 'vuex/dist/logger

Vue.use(Vuex)

export default new Vuex.Store({
    strict: config.debug ? true : false,
    state: {
        loading: true,
        title: "",
        auth: null,
    },
    mutations: {
        disableLoading: (state) => { state.loading = false },
        enableLoading: (state) => { state.loading = true }, // TODO: handle this better..
        setTitle: (state, val) => {
            state.title = val ? val.replace(" - ", " · ") : ""
            document.title = state.title != "" ? `${state.title} · ${config.application}` : config.application
        },
        set_auth(state, auth) { state.auth = auth }
    },
    actions: {
        get_auth: function ({ commit, state }, routeGuard) {
            return new Promise((resolve, reject) => {
                if (routeGuard && state.auth !== null) {
                    resolve()
                    return
                }

                api.http.get("/auth/self").then(resp => {
                    commit('set_auth', resp.data)
                    resolve()
                }).catch(err => {
                    commit('set_auth', false)
                    reject(err)
                })
            })
        },
        logout: function ({ commit }) {
            new Promise((resolve, reject) => {
                api.http.get("/auth/logout").then(resp => {
                    commit('set_auth', false)
                    resolve(resp)
                }).catch(err => {
                    reject(err)
                })
            })
        }
    },
    getters: {
        authed: (state) => { return !!state.auth.authenticated },
        admin: (state) => { return !!state.auth.admin },
        loading: (state) => { return state.loading }
    }
})
