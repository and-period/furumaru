<script setup lang="ts">
import type { Product } from '~/types/api'

interface Props {
  item: Product
}

const props = defineProps<Props>()

const sanitizeDescription = computed<string>(() => {
  // Remove HTML tags from the description
  return props.item.description.replace(/<("[^"]*"|'[^']*'|[^'">])*>/g, '')
})
</script>

<template>
  <div class="d-flex flex-column ga-3">
    <p class="text-subtitle-1">
      {{ item.name }}
    </p>
    <div class="description text-body-2">
      {{ sanitizeDescription }}
    </div>
  </div>
</template>

<style scoped>
.description {
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>
