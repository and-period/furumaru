<script lang="ts" setup>
import { useAuthStore } from '~/store/auth'
import { ApiBaseError } from '~/types/exception'
import { I18n } from '~/types/locales'

definePageMeta({
  layout: 'auth',
})

const authStore = useAuthStore()
const { verifyAuth } = authStore

const route = useRoute()

const i18n = useI18n()

const t = (str: keyof I18n['auth']['verify']) => {
  return i18n.t(`auth.verify.${str}`)
}

const id = computed<string>(() => {
  const id = route.query.id
  if (id) {
    return id as string
  } else {
    return ''
  }
})

const code = ref<string>('')
const errorMessage = ref<string>('')

const handleSubmit = async () => {
  try {
    await verifyAuth({
      verifyCode: code.value,
      id: id.value,
    })
  } catch (error) {
    if (error instanceof ApiBaseError) {
      errorMessage.value = error.message
    }
  }
}

onMounted(() => {
  if (id.value === '') {
    errorMessage.value = 'アカウント新規登録画面から操作を実施してください。'
  }
})

useSeoMeta({
  title: '認証コード入力',
})
</script>

<template>
  <the-verify-code-page
    v-model:code="code"
    :error-message="errorMessage"
    :page-name="t('pageName')"
    :button-text="t('btnText')"
    :message="t('message')"
    @submit="handleSubmit"
  />
</template>
