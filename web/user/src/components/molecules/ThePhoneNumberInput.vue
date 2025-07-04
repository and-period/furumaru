<script setup lang="ts">
import type { I18n } from '~/types/locales'

interface Props {
  modelValue: string
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

interface Emits {
  (e: 'update:modelValue', val: string): void
}

const emits = defineEmits<Emits>()

const i18n = useI18n()

const gt = (str: keyof I18n['purchase']['guest']) => {
  return i18n.t(`purchase.guest.${str}`)
}

const tel1 = ref<string>('')
const tel2 = ref<string>('')
const tel3 = ref<string>('')

watch([tel1, tel2, tel3], () => {
  emits('update:modelValue', `${tel1.value}-${tel2.value}-${tel3.value}`)
})

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

onMounted(() => {
  // 初期値を設定
  const tel = props.modelValue.split('-')
  if (tel.length === 3) {
    tel1.value = tel[0]
    tel2.value = tel[1]
    tel3.value = tel[2]
  }
})
</script>

<template>
  <div>
    <label class="text-[12px] tracking-[1.2px] text-main">
      {{ gt('phoneNumberLabel') }}
    </label>
    <div class="grid w-full grid-cols-11 items-center gap-2">
      <input
        v-model="tel1"
        type="tel"
        size="3"
        placeholder="000"
        class="col-span-3 block appearance-none rounded-none border-b border-main bg-transparent px-2 leading-10 outline-none ring-0 focus:outline-none"
        required
      >
      -
      <input
        v-model="tel2"
        type="tel"
        size="4"
        placeholder="0000"
        class="col-span-3 appearance-none rounded-none border-b border-main bg-transparent px-2 leading-10 outline-none ring-0 focus:outline-none"
        required
      >
      -
      <input
        v-model="tel3"
        type="tel"
        size="4"
        placeholder="0000"
        class="col-span-3 appearance-none rounded-none border-b border-main bg-transparent px-2 leading-10 outline-none ring-0 focus:outline-none"
        required
      >
    </div>
    <p :class="{ 'text-orange': hasError, 'text-left text-sm': true }">
      {{ viewMessage }}
    </p>
  </div>
</template>
