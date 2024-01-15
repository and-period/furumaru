import { useAuthStore } from '~/store/auth'

/**
 * 認証が必要なページにアクセスした際に、認証情報がない場合はログインページにリダイレクトするミドルウェア
 * グローバルなミドルウェアではないのでページごとに設定する必要がある
 */
export default defineNuxtRouteMiddleware(() => {
  const authStore = useAuthStore()

  if (!authStore.isAuthenticated) {
    return navigateTo('/signin')
  }
})
