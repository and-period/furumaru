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
  maxLength?: number
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
 * パスワードの表示切り替え
 * props.type === 'password'の場合のみ使用
 */
const showSecretValue = ref(false)

/**
 * フォームのtype属性
 */
const formType = computed(() => {
  if (props.type === 'password') {
    if (showSecretValue.value) {
      return 'text'
    } else {
      return 'password'
    }
  } else {
    return props.type
  }
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
      <div class="relative w-full">
        <input
          :id="id"
          v-model="value"
          :name="name"
          :placeholder="placeholder"
          :required="required"
          :type="formType"
          :pattern="pattern"
          :maxlength="maxLength"
          :class="{
            'block w-full appearance-none rounded-none border-b border-main bg-transparent px-2 leading-10 outline-none ring-0 focus:outline-none': true,
            'border-b-2 border-orange': hasError,
            'pr-8': type === 'password',
          }"
        />
        <button
          v-if="type === 'password'"
          class="absolute right-0 top-0 h-full px-2"
          type="button"
          @click="showSecretValue = !showSecretValue"
        >
          <template v-if="showSecretValue">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
              stroke-width="1.5"
              stroke="currentColor"
              class="h-5 w-5"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                d="M2.036 12.322a1.012 1.012 0 0 1 0-.639C3.423 7.51 7.36 4.5 12 4.5c4.638 0 8.573 3.007 9.963 7.178.07.207.07.431 0 .639C20.577 16.49 16.64 19.5 12 19.5c-4.638 0-8.573-3.007-9.963-7.178Z"
              />
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                d="M15 12a3 3 0 1 1-6 0 3 3 0 0 1 6 0Z"
              />
            </svg>
          </template>
          <template v-else>
            <svg
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
              stroke-width="1.5"
              stroke="currentColor"
              class="h-5 w-5"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                d="M3.98 8.223A10.477 10.477 0 0 0 1.934 12C3.226 16.338 7.244 19.5 12 19.5c.993 0 1.953-.138 2.863-.395M6.228 6.228A10.451 10.451 0 0 1 12 4.5c4.756 0 8.773 3.162 10.065 7.498a10.522 10.522 0 0 1-4.293 5.774M6.228 6.228 3 3m3.228 3.228 3.65 3.65m7.894 7.894L21 21m-3.228-3.228-3.65-3.65m0 0a3 3 0 1 0-4.243-4.243m4.242 4.242L9.88 9.88"
              />
            </svg>
          </template>
        </button>
      </div>
    </div>
    <p :class="{ 'text-orange': hasError, 'text-left text-sm': true }">
      {{ viewMessage }}
    </p>
  </div>
</template>
