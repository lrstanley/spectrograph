<route>
meta:
  title: 'An error occurred'
  layout: error
</route>

<template>
  <div class="text-center">
    <p class="text-lg font-semibold text-dnd-500">{{ code }}</p>
    <p class="mt-2 text-lg text-dnd-400">Error: {{ error || code }}</p>
  </div>
</template>

<script setup lang="ts">
const route = useRoute("/error")

const code = computed(() => {
  const err = route.query.code

  if (err === "404") {
    return "404: Page not found"
  } else if (err === "401") {
    return "401: Unauthorized"
  } else if (err === "403") {
    return "403: Forbidden"
  } else if (err === "500") {
    return "500: Internal server error"
  } else {
    return err || "Unknown error"
  }
})

const error = computed(() => {
  const err = route.query.e
  const code = route.query.code

  if (err) {
    return err
  } else if (code === "404") {
    return "The page you are looking for could not be found."
  } else if (code === "401") {
    return "Please login before viewing this page."
  } else if (code === "403") {
    return "You are not allowed to view this page."
  } else if (code === "500") {
    return "An internal server error occurred, please try again later."
  } else {
    return "An unknown error occurred, please try again later."
  }
})
</script>
