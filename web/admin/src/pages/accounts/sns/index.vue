<script lang="ts" setup>
import { useAuthStore } from '~/store'

const authStore = useAuthStore()

const authGoogleAccount = async (): Promise<void> => {
  try {
    const config = useRuntimeConfig()

    let redirectUri: string | undefined
    if (config.public.GOOGLE_CONNECT_REDIRECT_URI !== '') {
      redirectUri = config.public.GOOGLE_CONNECT_REDIRECT_URI
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
</script>

<template>
  <div>
    <v-btn @click="authGoogleAccount">
      Google連携
    </v-btn>
  </div>
</template>
