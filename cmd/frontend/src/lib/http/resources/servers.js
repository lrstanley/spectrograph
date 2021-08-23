import http from '@/lib/http/http'

function list(options) { return http.get('/servers', options) }
function get(id, options) { return http.get(`/servers/${id}`, options) }

export default {
    list,
    get,
}
