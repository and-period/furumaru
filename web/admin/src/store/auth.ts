import { defineStore } from 'pinia'

import { ApiClientFactory } from '.'

import Cookies from 'universal-cookie'

import { AuthApi, AuthResponse, SignInRequest } from '~/types/api'

import axios from 'axios'

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

        const cookies = new Cookies()
        cookies.set('refreshToken', this.user.refreshToken)
      } catch (err) {
        // TODO: エラーハンドリング
        throw new Error('Internal Server Error')
      }
    },
    async getAuthByRefreshToken(refreshToken: string): Promise<void> {
      try {
        const factory = new ApiClientFactory()
        const authApiClient = factory.create(AuthApi)
        const res = await authApiClient.v1RefreshAuthToken({
          refreshToken,
        })
        this.isAuthenticated = true
        this.user = res.data

        const cookies = new Cookies()
        cookies.set('refreshToken', this.user.refreshToken)
      } catch (error) {
        const cookies = new Cookies()
        cookies.remove('refreshToken')
        if (axios.isAxiosError(error)) {
          throw new Error(error.message)
        }
        throw new Error('Internal Server Error')
      }
    },
  },
})
