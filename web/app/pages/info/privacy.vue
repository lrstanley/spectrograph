<script setup lang="ts">
import doc from "@/content/privacy.md"

definePageMeta({ title: "Privacy policy", layout: "landing" })
</script>

<template>
  <MarkdownDoc :doc="doc" />
</template>
