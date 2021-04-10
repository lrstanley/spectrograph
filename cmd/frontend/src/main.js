import Vue from 'vue'

import app from '~/app.vue'
import router from '~/lib/core/router'
import state from '~/lib/core/state'
import vuetify from '~/lib/core/vuetify'

import '~/lib/core/progressbar'
import '~/lib/http/api'
import '~/lib/utils/global-prototypes'

Vue.config.productionTip = false

export default new Vue({
    vuetify: vuetify,
    router: router,
    store: state,
    el: '#app',
    render: h => h(app)
})
