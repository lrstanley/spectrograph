import { Promise } from 'q'
import http from '~/lib/http/http'

function self(options) { return http.get('/auth/self', options) }

function callback(code, state, options) {
    return http.get('/auth/callback', { ...options, params: {
        code: code,
        state: state,
    }})
}

function redirect(options) { return http.get('/auth/redirect', options) }

function logout(options) { return http.get('/auth/logout', options) }

export default {
    self,
    callback,
    redirect,
    logout,
}
