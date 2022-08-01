import axios from 'axios'
import { defineStore } from 'pinia'
import Cookies from 'universal-cookie'

import ApiClientFactory from '../plugins/factory'

import { useCommonStore } from './common'

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
        const commonStore = useCommonStore()
        commonStore.addSnackbar({
          message: 'パスワードを更新しました。',
          color: 'info',
        })
      } catch (error) {
        if (axios.isAxiosError(error)) {
          return Promise.reject(new Error(error.message))
        } else {
          return Promise.reject(new Error('不明なエラーが発生しました。'))
        }
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
        throw new Error('不明なエラーが発生しました。')
      }
    },

    setRedirectPath(payload: string) {
      this.redirectPath = payload
    },
  },
})
