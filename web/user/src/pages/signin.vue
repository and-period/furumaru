<script lang="ts" setup>
import { SignInRequest } from '~/types/api'
import { I18n } from '~/types/locales'

definePageMeta({
  layout: 'auth'
})

const i18n = useI18n()
const localePath = useLocalePath()

const t = (str: keyof I18n['auth']['signIn']): string => {
  return i18n.t(`auth.signIn.${str}`)
}

const handleClickGoogleSingInButton = () => {
  console.log('NOT IMPLEMENTED')
}
const handleClickFacebookSingInButton = () => {
  console.log('NOT IMPLEMENTED')
}
const handleClickLineSingInButton = () => {
  console.log('NOT IMPLEMENTED')
}

const formData = reactive<SignInRequest>({
  username: '',
  password: ''
})

const hasError = ref<boolean>(false)
</script>

<template>
  <the-sign-in-page
    v-model="formData"
    :page-name="t('pageName')"
    :button-text="t('signIn')"
    :has-error="hasError"
    :username-label="t('email')"
    :username-placeholder="t('email')"
    username-error-message=""
    :password-label="t('password')"
    :password-placeholder="t('password')"
    password-error-message=""
    :dont-have-account-text="t('dontHaveAccount')"
    :google-button-text="t('googleButtonText')"
    :facebook-button-text="t('facebookButtonText')"
    :line-button-text="t('lineButtonText')"
    :forget-password-link="{
      href: localePath('/'),
      text: t('forgetPasswordLink')
    }"
    :sign-up-link="{
      href: localePath('/signup'),
      text: t('signUpLink')
    }"
    @click:google-sing-in-button="handleClickGoogleSingInButton"
    @click:facebook-sing-in-button="handleClickFacebookSingInButton"
    @click:line-sing-in-button="handleClickLineSingInButton"
  />
</template>
