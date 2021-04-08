import config from '~/config.json'

if (!config.envs.hasOwnProperty(config.target)) {
    console.info("unknown environment target, defaulting to local:", config.target)
    config.target = "local"
}
var c = { ...config.envs[config.target], ...config }
c.target = config.target

c.debugLog = function (...data) { if (c.debug) { console.log(...data) } }

export default c
