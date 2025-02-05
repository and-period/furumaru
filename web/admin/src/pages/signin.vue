<script lang="ts" setup>
import { useAlert } from '~/lib/hooks'
import { useAuthStore } from '~/store'
import type { SignInRequest } from '~/types/api'

definePageMeta({
  layout: 'auth',
})

const config = useRuntimeConfig()
const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()
const { alertType, isShow, alertText, show } = useAlert('error')

const loading = ref<boolean>(false)
const formData = reactive<SignInRequest>({
  username: '',
  password: '',
})

// Google の認証ページにリダイレクト
const loginWithGoogle = async () => {
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

const handleSubmit = async () => {
  try {
    loading.value = true
    const path = await authStore.signIn(formData)
    router.push(path)
  }
  catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
  finally {
    loading.value = false
  }
}

try {
  const { error } = route.query as { error: string }
  if (error) {
    throw new Error(error)
  }
}
catch (err) {
  if (err instanceof Error) {
    show(err.message)
  }
  console.error(err)
}
</script>

<template>
  <templates-auth-sign-in
    v-model:form-data="formData"
    :loading="loading"
    :is-alert="isShow"
    :alert-type="alertType"
    :alert-text="alertText"
    @click:login-with-google="loginWithGoogle"
    @submit="handleSubmit"
  />
</template>
