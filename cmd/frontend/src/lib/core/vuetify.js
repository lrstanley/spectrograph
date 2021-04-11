import Vue from 'vue'
import Vuetify from 'vuetify/lib'
import '@mdi/font/css/materialdesignicons.css'
import colors from '~/lib/utils/discord-colors'

Vue.use(Vuetify)

export default new Vuetify({
    icons: {
        iconfont: 'mdi'
    },
    theme: {
        dark: true,
        default: 'dark',
        options: {
            themeCache: {
                get: (key) => localStorage.getItem(key),
                set: (key, value) => localStorage.setItem(key, value),
            },
        },
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
})
