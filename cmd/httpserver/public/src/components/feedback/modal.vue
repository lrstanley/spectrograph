<template>
  <Teleport to="body">
    <div
      class="transition-all duration-200 modal modal-bottom sm:modal-middle bg-channel-500/75"
      :class="props.open && 'modal-open'"
    >
      <div class="p-4 m-1 transition-all duration-200 modal-box md:m-0">
        <div class="flex flex-col items-start flex-auto w-full">
          <div v-if="!!$slots.header" class="flex flex-row items-center gap-3">
            <slot name="header" />
          </div>

          <div class="w-full text-left">
            <slot />

            <transition name="fade" appear out-in>
              <div v-if="props.error" class="w-full mt-2 alert alert-error">
                {{ props.error.message }}
              </div>
            </transition>
          </div>
        </div>

        <div class="modal-action">
          <slot name="actions" />
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
const props = defineProps<{
  open: boolean
  error?: Error
}>()
</script>
