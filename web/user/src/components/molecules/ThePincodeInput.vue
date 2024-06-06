<script lang="ts" setup>
interface Props {
  modelValue: string
  length?: number
}

const props = withDefaults(defineProps<Props>(), {
  length: 6,
})

const emits = defineEmits<{ (name: 'update:modelValue', val: string): void }>()

const divs = ref<any>([])
const values = ref<string[]>([])
const type = 'password'

onBeforeUpdate(() => {
  divs.value = []
})

onMounted(() => {
  values.value = [...Array(props.length)].map(() => '')
})

watch(
  values,
  () => {
    const val = values.value.map(item => item).join('')
    emits('update:modelValue', val)
  },
  { deep: true },
)

const handleKeyup = (event: KeyboardEvent, i: number) => {
  const key = event.key
  if (key === 'Backspace' && i !== 0) {
    divs.value[i].previousSibling.focus()
  }
  if (key.length === 1 && i !== props.length - 1) {
    divs.value[i].nextSibling.focus()
  }
}
</script>

<template>
  <div class="mt-5 flex flex-row justify-center px-2 text-center">
    <input
      v-for="item in [...Array(length)].map((_, i) => i)"
      :ref="
        (el) => {
          if (el) divs[item] = el
        }
      "
      :key="item"
      v-model="values[item]"
      class="m-2 h-12 w-10 rounded border border-main text-center focus:outline-current"
      :type="type"
      maxlength="1"
      @keyup="(e) => handleKeyup(e, item)"
    >
  </div>
</template>
