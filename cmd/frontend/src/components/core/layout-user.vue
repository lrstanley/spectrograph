<template>
    <v-container>
        <v-row>
            <v-col class="d-none d-lg-block" cols="3" xl="2">
                <v-sheet class="rounded-lg">
                    <nav-sidebar v-if="authed" :value="true" user-nav floating :mobile-breakpoint="0" width="100%" color="servers" class="rounded-lg" />
                </v-sheet>
            </v-col>

            <v-col>
                <v-sheet min-height="70vh" rounded="lg" class="d-flex align-stretch justify-start">
                    <v-overlay :value="loading" absolute>
                        <v-progress-circular indeterminate size="64" />
                    </v-overlay>
                    <router-view v-show="!loading" />
                </v-sheet>
            </v-col>
        </v-row>
    </v-container>
</template>

<script>
import { mapGetters } from "vuex"

export default {
    name: "default-layout",
    data: function () {
        return {}
    },
    computed: { ...mapGetters(["authed", "loading"]) },
}
</script>
