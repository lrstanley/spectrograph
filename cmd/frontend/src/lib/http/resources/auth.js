import { Promise } from 'q'
import http from '~/lib/http/http'

const self = {
    get: (options) => { return http.get('/auth/self', options) }
}

const callback = {
    get: (code, state, options) => {
        return http.get('/auth/callback', { ...options, params: {
            code: code,
            state: state,
        }})
    }
}

const redirect = {
    get: (options) => { return http.get('/auth/redirect', options) }
}

const logout = {
    get: (options) => { return http.get('/auth/logout', options) }
}

export default {
    self,
    callback,
    redirect,
    logout,
}
