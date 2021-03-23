<template>
    <v-container>
        <v-layout column align-center justify-center>
            <h1 class="display-3">Uhoh! Page not found.</h1>

            <p>
                The page you were looking for <span v-if="last" class="font-weight-bold">({{ last }})</span> doesn't exist. Maybe try heading back to
                the home page?
            </p>
            <router-link :to="{ name: 'index' }" tag="span">
                <v-btn color="light-blue" dark>Home</v-btn>
            </router-link>
        </v-layout>
    </v-container>
</template>

<script>
export default {
    name: "not-found",
    computed: {
        last() { return this.$route.redirectedFrom }
    }
}
</script>
