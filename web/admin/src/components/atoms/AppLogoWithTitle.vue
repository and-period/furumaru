<script lang="ts" setup>
const props = defineProps({
  width: {
    type: Number,
    default: 200,
  },
})
</script>

<template>
  <img
    class="logo"
    alt="ふるマル ロゴ"
    src="/app-logo-with-title.svg"
    :width="props.width"
  >
</template>

<style lang="scss" scoped>
.logo {
  height: auto;
  display: block;
  max-width: 80%;
}
</style>
