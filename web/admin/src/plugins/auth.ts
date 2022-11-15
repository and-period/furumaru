import { Plugin } from '@nuxt/types'

import { useAuthStore } from '~/store/auth'

/**
 * authStoreを注入するプラグイン
 */
const authPlugin: Plugin = (_ctx, inject) => {
  const authStore = useAuthStore()
  inject('auth', authStore)
}

export default authPlugin
