<script lang="ts" setup>
interface Props {
  modelValue: string | number
  placeholder: string
  type: string
  label?: string
  required?: boolean
  withLabel?: boolean
  message?: string
  error?: boolean
  errorMessage?: string
  name?: string
  id?: string
  pattern?: string
}

const props = withDefaults(defineProps<Props>(), {
  required: false,
  withLabel: true,
  error: false,
  message: undefined,
  errorMessage: '',
})
const emits = defineEmits<{
  (e: 'update:modelValue', val: string | number): void
}>()

const value = computed({
  get: () => props.modelValue,
  set: (val: string | number) => emits('update:modelValue', val),
})

/**
 * エラーの判定
 * errorMessageが渡されている場合はエラー状態にする
 */
const hasError = computed(() => {
  if (props.error) {
    return true
  }
  if (props.errorMessage !== '') {
    return true
  }
  return false
})

/**
 * メッセージエリアに表示する文字列
 * errorMessageを優先する
 */
const viewMessage = computed(() => {
  if (props.errorMessage !== '') {
    return props.errorMessage
  } else {
    return props.message
  }
})
</script>

<template>
  <div class="mb-1">
    <div class="w-full">
      <label v-if="withLabel" class="inline-block" :for="id">{{ label }}</label>
      <input
        :id="id"
        v-model="value"
        :name="name"
        :placeholder="placeholder"
        :required="required"
        :type="type"
        :pattern="pattern"
        :class="{
          'block w-full appearance-none border-b border-main bg-transparent px-2 leading-10 outline-none ring-0 focus:outline-none ': true,
          'border-b-2 border-orange': hasError,
        }"
      />
    </div>
    <p :class="{ 'text-orange': hasError, 'text-left text-sm': true }">
      {{ viewMessage }}
    </p>
  </div>
</template>
