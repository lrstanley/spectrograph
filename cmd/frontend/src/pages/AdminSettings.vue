<template>
  <v-container>
    <v-layout justify-center>
      <v-flex lg12 xl9>
        <admin-providers></admin-providers>
      </v-flex>
    </v-layout>
  </v-container>
</template>

<script>
// example: https://codesandbox.io/s/o29j95wx9
import Providers from '../components/AdminSettings/Providers.vue';

export default {
  name: "admin-settings",
  components: {
    "admin-providers": Providers,
  },
};
</script>

<style scoped>
.v-navigation-drawer {
  z-index: 2;
}
</style>
