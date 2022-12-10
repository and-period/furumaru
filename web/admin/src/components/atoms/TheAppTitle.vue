<template>
  <img class="title" alt="ふるマル" src="/app-title.svg" :width="width" />
</template>

<script>
import { defineComponent } from '@vue/composition-api'

export default defineComponent({
  props: {
    width: {
      type: Number,
      default: 140,
    },
  },
})
</script>

<style scoped lang="scss">
.title {
  height: auto;
}
</style>
