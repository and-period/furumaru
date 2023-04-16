import { useAuthStore } from '~/store'

/**
 * authStoreを注入するプラグイン
 */
export default defineNuxtPlugin(() => {
  return {
    provide: {
      auth: useAuthStore()
    }
  }
})
