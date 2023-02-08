// docs: https://pinia.vuejs.org/core-concepts/#option-stores
import { defineStore, acceptHMRUpdate } from 'pinia'
import { SignInRequest } from '~/types/api'

/**
 * 認証情報を管理するグローバルステート
 */
export const useAuthStore = defineStore('auth', {
  state: () => {
    return {
      isAuthenticated: false
    }
  },

  actions: {
    /**
     * ログインを実施する非同期関数
     * @param payload
     * @returns
     */
    async signIn (payload: SignInRequest) {
      try {
        await this.authApiClient().v1SignIn(payload)
        this.isAuthenticated = true
      } catch (error) {
        return this.errorHandler(error, {
          401: this.i18n.t('auth.signIn.authErrorMessage')
        })
      }
    }
  }
})

// ホットリロードを有効にする
if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(useAuthStore, import.meta.hot))
}
