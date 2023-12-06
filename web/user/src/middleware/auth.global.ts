import dayjs from 'dayjs'
import { useAuthStore } from '~/store/auth'

export default defineNuxtRouteMiddleware(async () => {
  const authStore = useAuthStore()

  if (authStore.isAuthenticated) {
    if (dayjs().isBefore(authStore.expiredAt)) {
      // ログイン済み && AccessTokenが有効
    } else {
      // 認証トークンの期限切れの場合
      // 認証トークンとユーザー情報の再取得を行う
      await authStore.refreshAccsessToken(authStore.refreshToken)
      await authStore.fetchUserInfo()
    }
  }
})
