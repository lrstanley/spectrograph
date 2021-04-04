<template>
    <v-app>
        <vue-progress-bar></vue-progress-bar>
        <router-view></router-view>
    </v-app>
</template>

<script>
export default {
    name: 'app',
    methods: {}
};
</script>

<style>
.clickable {
    cursor: pointer;
}

::-webkit-scrollbar {
    width: 12px;
    height: 6px;
}
::-webkit-scrollbar-track-piece {
    background-color: #2F3136;
    background-clip: padding-box;
}
::-webkit-scrollbar-thumb {
    background-color: #7289DA;
    background-clip: padding-box;
    border: 3px solid #36393F;
    border-radius: 7px;
}
::-webkit-scrollbar-thumb:window-inactive {
    background-color: #99AAB5;
}
</style>
