<script setup lang="ts">
import {
  Brush,
  Coffee,
  ExternalLink,
  Gamepad2,
  Hand,
  Layers,
  Mic,
  RefreshCcw,
  Regex,
  Sparkles,
  SquareTerminal,
  Trash2,
  UserPlus,
} from "@lucide/vue"
import imgNew from "@/assets/img/example-new-optimized.gif"
import imgOld from "@/assets/img/example-old.png"

definePageMeta({ title: "Spectrograph", layout: "landing" })

const { self } = useSelf()

const capabilities = [
  {
    icon: Mic,
    title: "Always an open channel",
    body: "Whether a group has 2 voice channels or 500, there's always exactly one empty channel to join. No fewer, no more.",
  },
  {
    icon: UserPlus,
    title: "Scales as people join",
    body: "When every channel in a group is taken, Spectrograph adds one at the bottom, copying the settings of the others.",
  },
  {
    icon: Trash2,
    title: "Cleans up after itself",
    body: "Too many empty channels? Spectrograph removes the spares, keeping the empty one at the bottom.",
  },
  {
    icon: RefreshCcw,
    title: "Keeps settings in sync",
    body: "Change the bitrate or user limit of a group's original channel and its clones follow. New clones copy its permissions too.",
  },
]

const extras = [
  { icon: Sparkles, title: "Grows with your community", body: "Voice channels scale up as your server does." },
  { icon: Coffee, title: "Easy as 1, 2, 3", body: "Actually, just 1. We'll scale to 2 and 3." },
  { icon: Brush, title: "No clutter", body: "Too many empty channels? We'll yeet them. We're good at that." },
  {
    icon: Layers,
    title: "Settings stay in sync",
    body: "Bitrate and user limits follow the original channel, because we know you'll forget.",
  },
  { icon: Gamepad2, title: "Made for game servers", body: "Cap channels at your party size, so you can try-hard even harder." },
  { icon: Regex, title: "Your naming, your rules", body: "Match your channel names with a regular expression. We won't judge." },
  {
    icon: SquareTerminal,
    title: "Configure it from Discord",
    body: "Run /spectrograph to check its status, edit settings, or pause it, without leaving your server.",
  },
  {
    icon: Hand,
    title: "Few permissions",
    body: "Channel access, plus the voice permissions clones copy. No admin, no Manage Roles.",
  },
]
</script>

