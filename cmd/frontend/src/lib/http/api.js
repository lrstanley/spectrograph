import http from '~/lib/http/http'
import auth from '~/lib/http/resources/auth'

function install(Vue, options) {
    http.interceptors.request.use((config) => {
        Vue.prototype.$Progress.start()
        return config
    })

    http.interceptors.response.use((response) => {
        if (response.data.error) {
            Vue.prototype.$Progress.fail()
        }

        Vue.prototype.$Progress.finish()
        return Promise.resolve(response)
    }, (error) => {
        Vue.prototype.$Progress.fail()
        return Promise.reject(error)
    })

    Vue.prototype.$api = api
    Vue.prototype.$http = http
    window.api = api
}

const api = {
    install: install,
    http: http,
    auth: auth,
}

export default api
