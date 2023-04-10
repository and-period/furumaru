import { useAuthStore, useMessageStore } from '~/store'

export default defineNuxtRouteMiddleware((to, _) => {
  const publicPages = ['/signin']
  if (publicPages.includes(to.path)) {
    return
  }

  const authStore = useAuthStore()
  if (!authStore.isAuthenticated) {
    return // 未ログイン
  }

  const messageStore = useMessageStore()
  if (messageStore.total > 0) {
    return // すでに取得済みのデータがある
  }

  messageStore.fetchMessages().catch((err: Error) => {
    console.log('failed to get messages', err)
  })
})
