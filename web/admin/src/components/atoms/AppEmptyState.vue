<script lang="ts" setup>
import { mdiInformationOutline } from '@mdi/js'

interface Props {
  message?: string
  icon?: string
}

const props = withDefaults(defineProps<Props>(), {
  message: 'データがありません',
  icon: mdiInformationOutline,
})
</script>

<template>
  <div class="app-empty-state d-flex flex-column align-center justify-center py-8">
    <v-icon
      :icon="props.icon"
      size="48"
      color="medium-emphasis"
      class="mb-4"
    />
    <div class="text-body-1 text-medium-emphasis mb-4">
      {{ props.message }}
    </div>
    <slot />
  </div>
</template>
