<template>
    <v-container>This is a test of {{ id }}</v-container>
</template>

<script>
export default {
    name: "user-server",
    title: "Manage Server",
    beforeRouteUpdate: function (to, from, next) {
        this.resetComponentData()
        next()
    },
    data: function () {
        return {
            id: this.$route.params.id,
        }
    },
}
</script>
