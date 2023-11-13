<script setup lang="ts">
interface Props {
  isShow: Boolean
  text: string
}

defineProps<Props>()

interface Emits {
  (e: 'update:isShow', val: boolean): void
}

const emits = defineEmits<Emits>()

const handelCloseButton = () => {
  emits('update:isShow', false)
}

onMounted(() => {
  setTimeout(() => emits('update:isShow', false), 3000)
})
</script>

<template>
  <Teleport to="#snackbar">
    <div
      v-if="isShow"
      class="mx-auto flex w-[400px] justify-between rounded-lg border-2 border-orange bg-white p-4 text-main shadow-md"
    >
      <div>{{ text }}</div>
      <button type="button" @click="handelCloseButton">
        <the-outline-close-icon />
      </button>
    </div>
  </Teleport>
</template>
