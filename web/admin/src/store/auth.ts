import dayjs from 'dayjs'
import type { Dayjs } from 'dayjs'
import { getToken, isSupported } from 'firebase/messaging'
import { messaging } from '~/plugins/firebase'
import { useApiClient } from '~/composables/useApiClient'
import { AdminType, AuthApi } from '~/types/api/v1'
import type { AuthProvider, AuthResponse, AuthUserResponse, Coordinator, ForgotAuthPasswordRequest, ResetAuthPasswordRequest, SignInRequest, UpdateAuthEmailRequest, UpdateAuthPasswordRequest, UpdateCoordinatorRequest, V1AuthCoordinatorPatchRequest, V1AuthDevicePostRequest, V1AuthEmailPatchRequest, V1AuthEmailVerifiedPostRequest, V1AuthForgotPasswordPostRequest, V1AuthForgotPasswordVerifiedPostRequest, V1AuthGoogleGetRequest, V1AuthGooglePostRequest, V1AuthLineGetRequest, V1AuthLinePostRequest, V1AuthPasswordPatchRequest, V1AuthPostRequest, V1AuthRefreshTokenPostRequest, VerifyAuthEmailRequest } from '~/types/api/v1'

interface FetchTokenResponse {
  access_token: string
  refresh_token: string
  id_token: string
  token_type: string
  expires_in: number
}

