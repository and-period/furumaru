import dayjs, { type Dayjs } from 'dayjs'
import { getToken, isSupported } from 'firebase/messaging'
import { defineStore } from 'pinia'
import Cookies from 'universal-cookie'

import { messaging } from '~/plugins/firebase'
import { apiClient } from '~/plugins/api-client'
import {
  AdminType,
  type AuthResponse,
  type AuthUserResponse,
  type Coordinator,
  type ForgotAuthPasswordRequest,
  type ResetAuthPasswordRequest,
  type Shipping,
  type SignInRequest,
  type UpdateAuthEmailRequest,
  type UpdateAuthPasswordRequest,
  type UpdateCoordinatorRequest,
  type UpsertShippingRequest,
  type VerifyAuthEmailRequest,
} from '~/types/api'
import { useProductTypeStore } from '~/store'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    redirectPath: '/',
    isAuthenticated: false,
    auth: undefined as AuthResponse | undefined,
    user: undefined as AuthUserResponse | undefined,
    coordinator: {} as Coordinator,
    shipping: {} as Shipping,
    expiredAt: undefined as Dayjs | undefined,
  }),

  getters: {
    adminId(state): string {
      return state.auth?.adminId || ''
    },
    accessToken(state): string | undefined {
      return state.auth?.accessToken
    },
    adminType(state): AdminType {
      return state.auth?.type || AdminType.UNKNOWN
    },
  },

  actions: {
    /**
     * サインイン
     * @param payload メールアドレス/パスワード
     * @returns 遷移先Path
     */
    async signIn(payload: SignInRequest): Promise<string> {
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
      }
      catch (err) {
        return this.errorHandler(err, { 401: 'ユーザー名またはパスワードが違います。' })
      }
    },

    /**
     * サインイン中管理者情報取得
     */
    async getUser(): Promise<void> {
      try {
        const res = await apiClient.authApi().v1GetAuthUser()
        this.user = res.data
      }
      catch (err) {
        return this.errorHandler(err, { 401: 'ユーザー名またはパスワードが違います。' })
      }
    },

    /**
     * メールアドレス更新
     * @param payload
     */
    async updateEmail(payload: UpdateAuthEmailRequest): Promise<void> {
      try {
        await apiClient.authApi().v1UpdateAuthEmail(payload)
      }
      catch (err) {
        return this.errorHandler(err, {
          409: 'このメールアドレスはすでに登録されているため、変更できません。',
          412: '変更前のメールアドレスと同じため、変更できません。',
        })
      }
    },

    /**
     * メールアドレス更新後の検証
     * @param payload
     */
    async verifyEmail(payload: VerifyAuthEmailRequest): Promise<void> {
      try {
        await apiClient.authApi().v1VerifyAuthEmail(payload)
      }
      catch (err) {
        return this.errorHandler(err, {
          401: '認証エラーです。再度検証をしてみてください',
          409: 'このメールアドレスはすでに利用されているため使用できません。',
        })
      }
    },

    /**
     * パスワード更新
     * @param payload
     */
    async updatePassword(payload: UpdateAuthPasswordRequest): Promise<void> {
      try {
        await apiClient.authApi().v1UpdateAuthPassword(payload)
      }
      catch (err) {
        return this.errorHandler(err, {
          400: '入力内容に誤りがあります。',
          401: '認証エラーです。再度試してみてください',
        })
      }
    },

    /**
     * パスワードリセットの検証
     */
    async forgotPassword(payload: ForgotAuthPasswordRequest): Promise<void> {
      try {
        await apiClient.authApi().v1ForgotAuthPassword(payload)
      }
      catch (err) {
        return this.errorHandler(err, { 400: '入力内容に誤りがあります。' })
      }
    },

    /**
     * パスワードリセット
     * @param payload
     */
    async resetPassword(payload: ResetAuthPasswordRequest): Promise<void> {
      try {
        await apiClient.authApi().v1ResetAuthPassword(payload)
      }
      catch (err) {
        return this.errorHandler(err, { 400: '入力内容に誤りがあります。' })
      }
    },

    /**
     * デバイス情報の登録
     * @param deviceToken デバイスID
     */
    async registerDeviceToken(deviceToken: string): Promise<void> {
      try {
        await apiClient.authApi().v1RegisterAuthDevice({ device: deviceToken })

        const cookies = new Cookies()
        cookies.set('deviceToken', deviceToken, { secure: true })
      }
      catch (err) {
        return this.errorHandler(err, {
          400: 'デバイス情報の登録に失敗しました。',
          401: '認証エラーです。再度ログインをしてください。',
        })
      }
    },

    /**
     * 認証情報の更新
     * @param refreshToken リフレッシュトークン
     */
    async getAuthByRefreshToken(refreshToken: string): Promise<void> {
      try {
        const res = await apiClient.authApi().v1RefreshAuthToken({
          refreshToken,
        })
        this.setExpiredAt(res.data)
        this.isAuthenticated = true
        this.auth = res.data
        this.auth.refreshToken = refreshToken
      }
      catch (err) {
        const cookies = new Cookies()
        cookies.remove('refreshToken')
        return this.errorHandler(err, { 401: '認証エラーです。再度ログインをしてください。' })
      }
    },

    /**
     * デバイス情報の取得
     * @returns デバイスID
     */
    async getDeviceToken(): Promise<string> {
      const runtimeConfig = useRuntimeConfig()

      const supported = await isSupported()
      if (!supported) {
        console.log('this browser does not support push notificatins.')
        return '' // Push通知未対応ブラウザ
      }

      return await getToken(messaging, {
        vapidKey: runtimeConfig.public.FIREBASE_VAPID_KEY,
      })
        .then((currentToken) => {
          return currentToken
        })
        .catch((err) => {
          console.log('failed to get device token', err)
          return ''
        })
    },

    /**
     * コーディネーターの詳細情報を取得する非同期関数
     */
    async getCoordinator(): Promise<void> {
      try {
        const res = await apiClient.authApi().v1GetAuthCoordinator()

        const productTypeStore = useProductTypeStore()
        this.coordinator = res.data.coordinator
        productTypeStore.productTypes = res.data.productTypes
      }
      catch (err) {
        return this.errorHandler(err, { 404: 'コーディネーター情報が見つかりません。' })
      }
    },

    /**
     * コーディネーターの情報を更新する非同期関数
     * @param payload
     */
    async updateCoordinator(payload: UpdateCoordinatorRequest): Promise<void> {
      try {
        await apiClient.authApi().v1UpdateAuthCoordinator(payload)
      }
      catch (err) {
        return this.errorHandler(err, { 400: '入力内容に誤りがあります。' })
      }
    },

    /**
     * 指定したコーディネーターの配送設定を取得する非同期関数
     * @param coordinatorId
     * @returns
     */
    async fetchShipping(): Promise<void> {
      try {
        const res = await apiClient.authApi().v1GetAuthShipping()
        this.shipping = res.data.shipping
      }
      catch (err) {
        return this.errorHandler(err, { 404: '配送設定が見つかりません。' })
      }
    },

    /**
     * 指定したコーディネーターの配送設定を変更する非同期関数
     * @param payload
     * @returns
     */
    async upsertShipping(payload: UpsertShippingRequest): Promise<void> {
      try {
        await apiClient.authApi().v1UpsertAuthShipping(payload)
      }
      catch (err) {
        return this.errorHandler(err, {
          400: '入力内容に誤りがあります。',
          404: '指定した配送設定が見つかりません。',
        })
      }
    },

    setRedirectPath(payload: string) {
      this.redirectPath = payload
    },

    setExpiredAt(auth: AuthResponse) {
      this.expiredAt = dayjs().add(auth.expiresIn, 'second')
    },

    async logout() {
      try {
        await apiClient.authApi().v1SignOut()
        const cookies = new Cookies()
        cookies.remove('refreshToken')
        this.$reset()
      }
      catch (error) {
        console.log('APIでエラーが発生しました。', error)
      }
      finally {
        this.isAuthenticated = false
        this.auth = undefined
        this.user = undefined
        this.expiredAt = undefined
      }
    },
  },
})
