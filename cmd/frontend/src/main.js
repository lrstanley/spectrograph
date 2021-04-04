import Vue from 'vue'
import VueProgressBar from 'vue-progressbar'

Vue.use(VueProgressBar, {
    color: colors.nitro,
    failedColor: 'red',
    thickness: '3px',
    location: 'top',
    transition: {
        speed: '0.3s',
        opacity: '0.6s',
        termination: 300
    },
    autoRevert: true
})

import Vuetify from 'vuetify'
import 'vuetify/dist/vuetify.min.css'
import '@mdi/font/css/materialdesignicons.css'
Vue.use(Vuetify)

import app from '~/app.vue'
import router from '~/lib/core/router'
import state from '~/lib/core/state'
import '~/lib/utils/global-prototypes'
import colors from '~/lib/utils/discord-colors'

Vue.config.productionTip = false

export default new Vue({
    vuetify: new Vuetify({
        icons: {
            iconfont: 'mdi'
        },
        theme: {
            dark: true,
            default: 'dark',
            themes: {
                dark: {
                    ...colors,
                    // primary: '#2196F3',
                    // secondary: '#424242',
                    // accent: '#FF4081',
                    // error: '#FF5252',
                    // info: '#2196F3',
                    // success: '#4CAF50',
                    // warning: '#FB8C00',

                    primary: colors.discord,
                    secondary: colors.chat,
                    accent: colors.nitro,
                    error: colors.dnd,
                    info: colors.bravery,
                    success: colors.online,
                    warning: colors.high,
                    anchor: '#8c9eff',
                }
            }
        }
    }),
    router,
    store: state,
    el: '#app',
    render: h => h(app)
})
