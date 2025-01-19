<script setup lang="ts">
import type { I18n } from '~/types/locales'

interface Props {
  required?: boolean
  error?: boolean
  errorMessage?: string
}

const props = withDefaults(defineProps<Props>(), {
  required: false,
  withLabel: true,
  error: false,
  message: undefined,
  errorMessage: '',
})

const modelValue = defineModel<string>({ required: true })

const i18n = useI18n()

const gt = (str: keyof I18n['purchase']['guest']) => {
  return i18n.t(`purchase.guest.${str}`)
}

/**
 * エラーの判定
 * errorMessageが渡されている場合はエラー状態にする
 */
const hasError = computed(() => {
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
  }
  else {
    return props.message
  }
})
</script>

<template>
  <div>
    <div class="w-full items-center gap-2">
      <input
        v-model="modelValue"
        type="tel"
        :placeholder="gt('phoneNumberLabel')"
        class="w-full block appearance-none rounded-none border-b border-main bg-transparent px-2 leading-10 outline-none ring-0 focus:outline-none"
        required
      >
    </div>
    <p :class="{ 'text-orange': hasError, 'text-left text-sm': true }">
      {{ viewMessage }}
    </p>
  </div>
</template>
