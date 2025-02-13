// docs: https://pinia.vuejs.org/core-concepts/#option-stores
import type { Dayjs } from 'dayjs'
import dayjs from 'dayjs'
import { defineStore, acceptHMRUpdate } from 'pinia'
import { fileUpload } from './helpter'
import type {
  AuthUserResponse,
  CreateAuthUserRequest,
  CreateAuthUserResponse,
  SignInRequest,
  UpdateAuthPasswordRequest,
  VerifyAuthUserRequest,
} from '~/types/api'
import { AuthError } from '~/types/exception'

interface FetchTokenResponse {
  access_token: string
  refresh_token: string
  id_token: string
  token_type: string
  expires_in: number
}

/**
 * 認証情報を管理するグローバルステート
 */
export const useAuthStore = defineStore('auth', {
  // 永続化の設定
  persist: {
    storage: persistedState.cookiesWithOptions({
      sameSite: 'strict',
    }),
  },

  state: () => {
    return {
      isAuthenticated: false,
      accessToken: '',
      refreshToken: '',
      expiredAt: undefined as Dayjs | undefined,
      user: undefined as AuthUserResponse | undefined,
    }
  },

  actions: {
    /**
     * メール認証ログインを実施する非同期関数
     * @param payload
     * @returns
     */
    async signIn(payload: SignInRequest) {
      try {
        const res = await this.authApiClient().v1SignIn({ body: payload })
        this.isAuthenticated = true
        this.accessToken = res.accessToken
        this.refreshToken = res.refreshToken
        this.setExpiredAt(res.expiresIn)
        await this.fetchUserInfo()
      }
      catch (error) {
        return this.errorHandler(error, {
          401: this.i18n.t('auth.signIn.authErrorMessage'),
        })
      }
    },

    /**
     * OAuthログインを実施する非同期関数
     * @param code OAuth認証コード
     * @param redirectUri リダイレクト先URI
     * @returns 遷移先Path
     */
    async signInWithOAuth(code: string, redirectUri: string): Promise<void> {
      try {
        const token = await this.fetchOAuthToken(code, redirectUri).catch((err) => {
          console.error('OAuthトークンの取得に失敗しました。', err)
          throw new Error('OAuthトークンの取得に失敗しました。')
        })

        this.isAuthenticated = true
        this.accessToken = token.access_token
        this.refreshToken = token.refresh_token
        this.setExpiredAt(token.expires_in)
        await this.fetchUserInfo()
      }
      catch (error) {
        return this.errorHandler(error, {
          401: this.i18n.t('auth.signIn.authErrorMessage'),
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
        client_id: config.public.COGNITO_CLIENT_ID as string,
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

    async signUp(
      payload: CreateAuthUserRequest,
    ): Promise<CreateAuthUserResponse> {
      try {
        const res = await this.authUserApiClient().v1CreateAuthUser({
          body: payload,
        })
        return res
      }
      catch (error) {
        return this.errorHandler(error, {
          409: '指定したメールアドレスはご利用できません。',
        })
      }
    },

    async verifyAuth(payload: VerifyAuthUserRequest) {
      try {
        await this.authUserApiClient().v1VerifyAuthUser({ body: payload })
      }
      catch (error) {
        return this.errorHandler(error)
      }
    },

    async fetchUserInfo() {
      const res = await this.authUserApiClient(this.accessToken).v1GetAuthUser()
      this.user = res
    },

    /**
     * フロントエンドの認証情報をリセットする関数
     */
    resetState() {
      this.$reset()
    },

    /**
     * ログアウト処理
     */
    async logout() {
      try {
        await this.authApiClient(this.accessToken).v1SignOut()
      }
      catch (error) {
        this.errorHandler(error)
      }
      finally {
        // stateを初期状態にリセット
        this.resetState()
      }
    },

    async refreshAccsessToken(refreshToken: string) {
      if (!refreshToken) {
        console.debug('リフレッシュトークンが存在しません')
        this.$reset()
        return Promise.reject(
          new AuthError('リフレッシュトークンが存在しません'),
        )
      }
      try {
        const res = await this.authApiClient().v1RefreshAuthToken({
          body: { refreshToken },
        })
        this.accessToken = res.accessToken
        this.refreshToken = res.refreshToken
        this.setExpiredAt(res.expiresIn)
      }
      catch (error) {
        return this.errorHandler(error)
      }
    },

    setExpiredAt(expiredAt: number) {
      this.expiredAt = dayjs().add(expiredAt, 'second')
    },

    /**
     * サムネイル変更
     * @param file
     * @returns
     */
    async updateThumbnail(file: File) {
      const mimeType = file.type

      const { key, url: uploadUrl } = await this.authUserApiClient(
        this.accessToken,
      ).v1GetUserThumbnailUploadUrl({
        body: { fileType: mimeType },
      })

      const url = await fileUpload(file, key, uploadUrl, this.accessToken)

      this.authUserApiClient(this.accessToken).v1UpdateAuthUserThumbnail({
        body: { thumbnailUrl: url },
      })
    },

    /**
     * ユーザー名変更
     * @param username
     * @returns
     */
    async updateUsername(username: string) {
      try {
        await this.authUserApiClient(this.accessToken).v1UpdateAuthUserUsername(
          {
            body: { username },
          },
        )
      }
      catch (error) {
        return this.errorHandler(error)
      }
    },

    /**
     * アカウントID（ユーザーID）変更
     * @param accountId
     */
    async updateAccountId(accountId: string) {
      try {
        await this.authUserApiClient(
          this.accessToken,
        ).v1UpdateAuthUserAccountId({
          body: { accountId },
        })
      }
      catch (error) {
        return this.errorHandler(error)
      }
    },

    /**
     * パスワード変更
     * @param payload
     * @returns
     */
    async updatePassword(payload: UpdateAuthPasswordRequest) {
      try {
        await this.authApiClient(this.accessToken).v1UpdateUserPassword({
          body: payload,
        })
      }
      catch (error) {
        return this.errorHandler(error, {
          401: '現在のパスワードが正しくありません。',
        })
      }
    },

    /**
     * メールアドレス更新
     * @param email
     */
    async updateEmail(email: string) {
      try {
        await this.authUserApiClient(this.accessToken).v1UpdateAuthUserEmail({
          body: { email },
        })
      }
      catch (error) {
        return this.errorHandler(error)
      }
    },

    /**
     * メール通知設定
     * @param enabled
     * @returns
     */
    async updateNotificationEnabled(enabled: boolean) {
      try {
        await this.authUserApiClient(
          this.accessToken,
        ).v1UpdateAuthUserNotification({ body: { enabled } })
      }
      catch (error) {
        return this.errorHandler(error)
      }
    },
  },
})

// ホットリロードを有効にする
if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(useAuthStore, import.meta.hot))
}
