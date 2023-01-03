<template>
  <div class="stack" :style="spaceStyle">
    <slot />
  </div>
</template>

<script lang="ts" setup>
const props = defineProps({
  space: {
    type: Number,
    default: 16,
  },
})

const spaceStyle = computed(() => {
  return {
    '--space': `${props.space}px`,
  }
})
</script>

<style lang="scss" scoped>
.stack {
  display: flex;
  flex-direction: column;
  justify-content: flex-start;
}

.stack > * {
  margin-top: 0;
  margin-bottom: 0;
}

.stack > * + * {
  margin-top: var(--space, 16px);
}
</style>
