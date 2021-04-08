<template>
    <v-navigation-drawer v-bind="$attrs" v-model="drawer">
        <v-list-item>
            <v-list-item-avatar>
                <v-img :src="$store.getters.user.avatar_url" alt="user avatar"></v-img>
            </v-list-item-avatar>
            <v-list-item-content>
                <v-list-item-title class="title">
                    {{ $store.getters.user.username }}
                </v-list-item-title>
                <v-list-item-subtitle class="grey--text text--lighten-1"> #{{ $store.getters.user.discriminator }} </v-list-item-subtitle>
            </v-list-item-content>
        </v-list-item>

        <v-list color="transparent" dense nav>
            <v-list-item v-for="n in 5" :key="n" link>
                <v-list-item-content>
                    <v-list-item-title>List Item {{ n }}</v-list-item-title>
                </v-list-item-content>
            </v-list-item>

            <!-- <v-divider class="my-2"></v-divider> -->
            <v-divider></v-divider>

            <v-list-item link>
                <v-list-item-content>
                    <v-list-item-title>Refresh</v-list-item-title>
                </v-list-item-content>
            </v-list-item>
        </v-list>
    </v-navigation-drawer>
</template>

<script>
export default {
    name: "nav-sidebar",
    inheritAttrs: false,
    props: ["value"],
    computed: {
        drawer: {
            get: function () {
                return this.value
            },
            set: function (value) {
                this.$emit("input", value)
            },
        },
    },
}
</script>
