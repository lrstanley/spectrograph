<template>
  <footer
    v-if="props.fullSize"
    class="p-10 rounded footer footer-center text-base-content"
    v-bind="$attrs"
  >
    <div class="grid grid-flow-row gap-4 md:gap-8 md:grid-flow-col">
      <a v-for="item in mainLinks" :key="item.name" :href="item.href" class="link link-hover">
        {{ item.name }}
      </a>
    </div>
    <div>
      <div class="grid grid-flow-col gap-4">
        <a
          v-for="item in socialLinks"
          :key="item.name"
          :href="item.href"
          class="text-gray-400 fill-current hover:text-gray-500"
        >
          <span class="sr-only">{{ item.name }}</span>
          <component :is="item.icon" class="w-6 h-6" aria-hidden="true" />
        </a>
      </div>
    </div>
    <div>
      <a href="https://liam.sh/" target="_blank" class="link link-hover">
        &copy; {{ new Date().getFullYear() }} Spectrograph
      </a>
    </div>
  </footer>

  <footer v-else class="py-4 mt-8 text-center md:flex md:justify-between" v-bind="$attrs">
    <p class="text-base text-gray-400">
      <a href="https://liam.sh/" target="_blank"> &copy; {{ new Date().getFullYear() }} Spectrograph </a>
    </p>
    <div class="flex justify-center mt-6 space-x-8 md:mt-0">
      <a
        v-for="link in socialLinks"
        :key="link.name"
        :href="link.href"
        class="inline-flex text-gray-400 hover:text-gray-500"
      >
        <span class="sr-only">{{ link.name }}</span>
        <component :is="link.icon" class="w-6 h-6" aria-hidden="true" />
      </a>
    </div>
  </footer>
</template>

<script setup lang="ts">
import { mainLinks, socialLinks } from "@/lib/core/navigation"

const props = defineProps<{
  fullSize?: boolean
}>()
</script>
