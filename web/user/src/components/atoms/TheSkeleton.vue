<script setup lang="ts">
interface Props {
  width?: string
  height?: string
  rounded?: boolean
  circle?: boolean
}

withDefaults(defineProps<Props>(), {
  width: '100%',
  height: '1rem',
  rounded: true,
  circle: false,
})
</script>

<template>
  <div
    aria-hidden="true"
    class="animate-pulse bg-gray-200"
    :class="{
      'rounded-md': rounded && !circle,
      'rounded-full': circle,
    }"
    :style="{ width, height }"
  />
</template>
