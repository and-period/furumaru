// docs: https://pinia.vuejs.org/core-concepts/#option-stores
import { defineStore, acceptHMRUpdate } from 'pinia'
import type {
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
      accessToken: '',
      refreshToken: '',
      user: undefined as AuthUserResponse | undefined,
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
        const res = await this.authApiClient().v1SignIn({ body: payload })
        this.isAuthenticated = true
        this.accessToken = res.accessToken
        this.refreshToken = res.refreshToken
        await this.fetchUserInfo()
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

    async fetchUserInfo() {
      const res = await this.authApiClient(this.accessToken).v1GetAuthUser()
      this.user = res
    },

    async logout() {
      try {
        await this.authApiClient(this.accessToken).v1SignOut()
      } catch (error) {
        this.errorHandler(error)
      } finally {
        // stateを初期状態にリセット
        this.$reset()
      }
    },
  },
})

// ホットリロードを有効にする
if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(useAuthStore, import.meta.hot))
}
