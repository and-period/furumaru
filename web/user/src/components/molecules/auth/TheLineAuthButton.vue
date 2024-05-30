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
    class="flex items-center bg-line"
    @click="handleClick"
  >
    <the-line-icon />
    <p class="w-full text-center text-white">
      {{ buttonText }}
    </p>
  </the-button>
</template>
