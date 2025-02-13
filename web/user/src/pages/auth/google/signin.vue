<script lang="ts" setup>
import { useAuthStore } from '~/store/auth'

const config = useRuntimeConfig()
const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()
const localePath = useLocalePath()

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

  const redirectUri = config.public.GOOGLE_SIGNIN_REDIRECT_URI as string
  await authStore.signInWithOAuth(code, redirectUri)
  router.push(localePath('/'))
}
catch (err) {
  console.error(err)

  const msg = 'Googleアカウントの認証に失敗しました'
  router.push(`/signin?error=${msg}`)
}
</script>
