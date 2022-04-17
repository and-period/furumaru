<template>
  <v-btn
    class="rounded-xl"
    color="accent"
    block
    type="submit"
    elevation="0"
    :height="isMobile ? 40 : 64"
    :x-large="!isMobile"
  >
    <slot />
  </v-btn>
</template>

<script lang="ts">
import { defineComponent } from '@nuxtjs/composition-api'

export default defineComponent({
  props: {
    isMobile: {
      type: Boolean,
      required: true,
    },
  },
})
</script>
