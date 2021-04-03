<template>
    <v-app>
        <v-container fill-height fluid>
            <v-row align="center" justify="center">
                <v-col cols="12" sm="12" md="8" lg="4" class="text-center">
                    <div v-if="error">
                        <span class="d-flex">
                            <v-btn text color="primary" class="mb-2" @click.up="handleNext()">
                                <v-icon>mdi-chevron-left</v-icon> go back
                            </v-btn>
                            <v-btn text color="primary" class="mb-2 ml-auto" @click.up="$router.push({ name: 'index' })">
                                <v-icon>mdi-home-outline</v-icon> home
                            </v-btn>
                        </span>
                        <v-alert border="left" colored-border type="error" elevation="2">
                            <v-row align="center">
                                <v-col class="grow">{{ error }}</v-col>
                                <v-col class="shrink">
                                    <v-btn color="error" @click.up="$router.push({ name: 'auth', params: { method: 'redirect' } })">Try again</v-btn>
                                </v-col>
                            </v-row>
                        </v-alert>
                    </div>
                    <v-progress-circular v-else :size="100" color="primary" class="align-self-center ma-10" indeterminate></v-progress-circular>
                    <div class="grey--text footer"><footer-metadata></footer-metadata></div>
                </v-col>
            </v-row>
        </v-container>
    </v-app>
</template>

<script>
import footerMetadata from '~/components/core/footer-metadata'

const NEXT_ROUTE = "next-route"

export default {
    name: "auth",
    title: "Login",
    components: { footerMetadata },
    props: {
        other_error: [String, null],
        next: [Object, null]
    },
    data: function() { return {
        error: null
    }},
    methods: {
        handleNext: function() {
            let next = window.localStorage.getItem(NEXT_ROUTE)
            window.localStorage.removeItem(NEXT_ROUTE)

            console.log(next)
            // TODO: test this.
            // if we previously stored a route to go to, but we had to intercept
            // and have the user login, try and redirect back to that route.
            if (!next) {
                this.$router.push({ name: 'index' })
            }

            next = JSON.parse(next)

            console.log("using next...")
            return this.$router.push({ path: next.path, query: next.query })
        },
        handle: function() {
            this.error = this.other_error

            if (this.next) {
                window.localStorage.setItem(NEXT_ROUTE, JSON.stringify({ path: this.next.path, query: this.next.query }))
            }

            if (this.$route.params.method == "redirect") {
                // TODO: grab from api package, intercept location and redirect
                // ourselves, so the users doesn't have the /api endpoint in their
                // back button.
                // window.location.replace(`${window.location.protocol}//${window.location.host}${this.$config.api_baseurl}/auth/redirect`)
                this.$api.auth.redirect.get().then((resp) => {
                    window.location.replace(resp.data.auth_redirect)
                }).catch((error) => {
                    this.error = error.message
                })
                return
            }

            if (this.$route.params.method == "logout") {
                this.$api.auth.logout.get().then((resp) => {
                    this.$store.commit('set_auth', false)
                    this.$router.push({ name: 'index' })
                }).catch((error) => {
                    this.error = error.message
                })
                return
            }

            if (this.$route.params.method == "callback") {
                this.$api.auth.callback.get(this.$route.query.code, this.$route.query.state).then((resp) => {
                    // if successful, we theoretically should be able to obtain our
                    // user info. if not, return the errors.
                    this.$store.dispatch('get_auth').then(() => {
                        this.handleNext()
                    }).catch((err) => {
                        this.error = err.message
                    })
                }).catch((err) => {
                    this.error = err.message
                })
            }
        }
    },
    beforeRouteUpdate: function(to, from, next) {
        // if this.$route.params.method changes.
        console.log(this)
        next()
        this.handle()
    },
    mounted: function() { return this.handle() }
}
</script>

<style scoped>
.footer {
    font-size: 14px;
    margin-top: 5px;
    text-align: center;
}
</style>
