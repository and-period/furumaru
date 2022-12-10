<template>
  <img
    class="logo"
    alt="ふるマル ロゴ"
    src="/app-logo-with-title.svg"
    :width="width"
  />
</template>

<script>
import { defineComponent } from '@vue/composition-api'

export default defineComponent({
  props: {
    width: {
      type: Number,
      default: 200,
    },
  },
})
</script>

<style scoped lang="scss">
.logo {
  height: auto;
  display: block;
}
</style>
