import dayjs from 'dayjs'
import Cookies from 'universal-cookie'

import { useAuthStore } from '~/store'

export default defineNuxtRouteMiddleware(async (to, _) => {
  const publicPages = ['/signin', '/recover', '/privacy', '/legal-notice', '/auth/youtube/callback', '/auth/youtube/complete']
  if (publicPages.includes(to.path)) {
    return
  }

  const store = useAuthStore()
  const cookies = new Cookies()

  // ログイン中 && AccessTokenの有効期限が切れていないかの検証
  if (store.isAuthenticated && dayjs().isBefore(store.expiredAt)) {
    return // ログイン済み && AccessTokenが有効
  }
  store.setRedirectPath(to.path)

  // RefreshTokenの有無検証
  const refreshToken: string = cookies.get('refreshToken')
  if (!refreshToken) {
    return navigateTo('/signin')
  }

  try {
    // AccessTokenの更新
    await store.getAuthByRefreshToken(refreshToken)
  } catch (err) {
    console.log('failed to refresh auth token', err)
    return navigateTo('/signin')
  }

  // ログインユーザーの情報取得
  store.getUser().catch((err) => {
    console.log('failed to get user', err)
  })

  // Push通知用のDeviceToken取得/登録
  store
    .getDeviceToken()
    .then((deviceToken) => {
      if (deviceToken === '') {
        return // Push通知が未有効化状態
      }
      const currentToken: string = cookies.get('deviceToken')
      if (deviceToken === currentToken) {
        return // API側へ登録済み
      }
      return store.registerDeviceToken(deviceToken)
    })
    .catch((err) => {
      console.log('push notifications are disabled.', err)
    })
})
