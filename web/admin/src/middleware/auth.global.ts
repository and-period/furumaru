import dayjs from 'dayjs'
import Cookies from 'universal-cookie'

import { useAuthStore } from '~/store'
import type { AuthResponse } from '~/types/api'

const publicPages = [
  '/health',
  '/signin',
  '/recover',
  '/privacy',
  '/legal-notice',
  '/auth/google/callback',
  '/auth/line/callback',
  '/auth/youtube/callback',
  '/auth/youtube/complete',
]

export default defineNuxtRouteMiddleware(async (to, _) => {
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

  try {
    const auth = cookies.get('auth')
    if (auth) {
      const parsed = JSON.parse(decodeURIComponent(auth)) as AuthResponse
      await store.setAuth(parsed)
    }
    else {
      // RefreshTokenの有無検証
      const refreshToken = cookies.get('refreshToken')
      if (!refreshToken) {
        console.log('refresh token is not found.')
        return navigateTo('/signin')
      }

      // AccessTokenの更新
      await store.getAuthByRefreshToken(refreshToken.value)

      // ログインユーザーの情報取得
      store.getUser().catch((err) => {
        console.log('failed to get user', err)
      })
    }
  }
  catch (err) {
    console.log('failed to refresh auth token', err)
    return navigateTo('/signin')
  }

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
