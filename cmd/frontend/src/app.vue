<template>
    <v-app>
        <vue-progress-bar></vue-progress-bar>
        <router-view></router-view>
    </v-app>
</template>

<script>
export default {
    name: "app",
}
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
    background-color: #2f3136;
    background-clip: padding-box;
}
::-webkit-scrollbar-thumb {
    background-color: #7289da;
    background-clip: padding-box;
    border: 3px solid #36393f;
    border-radius: 7px;
}
::-webkit-scrollbar-thumb:window-inactive {
    background-color: #99aab5;
}
</style>
