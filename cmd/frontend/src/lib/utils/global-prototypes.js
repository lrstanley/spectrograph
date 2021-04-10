import Vue from 'vue'
import formatDistance from 'date-fns/formatDistance'

Vue.prototype.humanizeISODate = function (isoDate) {
    return formatDistance(Date.parse(isoDate), new Date()).replace('about ', '') + ' ago'
}
