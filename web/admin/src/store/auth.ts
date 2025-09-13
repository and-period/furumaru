import dayjs from 'dayjs'
import type { Dayjs } from 'dayjs'
import { getToken, isSupported } from 'firebase/messaging'
import { messaging } from '~/plugins/firebase'
import { AdminType } from '~/types/api/v1'
import type { AuthProvider, AuthResponse, AuthUserResponse, ConnectGoogleAccountRequest, ConnectLINEAccountRequest, Coordinator, ForgotAuthPasswordRequest, ResetAuthPasswordRequest, Shipping, SignInRequest, UpdateAuthEmailRequest, UpdateAuthPasswordRequest, UpdateCoordinatorRequest, UpsertShippingRequest, V1AuthCoordinatorPatchRequest, V1AuthCoordinatorShippingsPatchRequest, V1AuthDevicePostRequest, V1AuthEmailPatchRequest, V1AuthEmailVerifiedPostRequest, V1AuthForgotPasswordPostRequest, V1AuthForgotPasswordVerifiedPostRequest, V1AuthGoogleGetRequest, V1AuthGooglePostRequest, V1AuthLineGetRequest, V1AuthLinePostRequest, V1AuthPasswordPatchRequest, V1AuthPostRequest, V1AuthRefreshTokenPostRequest, V1CoordinatorsGetRequest, VerifyAuthEmailRequest } from '~/types/api/v1'
import { useProductTypeStore } from '~/store'

interface FetchTokenResponse {
  access_token: string
  refresh_token: string
  id_token: string
  token_type: string
  expires_in: number
}

