<script lang="ts" setup>
interface Props {
  buttonText: string
}
defineProps<Props>()

const emits = defineEmits<{(e: 'click'): void}>()

const handleClick = (): void => {
  emits('click')
}
</script>

<template>
  <the-button class="bg-facebook flex items-center" @click="handleClick">
    <the-facebook-icon class="w-6" />
    <p class="text-center w-full text-white">
      {{ buttonText }}
    </p>
  </the-button>
</template>
