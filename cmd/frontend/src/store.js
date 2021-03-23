import Vue from 'vue'
import Vuex from 'vuex'
import axios from 'axios'

Vue.use(Vuex)

export default new Vuex.Store({
    state: {
        auth: null,
        check_auth: true,
    },
    mutations: {
        set_auth(state, auth) { state.auth = auth; },
    },
    actions: {
        get_auth({ commit, state }) {
            return new Promise((resolve, reject) => {

                axios("/api/v1/auth/self").then(resp => {
                    const user = resp.data.authenticated ? resp.data.user : false
                    commit('set_auth', user)
                    resolve(user)
                }).catch(err => {
                    commit('set_auth', false)
                    // TODO: if error field, send notification?
                    reject(err)
                })
            })
        },
        logout({ commit }) {
            new Promise((resolve, reject) => {
                axios("/api/v1/auth/logout").then(resp => {
                    commit('set_auth', false)
                    resolve(resp)
                }).catch(err => {
                    // TODO: if error field, send notification?
                    reject(err)
                })
            })
        }
    },
    getters: {
        // isAuthed: state => !!state.auth,
        // auth: state => state.auth,
    }
})
