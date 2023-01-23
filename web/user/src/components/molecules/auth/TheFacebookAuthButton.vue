<script lang="ts" setup>
const emits = defineEmits<{
  (e: 'click'): void
}>()

const handleClick = (): void => {
  emits('click')
}
</script>

<template>
  <the-button class="bg-facebook flex items-center" @click="handleClick">
    <p class="text-center w-full text-white">Log In with Facebook</p>
  </the-button>
</template>
