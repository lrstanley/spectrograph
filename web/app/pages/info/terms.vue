<script setup lang="ts">
import doc from "@/content/terms.md"

definePageMeta({ title: "Terms of service", layout: "landing" })
</script>

<template>
  <MarkdownDoc :doc="doc" />
</template>
