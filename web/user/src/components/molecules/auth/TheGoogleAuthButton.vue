<script lang="ts" setup>
interface Props {
  buttonText: string
}
defineProps<Props>()

const emits = defineEmits<{ (e: 'click'): void }>()

const handleClick = (): void => {
  emits('click')
}
</script>

<template>
  <the-button
    class="flex items-center border border-gray-400 text-typography"
    @click="handleClick"
  >
    <img
      src="~/assets/img/g-normal.png"
      height="24"
      width="24"
      alt="google icon"
      class="block"
    >
    <p class="w-full text-center">
      {{ buttonText }}
    </p>
  </the-button>
</template>