export const useAuthStore = defineStore('auth', () => {
  const { create, errorHandler } = useApiClient()
  const authApi = () => create(AuthApi)

  const redirectPath = ref('/')
  const isAuthenticated = ref(false)
  const auth = ref<AuthResponse | undefined>(undefined)
  const user = ref<AuthUserResponse | undefined>(undefined)
  const providers = ref<AuthProvider[]>([])
  const coordinator = ref<Coordinator>({} as Coordinator)
  const expiredAt = ref<Dayjs | undefined>(undefined)

  // getters
  const accessToken = computed<string | undefined>(() => auth.value?.accessToken)
  const adminId = computed<string>(() => user.value?.id || '')
  const shopIds = computed<string[]>(() => user.value?.shopIds || [])
  const adminType = computed<AdminType>(() => (user.value?.type as AdminType) || AdminType.AdminTypeUnknown)

  /**
   * トークンの保存
   */
  async function setAuth(authResponse: AuthResponse): Promise<string> {
    auth.value = authResponse
    isAuthenticated.value = true

    const refreshTokenExpires = dayjs().add(90, 'days')

    const cookie = useCookie('auth', { secure: true, maxAge: authResponse.expiresIn })
    const refreshToken = useCookie('refreshToken', { secure: true, expires: refreshTokenExpires.toDate() })

    cookie.value = encodeURIComponent(JSON.stringify(authResponse))
    refreshToken.value = authResponse.refreshToken

    await getUser()
    setExpiredAt(authResponse)
    acceptPushNotification()

    return redirectPath.value
  }

  /**
   * サインイン
   * @param payload メールアドレス/パスワード
   * @returns 遷移先Path
   */
  async function signIn(payload: SignInRequest): Promise<string> {
    try {
      const params: V1AuthPostRequest = {
        signInRequest: payload,
      }
      const res = await authApi().v1AuthPost(params)
      return await setAuth(res)
    }
    catch (err) {
      return errorHandler(err, { 401: 'ユーザー名またはパスワードが違います。' })
    }
  }

  /**
   * サインイン with OAuth
   * @param code 認証コード
   * @param redirectUri リダイレクト先URI
   * @returns 遷移先Path
   */
  async function signInWithOAuth(code: string, redirectUri: string): Promise<string> {
    try {
      const token = await fetchOAuthToken(code, redirectUri).catch((err) => {
        console.error('OAuthトークンの取得に失敗しました。', err)
        throw new Error('OAuthトークンの取得に失敗しました。')
      })

      const params: V1AuthRefreshTokenPostRequest = {
        refreshAuthTokenRequest: {
          refreshToken: token.refresh_token,
        },
      }
      const authRes = await authApi().v1AuthRefreshTokenPost(params).catch((err) => {
        console.error('OAuth認証に失敗しました。', err)
        throw new Error('OAuth認証に失敗しました。')
      })

      return await setAuth({ ...authRes, refreshToken: token.refresh_token })
    }
    catch (err) {
      return errorHandler(err, { 401: 'Googleアカウントでのログインに失敗しました。' })
    }
  }

  /**
   * サインイン中管理者情報取得
   */
  async function getUser(): Promise<void> {
    try {
      const res = await authApi().v1AuthUserGet()
      user.value = res
    }
    catch (err) {
      return errorHandler(err, { 401: 'ユーザー名またはパスワードが違います。' })
    }
  }

  /**
   * メールアドレス更新
   * @param payload
   */
  async function updateEmail(payload: UpdateAuthEmailRequest): Promise<void> {
    try {
      const params: V1AuthEmailPatchRequest = {
        updateAuthEmailRequest: payload,
      }
      await authApi().v1AuthEmailPatch(params)
    }
    catch (err) {
      return errorHandler(err, {
        409: 'このメールアドレスはすでに登録されているため、変更できません。',
        412: '変更前のメールアドレスと同じため、変更できません。',
      })
    }
  }

  /**
   * メールアドレス更新後の検証
   * @param payload
   */
  async function verifyEmail(payload: VerifyAuthEmailRequest): Promise<void> {
    try {
      const params: V1AuthEmailVerifiedPostRequest = {
        verifyAuthEmailRequest: payload,
      }
      await authApi().v1AuthEmailVerifiedPost(params)
    }
    catch (err) {
      return errorHandler(err, {
        401: '認証エラーです。再度検証をしてみてください',
        409: 'このメールアドレスはすでに利用されているため使用できません。',
      })
    }
  }

  /**
   * パスワード更新
   * @param payload
   */
  async function updatePassword(payload: UpdateAuthPasswordRequest): Promise<void> {
    try {
      const params: V1AuthPasswordPatchRequest = {
        updateAuthPasswordRequest: payload,
      }
      await authApi().v1AuthPasswordPatch(params)
    }
    catch (err) {
      return errorHandler(err, {
        400: '入力内容に誤りがあります。',
        401: '認証エラーです。再度試してみてください',
      })
    }
  }

  /**
   * パスワードリセットの検証
   */
  async function forgotPassword(payload: ForgotAuthPasswordRequest): Promise<void> {
    try {
      const params: V1AuthForgotPasswordPostRequest = {
        forgotAuthPasswordRequest: payload,
      }
      await authApi().v1AuthForgotPasswordPost(params)
    }
    catch (err) {
      return errorHandler(err, { 400: '入力内容に誤りがあります。' })
    }
  }

  /**
   * パスワードリセット
   * @param payload
   */
  async function resetPassword(payload: ResetAuthPasswordRequest): Promise<void> {
    try {
      const params: V1AuthForgotPasswordVerifiedPostRequest = {
        resetAuthPasswordRequest: payload,
      }
      await authApi().v1AuthForgotPasswordVerifiedPost(params)
    }
    catch (err) {
      return errorHandler(err, { 400: '入力内容に誤りがあります。' })
    }
  }

  /**
   * デバイス情報の登録
   * @param deviceToken デバイスID
   */
  async function registerDeviceToken(deviceToken: string): Promise<void> {
    try {
      const params: V1AuthDevicePostRequest = {
        registerAuthDeviceRequest: {
          device: deviceToken,
        },
      }
      await authApi().v1AuthDevicePost(params)

      const cookie = useCookie('deviceToken', { secure: true })
      cookie.value = deviceToken
    }
    catch (err) {
      return errorHandler(err, {
        400: 'デバイス情報の登録に失敗しました。',
        401: '認証エラーです。再度ログインをしてください。',
      })
    }
  }

  /**
   * 認証情報の更新
   * @param refreshToken リフレッシュトークン
   */
  async function getAuthByRefreshToken(refreshTokenStr: string): Promise<void> {
    try {
      const params: V1AuthRefreshTokenPostRequest = {
        refreshAuthTokenRequest: {
          refreshToken: refreshTokenStr,
        },
      }
      const res = await authApi().v1AuthRefreshTokenPost(params)
      setExpiredAt(res)
      isAuthenticated.value = true
      auth.value = res
      auth.value.refreshToken = refreshTokenStr
    }
    catch (err) {
      const authCookie = useCookie('auth', { secure: true })
      const refreshTokenCookie = useCookie('refreshToken', { secure: true })
      authCookie.value = undefined
      refreshTokenCookie.value = undefined
      return errorHandler(err, { 401: '認証エラーです。再度ログインをしてください。' })
    }
  }

  /**
   * デバイス情報の取得
   * @returns デバイスID
   */
  async function getDeviceToken(): Promise<string> {
    const runtimeConfig = useRuntimeConfig()

    const supported = await isSupported()
    if (!supported) {
      console.log('this browser does not support push notificatins.')
      return '' // Push通知未対応ブラウザ
    }

    // messaging が null の場合は空文字を返す
    if (!messaging) {
      console.log('Firebase Messaging is not initialized')
      return ''
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
  }

  /**
   * コーディネーターの詳細情報を取得する非同期関数
   */
  async function getCoordinator(): Promise<void> {
    try {
      const res = await authApi().v1AuthCoordinatorGet()
      coordinator.value = res.coordinator
    }
    catch (err) {
      return errorHandler(err, { 404: 'コーディネーター情報が見つかりません。' })
    }
  }

  /**
   * コーディネーターの情報を更新する非同期関数
   * @param payload
   */
  async function updateCoordinator(payload: UpdateCoordinatorRequest): Promise<void> {
    try {
      const params: V1AuthCoordinatorPatchRequest = {
        updateCoordinatorRequest: payload,
      }
      await authApi().v1AuthCoordinatorPatch(params)
    }
    catch (err) {
      return errorHandler(err, { 400: '入力内容に誤りがあります。' })
    }
  }

  /**
   * OAuth認証によるトークン発行
   * @param code 認証コード
   * @param redirectUri リダイレクト先URI
   * @returns
   */
  async function fetchOAuthToken(code: string, redirectUri: string): Promise<FetchTokenResponse> {
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
  }

  /**
   * Push通知の許可設定
   * @returns
   */
  async function acceptPushNotification(): Promise<void> {
    const supported = await isSupported()
    if (!supported) {
      console.log('this browser does not support push notificatins.')
      return // Push通知未対応ブラウザ
    }

    getDeviceToken()
      .then((deviceToken) => {
        if (deviceToken === '') {
          return // Push通知が未有効化状態
        }
        const cookie = useCookie('deviceToken', { secure: true })
        if (cookie.value === deviceToken) {
          return // API側へ登録済み
        }
        return registerDeviceToken(deviceToken)
      })
      .catch((err) => {
        console.log('push notifications are disabled.', err)
      })
  }

  /**
   * 認証済みプロバイダ一覧取得
   * @returns
   */
  async function listAuthProviders(): Promise<void> {
    try {
      const res = await authApi().v1AuthProvidersGet()
      providers.value = res.providers
    }
    catch (err) {
      return errorHandler(err, { 400: '入力内容に誤りがあります。' })
    }
  }

  /**
   * Google連携 - 認証ページへの遷移URL取得
   * @param state ランダム文字列
   * @param redirectUri リダイレクト先URI
   * @returns
   */
  async function getAuthGoogleUrl(state: string, redirectUri?: string): Promise<string> {
    try {
      const params: V1AuthGoogleGetRequest = {
        state,
        redirectUri,
      }
      const res = await authApi().v1AuthGoogleGet(params)
      return res.url
    }
    catch (err) {
      return errorHandler(err, { 400: '入力内容に誤りがあります。' })
    }
  }

  /**
   * Google連携 - アカウントの連携
   * @param code 認証コード
   * @param nonce ランダム文字列
   * @returns
   */
  async function linkGoogleAccount(code: string, nonce: string, redirectUri?: string): Promise<void> {
    try {
      const params: V1AuthGooglePostRequest = {
        connectGoogleAccountRequest: {
          code,
          nonce,
          redirectUri: redirectUri || '',
        },
      }
      await authApi().v1AuthGooglePost(params)
    }
    catch (err) {
      return errorHandler(err, { 400: '入力内容に誤りがあります。' })
    }
  }

  /**
   * LINE連携 - 認証ページへの遷移URL取得
   * @param state ランダム文字列
   * @param redirectUri リダイレクト先URI
   * @returns
   */
  async function getAuthLineUrl(state: string, redirectUri?: string): Promise<string> {
    try {
      const params: V1AuthLineGetRequest = {
        state,
        redirectUri,
      }
      const res = await authApi().v1AuthLineGet(params)
      return res.url
    }
    catch (err) {
      return errorHandler(err, { 400: '入力内容に誤りがあります。' })
    }
  }

  /**
   * LINE連携 - アカウントの連携
   * @param code 認証コード
   * @param nonce ランダム文字列
   * @returns
   */
  async function linkLineAccount(code: string, nonce: string, redirectUri?: string): Promise<void> {
    try {
      const params: V1AuthLinePostRequest = {
        connectLINEAccountRequest: {
          code,
          nonce,
          redirectUri: redirectUri || '',
        },
      }
      await authApi().v1AuthLinePost(params)
    }
    catch (err) {
      return errorHandler(err, { 400: '入力内容に誤りがあります。' })
    }
  }

  function setRedirectPath(payload: string) {
    redirectPath.value = payload
  }

  function setExpiredAt(authResponse: AuthResponse) {
    expiredAt.value = dayjs().add(authResponse.expiresIn, 'second')
  }

  async function logout() {
    try {
      await authApi().v1AuthDelete()

      const authCookie = useCookie('auth', { secure: true })
      const refreshTokenCookie = useCookie('refreshToken', { secure: true })
      authCookie.value = undefined
      refreshTokenCookie.value = undefined

      $reset()
    }
    catch (error) {
      console.log('APIでエラーが発生しました。', error)
    }
    finally {
      isAuthenticated.value = false
      auth.value = undefined
      user.value = undefined
      expiredAt.value = undefined
    }
  }

  function $reset() {
    redirectPath.value = '/'
    isAuthenticated.value = false
    auth.value = undefined
    user.value = undefined
    providers.value = []
    coordinator.value = {} as Coordinator
    expiredAt.value = undefined
  }

  return {
    // state
    redirectPath,
    isAuthenticated,
    auth,
    user,
    providers,
    coordinator,
    expiredAt,
    // getters
    accessToken,
    adminId,
    shopIds,
    adminType,
    // actions
    setAuth,
    signIn,
    signInWithOAuth,
    getUser,
    updateEmail,
    verifyEmail,
    updatePassword,
    forgotPassword,
    resetPassword,
    registerDeviceToken,
    getAuthByRefreshToken,
    getDeviceToken,
    getCoordinator,
    updateCoordinator,
    fetchOAuthToken,
    acceptPushNotification,
    listAuthProviders,
    getAuthGoogleUrl,
    linkGoogleAccount,
    getAuthLineUrl,
    linkLineAccount,
    setRedirectPath,
    setExpiredAt,
    logout,
    $reset,
  }
})
