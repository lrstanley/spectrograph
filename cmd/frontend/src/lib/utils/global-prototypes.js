import Vue from 'vue'
import formatDistance from 'date-fns/formatDistance'

import router from '~/lib/core/router'

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
