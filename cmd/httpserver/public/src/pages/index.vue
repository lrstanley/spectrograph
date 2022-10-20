<route lang="yaml">
meta:
  title: "Home"
  layout: landing
</route>

<template>
  <div>
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
                <a class="h-10 min-h-0 btn btn-primary" href="#more-info">
                  <i-fas-circle-info class="h-8 mr-2" />
                  More information
                </a>
                <a class="h-10 min-h-0 btn btn-secondary" href="/-/invite">
                  <i-fas-circle-plus class="h-8 mr-2" />
                  Add to your server
                </a>
              </div>
            </div>
          </div>
          <div class="mt-12 -mb-16 sm:-mb-48 lg:relative lg:m-0">
            <div
              class="max-w-2xl px-4 mx-auto mb-0 sm:mb-16 lg:mb-0 sm:max-w-md md:px-6 lg:max-w-none lg:px-0"
            >
              <img
                loading="lazy"
                alt="Image of a Discord client"
                class="w-full transition-all duration-100 lg:absolute lg:inset-y-0 lg:left-0 lg:h-full lg:w-auto lg:max-w-none"
                width="1256"
                height="856"
                :src="imgDiscord"
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
                channels that are empty. Why not be more efficient, without requiring mods/admins to step
                in and create spare channels?
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
              <img
                loading="lazy"
                alt="Image of many discord channels"
                class="relative mx-auto rounded-md"
                width="225"
                :src="imgExampleOld"
              />
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
                    voice channel group, copying the configuration from the other channels in the group.
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
                    permissions, etc), Spectrograph will persist these changes across all other channels
                    of the voice channel group.
                  </dd>
                </div>
              </dl>
            </div>

            <div class="relative mt-10 -mx-4 lg:mt-0" aria-hidden="true">
              <img
                loading="lazy"
                alt="Image of Spectrograph dynamically creating channels"
                class="relative mx-auto rounded-md"
                width="300"
                :src="imgExampleNew"
              />
            </div>
          </div>
        </div>
      </div>
    </div>

    <div class="bg-gradient-to-r from-nitro-800 to-bravery-700">
      <div class="max-w-4xl px-4 py-16 mx-auto sm:px-6 sm:pt-20 sm:pb-24 lg:max-w-7xl lg:px-8 lg:pt-24">
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
  </div>
</template>

<script setup lang="ts">
import imgDiscord from "@/assets/img/discord.png?format=webp&imagetools"
import imgExampleOld from "@/assets/img/example-old.png?format=webp&imagetools"
import imgExampleNew from "@/assets/img/example-new.png?format=webp&imagetools"
const features = [
  {
    name: "Auto-scaling magic",
    description: "Let us auto-scale your voice channels as your community scales.",
    icon: IEmojiSparkles,
  },
  {
    name: "Easy as 1, 2, 3",
    description: "Actually, just 1, because we'll auto-scale to 2 and 3.",
    icon: IEmojiHotBeverage,
  },
  {
    name: "Keep your channels clean",
    description: "Too many empty channels? No problem. We'll yeet them, we're good at that.",
    icon: IEmojiBroom,
  },
  {
    name: "Synchronize all the things",
    description: "All channels are synced (perms, user limit, etc), because we know you'll forget.",
    icon: IEmojiCounterclockwiseArrowsButton,
  },
  {
    name: "Perfect for game servers",
    description: "Limit channels to the size of a game party, so you can try-hard even harder.",
    icon: IEmojiAlienMonster,
  },
  {
    name: "Customizable regex matching",
    description: "Use regex to match your weird channel names. We won't judge, but we will say RTFM.",
    icon: IEmojiExplodingHead,
  },
  {
    name: "Few required permissions",
    description: "We're not the FBI or the chinese government. We don't need all those things.",
    icon: IEmojiManGesturingNo,
  },
]
</script>
