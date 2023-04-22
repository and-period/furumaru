import dayjs from 'dayjs'
import Cookies from 'universal-cookie'

import { useAuthStore } from '~/store'

export default defineNuxtRouteMiddleware(async (to, _) => {
  const publicPages = ['/signin']
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

  // AccessTokenの更新
  await store.getAuthByRefreshToken(refreshToken).catch((err) => {
    console.log('failed to refresh auth token', err)
    navigateTo('/signin')
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
      console.log('Push notifications are disabled.', err)
    })
})
