<template>
  <div class="flex flex-col flex-auto w-full">
    <TransitionRoot as="template" :show="sidebarOpen">
      <Dialog as="div" class="relative z-40 md:hidden" @close="sidebarOpen = false">
        <TransitionChild
          as="template"
          enter-from="opacity-0"
          enter-to="opacity-100"
          leave-from="opacity-100"
          leave-to="opacity-0"
        >
          <div class="fixed inset-0 transition-all duration-300 ease-linear bg-channel-600/75" />
        </TransitionChild>

        <div class="fixed inset-0 z-40 flex">
          <TransitionChild
            as="template"
            enter-from="-translate-x-full"
            enter-to="translate-x-0"
            leave-from="translate-x-0"
            leave-to="-translate-x-full"
          >
            <DialogPanel
              class="relative flex flex-col flex-1 w-full max-w-xs transition duration-300 ease-in-out"
            >
              <TransitionChild
                as="template"
                enter="ease-in-out duration-300"
                enter-from="opacity-0"
                enter-to="opacity-100"
                leave="ease-in-out duration-300"
                leave-from="opacity-100"
                leave-to="opacity-0"
              >
                <div class="absolute top-0 right-0 pt-2 -mr-12">
                  <button
                    type="button"
                    class="flex items-center justify-center w-10 h-10 ml-1 rounded-full focus:outline-none focus:ring-2 focus:ring-inset focus:ring-channel-400"
                    @click="sidebarOpen = false"
                  >
                    <i-fas-xmark class="w-6 h-6 text-channel-400" aria-hidden="true" />
                  </button>
                </div>
              </TransitionChild>

              <NavDashboard class="flex flex-col flex-1 overflow-y-auto" />
            </DialogPanel>
          </TransitionChild>

          <!-- Force sidebar to shrink to fit close icon -->
          <div class="w-14 shrink-0" />
        </div>
      </Dialog>
    </TransitionRoot>

    <NavDashboard class="hidden md:fixed md:inset-y-0 md:flex md:w-64 md:flex-col" />

    <div class="flex flex-col flex-1 md:pl-64">
      <div class="sticky top-0 z-10 pt-1 pl-1 bg-channel-500 sm:pl-3 sm:pt-3 md:hidden">
        <button
          type="button"
          class="-ml-0.5 -mt-0.5 inline-flex h-12 w-12 items-center justify-center rounded-md text-channel-400 hover:text-gray-900 focus:outline-none focus:ring-2 focus:ring-inset focus:ring-bravery-500"
          @click="sidebarOpen = true"
        >
          <i-fas-bars class="w-6 h-6" aria-hidden="true" />
        </button>
      </div>

      <main class="flex-1 bg-chat-500">
        <div class="flex flex-col h-full px-4 py-8 mx-auto md:px-6 max-w-7xl">
          <router-view v-slot="{ Component, route }">
            <transition name="fade" mode="out-in" appear>
              <Suspense>
                <component :is="Component" :key="route.path" />

                <template #fallback>
                  <div class="flex flex-col h-full gap-4 mx-auto place-content-center">
                    <i-fas-circle-notch
                      class="h-12 text-4xl align-middle animate-spin text-discord-500"
                    />
                  </div>
                </template>
              </Suspense>
            </transition>
          </router-view>
        </div>
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
const sidebarOpen = ref(false)
</script>
