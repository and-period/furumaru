<script lang="ts" setup>
interface Props {
  label: string
  color: string
}

defineProps<Props>()
</script>

<template>
  <v-chip
    :color="color"
    :aria-label="label"
  >
    {{ label }}
  </v-chip>
</template>
