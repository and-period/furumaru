// docs: https://pinia.vuejs.org/core-concepts/#option-stores
import { defineStore, acceptHMRUpdate } from 'pinia'
import {
  AuthUserResponse,
  CreateAuthRequest,
  SignInRequest,
  VerifyAuthRequest,
} from '~/types/api'

/**
 * 認証情報を管理するグローバルステート
 */
export const useAuthStore = defineStore('auth', {
  // 永続化の設定
  persist: {
    storage: persistedState.cookiesWithOptions({
      sameSite: 'strict',
    }),
  },

  state: () => {
    return {
      isAuthenticated: false,
    }
  },

  actions: {
    /**
     * ログインを実施する非同期関数
     * @param payload
     * @returns
     */
    async signIn(payload: SignInRequest) {
      try {
        await this.authApiClient().v1SignIn({ body: payload })
        this.isAuthenticated = true
      } catch (error) {
        return this.errorHandler(error, {
          401: this.i18n.t('auth.signIn.authErrorMessage'),
        })
      }
    },

    async signUp(payload: CreateAuthRequest): Promise<AuthUserResponse> {
      try {
        const res = await this.authApiClient().v1CreateAuth({ body: payload })
        return res
      } catch (error) {
        return this.errorHandler(error, {
          409: '指定したメールアドレスはご利用できません。',
        })
      }
    },

    async verifyAuth(payload: VerifyAuthRequest) {
      try {
        await this.authApiClient().v1VerifyAuth({ body: payload })
      } catch (error) {
        return this.errorHandler(error)
      }
    },
  },
})

// ホットリロードを有効にする
if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(useAuthStore, import.meta.hot))
}
