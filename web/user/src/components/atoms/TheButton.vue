<script lang="ts" setup>
const emits = defineEmits<{ (e: 'click'): void }>()

const handleClick = (): void => {
  emits('click')
}
</script>

<template>
  <button
    class="p-3"
    @click="handleClick"
  >
    <slot />
  </button>
</template>
