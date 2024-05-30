<script lang="ts" setup>
const props = defineProps({
  width: {
    type: Number,
    default: 140,
  },
})
</script>

<template>
  <img
    class="title"
    alt="ふるマル"
    src="/app-title.svg"
    :width="props.width"
  >
</template>

<style lang="scss" scoped>
.title {
  height: auto;
}
</style>
