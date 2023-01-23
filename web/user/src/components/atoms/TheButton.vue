<script lang="ts" setup>
const emits = defineEmits<{
  (e: 'click'): void
}>()

const handleClick = (): void => {
  emits('click')
}
</script>

<template>
  <button class="p-2" @click="handleClick">
    <slot />
  </button>
</template>