export const useAuthStore = defineStore('auth', {
  state: () => ({
    redirectPath: '/',
    isAuthenticated: false,
    auth: undefined as AuthResponse | undefined,
    user: undefined as AuthUserResponse | undefined,
    providers: [] as AuthProvider[],
    coordinator: {} as Coordinator,
    shipping: {} as Shipping,
    expiredAt: undefined as Dayjs | undefined,
  }),

  getters: {
    accessToken(state): string | undefined {
      return state.auth?.accessToken
    },
    adminId(state): string {
      return state.user?.id || ''
    },
    shopIds(state): string[] {
      return state.user?.shopIds || []
    },
    adminType(state): AdminType {
      return state.user?.type as AdminType || AdminType.AdminTypeUnknown
    },
  },

  actions: {
    /**
     * トークンの保存
     */
    async setAuth(auth: AuthResponse): Promise<string> {
      this.auth = auth
      this.isAuthenticated = true

      const refreshTokenExpires = dayjs().add(90, 'days')

      const cookie = useCookie('auth', { secure: true, maxAge: auth.expiresIn })
      const refreshToken = useCookie('refreshToken', { secure: true, expires: refreshTokenExpires.toDate() })

      cookie.value = encodeURIComponent(JSON.stringify(auth))
      refreshToken.value = auth.refreshToken

      await this.getUser()
      this.setExpiredAt(auth)
      this.acceptPushNotification()

      return this.redirectPath
    },

    /**
     * サインイン
     * @param payload メールアドレス/パスワード
     * @returns 遷移先Path
     */
    async signIn(payload: SignInRequest): Promise<string> {
      try {
        const params: V1AuthPostRequest = {
          signInRequest: payload,
        }
        const res = await this.authApi().v1AuthPost(params)
        return await this.setAuth(res)
      }
      catch (err) {
        return this.errorHandler(err, { 401: 'ユーザー名またはパスワードが違います。' })
      }
    },

    /**
     * サインイン with OAuth
     * @param code 認証コード
     * @param redirectUri リダイレクト先URI
     * @returns 遷移先Path
     */
    async signInWithOAuth(code: string, redirectUri: string): Promise<string> {
      try {
        const token = await this.fetchOAuthToken(code, redirectUri).catch((err) => {
          console.error('OAuthトークンの取得に失敗しました。', err)
          throw new Error('OAuthトークンの取得に失敗しました。')
        })

        const params: V1AuthRefreshTokenPostRequest = {
          refreshAuthTokenRequest: {
            refreshToken: token.refresh_token,
          },
        }
        const auth = await this.authApi().v1AuthRefreshTokenPost(params).catch((err) => {
          console.error('OAuth認証に失敗しました。', err)
          throw new Error('OAuth認証に失敗しました。')
        })

        return await this.setAuth({ ...auth, refreshToken: token.refresh_token })
      }
      catch (err) {
        return this.errorHandler(err, { 401: 'Googleアカウントでのログインに失敗しました。' })
      }
    },

    /**
     * サインイン中管理者情報取得
     */
    async getUser(): Promise<void> {
      try {
        const res = await this.authApi().v1AuthUserGet()
        this.user = res
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
        const params: V1AuthEmailPatchRequest = {
          updateAuthEmailRequest: payload,
        }
        await this.authApi().v1AuthEmailPatch(params)
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
        const params: V1AuthEmailVerifiedPostRequest = {
          verifyAuthEmailRequest: payload,
        }
        await this.authApi().v1AuthEmailVerifiedPost(params)
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
        const params: V1AuthPasswordPatchRequest = {
          updateAuthPasswordRequest: payload,
        }
        await this.authApi().v1AuthPasswordPatch(params)
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
        const params: V1AuthForgotPasswordPostRequest = {
          forgotAuthPasswordRequest: payload,
        }
        await this.authApi().v1AuthForgotPasswordPost(params)
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
        const params: V1AuthForgotPasswordVerifiedPostRequest = {
          resetAuthPasswordRequest: payload,
        }
        await this.authApi().v1AuthForgotPasswordVerifiedPost(params)
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
        const params: V1AuthDevicePostRequest = {
          registerAuthDeviceRequest: {
            device: deviceToken,
          },
        }
        await this.authApi().v1AuthDevicePost(params)

        const cookie = useCookie('deviceToken', { secure: true })
        cookie.value = deviceToken
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
        const params: V1AuthRefreshTokenPostRequest = {
          refreshAuthTokenRequest: {
            refreshToken,
          },
        }
        const res = await this.authApi().v1AuthRefreshTokenPost(params)
        this.setExpiredAt(res)
        this.isAuthenticated = true
        this.auth = res
        this.auth.refreshToken = refreshToken
      }
      catch (err) {
        const auth = useCookie('auth', { secure: true })
        const refreshToken = useCookie('refreshToken', { secure: true })
        auth.value = undefined
        refreshToken.value = undefined
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
        const res = await this.authApi().v1AuthCoordinatorGet()

        const productTypeStore = useProductTypeStore()
        this.coordinator = res.coordinator
        productTypeStore.productTypes = res.productTypes
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
        const params: V1AuthCoordinatorPatchRequest = {
          updateCoordinatorRequest: payload,
        }
        await this.authApi().v1AuthCoordinatorPatch(params)
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
        const res = await this.authApi().v1AuthCoordinatorShippingsGet()
        this.shipping = res.shipping
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
        const params: V1AuthCoordinatorShippingsPatchRequest = {
          upsertShippingRequest: payload,
        }
        await this.authApi().v1AuthCoordinatorShippingsPatch(params)
      }
      catch (err) {
        return this.errorHandler(err, {
          400: '入力内容に誤りがあります。',
          404: '指定した配送設定が見つかりません。',
        })
      }
    },

    /**
     * OAuth認証によるトークン発行
     * @param code 認証コード
     * @param redirectUri リダイレクト先URI
     * @returns
     */
    async fetchOAuthToken(code: string, redirectUri: string): Promise<FetchTokenResponse> {
      if (code === '' || redirectUri === '') {
        throw new Error('code or redirectUri is empty.')
      }

      const config = useRuntimeConfig()

      const params = new URLSearchParams({
        grant_type: 'authorization_code',
        client_id: config.public.COGNITO_CLIENT_ID || '',
        redirect_uri: redirectUri,
        code,
      })

      const res = await $fetch<FetchTokenResponse>(
        `https://${config.public.COGNITO_AUTH_DOMAIN}/oauth2/token`,
        {
          method: 'POST',
          headers: {
            'Content-Type': 'application/x-www-form-urlencoded',
          },
          body: params.toString(),
        },
      )

      return res
    },

    /**
     * Push通知の許可設定
     * @returns
     */
    async acceptPushNotification(): Promise<void> {
      const supported = await isSupported()
      if (!supported) {
        console.log('this browser does not support push notificatins.')
        return // Push通知未対応ブラウザ
      }

      this.getDeviceToken()
        .then((deviceToken) => {
          if (deviceToken === '') {
            return // Push通知が未有効化状態
          }
          const cookie = useCookie('deviceToken', { secure: true })
          if (cookie.value === deviceToken) {
            return // API側へ登録済み
          }
          return this.registerDeviceToken(deviceToken)
        })
        .catch((err) => {
          console.log('push notifications are disabled.', err)
        })
    },

    /**
     * 認証済みプロバイダ一覧取得
     * @returns
     */
    async listAuthProviders(): Promise<void> {
      try {
        const res = await this.authApi().v1AuthProvidersGet()
        this.providers = res.providers
      }
      catch (err) {
        return this.errorHandler(err, { 400: '入力内容に誤りがあります。' })
      }
    },

    /**
     * Google連携 - 認証ページへの遷移URL取得
     * @param state ランダム文字列
     * @param redirectUri リダイレクト先URI
     * @returns
     */
    async getAuthGoogleUrl(state: string, redirectUri?: string): Promise<string> {
      try {
        const params: V1AuthGoogleGetRequest = {
          state,
          redirectUri,
        }
        const res = await this.authApi().v1AuthGoogleGet(params)
        return res.url
      }
      catch (err) {
        return this.errorHandler(err, { 400: '入力内容に誤りがあります。' })
      }
    },

    /**
     * Google連携 - アカウントの連携
     * @param code 認証コード
     * @param nonce ランダム文字列
     * @returns
     */
    async linkGoogleAccount(code: string, nonce: string, redirectUri?: string): Promise<void> {
      try {
        const params: V1AuthGooglePostRequest = {
          connectGoogleAccountRequest: {
            code,
            nonce,
            redirectUri: redirectUri || '',
          },
        }
        await this.authApi().v1AuthGooglePost(params)
      }
      catch (err) {
        return this.errorHandler(err, { 400: '入力内容に誤りがあります。' })
      }
    },

    /**
     * LINE連携 - 認証ページへの遷移URL取得
     * @param state ランダム文字列
     * @param redirectUri リダイレクト先URI
     * @returns
     */
    async getAuthLineUrl(state: string, redirectUri?: string): Promise<string> {
      try {
        const params: V1AuthLineGetRequest = {
          state,
          redirectUri,
        }
        const res = await this.authApi().v1AuthLineGet(params)
        return res.url
      }
      catch (err) {
        return this.errorHandler(err, { 400: '入力内容に誤りがあります。' })
      }
    },

    /**
     * LINE連携 - アカウントの連携
     * @param code 認証コード
     * @param nonce ランダム文字列
     * @returns
     */
    async linkLineAccount(code: string, nonce: string, redirectUri?: string): Promise<void> {
      try {
        const params: V1AuthLinePostRequest = {
          connectLINEAccountRequest: {
            code,
            nonce,
            redirectUri: redirectUri || '',
          },
        }
        await this.authApi().v1AuthLinePost(params)
      }
      catch (err) {
        return this.errorHandler(err, { 400: '入力内容に誤りがあります。' })
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
        await this.authApi().v1AuthDelete()

        const auth = useCookie('auth', { secure: true })
        const refreshToken = useCookie('refreshToken', { secure: true })
        auth.value = undefined
        refreshToken.value = undefined

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
