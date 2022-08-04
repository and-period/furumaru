import axios from 'axios'
import { defineStore } from 'pinia'
import Cookies from 'universal-cookie'

import ApiClientFactory from '../plugins/factory'

import { useCommonStore } from './common'

import {
  AuthApi,
  AuthResponse,
  SignInRequest,
  UpdateAuthEmailRequest,
  UpdateAuthPasswordRequest,
  VerifyAuthEmailRequest,
} from '~/types/api'
import {
  AuthError,
  ConnectionError,
  InternalServerError,
  ValidationError,
} from '~/types/exception'

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
        if (axios.isAxiosError(err)) {
          if (!err.response) {
            return Promise.reject(new ConnectionError(err))
          }
          switch (err.response.status) {
            case 400:
            case 401:
              return Promise.reject(
                new ValidationError(
                  err.response.status,
                  'ユーザー名またはパスワードが違います。',
                  err
                )
              )
            default:
              return Promise.reject(new InternalServerError(err))
          }
        }
        throw new InternalServerError(err)
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
      } catch (err) {
        if (axios.isAxiosError(err)) {
          if (!err.response) {
            return Promise.reject(new ConnectionError(err))
          }
          const statusCode = err.response.status
          switch (statusCode) {
            case 401:
              return Promise.reject(
                new AuthError(
                  statusCode,
                  '認証エラー。再度ログインをしてください。',
                  err
                )
              )
            case 400:
              return Promise.reject(
                new ValidationError(statusCode, '入力値に誤りがあります。', err)
              )
            default:
              return Promise.reject(new InternalServerError(err))
          }
        }
        throw new InternalServerError(err)
      }
    },

    async emailUpdate(payload: UpdateAuthEmailRequest): Promise<void> {
      try {
        const factory = new ApiClientFactory()
        const authApiClient = factory.create(AuthApi, this.user?.accessToken)
        await authApiClient.v1UpdateAuthEmail(payload)
        const commonStore = useCommonStore()
        commonStore.addSnackbar({
          message: '認証コードを送信しました。',
          color: 'info',
        })
      } catch (err) {
        console.log(err)
        throw new Error('Internal Server Error')
      }
    },

    async codeVerify(payload: VerifyAuthEmailRequest): Promise<void> {
      try {
        const factory = new ApiClientFactory()
        const authApiClient = factory.create(AuthApi, this.user?.accessToken)
        await authApiClient.v1VerifyAuthEmail(payload)
        const commonStore = useCommonStore()
        commonStore.addSnackbar({
          message: 'メールアドレスが変更されました。',
          color: 'info',
        })
      } catch (err) {
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
          if (!error.response) {
            return Promise.reject(new ConnectionError(error))
          }
          if (error.response.status === 401) {
            return Promise.reject(
              new AuthError(
                error.response.status,
                '認証エラー。再度ログインをしてください。',
                error
              )
            )
          } else {
            return Promise.reject(new InternalServerError(error))
          }
        }
        throw new InternalServerError(error)
      }
    },

    setRedirectPath(payload: string) {
      this.redirectPath = payload
    },
  },
})
