import { defineStore } from 'pinia'

import { ApiClientFactory } from '.'

import { AuthApi, AuthResponse, SignInRequest, UpdateAuthPasswordRequest } from '~/types/api'

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
        const factory = new ApiClientFactory()
        const authApiClient = factory.create(AuthApi)
        const res = await authApiClient.v1SignIn(payload)
        this.isAuthenticated = true
        this.user = res.data
      } catch (err) {
        // TODO: エラーハンドリング
        console.log(err)
        throw new Error('Internal Server Error')
      }
    },
    async passwordUpdate(payload: UpdateAuthPasswordRequest): Promise<void> {
      try {
        const factory = new ApiClientFactory()
        console.log(this.user?.accessToken)
        const authApiClient = factory.create(AuthApi, this.user?.accessToken)
        await authApiClient.v1UpdateAuthPassword(payload)
        this.isAuthenticated = true
      } catch (err) {
        // TODO: エラーハンドリング
        console.log(err)
        throw new Error('Internal Server Error')
      }
    }
  },
})
