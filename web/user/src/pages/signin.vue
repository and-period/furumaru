<script lang="ts" setup>
import { useAuthStore } from '~/store/auth'
import type { SignInRequest } from '~/types/api'
import { ApiBaseError } from '~/types/exception'
import type { I18n } from '~/types/locales'

definePageMeta({
  layout: 'auth',
})

const config = useRuntimeConfig()
const route = useRoute()
const router = useRouter()
const i18n = useI18n()

const { signIn } = useAuthStore()
const localePath = useLocalePath()

const t = (str: keyof I18n['auth']['signIn']): string => {
  return i18n.t(`auth.signIn.${str}`)
}

const handleClickGoogleSingInButton = async () => {
  const state = crypto.randomUUID()
  sessionStorage.setItem('oauth_state', state)

  const url = `https://${config.public.COGNITO_AUTH_DOMAIN}/oauth2/authorize`
    + `?response_type=CODE`
    + `&client_id=${config.public.COGNITO_CLIENT_ID}`
    + `&redirect_uri=${config.public.GOOGLE_SIGNIN_REDIRECT_URI}`
    + `&state=${state}`
    + `&identity_provider=Google`
    + `&scope=openid email aws.cognito.signin.user.admin`
  await navigateTo(url, { external: true })
}
const handleClickLineSingInButton = async () => {
  const state = crypto.randomUUID()
  sessionStorage.setItem('oauth_state', state)

  const url = `https://${config.public.COGNITO_AUTH_DOMAIN}/oauth2/authorize`
    + `?response_type=CODE`
    + `&client_id=${config.public.COGNITO_CLIENT_ID}`
    + `&redirect_uri=${config.public.LINE_SIGNIN_REDIRECT_URI}`
    + `&state=${state}`
    + `&identity_provider=LINE`
    + `&scope=openid email profile aws.cognito.signin.user.admin`
  await navigateTo(url, { external: true })
}

const formData = reactive<SignInRequest>({
  username: '',
  password: '',
})

const hasError = ref<boolean>(false)
const errorMessage = ref<string>('')

const handleSubmit = async () => {
  try {
    await signIn(formData)
    router.push(localePath('/'))
  }
  catch (error) {
    hasError.value = true
    if (error instanceof ApiBaseError) {
      errorMessage.value = error.message
    }
  }
}

useSeoMeta({
  title: '新規アカウント登録',
})

try {
  const { error } = route.query as { error: string }
  if (error) {
    throw new Error(error)
  }
}
catch (error) {
  errorMessage.value = error as string
  console.error(error)
}
</script>

<template>
  <the-sign-in-page
    v-model="formData"
    :page-name="t('pageName')"
    :button-text="t('signIn')"
    :has-error="hasError"
    :error-message="errorMessage"
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
      text: t('forgetPasswordLink'),
    }"
    :sign-up-link="{
      href: localePath('/signup'),
      text: t('signUpLink'),
    }"
    @submit="handleSubmit"
    @click:google-sing-in-button="handleClickGoogleSingInButton"
    @click:line-sing-in-button="handleClickLineSingInButton"
  />
</template>
