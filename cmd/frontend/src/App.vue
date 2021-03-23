<template>
    <v-app>
        <vue-progress-bar></vue-progress-bar>
        <v-toolbar color="light-blue lighten-1" app fixed clipped-left>
            <v-toolbar-side-icon class="hidden-md-and-up" @click="drawer = !drawer"></v-toolbar-side-icon>
            <v-icon class="hidden-sm-and-down">code</v-icon>
            <router-link :to="{ name: 'index' }" tag="span" class="title ml-3 mr-5 clickable">src<span class="font-weight-light">-todo</router-link>
            <v-text-field solo-inverted flat hide-details label="Search" prepend-inner-icon="search"></v-text-field>
            <v-spacer></v-spacer>

            <v-btn v-if="auth == false" flat light href="/api/v1/auth/github/redirect">Login</v-btn>

            <v-progress-circular v-if="auth == null" indeterminate color="purple"></v-progress-circular>

            <!-- :nudge-width="200"; TODO: https://github.com/vuetifyjs/vuetify/issues/3438 -->
            <v-menu v-if="auth" offset-y :nudge-top="-15">
                <!-- <v-btn slot="activator" flat light>Hello, Liam</v-btn> -->
                <v-avatar slot="activator" size="36px">
                    <img class="user-avatar" :src="auth.avatar_url" :alt="auth.name">
                </v-avatar>

                <v-card>
                    <v-list>
                        <v-list-tile avatar>
                            <v-list-tile-avatar>
                                <img :src="auth.avatar_url" :alt="auth.name">
                            </v-list-tile-avatar>

                            <v-list-tile-content>
                                <v-list-tile-title>{{ auth.name }}</v-list-tile-title>
                                <v-list-tile-sub-title>{{ auth.username }}</v-list-tile-sub-title>
                            </v-list-tile-content>
                        </v-list-tile>
                    </v-list>

                    <v-divider></v-divider>

                    <v-list>
                        <v-list-tile to="/admin/settings">
                            <v-list-tile-action>
                                <v-icon>security</v-icon>
                            </v-list-tile-action>
                            <v-list-tile-title>Admin settings</v-list-tile-title>
                        </v-list-tile>
                        <v-list-tile href="/api/v1/auth/github/manage">
                            <v-list-tile-action>
                                <v-icon>settings</v-icon>
                            </v-list-tile-action>
                            <v-list-tile-title>Manage Github permissions</v-list-tile-title>
                        </v-list-tile>
                        <v-list-tile @click="logout()">
                            <v-list-tile-action>
                                <v-icon>lock</v-icon>
                            </v-list-tile-action>
                            <v-list-tile-title>Sign out</v-list-tile-title>
                        </v-list-tile>
                    </v-list>
                </v-card>
            </v-menu>
        </v-toolbar>
        <v-content>
            <v-container fluid fill-height class="grey lighten-4">
                <router-view></router-view>
            </v-container>
        </v-content>
    </v-app>
</template>

<script>
export default {
    name: 'app',
    computed: {
        auth() { return this.$store.state.auth; }
    },
    methods: {
        test: function() {
            this.$http.get("/api/v1/test")
        },
        logout: function() {
            this.$store.dispatch('logout').then(() => {
                this.$router.push('/')
            })
        }
    },
    created: function() {
        this.$Progress.start()
        this.$router.beforeEach((to, from, next) => {
            if (to.meta.progress !== undefined) {
                let meta = to.meta.progress
                this.$Progress.parseMeta(meta)
            }

            // console.log(to.name)
            // if ((to.meta.auth_required === true || to.name.startswith('admin-')) && this.$store.state.auth === null) {
            //     window.location.replace("/api/v1/auth/github/redirect")
            // }

            this.$Progress.start()
            next()
        })

        this.$router.afterEach((to, from) => {
            this.$Progress.finish()
        })

        this.$http.interceptors.request.use((config) => {
            this.$Progress.start()
            return config
        }, (error) => {
            this.$Progress.fail()
            return Promise.reject(error)
        })

        this.$http.interceptors.response.use((response) => {
            this.$Progress.finish()
            return response
        }, (error) => {
            this.$Progress.fail()
            return Promise.reject(error)
        })
    },
    mounted: function() {
        this.$Progress.finish()
    }
};
</script>

<style>
.clickable {
    cursor: pointer;
}
.user-avatar:hover {
    box-shadow: 0 0 8px white;
    transition: box-shadow 0.2s ease-in-out;
}

.v-snack [role="progressbar"] {
    display: none;
}

::-webkit-scrollbar {
    width: 10px;
    height: 6px;
}
::-webkit-scrollbar-track-piece {
    background-color: #F5F5F5;
    background-clip: padding-box;
}
::-webkit-scrollbar-thumb {
    background-color: #1678c2;
    background-clip: padding-box;
    border: 2px solid #FFFFFF;
    border-radius: 6px;
}
::-webkit-scrollbar-thumb:window-inactive {
    background-color: #1678c2;
}
</style>
