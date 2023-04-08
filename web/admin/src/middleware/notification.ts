import { Context } from '@nuxt/types'

import { useAuthStore } from '~/store'
import { useMessageStore } from '~/store'

export default ({ route }: Context) => {
  const publicPages = ['/signin']
  if (publicPages.includes(route.path)) {
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
}
