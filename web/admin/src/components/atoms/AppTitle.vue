<script lang="ts" setup>
const props = defineProps({
  width: {
    type: Number,
    default: 140,
  },
})

const titleSrc = '/app-title.svg'
</script>

<template>
  <img
    class="title"
    alt="ふるマル"
    :src="titleSrc"
    :width="props.width"
  >
</template>

<style lang="scss" scoped>
.title {
  height: auto;
}
</style>