<template>
  <section class="relative overflow-hidden border-b">
    <div
      aria-hidden="true"
      class="pointer-events-none absolute inset-x-0 bottom-0 h-1 bg-linear-to-r from-magma-1 via-magma-3 to-magma-5"
    />
    <div class="mx-auto grid max-w-6xl items-center gap-12 px-4 py-16 sm:px-6 lg:grid-cols-[1.15fr_1fr] lg:py-24">
      <div class="flex flex-col gap-6">
        <h1 class="font-heading text-5xl leading-[0.95] font-bold tracking-[-0.035em] text-balance sm:text-6xl lg:text-7xl">
          Voice channels that make room for everyone.
        </h1>
        <p class="max-w-[52ch] text-lg text-muted-foreground">
          Spectrograph is a Discord bot that clones a voice channel when every one in a group is taken, and removes the
          spares when people leave. There's always one empty channel to join, without a moderator stepping in.
        </p>
        <div class="flex flex-wrap gap-3">
          <Button as-child size="lg">
            <a :href="inviteURL()" target="_blank" rel="noopener">
              Add to Discord
              <ExternalLink data-icon="inline-end" />
            </a>
          </Button>
          <Button as-child size="lg" variant="outline">
            <NuxtLink :to="self ? '/dashboard' : '/login'">{{ self ? "Open dashboard" : "Log in to configure" }}</NuxtLink>
          </Button>
        </div>
        <p class="text-sm text-muted-foreground">Free to use. Works with any voice channel whose name ends in " +".</p>
      </div>
      <ChannelDemo class="w-full max-w-md justify-self-center lg:justify-self-end" />
    </div>
  </section>

  <section id="how-it-works" class="scroll-mt-16 border-b">
    <div class="mx-auto grid max-w-6xl gap-16 px-4 py-20 sm:px-6">
      <div class="grid items-center gap-10 md:grid-cols-[auto_1fr]">
        <figure class="mx-auto flex flex-col items-center gap-3">
          <img
            :src="imgOld"
            alt="A Discord channel list with ten empty voice channels, Free Roam 6 to 15"
            width="179"
            height="355"
            loading="lazy"
            class="border opacity-80 grayscale"
          />
          <figcaption class="text-xs text-muted-foreground">Ten channels, nobody in them</figcaption>
        </figure>
        <div class="flex max-w-xl flex-col gap-3">
          <h2 class="font-heading text-3xl font-semibold tracking-tight">Out with the old</h2>
          <p class="text-muted-foreground">
            Larger servers create a pile of numbered voice channels up front, just in case. Most sit empty, and when
            they all fill up, someone has to ask a moderator for another one.
          </p>
        </div>
      </div>

      <div class="grid items-center gap-10 md:grid-cols-[1fr_auto]">
        <div class="flex flex-col gap-6">
          <div class="flex max-w-xl flex-col gap-3">
            <h2 class="font-heading text-3xl font-semibold tracking-tight">In with the new, and simple</h2>
            <p class="text-muted-foreground">
              Spectrograph watches who joins and leaves voice channels, and how channels are created, deleted and
              edited, then adjusts the group to fit.
            </p>
          </div>
          <dl class="grid gap-6 sm:grid-cols-2">
            <div v-for="item in capabilities" :key="item.title" class="flex flex-col gap-2">
              <dt class="flex items-center gap-2 font-semibold">
                <component :is="item.icon" class="size-4 text-primary" aria-hidden="true" />
                {{ item.title }}
              </dt>
              <dd class="text-sm text-muted-foreground">{{ item.body }}</dd>
            </div>
          </dl>
        </div>
        <figure class="mx-auto flex flex-col items-center gap-3 md:order-last">
          <img
            :src="imgNew"
            alt="Recording of Spectrograph adding a voice channel as people join"
            width="238"
            height="267"
            loading="lazy"
            class="border"
          />
          <figcaption class="text-xs text-muted-foreground">One spare channel, added as needed</figcaption>
        </figure>
      </div>
    </div>
  </section>

  <section class="border-b bg-card">
    <div class="mx-auto flex max-w-6xl flex-col gap-10 px-4 py-20 sm:px-6">
      <div class="flex max-w-2xl flex-col gap-3">
        <h2 class="font-heading text-3xl font-semibold tracking-tight">Let Spectrograph do the work for you</h2>
        <p class="text-muted-foreground">
          I'm lazy. You're lazy. We're all lazy. Let the bot handle channels, so you can focus on the important things,
          like sending memes or arguing about which anime is best.
        </p>
      </div>
      <ul class="grid gap-x-10 gap-y-8 sm:grid-cols-2 lg:grid-cols-3">
        <li v-for="item in extras" :key="item.title" class="flex gap-3">
          <component :is="item.icon" class="mt-0.5 size-5 shrink-0 text-magma-4" aria-hidden="true" />
          <div class="flex flex-col gap-1">
            <h3 class="font-semibold">{{ item.title }}</h3>
            <p class="text-sm text-muted-foreground">{{ item.body }}</p>
          </div>
        </li>
      </ul>
    </div>
  </section>

  <section id="setup" class="scroll-mt-16">
    <div class="mx-auto flex max-w-6xl flex-col gap-10 px-4 py-20 sm:px-6">
      <h2 class="font-heading text-3xl font-semibold tracking-tight">Set up in a minute</h2>
      <ol class="grid gap-8 md:grid-cols-3">
        <li class="flex flex-col gap-2 border-t-2 border-magma-2 pt-4">
          <span class="font-heading text-4xl font-bold text-magma-2" aria-hidden="true">1</span>
          <h3 class="font-semibold">Add the bot</h3>
          <p class="text-sm text-muted-foreground">
            Invite Spectrograph to your server. It asks to view and manage channels, plus the voice permissions
            (connect, speak, video, voice activity) it copies onto cloned channels.
          </p>
        </li>
        <li class="flex flex-col gap-2 border-t-2 border-magma-3 pt-4">
          <span class="font-heading text-4xl font-bold text-magma-3" aria-hidden="true">2</span>
          <h3 class="font-semibold">Name a channel</h3>
          <p class="text-sm text-muted-foreground">
            End a voice channel's name with a space and a plus, like <code class="font-mono text-foreground">Squad +</code>.
            Prefer your own naming? Set a pattern in the dashboard, or with
            <code class="font-mono text-foreground">/spectrograph settings</code> in Discord.
          </p>
        </li>
        <li class="flex flex-col gap-2 border-t-2 border-magma-5 pt-4">
          <span class="font-heading text-4xl font-bold text-magma-5" aria-hidden="true">3</span>
          <h3 class="font-semibold">Join it</h3>
          <p class="text-sm text-muted-foreground">
            When someone joins, an empty copy appears below. When they leave, the spare goes away.
          </p>
        </li>
      </ol>
      <div class="flex flex-wrap gap-3">
        <Button as-child size="lg">
          <a :href="inviteURL()" target="_blank" rel="noopener">
            Add to Discord
            <ExternalLink data-icon="inline-end" />
          </a>
        </Button>
      </div>
    </div>
  </section>
</template>
