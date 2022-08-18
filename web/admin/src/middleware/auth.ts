import { Context, Middleware } from '@nuxt/types'
import Cookies from 'universal-cookie'

import { useAuthStore } from '~/store/auth'

const routing: Middleware = async ({ route, redirect }: Context) => {
  const publicPages = ['/signin']
  if (publicPages.includes(route.path)) {
    return
  }

  const store = useAuthStore()
  const cookies = new Cookies()

  // TODO: AccessTokenの有効期限確認も追加
  if (store.isAuthenticated) {
    return // ログイン済み
  }
  store.setRedirectPath(route.path)

  // RefreshTokenの有無検証
  const refreshToken: string = cookies.get('refreshToken')
  if (!refreshToken) {
    redirect('/signin')
    return
  }

  // AccessTokenの更新
  await store.getAuthByRefreshToken(refreshToken).catch(() => {
    redirect('/signin')
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
}

export default routing
