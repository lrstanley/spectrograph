<template>
  <div class="col-span-3 mt-5 md:col-span-2 md:mt-0">
    <div class="border rounded shadow sm:overflow-hidden border-chat-800">
      <div
        v-if="!!$slots.header"
        class="flex flex-row items-center justify-between flex-1 px-6 py-3 bg-chat-900"
      >
        <slot name="header" />
      </div>

      <div class="bg-channel-600">
        <slot />
      </div>

      <div
        v-if="!!$slots.actions"
        class="flex flex-row items-center justify-between flex-1 px-6 py-3 bg-chat-800"
      >
        <slot name="actions" />
      </div>
    </div>
  </div>
</template>
