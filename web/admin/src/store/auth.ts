import dayjs, { Dayjs } from 'dayjs'
import { getToken, isSupported } from 'firebase/messaging'
import { defineStore } from 'pinia'
import Cookies from 'universal-cookie'

import { messaging } from '~/plugins/firebase'
import { apiClient } from '~/plugins/api-client'
import {
  AdminRole,
  type AuthResponse,
  type AuthUserResponse,
  type ForgotAuthPasswordRequest,
  type ResetAuthPasswordRequest,
  type SignInRequest,
  type UpdateAuthEmailRequest,
  type UpdateAuthPasswordRequest,
  type VerifyAuthEmailRequest
} from '~/types/api'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    redirectPath: '/',
    isAuthenticated: false,
    auth: undefined as AuthResponse | undefined,
    user: undefined as AuthUserResponse | undefined,
    expiredAt: undefined as Dayjs | undefined
  }),

  getters: {
    adminId (state): string {
      return state.auth?.adminId || ''
    },
    accessToken (state): string | undefined {
      return state.auth?.accessToken
    },
    role (state): AdminRole {
      return state.auth?.role || AdminRole.UNKNOWN
    }
  },

  actions: {
    /**
     * サインイン
     * @param payload メールアドレス/パスワード
     * @returns 遷移先Path
     */
    async signIn (payload: SignInRequest): Promise<string> {
      try {
        const res = await apiClient.authApi().v1SignIn(payload)
        this.auth = res.data
        this.isAuthenticated = true

        this.getUser()
        this.setExpiredAt(res.data)

        const cookies = new Cookies()
        cookies.set('refreshToken', this.auth.refreshToken, { secure: true })

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

    /**
     * サインイン中管理者情報取得
     */
    async getUser (): Promise<void> {
      try {
        const res = await apiClient.authApi().v1GetAuthUser()
        this.user = res.data
      } catch (err) {
        return this.errorHandler(err)
      }
    },

    /**
     * メールアドレス更新
     * @param payload
     */
    async updateEmail (payload: UpdateAuthEmailRequest): Promise<void> {
      try {
        await apiClient.authApi().v1UpdateAuthEmail(payload)
      } catch (err) {
        return this.errorHandler(err, {
          409: 'このメールアドレスはすでに登録されているため、変更できません。',
          412: '変更前のメールアドレスと同じため、変更できません。'
        })
      }
    },

    /**
     * メールアドレス更新後の検証
     * @param payload
     */
    async verifyEmail (payload: VerifyAuthEmailRequest): Promise<void> {
      try {
        await apiClient.authApi().v1VerifyAuthEmail(payload)
      } catch (err) {
        return this.errorHandler(err, { 409: 'このメールアドレスはすでに利用されているため使用できません。' })
      }
    },

    /**
     * パスワード更新
     * @param payload
     */
    async updatePassword (payload: UpdateAuthPasswordRequest): Promise<void> {
      try {
        await apiClient.authApi().v1UpdateAuthPassword(payload)
      } catch (err) {
        return this.errorHandler(err)
      }
    },

    /**
     * パスワードリセットの検証
     */
    async forgotPassword (payload: ForgotAuthPasswordRequest): Promise<void> {
      try {
        await apiClient.authApi().v1ForgotAuthPassword(payload)
      } catch (err) {
        return this.errorHandler(err)
      }
    },

    /**
     * パスワードリセット
     * @param payload
     */
    async resetPassword (payload: ResetAuthPasswordRequest): Promise<void> {
      try {
        await apiClient.authApi().v1ResetAuthPassword(payload)
      } catch (err) {
        return this.errorHandler(err)
      }
    },

    /**
     * デバイス情報の登録
     * @param deviceToken デバイスID
     */
    async registerDeviceToken (deviceToken: string): Promise<void> {
      try {
        await apiClient.authApi().v1RegisterAuthDevice({ device: deviceToken })

        const cookies = new Cookies()
        cookies.set('deviceToken', deviceToken, { secure: true })
      } catch (err) {
        return this.errorHandler(err)
      }
    },

    /**
     * 認証情報の更新
     * @param refreshToken リフレッシュトークン
     */
    async getAuthByRefreshToken (refreshToken: string): Promise<void> {
      try {
        const res = await apiClient.authApi().v1RefreshAuthToken({
          refreshToken
        })
        this.setExpiredAt(res.data)
        this.isAuthenticated = true
        this.auth = res.data
        this.auth.refreshToken = refreshToken
      } catch (err) {
        const cookies = new Cookies()
        cookies.remove('refreshToken')
        return this.errorHandler(err)
      }
    },

    /**
     * デバイス情報の取得
     * @returns デバイスID
     */
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

    setExpiredAt (auth: AuthResponse) {
      this.expiredAt = dayjs().add(auth.expiresIn, 'second')
    },

    async logout () {
      try {
        await apiClient.authApi().v1SignOut()
        const cookies = new Cookies()
        cookies.remove('refreshToken')
        this.$reset()
      } catch (error) {
        console.log('APIでエラーが発生しました。', error)
      } finally {
        this.isAuthenticated = false
        this.auth = undefined
        this.user = undefined
        this.expiredAt = undefined
      }
    }
  }
})
