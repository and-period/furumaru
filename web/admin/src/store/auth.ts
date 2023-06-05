import { getToken, isSupported } from 'firebase/messaging'
import { defineStore } from 'pinia'
import Cookies from 'universal-cookie'

import dayjs, { Dayjs } from 'dayjs'
import { useCommonStore } from './common'
import { messaging } from '~/plugins/firebase'
import {
  AuthResponse,
  SignInRequest,
  UpdateAuthEmailRequest,
  UpdateAuthPasswordRequest,
  VerifyAuthEmailRequest
} from '~/types/api'
import { apiClient } from '~/plugins/api-client'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    redirectPath: '/',
    isAuthenticated: false,
    user: undefined as AuthResponse | undefined,
    expiredAt: undefined as Dayjs | undefined
  }),

  getters: {
    accessToken (state): string | undefined {
      return state.user?.accessToken
    }
  },

  actions: {
    async signIn (payload: SignInRequest): Promise<string> {
      try {
        const res = await apiClient.authApi().v1SignIn(payload)
        this.setExpiredAt(res.data)
        this.isAuthenticated = true
        this.user = res.data

        const cookies = new Cookies()
        cookies.set('refreshToken', this.user.refreshToken, { secure: true })

        // Push通知の許可設定
        this.getDeviceToken()
          .then((deviceToken) => {
            if (deviceToken === '') {
              return // Push通知が未有効化状態
            }
            const currentToken: string = cookies.get('deviceToken')
            if (deviceToken === currentToken) {
              return // API側へ登録済み
            }
            return this.registerDeviceToken(deviceToken)
          })
          .catch((err) => {
            console.log('push notifications are disabled.', err)
          })

        return this.redirectPath
      } catch (err) {
        return this.errorHandler(err, { 401: 'ユーザー名またはパスワードが違います。' })
      }
    },

    async updatePassword (payload: UpdateAuthPasswordRequest): Promise<void> {
      try {
        await apiClient.authApi().v1UpdateAuthPassword(payload)
        const commonStore = useCommonStore()
        commonStore.addSnackbar({
          message: 'パスワードを更新しました。',
          color: 'info'
        })
      } catch (err) {
        return this.errorHandler(err)
      }
    },

    async emailUpdate (payload: UpdateAuthEmailRequest): Promise<void> {
      try {
        await apiClient.authApi().v1UpdateAuthEmail(payload)
        const commonStore = useCommonStore()
        commonStore.addSnackbar({
          message: '認証コードを送信しました。',
          color: 'info'
        })
      } catch (err) {
        return this.errorHandler(err, {
          409: 'このメールアドレスはすでに登録されているため、変更できません。',
          412: '変更前のメールアドレスと同じため、変更できません。'
        })
      }
    },

    async codeVerify (payload: VerifyAuthEmailRequest): Promise<void> {
      try {
        await apiClient.authApi().v1VerifyAuthEmail(payload)
        const commonStore = useCommonStore()
        commonStore.addSnackbar({
          message: 'メールアドレスが変更されました。',
          color: 'info'
        })
      } catch (err) {
        return this.errorHandler(err, { 409: 'このメールアドレスはすでに利用されているため使用できません。' })
      }
    },

    async getAuthByRefreshToken (refreshToken: string): Promise<void> {
      try {
        const res = await apiClient.authApi().v1RefreshAuthToken({
          refreshToken
        })
        this.setExpiredAt(res.data)
        this.isAuthenticated = true
        this.user = res.data
        this.user.refreshToken = refreshToken
      } catch (err) {
        const cookies = new Cookies()
        cookies.remove('refreshToken')
        return this.errorHandler(err)
      }
    },

    async registerDeviceToken (deviceToken: string): Promise<void> {
      try {
        await apiClient.authApi().v1RegisterAuthDevice({ device: deviceToken })

        const cookies = new Cookies()
        cookies.set('deviceToken', deviceToken, { secure: true })
      } catch (err) {
        return this.errorHandler(err)
      }
    },

    async getDeviceToken (): Promise<string> {
      const runtimeConfig = useRuntimeConfig()

      const supported = await isSupported()
      if (!supported) {
        console.log('this browser does not support push notificatins.')
        return '' // Push通知未対応ブラウザ
      }

      return await getToken(messaging, {
        vapidKey: runtimeConfig.public.FIREBASE_VAPID_KEY
      })
        .then((currentToken) => {
          return currentToken
        })
        .catch((err) => {
          console.log('failed to get device token', err)
          return ''
        })
    },

    setRedirectPath (payload: string) {
      this.redirectPath = payload
    },

    setExpiredAt (user: AuthResponse) {
      this.expiredAt = dayjs().add(user.expiresIn, 'second')
    },

    logout () {
      try {
        apiClient.authApi().v1SignOut()
        const cookies = new Cookies()
        cookies.remove('refreshToken')
        this.$reset()
      } catch (error) {
        console.log('APIでエラーが発生しました。', error)
      }
    }
  }
})
