<script lang="ts" setup>
import { useAuthStore } from '~/store'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

try {
  const { code, state } = route.query as { code: string, state: string }
  if (!code || code === '') {
    throw new Error('code is empty')
  }

  const oauthState = sessionStorage.getItem('oauth_state')
  const oauthNonce = sessionStorage.getItem('oauth_nonce') || ''

  sessionStorage.removeItem('oauth_state')
  sessionStorage.removeItem('oauth_nonce')

  if (!oauthState || oauthState !== state) {
    throw new Error('state is invalid')
  }

  const config = useRuntimeConfig()

  let redirectUri: string | undefined
  if (config.public.GOOGLE_CONNECT_REDIRECT_URI !== '') {
    redirectUri = config.public.GOOGLE_CONNECT_REDIRECT_URI
  }

  await authStore.linkGoogleAccount(code, oauthNonce, redirectUri)

  // TODO: 成功のメッセージを表示するように
  router.push('/accounts/sns')
}
catch (err) {
  console.error(err)
}
</script>
