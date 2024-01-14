<script lang="ts" setup>
import { useAuthStore } from '~/store/auth'
import { ApiBaseError } from '~/types/exception'
import type { I18n } from '~/types/locales'

definePageMeta({
  layout: 'auth',
})

const authStore = useAuthStore()
const { verifyAuth } = authStore

const route = useRoute()
const router = useRouter()

const i18n = useI18n()

const t = (str: keyof I18n['auth']['verify']) => {
  return i18n.t(`auth.verify.${str}`)
}

// 新規登録時に発行されたID
const id = computed<string>(() => {
  const id = route.query.id
  if (id) {
    return id as string
  } else {
    return ''
  }
})

// 買い物カゴ画面から認証に飛ばされたかのフラグ
const redirectToPurchase = computed<boolean>(() => {
  const redirectToPurchaseParam = route.query.redirect_to_purchase
  if (redirectToPurchaseParam) {
    return Boolean(redirectToPurchaseParam)
  } else {
    return false
  }
})

// コーディネーターID
const coordinatorId = computed<string>(() => {
  const id = route.query.coordinatorId
  if (id) {
    return String(id)
  } else {
    return ''
  }
})

// カート番号
const cartNumber = computed<number | undefined>(() => {
  const id = route.query.cartNumber
  const idNumber = Number(id)
  if (idNumber === 0) {
    return undefined
  }
  if (isNaN(idNumber)) {
    return undefined
  }
  return idNumber
})

const code = ref<string>('')
const errorMessage = ref<string>('')

const handleSubmit = async () => {
  try {
    await verifyAuth({
      verifyCode: code.value,
      id: id.value,
    })
    if (redirectToPurchase.value) {
      router.push({
        path: '/v1/purchase/auth',
        query: {
          from_new_accounet: true,
          coordinatorId: coordinatorId.value,
          cartNumber: cartNumber.value,
        },
      })
    } else {
      router.push('/')
    }
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
