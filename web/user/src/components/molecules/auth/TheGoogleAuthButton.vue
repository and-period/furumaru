<script lang="ts" setup>
const emits = defineEmits<{
  (e: 'click'): void
}>()

const handleClick = (): void => {
  emits('click')
}
</script>

<template>
  <the-button class="border border-gray-400 flex items-center" @click="handleClick">
    <img src="~/assets/img/g-normal.png" height="18" width="18" alt="google icon" class="block" />
    <p class="text-center w-full">Log In with Google</p>
  </the-button>
</template>
