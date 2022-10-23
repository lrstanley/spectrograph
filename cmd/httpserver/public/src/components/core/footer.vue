<template>
  <footer
    v-if="props.fullSize"
    class="p-10 rounded footer footer-center text-base-content"
    v-bind="$attrs"
    aria-label="page footer"
  >
    <nav class="grid grid-flow-row gap-4 md:gap-8 md:grid-flow-col" aria-label="service links">
      <template v-for="link in mainLinks" :key="link.name">
        <component
          :is="link.to ? 'router-link' : 'a'"
          v-bind="link.to ? { to: link.to } : { href: link.href }"
          class="link link-hover"
        >
          {{ link.name }}
        </component>
      </template>
    </nav>
    <nav class="grid grid-flow-col gap-4" aria-label="social links">
      <a
        v-for="item in socialLinks"
        :key="item.name"
        :href="item.href"
        class="text-gray-400 fill-current hover:text-gray-500"
      >
        <span class="sr-only">{{ item.name }}</span>
        <component :is="item.icon" class="w-6 h-6" aria-hidden="true" />
      </a>
    </nav>
    <div aria-hidden="true">
      <a href="https://liam.sh/" target="_blank" class="link link-hover">
        &copy; {{ new Date().getFullYear() }} Spectrograph
      </a>
    </div>
  </footer>

  <footer
    v-else
    class="py-4 mt-8 text-center md:flex md:justify-between"
    v-bind="$attrs"
    aria-label="page footer"
  >
    <p class="text-base text-gray-400" aria-hidden="true">
      <a href="https://liam.sh/" target="_blank"> &copy; {{ new Date().getFullYear() }} Spectrograph </a>
    </p>
    <nav class="flex justify-center mt-6 space-x-8 md:mt-0" aria-label="social links">
      <a
        v-for="link in socialLinks"
        :key="link.name"
        :href="link.href"
        class="inline-flex text-gray-400 hover:text-gray-500"
      >
        <span class="sr-only">{{ link.name }}</span>
        <component :is="link.icon" class="w-6 h-6" aria-hidden="true" />
      </a>
    </nav>
  </footer>
</template>

<script setup lang="ts">
import { mainLinks, socialLinks } from "@/lib/core/navigation"

const props = defineProps<{
  fullSize?: boolean
}>()
</script>
