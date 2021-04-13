<template>
    <v-container fill-height fluid>
        <v-row align="center" justify="center">
            <v-col cols="12" sm="12" md="8" lg="4" class="text-center">
                <div v-if="error">
                    <span class="d-flex">
                        <v-btn text color="primary" class="mb-2" @click.up="$router.go(-2)">
                            <v-icon>{{ mdiChevronLeft }}</v-icon> go back
                        </v-btn>
                        <v-btn text color="primary" class="mb-2 ml-auto" @click.up="$router.push({ name: 'index' })">
                            <v-icon>{{ mdiHomeOutline }}</v-icon> home
                        </v-btn>
                    </span>
                    <v-alert border="left" colored-border type="error" elevation="2">
                        <v-row align="center">
                            <v-col class="grow">{{ error }}</v-col>
                            <v-col class="shrink">
                                <v-btn color="error" @click.up="$router.replace({ name: 'auth', params: { method: 'redirect' } })">Try again</v-btn>
                            </v-col>
                        </v-row>
                    </v-alert>
                </div>
                <v-progress-circular v-else :size="100" color="primary" class="align-self-center ma-10" indeterminate />
                <div class="grey--text footer">
                    <footer-metadata />
                </div>
            </v-col>
        </v-row>
    </v-container>
</template>

<script>
import { mdiHomeOutline, mdiChevronLeft } from "@mdi/js"
import next from "@/lib/utils/next"

export default {
    name: "auth",
    title: "Login",
    beforeRouteUpdate: function (to, from, next) {
        // if this.$route.params.method changes.
        next()
        this.handle()
    },
    props: {
        other_error: [String, null],
    },
    data: function () {
        return {
            mdiHomeOutline,
            mdiChevronLeft,
            error: null,
        }
    },
    mounted: function () {
        return this.handle()
    },
    methods: {
        handle: function () {
            this.error = this.other_error
            next.store(this.$route.params.next)

            if (this[this.$route.params.method]) {
                this[this.$route.params.method]()
                return
            }

            this.$router.push({ name: "not-found" })
        },
        redirect: function () {
            this.$api.auth
                .redirect()
                .then((resp) => {
                    window.location.replace(resp.data.auth_redirect)
                })
                .catch((error) => {
                    this.error = error.message
                })
            return
        },
        callback: function () {
            this.$api.auth
                .callback(this.$route.query.code, this.$route.query.state)
                .then(() => {
                    // if successful, we theoretically should be able to obtain our
                    // user info. if not, return the errors.
                    this.$store
                        .dispatch("get_auth")
                        .then(() => {
                            next.restore(this.$router, { name: "index" })
                        })
                        .catch((err) => {
                            this.error = err.message
                        })
                })
                .catch((err) => {
                    this.error = err.message
                })
        },
        logout: function () {
            this.$store
                .dispatch("logout")
                .then(() => {
                    this.$router.push({ name: "index" })
                })
                .catch((error) => {
                    this.error = error.message
                })
            return
        },
    },
}
</script>

<style scoped>
.footer {
    font-size: 14px;
    margin-top: 5px;
    text-align: center;
}
</style>
