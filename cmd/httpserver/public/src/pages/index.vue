<route lang="yaml">
meta:
  title: "Home"
</route>

<template>
  <div class="flex flex-col flex-auto max-w-full">
    <Popover as="header" class="relative">
      <div class="pt-6">
        <nav
          class="relative flex items-center justify-between px-4 mx-auto max-w-7xl sm:px-6"
          aria-label="Global"
        >
          <div class="flex items-center flex-1">
            <div class="flex items-center justify-between w-full md:w-auto">
              <a href="#" class="inline-flex items-center">
                <img class="w-auto h-12" src="/img/mic.png" />
                <h2 class="text-xl text-gradient bg-gradient-to-r from-nitro-500 to-dnd-400">
                  Spectrograph
                </h2>
              </a>
              <div class="flex items-center -mr-2 md:hidden">
                <PopoverButton
                  class="inline-flex items-center justify-center p-2 text-white rounded-md bg-bravery-500 focus-ring-inset hover:bg-bravery-600 focus:outline-none"
                >
                  <span class="sr-only">Open main menu</span>
                  <i-fas-bars class="w-6 h-6" aria-hidden="true" />
                </PopoverButton>
              </div>
            </div>
            <div class="hidden space-x-8 md:ml-10 md:flex">
              <a
                v-for="link in headerLinks"
                :key="link.name"
                :href="link.href"
                class="text-base font-medium text-white hover:text-gray-300 focus:ring-1 ring-white"
              >
                {{ link.name }}
              </a>
            </div>
          </div>
          <div class="hidden md:flex md:items-center md:space-x-6">
            <RouterButton
              to="dashboard"
              block
              type="default"
              class="bg-gradient-to-r from-nitro-700/80 to-bravery-700/80 hover:from-nitro-600/80 hover:to-bravery-600/80"
            >
              {{ state.base.self ? "Go to Dashboard" : "Login" }}
            </RouterButton>
          </div>
        </nav>
      </div>

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
            <div class="flex items-center justify-between px-5 pt-4">
              <div>
                <a href="#" class="inline-flex items-center">
                  <img class="w-auto h-12" src="/img/mic.png" />
                  <h2 class="text-xl text-gradient bg-gradient-to-r from-nitro-500 to-dnd-400">
                    Spectrograph
                  </h2>
                </a>
              </div>
              <div class="-mr-2">
                <PopoverButton
                  class="inline-flex items-center justify-center p-2 text-white rounded-md bg-bravery-500 hover:bg-gray-100 focus:outline-none focus:ring-2 focus:ring-inset focus:ring-indigo-600"
                >
                  <span class="sr-only">Close menu</span>
                  <i-fas-xmark class="w-6 h-6" aria-hidden="true" />
                </PopoverButton>
              </div>
            </div>
            <div class="pt-5 pb-6">
              <div class="px-2 space-y-1">
                <a
                  v-for="link in headerLinks"
                  :key="link.name"
                  :href="link.href"
                  class="block px-3 py-2 text-base font-medium text-white rounded-md hover:bg-channel-600"
                >
                  {{ link.name }}
                </a>
              </div>
              <div class="px-5 mt-6">
                <RouterButton
                  to="dashboard"
                  block
                  type="default"
                  class="bg-gradient-to-r from-nitro-500 to-bravery-500"
                >
                  {{ state.base.self ? "Go to Dashboard" : "Login" }}
                </RouterButton>
              </div>
            </div>
          </div>
        </PopoverPanel>
      </transition>
    </Popover>

    <main>
      <div class="pt-10 lg:overflow-hidden lg:py-16">
        <div class="mx-auto max-w-7xl lg:px-8">
          <div class="lg:grid lg:grid-cols-2 lg:gap-8">
            <div
              class="max-w-md px-4 mx-auto sm:max-w-2xl sm:px-6 sm:text-center lg:flex lg:items-center lg:px-0 lg:text-left"
            >
              <div class="lg:py-24">
                <h1 class="mt-4 text-4xl font-bold tracking-tight text-white lg:mt-6 xl:text-6xl">
                  <span class="block">It's time to use</span>
                  <span class="block pb-1.5 text-gradient bg-gradient-to-r from-nitro-500 to-dnd-500">
                    Spectrograph
                  </span>
                  <span class="block">in your server.</span>
                </h1>
                <p class="mt-3 text-base text-gray-300 sm:mt-5 sm:text-xl lg:text-lg xl:text-xl">
                  Spectrograph is a Discord bot that manages creating/destroying voice channels depending
                  on occupancy, <span class="text-dnd-400">so you don't have to</span>.
                </p>
                <div class="flex flex-col flex-auto w-full gap-2 mt-10 lg:flex-row sm:mt-12">
                  <BaseButton href="#more-info" size="lg" type="info" block>
                    More information
                  </BaseButton>
                  <BaseButton href="/-/invite" size="lg" type="success" block transparent>
                    Add to your server
                  </BaseButton>
                </div>
              </div>
            </div>
            <div class="mt-12 -mb-16 sm:-mb-48 lg:relative lg:m-0">
              <div
                class="max-w-2xl px-4 mx-auto mb-0 sm:mb-16 lg:mb-0 sm:max-w-md md:px-6 lg:max-w-none lg:px-0"
              >
                <img
                  class="w-full lg:absolute lg:inset-y-0 lg:left-0 lg:h-full lg:w-auto lg:max-w-none"
                  src="/img/discord.png"
                />
              </div>
            </div>
          </div>
        </div>
      </div>

      <div id="more-info" class="relative pt-16 pb-32 overflow-hidden">
        <div class="relative max-w-xl px-4 mx-auto sm:px-6 lg:max-w-7xl lg:px-8">
          <div class="relative mt-12 sm:mt-16 lg:mt-24">
            <div class="lg:grid lg:grid-flow-row-dense lg:grid-cols-2 lg:items-center lg:gap-8">
              <div class="lg:col-start-2">
                <h3 class="text-2xl font-bold tracking-tight text-idle-500 sm:text-3xl">
                  Out with the old.
                </h3>
                <p class="mt-3 text-lg text-gray-500">
                  Larger servers have these types of channels setup. This can get annoying -- so many
                  channels that are empty. Why not be more efficient, without requiring mods/admins to
                  step in and create spare channels?
                </p>

                <dl class="mt-10 space-y-10">
                  <div class="relative">
                    <dt>
                      <div
                        class="absolute flex items-center justify-center w-12 h-12 text-white rounded-md bg-discord-500"
                      >
                        <i-fas-microphone class="w-6 h-6" aria-hidden="true" />
                      </div>
                      <p class="ml-16 text-lg font-medium leading-6 text-bravery-500">
                        Always-available voice channels
                      </p>
                    </dt>
                    <dd class="mt-2 ml-16 text-base text-gray-500">
                      Always have one empty channel with Spectrograph. If there are 2 (or 500) voice
                      channels, there is always an available one to join. No less, no more.
                    </dd>
                  </div>
                </dl>
              </div>

              <div class="relative mt-10 -mx-4 lg:col-start-1 lg:mt-0">
                <img class="relative mx-auto rounded-md" width="225" src="/img/example-bad.png" />
              </div>
            </div>

            <div class="relative mt-12 lg:mt-32 lg:grid lg:grid-cols-2 lg:items-center lg:gap-8">
              <div class="relative">
                <h3 class="text-2xl font-bold tracking-tight text-nitro-500 sm:text-3xl">
                  In with the new (and simple).
                </h3>
                <p class="mt-3 text-lg text-gray-500">
                  Spectrograph monitors voice channel member activity, as well as the
                  creation/deletion/editing of channels.
                </p>

                <dl class="mt-10 space-y-10">
                  <div class="relative">
                    <dt>
                      <div
                        class="absolute flex items-center justify-center w-12 h-12 text-white rounded-md bg-discord-500"
                      >
                        <i-fas-trash-can class="w-6 h-6" aria-hidden="true" />
                      </div>
                      <p class="ml-16 text-lg font-medium leading-6 text-bravery-500">
                        Cleanup of old channels
                      </p>
                    </dt>
                    <dd class="mt-2 ml-16 text-base text-gray-500">
                      Too many empty channels? No problem. Spectrograph will remove them, ensuring the
                      empty channel is always at the bottom.
                    </dd>
                  </div>

                  <div class="relative">
                    <dt>
                      <div
                        class="absolute flex items-center justify-center w-12 h-12 text-white rounded-md bg-discord-500"
                      >
                        <i-fas-user-plus class="w-6 h-6" aria-hidden="true" />
                      </div>
                      <p class="ml-16 text-lg font-medium leading-6 text-bravery-500">
                        Auto-scale channels
                      </p>
                    </dt>
                    <dd class="mt-2 ml-16 text-base text-gray-500">
                      All channels occupied? Spectrograph will create a new channel at the bottom of the
                      voice channel group, copying the configuration from the other channels in the
                      group.
                    </dd>
                  </div>

                  <div class="relative">
                    <dt>
                      <div
                        class="absolute flex items-center justify-center w-12 h-12 text-white rounded-md bg-discord-500"
                      >
                        <i-fas-gear class="w-6 h-6" aria-hidden="true" />
                      </div>
                      <p class="ml-16 text-lg font-medium leading-6 text-bravery-500">
                        Synchronized channel configuration
                      </p>
                    </dt>
                    <dd class="mt-2 ml-16 text-base text-gray-500">
                      When the first voice channel in the given group is changed (bitrate, user limit,
                      permissions, etc), Spectrograph will persist these changes across all other
                      channels of the voice channel group.
                    </dd>
                  </div>
                </dl>
              </div>

              <div class="relative mt-10 -mx-4 lg:mt-0" aria-hidden="true">
                <img class="relative mx-auto rounded-md" width="300" src="/img/example-new.png" />
              </div>
            </div>
          </div>
        </div>
      </div>

      <div class="bg-gradient-to-r from-nitro-800 to-bravery-700">
        <div
          class="max-w-4xl px-4 py-16 mx-auto sm:px-6 sm:pt-20 sm:pb-24 lg:max-w-7xl lg:px-8 lg:pt-24"
        >
          <h2 class="text-3xl font-bold tracking-tight text-white">
            Let Spectrograph do the work for you.
          </h2>
          <p class="max-w-3xl mt-4 text-lg text-gray-200">
            I'm lazy. You're lazy. We're all lazy. Let us do the work for you, so you can focus on the
            important things, like sending dank memes or having a heated debate about which anime is the
            best.
          </p>
          <div
            class="grid grid-cols-1 mt-12 gap-x-6 gap-y-12 sm:grid-cols-2 lg:mt-16 lg:grid-cols-4 lg:gap-x-8 lg:gap-y-16"
          >
            <div v-for="feature in features" :key="feature.name">
              <div>
                <span class="flex items-center justify-center w-12 h-12 rounded-md bg-white/10">
                  <component :is="feature.icon" class="w-6 h-6 text-white" aria-hidden="true" />
                </span>
              </div>
              <div class="mt-6">
                <h3 class="text-lg font-medium text-white">{{ feature.name }}</h3>
                <p class="mt-2 text-base text-gray-200">{{ feature.description }}</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </main>

    <footer>
      <div class="px-4 py-12 mx-auto overflow-hidden max-w-7xl sm:px-6 lg:px-8">
        <nav class="flex flex-wrap justify-center -mx-5 -my-2" aria-label="Footer">
          <div v-for="item in mainLinks" :key="item.name" class="px-5 py-2">
            <a :href="item.href" class="text-base text-gray-400 hover:text-gray-600">{{ item.name }}</a>
          </div>
        </nav>

        <div class="flex justify-center mt-8 space-x-6">
          <a
            v-for="item in socialLinks"
            :key="item.name"
            :href="item.href"
            class="text-gray-400 hover:text-gray-500"
          >
            <span class="sr-only">{{ item.name }}</span>
            <component :is="item.icon" class="w-6 h-6" aria-hidden="true" />
          </a>
        </div>
        <p class="mt-8 text-base text-center text-gray-400">
          <a href="https://liam.sh/" target="_blank">
            &copy; {{ new Date().getFullYear() }} Spectrograph
          </a>
        </p>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { headerLinks, mainLinks, socialLinks } from "@/lib/core/navigation"

const state = useState()

const features = [
  {
    name: "Auto-scaling magic",
    description: "Let us auto-scale your voice channels as your community scales.",
    icon: IconEmojiSparkles,
  },
  {
    name: "Easy as 1, 2, 3",
    description: "Actually, just 1, because we'll auto-scale to 2 and 3.",
    icon: IconEmojiHotBeverage,
  },
  {
    name: "Keep your channels clean",
    description: "Too many empty channels? No problem. We'll yeet them, we're good at that.",
    icon: IconEmojiBroom,
  },
  {
    name: "Synchronize all the things",
    description: "All channels are synced (perms, user limit, etc), because we know you'll forget.",
    icon: IconEmojiCounterclockwiseArrowsButton,
  },
  {
    name: "Perfect for game servers",
    description: "Limit channels to the size of a game party, so you can try-hard even harder.",
    icon: IconEmojiAlienMonster,
  },
  {
    name: "Customizable regex matching",
    description: "Use regex to match your weird channel names. We won't judge, but we will say RTFM.",
    icon: IconEmojiExplodingHead,
  },
  {
    name: "Few required permissions",
    description: "We're not the FBI or the chinese government. We don't need all those things.",
    icon: IconEmojiManGesturingNo,
  },
]
</script>
