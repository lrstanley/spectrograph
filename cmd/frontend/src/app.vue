<template>
    <router-view></router-view>
</template>

<script>
export default {
    name: 'app',
    methods: {},
    created: function() {
        this.$Progress.start()
        this.$router.beforeEach((to, from, next) => {
            if (to.meta.progress !== undefined) {
                let meta = to.meta.progress
                this.$Progress.parseMeta(meta)
            }

            this.$Progress.start()
            next()
        })

        this.$router.afterEach((to, from) => {
            this.$Progress.finish()
        })

        this.$http.interceptors.request.use((config) => {
            this.$Progress.start()
            return config
        }, (error) => {
            this.$Progress.fail()
            return Promise.reject(error)
        })

        this.$http.interceptors.response.use((response) => {
            this.$Progress.finish()
            return response
        }, (error) => {
            this.$Progress.fail()
            return Promise.reject(error)
        })
    },
    mounted: function() {
        this.$Progress.finish()
    }
};
</script>

<style>
.clickable {
    cursor: pointer;
}
.user-avatar:hover {
    box-shadow: 0 0 8px white;
    transition: box-shadow 0.2s ease-in-out;
}

::-webkit-scrollbar {
    width: 10px;
    height: 6px;
}
::-webkit-scrollbar-track-piece {
    background-color: #F5F5F5;
    background-clip: padding-box;
}
::-webkit-scrollbar-thumb {
    background-color: #1678c2;
    background-clip: padding-box;
    border: 2px solid #FFFFFF;
    border-radius: 6px;
}
::-webkit-scrollbar-thumb:window-inactive {
    background-color: #1678c2;
}
</style>
