import axios from 'axios'
import { defineStore } from 'pinia'
import Cookies from 'universal-cookie'

import { ApiClientFactory } from '.'

import {
  AuthApi,
  AuthResponse,
  SignInRequest,
  UpdateAuthPasswordRequest,
} from '~/types/api'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    redirectPath: '/',
    isAuthenticated: false,
    user: undefined as AuthResponse | undefined,
  }),

  getters: {
    accessToken(state): string | undefined {
      return state.user?.accessToken
    },
  },

  actions: {
    async signIn(payload: SignInRequest): Promise<string> {
      try {
        const factory = new ApiClientFactory()
        const authApiClient = factory.create(AuthApi)
        const res = await authApiClient.v1SignIn(payload)
        this.isAuthenticated = true
        this.user = res.data

        const cookies = new Cookies()
        cookies.set('refreshToken', this.user.refreshToken)
        return this.redirectPath
      } catch (err) {
        console.log(err)
        if (axios.isAxiosError(err)) {
          if (!err.response) {
            return Promise.reject(
              new Error(
                '現在、システムが停止中です。時間をおいてから再度アクセスしてください。'
              )
            )
          }
          switch (err.response?.status) {
            case 400:
            case 401:
              return Promise.reject(
                new Error('ユーザー名またはパスワードが違います。')
              )
            default:
              return Promise.reject(
                new Error(
                  '現在、システムが停止中です。時間をおいてから再度アクセスしてください。'
                )
              )
          }
        }
        throw new Error(
          '不明なエラーが発生しました。管理者にお問い合わせください。'
        )
      }
    },

    async passwordUpdate(payload: UpdateAuthPasswordRequest): Promise<void> {
      try {
        const factory = new ApiClientFactory()
        const authApiClient = factory.create(AuthApi, this.user?.accessToken)
        await authApiClient.v1UpdateAuthPassword(payload)
      } catch (err) {
        // TODO: エラーハンドリング
        console.log(err)
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
        this.user.refreshToken = refreshToken
      } catch (error) {
        const cookies = new Cookies()
        cookies.remove('refreshToken')
        if (axios.isAxiosError(error)) {
          throw new Error(error.message)
        }
        throw new Error('Internal Server Error')
      }
    },

    setRedirectPath(payload: string) {
      this.redirectPath = payload
    },
  },
})
