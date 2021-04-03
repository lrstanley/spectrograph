<template>
    <v-app id="inspire" dark>
        <v-app-bar app darkcolor="white" flat>
            <v-app-bar-nav-icon class="d-lg-none d-xl-none" @click="drawer = !drawer"></v-app-bar-nav-icon>
            <v-container class="py-0 fill-height">
                <v-avatar class="mr-10" color="grey darken-1" size="32"></v-avatar>

                <v-btn v-for="n in 5" :key="n" text>
                    link {{ n }}
                </v-btn>

                <v-btn @click="$router.push({ name: 'not-found' })" text>not found</v-btn>

                <v-spacer></v-spacer>

                <v-responsive max-width="260" class="d-none d-md-block">
                    <v-text-field dense flat hide-details rounded solo></v-text-field>
                </v-responsive>
            </v-container>
        </v-app-bar>

        <v-navigation-drawer v-model="drawer" fixed temporary>
            <nav-sidebar></nav-sidebar>
        </v-navigation-drawer>

        <!-- class="grey lighten-3" -->
        <v-main>
            <v-container>
                <v-row>
                    <v-col class="d-none d-lg-block" cols="3" xl="2">
                        <v-sheet rounded="lg">
                            <nav-sidebar></nav-sidebar>
                        </v-sheet>
                    </v-col>

                    <v-col>
                        <v-sheet min-height="70vh" rounded="lg" class="d-flex flex-column align-stretch justify-start">
                            <v-overlay :value="$store.getters.loading" absolute>
                                <v-progress-circular indeterminate size="64"></v-progress-circular>
                            </v-overlay>
                            <router-view v-show="!$store.getters.loading"></router-view>
                        </v-sheet>
                    </v-col>
                </v-row>
            </v-container>
        </v-main>

        <v-footer app padless inset>
            <div class="flex-grow-1"></div>
            <div><footer-metadata></footer-metadata></div>
        </v-footer>
    </v-app>
</template>

<script>
import navSidebar from '~/components/core/nav-sidebar'
import footerMetadata from '~/components/core/footer-metadata'

export default {
    name: 'default-layout',
    components: { navSidebar, footerMetadata },
    data: function() { return {
        drawer: false,
    }},
};
</script>

<style scoped>
.v-footer {
    padding: 0 8px;
    font-size: 14px;
}
.logo {
    max-height: 30px;
}
</style>
