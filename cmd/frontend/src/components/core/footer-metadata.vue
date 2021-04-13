<template>
    <span>
        Made with <v-icon color="error" small>{{ mdiHeart }}</v-icon> by
        <a href="https://liam.sh/" target="_blank">Liam Stanley</a>
    </span>
</template>

<script>
import { mdiHeart } from "@mdi/js"

export default {
    name: "footer-metadata",
    data: function () {
        return {
            mdiHeart,
        }
    },
}
</script>

<style scoped>
a:visited,
a:link,
a:active {
    text-decoration: none;
}
</style>
