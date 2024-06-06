<script setup lang="ts">
interface Props {
  isShow: boolean
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
      class="mx-auto flex w-[90%] justify-between rounded-lg border-2 border-orange bg-white p-4 text-main shadow-md md:w-[400px]"
    >
      <div>{{ text }}</div>
      <button
        type="button"
        class="ml-4"
        @click="handelCloseButton"
      >
        <the-outline-close-icon />
      </button>
    </div>
  </Teleport>
</template>
