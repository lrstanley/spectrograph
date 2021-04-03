import Vue from 'vue'
import VueProgressBar from 'vue-progressbar'

Vue.use(VueProgressBar, {
    color: 'rgb(120, 96, 255)',
    failedColor: 'red',
    thickness: '4px'
})

import Vuetify from 'vuetify'
import 'vuetify/dist/vuetify.min.css'
import '@mdi/font/css/materialdesignicons.css'
Vue.use(Vuetify)

import app from '~/app.vue'
import router from '~/lib/core/router'
import state from '~/lib/core/state'
import '~/lib/utils/global-prototypes'

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
                    primary: '#2196F3',
                    secondary: '#424242',
                    accent: '#FF4081',
                    error: '#FF5252',
                    info: '#2196F3',
                    success: '#4CAF50',
                    warning: '#FB8C00',
                    // Discord:    #7289DA
                    // Bravery:    #9B84EE
                    // Online:     #34B581
                    // Balance:    #44DDBF
                    // DND:        #F04747
                    // Brilliance: #F47B68
                    // Idle:       #FAA61A
                    // High:       #F57731
                    // Nitro:      #FF73FA
                    // Skin:       #F9C9A9
                    // White:      #FFFFFF
                    // Grey:       #99AAB5
                    // Chat:       #36393F
                    // Channels:   #2F3136
                    // Servers:    #2F3136
                }
            }
        }
    }),
    router,
    store: state,
    el: '#app',
    render: h=> h(app)
})
