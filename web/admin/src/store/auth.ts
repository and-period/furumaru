import { defineStore } from 'pinia'

import { AuthApi, AuthResponse, SignInRequest } from '~/types/api'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    isAuthenticated: false,
    user: undefined as AuthResponse | undefined,
  }),
  getters: {
    accessToken(state): string | undefined {
      return state.user?.accessToken
    },
  },
  actions: {
    async signIn(payload: SignInRequest): Promise<void> {
      try {
        const authApiClient = new AuthApi()
        const res = await authApiClient.v1SignIn(payload)
        this.isAuthenticated = true
        this.user = res.data
      } catch (err) {
        // TODO: エラーハンドリング
        throw new Error('Internal Server Error')
      }
    },
  },
})
