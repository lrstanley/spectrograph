import axios from 'axios'
import config from '~/lib/core/config'
import utils from '~/lib/http/utils'

import app from '~/main'

var http = axios.create({
    baseURL: config.api_baseurl,
    headers: {
        common: {
            "Content-Type": "application/json"
        },
    },
    responseType: "json",
    maxRedirects: 0,
})

http.interceptors.request.use((config) => {
    app.$Progress.start()
    return config
})

http.interceptors.response.use((response) => {
    if (response.config.sort === true && Array.isArray(response.data)) {
        response.data = response.data.sort((a, b) => { return a.name > b.name ? 1 : -1 })
    }

    if (response.data.error) {
        let error = new Error()
        error.response = response

        app.$Progress.fail()
        return Promise.reject(utils.formatError(error))
    }

    app.$Progress.finish()
    return Promise.resolve(response)
}, (error) => {
    app.$Progress.fail()
    return Promise.reject(utils.formatError(error))
})

export default http
