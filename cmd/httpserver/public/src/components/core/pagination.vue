<template>
  <div class="btn-group btn-group-vertical lg:btn-group-horizontal">
    <button
      class="rounded-sm btn btn-sm btn-primary"
      :disabled="!page?.hasPreviousPage"
      @click="cursor = 'b.' + page?.startCursor"
    >
      Previous
    </button>
    <button
      class="rounded-sm btn btn-sm btn-primary"
      :disabled="!page?.hasNextPage"
      @click="cursor = 'a.' + page?.endCursor"
    >
      Next
    </button>
  </div>
</template>

<script setup lang="ts">
import type { PageInfo } from "@/lib/api"

const props = defineProps<{
  modelValue?: string
  pageInfo?: PageInfo
}>()

const emit = defineEmits(["update:modelValue"])

const cursor = useVModel(props, "modelValue", emit)
const page = computed(() => props.pageInfo)
</script>
