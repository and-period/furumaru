<template>
  <v-btn
    class="rounded-xl"
    color="accent"
    block
    type="submit"
    elevation="0"
    :height="props.isMobile ? 40 : 64"
    :size="!props.isMobile && 'x-large'"
  >
    <slot />
  </v-btn>
</template>

<script lang="ts" setup>
const props = defineProps({
  isMobile: {
    type: Boolean,
    default: false,
  },
})
</script>
