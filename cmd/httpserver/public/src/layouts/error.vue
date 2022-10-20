<template>
  <div class="flex flex-col flex-auto">
    <main class="w-full px-6 mx-auto max-w-7xl md:px-8">
      <div class="pt-8 shrink-0">
        <img loading="eager" alt="Spectrograph logo" class="w-auto h-20 mx-auto" :src="imgMicWithBg" />
      </div>

      <div class="max-w-xl py-8 mx-auto">
        <router-view v-slot="{ Component, route }">
          <transition name="fade" mode="out-in" appear>
            <Suspense>
              <component :is="Component" :key="route.path" />

              <template #fallback>
                <div class="flex flex-col h-full gap-4 mx-auto place-content-center">
                  <i-fas-circle-notch class="h-12 text-4xl align-middle animate-spin text-discord-500" />
                </div>
              </template>
            </Suspense>
          </transition>
        </router-view>

        <div class="mt-12">
          <ul role="list" class="mt-4 divide-y divide-channel-600 border-channel-600 border-y">
            <li
              v-for="link in mainLinks"
              :key="link.name"
              class="relative flex items-start py-4 space-x-4"
            >
              <div class="shrink-0">
                <span class="flex items-center justify-center w-12 h-12 rounded-lg bg-discord-600">
                  <component :is="link.icon" class="w-6 h-6 text-white" aria-hidden="true" />
                </span>
              </div>
              <div class="flex-1 min-w-0">
                <h3 class="text-base font-medium text-nitro-500">
                  <span
                    class="rounded-sm focus-within:ring-2 focus-within:ring-bravery-500 focus-within:ring-offset-2"
                  >
                    <a href="#" class="focus:outline-none">
                      <span class="absolute inset-0" aria-hidden="true" />
                      {{ link.name }}
                    </a>
                  </span>
                </h3>
                <p class="text-base text-white">{{ link.description ?? "---" }}</p>
              </div>
              <div class="self-center shrink-0">
                <i-fas-chevron-right class="w-5 h-5 text-gray-400" aria-hidden="true" />
              </div>
            </li>
          </ul>
          <div class="mt-8">
            <a
              href="/"
              class="inline-flex text-base font-medium text-bravery-600 hover:text-bravery-500"
            >
              <div class="self-center shrink-0">
                <i-fas-chevron-left class="w-4 h-4 mr-2 text-gray-400" aria-hidden="true" />
              </div>

              Or go back home
            </a>
          </div>

          <CoreFooter class="border-t border-channel-600" />
        </div>
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import imgMicWithBg from "@/assets/img/mic-with-bg.png?format=webp&imagetools"
import { mainLinks } from "@/lib/core/navigation"
</script>
