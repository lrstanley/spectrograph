import Vue from 'vue'
import axios from 'axios'
import formatDistance from 'date-fns/formatDistance'

import config from '~/lib/core/config'
import router from '~/lib/core/router'
import api from '~/lib/http/api'

Vue.prototype.$api = api
Vue.prototype.$http = axios
Vue.prototype.$config = config
Vue.prototype.debug = config.debugLog

window.api = api

Vue.prototype.updateQuery = function (update) {
    let query = Object.assign({}, router.app.$route.query, update)
    for (let key in query) {
        // if a GET param is false, just remote it.
        if (!query[key]) { delete query[key] }
    }
    router.app.$router.push({ query: query }).catch(() => { })
}

Vue.prototype.resetComponentData = function () {
    Object.assign(this.$data, this.$options.data.apply(this))
}

Vue.prototype.humanizeISODate = function (isoDate) {
    return formatDistance(Date.parse(isoDate), new Date()).replace('about ', '') + ' ago'
}
