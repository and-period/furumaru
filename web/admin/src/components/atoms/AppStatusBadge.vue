<script lang="ts" setup>
interface Props {
  label: string
  color: string
}

defineProps<Props>()
</script>

<template>
  <v-chip :color="color">
    {{ label }}
  </v-chip>
</template>
