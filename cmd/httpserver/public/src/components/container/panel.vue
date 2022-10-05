<template>
  <div class="col-span-3 mt-5 md:col-span-2 md:mt-0">
    <div class="border rounded shadow sm:overflow-hidden border-chat-900">
      <div class="px-6 py-5 space-y-6 bg-chat-600">
        <div class="grid grid-cols-2 gap-6">
          <slot />
        </div>
      </div>
      <div class="flex flex-row items-center justify-between flex-1 px-6 py-3 bg-chat-700">
        <slot name="actions" />
      </div>
    </div>
  </div>
</template>
