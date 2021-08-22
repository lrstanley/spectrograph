<template>
    <!-- TODO: split user stuff out of this? -->
    <v-navigation-drawer v-model="drawer" v-bind="$attrs">
        <v-list-item v-if="authed">
            <v-list-item-avatar>
                <v-img :src="user.avatar_url" alt="user avatar" />
            </v-list-item-avatar>
            <v-list-item-content>
                <v-list-item-title class="title text-truncate" :title="user.username">
                    {{ user.username }}
                </v-list-item-title>
                <v-list-item-subtitle class="grey--text text--lighten-1"> #{{ user.discriminator }} </v-list-item-subtitle>
            </v-list-item-content>
        </v-list-item>

        <v-list v-if="!userNav" color="transparent" dense nav>
            <v-list-item exact :to="{ name: 'index' }">
                <v-list-item-icon>
                    <v-icon>{{ mdiHome }}</v-icon>
                </v-list-item-icon>
                <v-list-item-title>Home</v-list-item-title>
            </v-list-item>
            <v-list-item exact :to="{ name: 'user-details' }">
                <v-list-item-icon>
                    <v-icon>{{ mdiServer }}</v-icon>
                </v-list-item-icon>
                <v-list-item-title>Manage Servers</v-list-item-title>
            </v-list-item>
            <v-list-item exact link>
                <v-list-item-icon>
                    <v-icon>{{ mdiFileDocumentEdit }}</v-icon>
                </v-list-item-icon>
                <v-list-item-title>Documentation</v-list-item-title>
            </v-list-item>

            <v-divider class="my-4" />
        </v-list>
        <v-list v-if="authed" color="transparent" dense nav>
            <v-list-item v-for="server in user.joined_servers" :key="server.id" link>
                <v-list-item-avatar>
                    <v-img v-if="server.icon_url" :src="server.icon_url" alt="server icon" />
                    <v-avatar v-else color="#1E1E1E">
                        <span class="white--text headline">{{ serverInitials(server.name) }}</span>
                    </v-avatar>
                </v-list-item-avatar>
                <v-list-item-content>
                    <v-list-item-title>{{ server.name }}</v-list-item-title>
                </v-list-item-content>
                <v-icon :color="true ? 'success' : 'error'">{{ true ? mdiCheck : mdiCloseCircleOutline }}</v-icon>
            </v-list-item>
        </v-list>

        <template v-slot:append>
            <!-- <v-list-item exact :to="{ name: 'auth', params: { method: 'logout' } }">
                    <v-list-item-icon><v-icon>{{ mdiLockRemove }}</v-icon></v-list-item-icon>
                    <v-list-item-title>Logout</v-list-item-title>
                </v-list-item> -->
            <div class="pa-2">
                <v-btn color="success" block :href="$config.bot_auth_url" target="_blank">Add to server</v-btn>
            </div>
            <div v-if="!userNav" class="pa-2">
                <v-btn v-if="authed" block exact :to="{ name: 'auth', params: { method: 'logout' } }"> Logout </v-btn>
                <v-btn v-else block exact :to="{ name: 'auth', params: { method: 'redirect' } }">Login</v-btn>
            </div>
        </template>
    </v-navigation-drawer>
</template>

<script>
import { mdiHome, mdiServer, mdiFileDocumentEdit, mdiCheck, mdiCloseCircleOutline } from "@mdi/js"
import { mapGetters } from "vuex"

export default {
    name: "nav-sidebar-user",
    inheritAttrs: false,
    props: {
        value: {
            type: Boolean,
            default: false,
        },
        userNav: {
            type: Boolean,
            default: false,
        },
    },
    data: function () {
        return {
            mdiHome,
            mdiServer,
            mdiFileDocumentEdit,
            mdiCheck,
            mdiCloseCircleOutline,
        }
    },
    computed: {
        ...mapGetters(["authed", "user"]),
        drawer: {
            get: function () {
                return this.value
            },
            set: function (value) {
                this.$emit("input", value)
            },
        },
    },
    methods: {
        serverInitials: (name) => {
            if (!name || name.length < 1) return "?"
            let letter = name.match(/([a-zA-Z0-9])/g)
            if (letter !== null) {
                return letter[0]
            }
            return name[0]
        },
    },
}
</script>
