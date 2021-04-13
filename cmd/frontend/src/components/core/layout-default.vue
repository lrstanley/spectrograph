<template>
    <div>
        <v-app-bar app color="secondary" flat>
            <!-- TODO: handle side navigation via route metadata field? -->
            <v-app-bar-nav-icon class="d-lg-none d-xl-none" @click="drawer = !drawer" />
            <v-container class="py-0 fill-height">
                <v-avatar size="45" @click.up="$router.push({ name: authed ? 'user-details' : 'index' })">
                    <v-img src="/src/static/img/mic.png" />
                </v-avatar>

                <v-toolbar-title @click.up="$router.push({ name: authed ? 'user-details' : 'index' })">
                    {{ $config.application }}
                </v-toolbar-title>

                <nav-appbar v-if="$vuetify.breakpoint.lgAndUp" />
                <template v-else>
                    <v-spacer />
                </template>
                <!-- <v-avatar class="hidden-sm-and-down" color="grey darken-1 shrink" size="32"></v-avatar> -->
                <v-btn v-if="!authed" text exact :to="{ name: 'auth', params: { method: 'redirect' } }" class="pr-xs-0">Login</v-btn>
                <v-btn v-if="authed" text exact :to="{ name: 'auth', params: { method: 'logout' } }" class="pr-xs-0">Logout</v-btn>
            </v-container>
        </v-app-bar>

        <nav-sidebar v-model="drawer" fixed temporary />

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
import { mapGetters } from "vuex"

export default {
    name: "default-layout",
    data: function () {
        return {
            drawer: false,
        }
    },
    computed: { ...mapGetters(["authed", "loading"]) },
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
