<template>
  <Landing>
    <div class="grid lg:w-[72rem] gap-5 px-6 mx-auto mt-14 md:gap-10">
      <div class="flex flex-col self-start flex-auto gap-4 md:sticky md:left-0 md:top-10">
        <!-- throw this out? it's a PITA? -->
        <NavTableOfContents :element="el" />
      </div>

      <main>
        <router-view v-slot="{ Component, route }">
          <transition name="fade" mode="out-in" appear>
            <Suspense ref="el">
              <component :is="Component" :key="route.path" />

              <template #fallback>
                <div class="flex flex-col h-full gap-4 mx-auto place-content-center">
                  <i-fas-circle-notch class="h-12 text-4xl align-middle animate-spin text-discord-500" />
                </div>
              </template>
            </Suspense>
          </transition>
        </router-view>
      </main>
    </div>
  </Landing>
</template>

<script setup lang="ts">
import Landing from "@/layouts/landing.vue"

const el = ref<HTMLElement>()
</script>

<style scoped>
.grid {
  @apply grid-cols-1 md:grid-cols-[240px,1fr];
}
</style>
