<template>
    <v-row no-gutters class="flex-sm-row flex-column">
        <v-col v-if="loading" style="position: relative">
            <v-overlay absolute>
                <v-progress-circular indeterminate size="64" />
            </v-overlay>
        </v-col>
        <v-col v-if="!loading && !error && server" class="flex-grow-0">
            <v-navigation-drawer color="servers" class="accent-4 elevation-2" style="border-radius: 8px 0 0 8px" floating permanent>
                <v-list-item class="px-2">
                    <v-list-item-avatar>
                        <v-img src="https://randomuser.me/api/portraits/men/85.jpg" />
                    </v-list-item-avatar>

                    <v-list-item-title>{{ server.discord.name }}</v-list-item-title>

                    <v-btn icon @click.stop="mini = !mini">
                        <v-icon>mdi-chevron-left</v-icon>
                    </v-btn>
                </v-list-item>

                <v-divider />

                <v-list nav transparent>
                    <v-list-item v-for="item in items" :key="item.title" link color="transparent">
                        <v-list-item-icon>
                            <v-icon>{{ item.icon }}</v-icon>
                        </v-list-item-icon>

                        <v-list-item-title>{{ item.title }}</v-list-item-title>
                    </v-list-item>
                </v-list>

                <!-- <template #append>
                            <div class="pa-2">
                                <v-btn block> Logout </v-btn>
                            </div>
                        </template> -->
            </v-navigation-drawer>
        </v-col>
        <v-col v-if="!loading && !error && server" class="mr-auto">
            <v-container class="flex-grow-1"> This is a test </v-container>
        </v-col>
    </v-row>
</template>

<script>
export default {
    name: "user-server",
    title: "Manage Server",
    beforeRouteUpdate: function (to, from, next) {
        this.resetComponentData()
        next()
        this.fetch()
    },
    data: function () {
        return {
            items: [
                { title: "Dashboardddddddddddddddddddddd", icon: "mdi-view-dashboard" },
                { title: "Account", icon: "mdi-account-box" },
                { title: "Admin", icon: "mdi-gavel" },
            ],
            loading: true,
            error: null,
            server: null,
        }
    },
    mounted: function () {
        this.fetch()
    },
    methods: {
        fetch: async function () {
            try {
                this.loading = true
                let resp = await this.$api.servers.get(this.$route.params.id)

                this.server = resp.data.server
            } catch (err) {
                this.error = err
            } finally {
                setTimeout(() => {
                    this.loading = false
                }, 200)
            }
        },
    },
}
</script>
