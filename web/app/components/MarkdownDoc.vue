<script setup lang="ts">
const props = defineProps<{
  doc: { html: string; headings: { id: string; text: string; depth: number }[] }
}>()

const toc = computed(() => props.doc.headings.filter((h) => h.depth === 2))

// Rendered markdown uses plain anchors; route same-site links through the
// router instead of reloading the page.
function onClick(event: MouseEvent) {
  const anchor = (event.target as HTMLElement).closest("a")
  const href = anchor?.getAttribute("href")
  if (!anchor || !href || anchor.target || event.metaKey || event.ctrlKey || event.shiftKey) {
    return
  }
  if (href.startsWith("/")) {
    event.preventDefault()
    navigateTo(href)
  }
}
</script>

<template>
  <div class="mx-auto grid max-w-6xl gap-10 px-4 py-12 sm:px-6 lg:grid-cols-[minmax(0,1fr)_14rem]">
    <!-- Trusted HTML, rendered from the repository's markdown at build time. Clicks are delegated to its links. -->
    <!-- biome-ignore lint/a11y/useKeyWithClickEvents: links inside handle keyboard activation natively. -->
    <article class="prose min-w-0" @click="onClick" v-html="doc.html" />
    <aside class="hidden lg:block">
      <nav aria-label="On this page" class="sticky top-20 flex flex-col gap-1 border-l pl-4 text-sm">
        <p class="mb-2 font-semibold">On this page</p>
        <a v-for="h in toc" :key="h.id" :href="`#${h.id}`" class="py-0.5 text-muted-foreground hover:text-foreground">
          {{ h.text }}
        </a>
      </nav>
    </aside>
  </div>
</template>
