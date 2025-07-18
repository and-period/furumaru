<script lang="ts" setup>
const props = defineProps({
  width: {
    type: Number,
    default: 200,
  },
})

const logoSrc = '/app-logo-with-title.svg'
</script>

<template>
  <img
    class="logo"
    alt="ふるマル ロゴ"
    :src="logoSrc"
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
