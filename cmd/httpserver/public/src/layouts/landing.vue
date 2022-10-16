<template>
  <div class="flex flex-col flex-auto max-w-full">
    <Popover as="header" class="relative">
      <NavMain class="relative px-4 pt-6 mx-auto max-w-7xl sm:px-6" />

      <transition
        enter-active-class="duration-150 ease-out"
        enter-from-class="scale-95 opacity-0"
        enter-to-class="scale-100 opacity-100"
        leave-active-class="duration-100 ease-in"
        leave-from-class="scale-100 opacity-100"
        leave-to-class="scale-95 opacity-0"
      >
        <PopoverPanel focus class="absolute inset-x-0 top-0 z-10 p-2 transition origin-top md:hidden">
          <div class="overflow-hidden rounded-lg shadow-md bg-channel-500 ring-1 ring-black/5">
            <NavMain class="px-5 py-4" is-panel />
          </div>
        </PopoverPanel>
      </transition>
    </Popover>

    <main>
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
    </main>

    <CoreFooter full-size />
  </div>
</template>
