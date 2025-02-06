<script lang="ts" setup>
import { useAlert } from '~/lib/hooks'
import { useAuthStore, useCommonStore } from '~/store'

const route = useRoute()
const authStore = useAuthStore()
const commonStore = useCommonStore()
const { isShow, alertText, alertType, show } = useAlert('error')

const { providers } = storeToRefs(authStore)

const loading = ref<boolean>(false)

const fetchState = useAsyncData(async (): Promise<void> => {
  await fetchAuthProviders()
})

const fetchAuthProviders = async (): Promise<void> => {
  try {
    await authStore.listAuthProviders()
  }
  catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.error(err)
  }
}

const isLoading = (): boolean => {
  return fetchState?.pending?.value || loading.value
}

const authGoogleAccount = async (): Promise<void> => {
  try {
    const config = useRuntimeConfig()

    let redirectUri: string | undefined
    if (config.public.GOOGLE_CONNECT_REDIRECT_URI !== '') {
      redirectUri = config.public.GOOGLE_CONNECT_REDIRECT_URI as string
    }

    const state = crypto.randomUUID()
    const authUrl = await authStore.getAuthGoogleUrl(state, redirectUri)
    const parsed = new URL(authUrl)

    sessionStorage.setItem('oauth_state', state)
    sessionStorage.setItem('oauth_nonce', parsed.searchParams.get('nonce') || '')

    await navigateTo(authUrl, { external: true })
  }
  catch (err) {
    console.error(err)
  }
}

const authLineAccount = async (): Promise<void> => {
  try {
    const config = useRuntimeConfig()

    let redirectUri: string | undefined
    if (config.public.LINE_CONNECT_REDIRECT_URI !== '') {
      redirectUri = config.public.LINE_CONNECT_REDIRECT_URI as string
    }

    const state = crypto.randomUUID()
    const authUrl = await authStore.getAuthLineUrl(state, redirectUri)
    const parsed = new URL(authUrl)

    sessionStorage.setItem('oauth_state', state)
    sessionStorage.setItem('oauth_nonce', parsed.searchParams.get('nonce') || '')

    await navigateTo(authUrl, { external: true })
  }
  catch (err) {
    console.error(err)
  }
}

try {
  const { success, error } = route.query as { success: string, error: string }
  if (error) {
    throw new Error(error)
  }

  if (success) {
    commonStore.addSnackbar({ color: 'success', message: success })
  }

  await fetchState.execute()
}
catch (err) {
  if (err instanceof Error) {
    show(err.message)
  }
  console.error(err)
}
</script>

<template>
  <templates-auth-provider
    :loading="isLoading"
    :is-alert="isShow"
    :alert-type="alertType"
    :alert-text="alertText"
    :providers="providers"
    @click:google="authGoogleAccount"
    @click:line="authLineAccount"
  />
</template>
