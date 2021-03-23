<template>
    <v-layout justify-center align-center>
        <v-progress-circular v-if="!error" :size="90" :width="7" color="light-blue" indeterminate></v-progress-circular>
        <v-alert :value="error" color="error" icon="warning" outline>{{ error }}</v-alert>
    </v-layout>
</template>

<script>
export default {
    name: "auth",
    data: () => ({
        error: false,
    }),
    mounted: function() {
        // this.$store.commit('disable_auth', true);
        // Send the oauth response GET data to the backend. We're doing
        // this on the frontend, so that way if there is an error, we can
        // better show it to the user. If we do it all on the backend,
        // we'd have to figure out some way of notifying the frontend,
        // which is... messy.
        this.$http.get(
            '/api/v1/auth/github/callback', {
                params: {
                    code: this.$route.query.code,
                    state: this.$route.query.state
                }
            }
        ).catch((error) => {
            this.error = error.response.data.error || error;
        }).finally(() => {
            this.$store.dispatch('get_auth').then(() => {
                this.$router.push('/');
            })
        })
    }
}
</script>
