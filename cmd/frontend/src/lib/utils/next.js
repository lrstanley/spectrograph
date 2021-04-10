import Promise from 'q'

const NEXT_ROUTE = "next-route"

function store(next) {
    if (!next || !next.path) { return }
    window.localStorage.setItem(NEXT_ROUTE, JSON.stringify({ path: next.path, query: next.query }))
}

function restore(router, defaultRoute) {
    if (!router) {
        throw 'Invalid router object provided'
    }

    return new Promise((resolve, reject) => {
        let next = window.localStorage.getItem(NEXT_ROUTE)
        window.localStorage.removeItem(NEXT_ROUTE)

        // if we previously stored a route to go to, but we had to intercept
        // and have the user login, try and redirect back to that route.
        if (!next) {
            if (defaultRoute) {
                router.push(defaultRoute)
            }

            reject()
            return
        }

        next = JSON.parse(next)

        router.push({ path: next.path, query: next.query })
        resolve()
    })
}

export default {
    store: store,
    restore: restore
}
