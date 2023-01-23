<script lang="ts" setup>
interface Props {
  modelValue: string | number
  label: string
  placeholder: string
  type: string
  required?: boolean
  withLabel?: boolean
  message?: string
  error?: boolean
  errorMessage?: string
}

const props = withDefaults(defineProps<Props>(), {
  required: false,
  withLabel: true,
  error: false,
  errorMessage: '',
})
const emits = defineEmits<{ (e: 'update:modelValue', val: string | number): void }>()

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
      <label v-if="withLabel" class="form-label inline-block">{{ label }}</label>
      <input
        v-model="value"
        :placeholder="placeholder"
        :required="required"
        :type="type"
        :class="{
          'form-control block w-full px-2 bg-transparent border-b border-text-main focus:outline-none leading-10': true,
          'border-b-2 border-orange': hasError,
        }"
      />
    </div>
    <p :class="{ 'text-orange': hasError, 'text-left text-sm': true }">
      {{ viewMessage }}
    </p>
  </div>
</template>
