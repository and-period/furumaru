<script setup lang="ts">
interface Props {
  boxSize: number
  useRate: number
}

defineProps<Props>()
</script>

<template>
  <div>
    <the-60-cardboard v-if="boxSize === 60" :use-rate="useRate" />
    <the-80-cardboard v-if="boxSize === 80" :use-rate="useRate" />
    <the-100-cardboard v-if="boxSize === 100" :use-rate="useRate" />
  </div>
</template>
