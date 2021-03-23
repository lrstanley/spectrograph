import Vue from 'vue'
import Axios from 'axios'
import VueAxios from 'vue-axios'
import VueProgressBar from 'vue-progressbar'
import SnackbarStackPlugin from 'snackbarstack'
import Vuetify from 'vuetify'
import 'vuetify/dist/vuetify.min.css'

import App from './App.vue'
import router from './router'
import store from './store'

Vue.prototype.$http = Axios;
Vue.config.productionTip = false

Vue.use(VueProgressBar, {
    color: 'rgb(120, 96, 255)',
    failedColor: 'red',
    thickness: '4px'
  })
Vue.use(Vuetify)
Vue.use(SnackbarStackPlugin)
Vue.use(VueAxios, Axios)

new Vue({
    router,
    store,
    el: '#vue',
    render: h => h(App)
})
