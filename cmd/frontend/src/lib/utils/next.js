const NEXT_ROUTE = "next-route"

function store(route) {
    if (!route || !route.path) { return }
    window.localStorage.setItem(NEXT_ROUTE, JSON.stringify(route))
}

function restore(router, defaultRoute) {
    if (!router) {
        throw 'Invalid router object provided'
    }

    let route = window.localStorage.getItem(NEXT_ROUTE)
    window.localStorage.removeItem(NEXT_ROUTE)

    // if we previously stored a route to go to, but we had to intercept
    // and have the user login, try and redirect back to that route.
    if (!route) {
        if (defaultRoute) {
            return router.push(defaultRoute)
        }
        return
    }

    route = JSON.parse(route)
    return router.push({ path: route.path, query: route.query })
}

export default {
    store: store,
    restore: restore
}
