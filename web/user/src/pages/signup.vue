<script lang="ts" setup>
import { CreateAuthRequest } from '~/types/api'
import { I18n } from '~/types/locales'

definePageMeta({
  layout: 'auth'
})

const router = useRouter()
const i18n = useI18n()
const localePath = useLocalePath()

const t = (str: keyof I18n['auth']['signUp']) => {
  return i18n.t(`auth.signUp.${str}`)
}

const formData = reactive<CreateAuthRequest>({
  email: '',
  phoneNumber: '',
  password: '',
  passwordConfirmation: ''
})

const handleSubmit = () => {
  router.push('/verify')
}
</script>

<template>
  <the-sign-up-page
    v-model="formData"
    :page-name="t('pageName')"
    :button-text="t('signUp')"
    :tel-label="t('tel')"
    :tel-placeholder="t('tel')"
    tel-error-message=""
    :email-label="t('email')"
    :email-placeholder="t('email')"
    email-error-message=""
    :password-label="t('password')"
    :password-placeholder="t('password')"
    password-error-message=""
    :password-confirm-label="t('passwordConfirm')"
    :password-confirm-placeholder="t('passwordConfirm')"
    password-confirm-error-message=""
    :already-has-link="{
      href: localePath('/signin'),
      text: t('alreadyHas')
    }"
    @submit="handleSubmit"
  />
</template>
