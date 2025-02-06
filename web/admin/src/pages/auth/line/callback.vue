<script lang="ts" setup>
import { useAuthStore } from '~/store'

const config = useRuntimeConfig()
const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

try {
  const { code, state } = route.query as { code: string, state: string }
  if (!code || code === '') {
    throw new Error('code is empty')
  }

  const oauthState = sessionStorage.getItem('oauth_state')
  sessionStorage.removeItem('oauth_state')

  if (!oauthState || oauthState !== state) {
    throw new Error('state is invalid')
  }

  const redirectUri = config.public.LINE_SIGNIN_REDIRECT_URI
  const path = await authStore.signInWithOAuth(code, redirectUri)
  router.push(path)
}
catch (err) {
  const msg = 'LINEアカウントの認証に失敗しました'
  router.push(`/auth/signin?error=${msg}`)
}
</script>
