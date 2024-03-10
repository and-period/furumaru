<script setup lang="ts">
interface Props {
  isAuthenticated: boolean
}

interface Emits {
  (e: 'submit', value: string): void
}

defineProps<Props>()

const emits = defineEmits<Emits>()

const modelValue = defineModel<string>('modelValue', {
  default: '',
  required: true,
})

const canSubmit = computed<boolean>(() => modelValue.value.length > 0)

const handleSubmit = () => {
  emits('submit', modelValue.value)
}
</script>

<template>
  <form class="flex w-full flex-col" @submit.prevent="handleSubmit">
    <div class="flex w-full items-center gap-4">
      <input
        v-model="modelValue"
        type="text"
        class="block w-full border-typography p-2 shadow-[0_1px_0_0] focus:border-0 focus:border-main focus:shadow-[0_2px_0_0] focus:outline-none focus:ring-0"
        placeholder="コメントする…"
      />
      <button
        class="whitespace-nowrap rounded-lg bg-main px-4 py-2 text-white disabled:bg-main/50"
        type="submit"
        :disabled="!canSubmit"
      >
        送信
      </button>
    </div>
    <span class="mt-1 text-[12px] tracking-[10%]">
      <template v-if="!isAuthenticated">※ゲストとしてコメントします </template>
    </span>
  </form>
</template>
