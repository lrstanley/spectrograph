<template>
    <div>
        <v-app-bar app color="secondary" flat>
            <!-- TODO: handle side navigation via route metadata field? -->
            <v-app-bar-nav-icon v-if="$route.path.startsWith('/user')" class="d-lg-none d-xl-none" @click="drawer = !drawer" />
            <v-container class="py-0 fill-height">
                <v-avatar size="45" class="mr-5">
                    <v-img src="/src/static/img/mic.png" />
                </v-avatar>

                <v-btn text exact :to="{ name: 'index' }">Home</v-btn>
                <v-btn text :to="{ path: '/user/details' }">Manage Servers</v-btn>
                <v-spacer />
                <v-btn text>Documentation</v-btn>
                <v-badge color="accent" icon="mdi-discord" overlap>
                    <v-btn color="discord" href="https://liam.sh/chat" target="_blank">Support</v-btn>
                </v-badge>
                <!-- <v-btn color="accent">
                    <v-icon class="mr-1">mdi-help-circle</v-icon> Support
                </v-btn> -->
                <!-- <v-btn v-for="n in 5" :key="n" text> link {{ n }} </v-btn> -->

                <!-- <v-btn @click="$router.push({ name: 'not-found' })" text>not found</v-btn> -->

                <!-- <v-spacer></v-spacer> -->

                <v-responsive max-width="260" class="d-none d-md-block ml-10 mr-5">
                    <v-text-field dense flat hide-details rounded solo prepend-inner-icon="mdi-magnify" placeholder="Search servers" />
                </v-responsive>
                <!-- <v-avatar class="hidden-sm-and-down" color="grey darken-1 shrink" size="32"></v-avatar> -->
                <v-btn v-if="!$store.getters.authed" text @click="$router.push({ name: 'auth', params: { method: 'redirect' } })">Login</v-btn>
                <v-btn v-if="$store.getters.authed" text @click="$router.push({ name: 'auth', params: { method: 'logout' } })">Logout</v-btn>
            </v-container>
        </v-app-bar>

        <nav-user-sidebar v-if="$store.getters.authed" v-model="drawer" fixed temporary />

        <v-main>
            <router-view />
        </v-main>

        <v-footer app padless inset>
            <div class="flex-grow-1" />
            <footer-metadata />
        </v-footer>
    </div>
</template>

<script>
import navUserSidebar from "~/components/core/nav-user-sidebar.vue"
import footerMetadata from "~/components/core/footer-metadata.vue"

export default {
    name: "default-layout",
    components: { navUserSidebar, footerMetadata },
    data: function () {
        return {
            drawer: false,
        }
    },
}
</script>

<style scoped>
.nav-radius {
    border-radius: 8px;
}
.v-footer {
    padding: 0 8px;
    font-size: 14px;
}
.logo {
    max-height: 30px;
}
div.v-text-field >>> div.v-input__slot {
    padding: 0 15px 0 12px;
}
</style>
