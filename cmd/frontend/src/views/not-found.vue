<template>
    <v-container fill-height fluid class="d-flex flex-column align-center justify-center">
        <h1 class="display-3">Uhoh! Page not found.</h1>

        <p>The page you were looking for doesn't exist.</p>

        <span>
            <v-btn color="warning" @click.up="$router.back()">
                <v-icon class="mr-1">{{ mdiArrowLeft }}</v-icon> Back
            </v-btn>
            <v-btn color="primary" @click.up="$router.push({ name: 'index' })">
                <v-icon class="mr-1">{{ mdiHome }}</v-icon> Home
            </v-btn>
        </span>
    </v-container>
</template>

<script>
import { mdiArrowLeft, mdiHome } from "@mdi/js"

export default {
    name: "not-found",
    data: function () {
        return {
            mdiArrowLeft,
            mdiHome,
        }
    },
}
</script>
