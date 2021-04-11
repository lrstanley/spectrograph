<template>
    <v-container>
        <v-row>
            <v-col class="d-none d-lg-block" cols="3" xl="2">
                <v-sheet class="nav-radius">
                    <nav-user-sidebar
                        v-if="$store.getters.authed"
                        :value="true"
                        floating
                        :mobile-breakpoint="0"
                        width="100%"
                        color="servers"
                        class="nav-radius"
                    />
                </v-sheet>
            </v-col>

            <v-col>
                <v-sheet min-height="70vh" rounded="lg" class="d-flex flex-column align-stretch justify-start pa-4">
                    <v-overlay :value="$store.getters.loading" absolute>
                        <v-progress-circular indeterminate size="64" />
                    </v-overlay>
                    <router-view v-show="!$store.getters.loading" />
                </v-sheet>
            </v-col>
        </v-row>
    </v-container>
</template>

<script>
import navUserSidebar from "~/components/core/nav-user-sidebar.vue"

export default {
    name: "default-layout",
    components: { navUserSidebar },
    data: function () {
        return {}
    },
}
</script>

<style scoped>
.nav-radius {
    border-radius: 8px;
}
</style>
