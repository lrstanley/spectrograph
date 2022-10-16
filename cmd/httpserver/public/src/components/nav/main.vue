<template>
  <nav class="flex flex-col justify-between md:items-center md:flex-row" aria-label="Main navigation">
    <div class="flex flex-col flex-1 md:flex-row" :class="props.isPanel ? '' : 'items-center'">
      <div class="flex items-center justify-between w-full md:w-auto">
        <a href="#" class="inline-flex items-center">
          <img alt="Spectrograph logo" class="w-12 h-12" :src="imgLogo" />
          <router-link to="/">
            <h2 class="text-xl text-gradient bg-gradient-to-r from-nitro-500 to-dnd-400">
              Spectrograph
            </h2>
          </router-link>
        </a>
        <div class="flex items-center -mr-2 md:hidden">
          <PopoverButton
            class="inline-flex items-center justify-center p-2 text-white rounded-md bg-bravery-500 focus-ring-inset hover:bg-bravery-600 focus:outline-none"
            :class="props.isPanel ? '' : 'mr-4'"
          >
            <span v-if="props.isPanel" class="sr-only">Close main menu</span>
            <span v-else class="sr-only">Open main menu</span>

            <i-fas-xmark v-if="props.isPanel" class="w-6 h-6" aria-hidden="true" />
            <i-fas-bars v-else class="w-6 h-6" aria-hidden="true" />
          </PopoverButton>
        </div>
      </div>
      <div class="md:ml-10" :class="props.isPanel ? 'flex flex-col' : 'hidden md:flex space-x-8'">
        <a
          v-for="link in headerLinks"
          :key="link.name"
          :href="link.href"
          class="text-lg font-medium text-white md:text-base hover:text-gray-300 focus:ring-1 ring-white"
          :class="props.isPanel ? 'px-4 py-1 my-1 bg-channel-600 rounded' : ''"
        >
          {{ link.name }}
        </a>
      </div>
    </div>
    <div class="items-center mt-2 space-x-6 md:mt-0" :class="props.isPanel ? 'flex' : 'hidden md:flex'">
      <router-link
        class="h-10 min-h-0 text-white md:grow-0 grow btn bg-gradient-to-r from-nitro-700/80 to-bravery-700/80 hover:from-nitro-600/80 hover:to-bravery-600/80"
        to="/dashboard"
      >
        {{ state.base.self ? "Go to Dashboard" : "Login" }}
      </router-link>
    </div>
  </nav>
</template>

<script setup lang="ts">
import imgLogo from "@/assets/img/mic.png?format=webp&imagetools"
import { headerLinks } from "@/lib/core/navigation"

const state = useState()

const props = defineProps<{
  isPanel?: boolean
}>()
</script>
