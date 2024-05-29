import dayjs from 'dayjs'
import { useAuthStore } from '~/store/auth'
import { ApiBaseError } from '~/types/exception'

export default defineNuxtRouteMiddleware(async () => {
  const authStore = useAuthStore()

  if (authStore.isAuthenticated) {
    if (dayjs().isBefore(authStore.expiredAt)) {
      // ログイン済み && AccessTokenが有効
    }
    else {
      // 認証トークンの期限切れの場合
      // 認証トークンとユーザー情報の再取得を行う
      try {
        await authStore.refreshAccsessToken(authStore.refreshToken)
      }
      catch (error) {
        // リフレッシュトークンが期限切れの場合
        if (error instanceof ApiBaseError) {
          console.log('リフレッシュトークンが期限切れです', error)
        }
        // フロントエンドの認証情報をクリアする
        authStore.resetState()
        return
      }
      try {
        await authStore.fetchUserInfo()
      }
      catch (error) {
        if (error instanceof ApiBaseError) {
          console.log('認証エラー', error)
        }
      }
    }
  }
})
