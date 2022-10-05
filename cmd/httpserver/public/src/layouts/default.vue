<template>
  <!-- <div class="flex flex-col flex-auto">
    <div class="flex flex-col flex-auto sm:container sm:mx-auto">
      <template v-if="props.loading"> Loading... </template>
      <FeedbackAlert
        v-else-if="props.error"
        v-motion-fade
        type="error"
        class="m-2 md:m-6"
      >
        An error occurred: {{ props.error }}
      </FeedbackAlert>
      <template v-else>
        <slot><router-view /></slot>
      </template>
    </div>
  </div> -->
  <router-view />
</template>
