<script lang="ts" setup>
import { useAuthStore } from '~/store/auth'
import type { CreateAuthRequest } from '~/types/api'
import type { I18n } from '~/types/locales'
import { convertJapaneseToI18nPhoneNumber } from '~/lib/phone-number'
import { ApiBaseError } from '~/types/exception'

definePageMeta({
  layout: 'auth',
})

const route = useRoute()
const router = useRouter()
const i18n = useI18n()
const localePath = useLocalePath()

const { signUp } = useAuthStore()

// 買い物カゴ画面から飛ばされたかのフラグ
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

const t = (str: keyof I18n['auth']['signUp']) => {
  return i18n.t(`auth.signUp.${str}`)
}

const formData = reactive<CreateAuthRequest>({
  username: '',
  accountId: '',
  lastname: '',
  firstname: '',
  lastnameKana: '',
  firstnameKana: '',
  email: '',
  phoneNumber: '',
  password: '',
  passwordConfirmation: '',
})

const hasError = ref<boolean>(false)
const telErrorMessage = ref<string>('')
const passwordConfirmErrorMessage = ref<string>('')

const apiErrorMessage = ref<string>('')

const validate = () => {
  hasError.value = false
  telErrorMessage.value = ''
  passwordConfirmErrorMessage.value = ''

  if (!formData.phoneNumber.startsWith('0')) {
    telErrorMessage.value = '電話番号は0から始まる値を入力してください'
    hasError.value = true
  }

  if (formData.password !== formData.passwordConfirmation) {
    passwordConfirmErrorMessage.value = 'パスワードが一致しません。'
    hasError.value = true
  }

  return hasError.value
}

const handleSubmit = async () => {
  try {
    if (validate()) {
      return
    }

    const result = await signUp({
      ...formData,
      phoneNumber: convertJapaneseToI18nPhoneNumber(formData.phoneNumber),
    })

    router.push({
      path: 'verify',
      query: {
        id: result.id,
        redirect_to_purchase: redirectToPurchase.value,
        coordinatorId: coordinatorId.value,
        cartNumber: cartNumber.value,
      },
    })
  } catch (error) {
    if (error instanceof ApiBaseError) {
      apiErrorMessage.value = error.message
    }
  }
}

useSeoMeta({
  title: '新規アカウント登録',
})
</script>

<template>
  <the-sign-up-page
    v-model="formData"
    :page-name="t('pageName')"
    :error-message="apiErrorMessage"
    :button-text="t('signUp')"
    :tel-label="t('tel')"
    :tel-placeholder="t('tel')"
    :tel-error-message="telErrorMessage"
    :email-label="t('email')"
    :email-placeholder="t('email')"
    email-error-message=""
    :password-label="t('password')"
    :password-placeholder="t('password')"
    password-error-message=""
    :password-confirm-label="t('passwordConfirm')"
    :password-confirm-placeholder="t('passwordConfirm')"
    :password-confirm-error-message="passwordConfirmErrorMessage"
    :already-has-link="{
      href: localePath('/signin'),
      text: t('alreadyHas'),
    }"
    @submit="handleSubmit"
  />
</template>
