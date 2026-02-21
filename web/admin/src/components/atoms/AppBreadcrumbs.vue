<script lang="ts" setup>
import { mdiChevronRight } from '@mdi/js'

defineProps<{
  items: Array<{
    title: string
    to?: string
    disabled?: boolean
  }>
}>()
</script>

<template>
  <v-breadcrumbs
    :items="items"
    class="app-breadcrumbs px-0"
  >
    <template #divider>
      <v-icon
        :icon="mdiChevronRight"
        size="small"
      />
    </template>
  </v-breadcrumbs>
</template>

<style lang="scss" scoped>
.app-breadcrumbs {
  font-size: 0.875rem;
  color: var(--color-text-secondary);
  padding-block: 8px;
}
</style>
