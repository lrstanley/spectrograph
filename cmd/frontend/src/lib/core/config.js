import Vue from 'vue'
import config from '~/config.json'

if (!Object.prototype.hasOwnProperty.call(config.envs, config.target)) {
    console.warn(`[${config.application}] unknown environment target, defaulting to local:`, config.target)
    config.target = "local"
}
var c = { ...config.envs[config.target], ...config }
c.target = config.target

c.debugLog = function (...data) { if (c.debug) { console.warn(`[${config.application}:debug]`, ...data) } }

// register global prototypes.
Vue.prototype.$config = c
Vue.prototype.debug = c.debugLog

export default c
