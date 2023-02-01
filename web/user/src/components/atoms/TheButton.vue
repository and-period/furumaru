<script lang="ts" setup>
const emits = defineEmits<{(e: 'click'): void}>()

const handleClick = (): void => {
  emits('click')
}
</script>

<template>
  <button class="px-2 py-3" @click="handleClick">
    <slot />
  </button>
</template>
