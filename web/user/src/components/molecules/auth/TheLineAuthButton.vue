<script lang="ts" setup>
const emits = defineEmits<{
  (e: 'click'): void
}>()

const handleClick = (): void => {
  emits('click')
}
</script>

<template>
  <the-button class="bg-line flex items-center" @click="handleClick">
    <p class="text-center w-full text-white">Log In with LINE</p>
  </the-button>
</template>
