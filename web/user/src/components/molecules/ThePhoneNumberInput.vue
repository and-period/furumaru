<script setup lang="ts">
interface Props {
  modelValue: string
  required: boolean
}

defineProps<Props>()

interface Emits {
  (e: 'update:modelValue', val: string): void
}

const emits = defineEmits<Emits>()

const tel1 = ref<string>('')
const tel2 = ref<string>('')
const tel3 = ref<string>('')

watch([tel1, tel2, tel3], () => {
  emits('update:modelValue', `${tel1.value}-${tel2.value}-${tel3.value}`)
})
</script>

<template>
  <div>
    <label class="text-[12px] tracking-[1.2px] text-main">電話番号</label>
    <div class="grid w-full grid-cols-11 items-center gap-2">
      <input
        v-model="tel1"
        type="tel"
        size="3"
        placeholder="000"
        class="col-span-3 block appearance-none rounded-none border-b border-main bg-transparent px-2 leading-10 outline-none ring-0 focus:outline-none"
        required
      />
      -
      <input
        v-model="tel2"
        type="tel"
        size="4"
        placeholder="0000"
        class="col-span-3 appearance-none rounded-none border-b border-main bg-transparent px-2 leading-10 outline-none ring-0 focus:outline-none"
        required
      />
      -
      <input
        v-model="tel3"
        type="tel"
        size="4"
        placeholder="0000"
        class="col-span-3 appearance-none rounded-none border-b border-main bg-transparent px-2 leading-10 outline-none ring-0 focus:outline-none"
        required
      />
    </div>
  </div>
</template>
