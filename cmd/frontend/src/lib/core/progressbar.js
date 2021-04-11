import Vue from 'vue'
import VueProgressBar from 'vue-progressbar'

import colors from '@/lib/utils/discord-colors'

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
