<script setup lang="ts">
import { Volume2 } from "@lucide/vue"
import { cn } from "@/lib/utils"

// A small simulation of what Spectrograph does to a channel group: people join
// and leave, a channel is cloned when every channel is taken, and spare empty
// channels are removed.

type Channel = { id: number; members: number[]; fresh: boolean }

const limit = 4
const maxChannels = 5
const reducedMotion = usePreferredReducedMotion()

let nextChannel = 1
let nextMember = 1

function channel(members: number[]): Channel {
  return { id: nextChannel++, members, fresh: false }
}

const channels = ref<Channel[]>([channel([nextMember++, nextMember++]), channel([nextMember++]), channel([])])
const status = ref("Watching 3 channels")

function reconcile() {
  const empty = channels.value.filter((c) => c.members.length === 0)
  if (empty.length === 0 && channels.value.length < maxChannels) {
    channels.value.push({ ...channel([]), fresh: true })
    status.value = "Every channel is in use, so one was cloned"
    return
  }
  if (empty.length > 1) {
    // Keep the last empty channel, remove the other spares.
    const keep = empty[empty.length - 1]
    channels.value = channels.value.filter((c) => c.members.length > 0 || c === keep || c === channels.value[0])
    status.value = "Removed a spare empty channel"
  }
}

function tick() {
  for (const c of channels.value) {
    c.fresh = false
  }
  const total = channels.value.reduce((sum, c) => sum + c.members.length, 0)
  const open = channels.value.filter((c) => c.members.length < limit)
  const join = total < 3 || (total < 12 && Math.random() < 0.62)

  if (join && open.length) {
    // People prefer channels that already have someone in them.
    const target = open.find((c) => c.members.length > 0) ?? open[0]
    if (target) {
      target.members.push(nextMember++)
      status.value = "Someone joined a channel"
    }
  } else {
    const occupied = channels.value.filter((c) => c.members.length > 0)
    const target = occupied[Math.floor(Math.random() * occupied.length)]
    if (target) {
      target.members.splice(Math.floor(Math.random() * target.members.length), 1)
      status.value = "Someone left a channel"
    }
  }
  reconcile()
}

const { pause, resume } = useIntervalFn(tick, 1600, { immediate: false })

watch(
  reducedMotion,
  (value) => {
    if (value === "reduce") pause()
    else resume()
  },
  { immediate: true },
)

const tones = ["bg-magma-1", "bg-magma-2", "bg-magma-3", "bg-magma-4", "bg-magma-5"]
</script>

<template>
  <figure class="relative">
    <div class="border bg-card shadow-[0_30px_80px_-30px] shadow-magma-2/50">
      <div class="flex items-center justify-between gap-3 border-b px-4 py-3">
        <span class="text-sm font-semibold">Voice channels</span>
        <code class="font-mono text-[11px] text-muted-foreground" title="Channels matching this pattern are managed">^.* \+$</code>
      </div>
      <TransitionGroup tag="ul" name="channel" class="relative flex flex-col gap-1 p-2" aria-hidden="true">
        <li
          v-for="(c, index) in channels"
          :key="c.id"
          :class="cn('flex min-h-11 items-center gap-3 px-3 py-2 transition-colors duration-700', c.fresh ? 'bg-primary/15' : 'bg-transparent')"
        >
          <Volume2 class="size-4 shrink-0 text-muted-foreground" />
          <span class="min-w-0 flex-1 truncate text-sm">
            Squad{{ index > 0 ? ` ${index + 1}` : "" }} +
          </span>
          <span class="flex items-center gap-1">
            <span
              v-for="m in c.members"
              :key="m"
              :class="cn('size-5 rounded-full ring-2 ring-card', tones[m % tones.length])"
            />
          </span>
          <span class="w-8 text-right font-mono text-[11px] tabular-nums text-muted-foreground">
            {{ c.members.length }}/{{ limit }}
          </span>
        </li>
      </TransitionGroup>
      <div class="border-t px-4 py-2.5 text-xs text-muted-foreground" aria-live="off">
        {{ status }}
      </div>
    </div>
    <figcaption class="sr-only">
      An animated example of a voice channel group. When every channel is in use, Spectrograph clones another one, and
      it removes spare empty channels when people leave.
    </figcaption>
  </figure>
</template>

<style scoped>
.channel-move,
.channel-enter-active,
.channel-leave-active {
  transition: all 0.45s cubic-bezier(0.2, 0.7, 0.2, 1);
}
.channel-enter-from {
  opacity: 0;
  transform: translateY(-6px);
}
.channel-leave-to {
  opacity: 0;
  transform: translateX(12px);
}
@media (prefers-reduced-motion: reduce) {
  .channel-move,
  .channel-enter-active,
  .channel-leave-active {
    transition: none;
  }
}
.channel-leave-active {
  position: absolute;
  inset-inline: 0.5rem;
}
</style>
