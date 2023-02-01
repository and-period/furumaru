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
  <the-button class="bg-line flex items-center" @click="handleClick">
    <p class="text-center w-full text-white">
      {{ buttonText }}
    </p>
  </the-button>
</template>
