<script lang="ts" setup>
const props = defineProps({
  width: {
    type: Number,
    default: 200,
  },
})
</script>

<style lang="scss" scoped>
.logo {
  height: auto;
  display: block;
}
</style>

<template>
  <img
    class="logo"
    alt="ふるマル ロゴ"
    src="/app-logo-with-title.svg"
    :width="props.width"
  />
</template>
