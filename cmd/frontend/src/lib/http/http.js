import axios from 'axios'
import config from '@/lib/core/config'
import utils from '@/lib/http/utils'

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

http.interceptors.response.use((response) => {
    if (response.config.sort === true && Array.isArray(response.data)) {
        response.data = response.data.sort((a, b) => { return a.name > b.name ? 1 : -1 })
    }

    if (response.data.error) {
        let error = new Error()
        error.response = response

        return Promise.reject(utils.formatError(error))
    }

    return Promise.resolve(response)
}, (error) => {
    return Promise.reject(utils.formatError(error))
})

export default http
