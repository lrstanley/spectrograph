<template>
    <div class="d-flex align-center flex-grow-1">
        <v-btn text exact :to="{ name: 'index' }" class="ml-5">Home</v-btn>
        <v-btn text :to="{ name: 'user-details' }">Manage Servers</v-btn>
        <v-spacer />
        <v-btn text>Documentation</v-btn>
        <v-badge color="accent" :icon="mdiDiscord" overlap>
            <v-btn color="discord" href="https://liam.sh/chat" target="_blank">Support</v-btn>
        </v-badge>
        <v-responsive max-width="260" class="d-none d-md-block ml-10 mr-5">
            <v-text-field dense flat hide-details rounded solo :prepend-inner-icon="mdiMagnify" placeholder="Search servers" />
        </v-responsive>
    </div>
</template>

<script>
import { mdiDiscord, mdiMagnify } from "@mdi/js"

export default {
    name: "nav-appbar",
    data: function () {
        return {
            mdiDiscord,
            mdiMagnify,
        }
    },
}
</script>

<style scoped>
</style>
