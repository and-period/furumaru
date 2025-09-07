import { defineStore } from 'pinia';
import ApiClientFactory from '@/plugins/helper/factory';
import {
  AuthUserApi,
  OtherApi,
  type AuthUserResponse,
  type CreateAuthUserRequest,
  type CreateAuthUserResponse,
  type UpdateAuthUserEmailRequest,
  type UpdateAuthUserUsernameRequest,
  type UpdateAuthUserAccountIdRequest,
  type VerifyAuthUserRequest,
  type VerifyAuthUserEmailRequest,
  UploadStatus,
} from '@/types/api';
import { useAuthStore } from '@/stores/auth';

interface UserState {
  user: AuthUserResponse | null;
  isLoading: boolean;
  error: string | null;
}

/**
 * ユーザー登録・プロフィール関連のAPIを扱うストア (LIFF用)
 */
export const useUserStore = defineStore('user', {
  state: (): UserState => ({
    user: null,
    isLoading: false,
    error: null,
  }),

  actions: {
    // 内部ユーティリティ: APIクライアント生成
    authUserApiClient(token?: string) {
      const factory = new ApiClientFactory();
      return factory.create<AuthUserApi>(AuthUserApi, token);
    },

    otherApiClient(token?: string) {
      const factory = new ApiClientFactory();
      return factory.create<OtherApi>(OtherApi, token);
    },

    // 購入者登録 (メール/SMS認証)
    async signUp(payload: CreateAuthUserRequest): Promise<CreateAuthUserResponse> {
      this.isLoading = true;
      this.error = null;
      try {
        const api = this.authUserApiClient();
        const res = await api.v1CreateAuthUser({ body: payload });
        return res;
      }
      catch (e) {
        const message = e instanceof Error ? e.message : 'Failed to sign up';
        this.error = message;
        console.error('User signUp failed:', e);
        throw e;
      }
      finally {
        this.isLoading = false;
      }
    },

    // 購入者登録 - コード検証 (メール/SMS認証)
    async verifyAuth(payload: VerifyAuthUserRequest): Promise<void> {
      this.isLoading = true;
      this.error = null;
      try {
        const api = this.authUserApiClient();
        await api.v1VerifyAuthUser({ body: payload });
      }
      catch (e) {
        const message = e instanceof Error ? e.message : 'Failed to verify user';
        this.error = message;
        console.error('User verifyAuth failed:', e);
        throw e;
      }
      finally {
        this.isLoading = false;
      }
    },

    // メールアドレス更新 - コード検証
    async verifyEmail(payload: VerifyAuthUserEmailRequest): Promise<void> {
      this.isLoading = true;
      this.error = null;
      try {
        const token = useAuthStore().token?.accessToken;
        const api = this.authUserApiClient(token);
        await api.v1VerifyAuthUserEmail({ body: payload });
      }
      catch (e) {
        const message = e instanceof Error ? e.message : 'Failed to verify email';
        this.error = message;
        console.error('User verifyEmail failed:', e);
        throw e;
      }
      finally {
        this.isLoading = false;
      }
    },

    // 自ユーザー取得
    async fetchMe(): Promise<AuthUserResponse | null> {
      this.isLoading = true;
      this.error = null;
      try {
        const token = useAuthStore().token?.accessToken;
        if (!token) {
          // 未認証時はnullで返す
          this.user = null;
          return null;
        }
        const api = this.authUserApiClient(token);
        const res = await api.v1GetAuthUser();
        this.user = res;
        return res;
      }
      catch (e) {
        const message = e instanceof Error ? e.message : 'Failed to fetch user';
        this.error = message;
        console.error('User fetchMe failed:', e);
        throw e;
      }
      finally {
        this.isLoading = false;
      }
    },

    async updateUsername(username: string): Promise<void> {
      this.isLoading = true;
      this.error = null;
      try {
        const token = useAuthStore().token?.accessToken;
        if (!token) throw new Error('Not authenticated');
        const api = this.authUserApiClient(token);
        await api.v1UpdateAuthUserUsername({ body: { username } as UpdateAuthUserUsernameRequest });
        await this.fetchMe();
      }
      catch (e) {
        const message = e instanceof Error ? e.message : 'Failed to update username';
        this.error = message;
        console.error('User updateUsername failed:', e);
        throw e;
      }
      finally {
        this.isLoading = false;
      }
    },

    async updateAccountId(accountId: string): Promise<void> {
      this.isLoading = true;
      this.error = null;
      try {
        const token = useAuthStore().token?.accessToken;
        if (!token) throw new Error('Not authenticated');
        const api = this.authUserApiClient(token);
        await api.v1UpdateAuthUserAccountId({ body: { accountId } as UpdateAuthUserAccountIdRequest });
        await this.fetchMe();
      }
      catch (e) {
        const message = e instanceof Error ? e.message : 'Failed to update account id';
        this.error = message;
        console.error('User updateAccountId failed:', e);
        throw e;
      }
      finally {
        this.isLoading = false;
      }
    },

    async updateEmail(email: string): Promise<void> {
      this.isLoading = true;
      this.error = null;
      try {
        const token = useAuthStore().token?.accessToken;
        if (!token) throw new Error('Not authenticated');
        const api = this.authUserApiClient(token);
        await api.v1UpdateAuthUserEmail({ body: { email } as UpdateAuthUserEmailRequest });
      }
      catch (e) {
        const message = e instanceof Error ? e.message : 'Failed to update email';
        this.error = message;
        console.error('User updateEmail failed:', e);
        throw e;
      }
      finally {
        this.isLoading = false;
      }
    },

    async updateNotificationEnabled(enabled: boolean): Promise<void> {
      this.isLoading = true;
      this.error = null;
      try {
        const token = useAuthStore().token?.accessToken;
        if (!token) throw new Error('Not authenticated');
        const api = this.authUserApiClient(token);
        await api.v1UpdateAuthUserNotification({ body: { enabled } });
        await this.fetchMe();
      }
      catch (e) {
        const message = e instanceof Error ? e.message : 'Failed to update notification settings';
        this.error = message;
        console.error('User updateNotificationEnabled failed:', e);
        throw e;
      }
      finally {
        this.isLoading = false;
      }
    },

    // サムネイル更新（アップロード〜反映まで）
    async updateThumbnail(file: File): Promise<void> {
      this.isLoading = true;
      this.error = null;
      try {
        const token = useAuthStore().token?.accessToken;
        if (!token) throw new Error('Not authenticated');

        const mimeType = file.type;
        const authApi = this.authUserApiClient(token);
        const { key, url: uploadUrl } = await authApi.v1GetUserThumbnailUploadUrl({
          body: { fileType: mimeType },
        });

        // 署名付きURLにアップロード
        await fetch(uploadUrl, {
          method: 'PUT',
          body: file,
          headers: {
            'Content-Type': file.type,
          },
        });

        // アップロードの完了をポーリング
        const otherApi = this.otherApiClient(token);
        let url: string | undefined;
        // eslint-disable-next-line no-constant-condition
        while (true) {
          await new Promise(resolve => setTimeout(resolve, 200));
          const state = await otherApi.v1GetUploadState({ key });
          if (state.status === UploadStatus.SUCEEDED) {
            url = state.url;
            break;
          }
          if (state.status === UploadStatus.FAILED) {
            throw new Error('ファイルのアップロードに失敗しました。');
          }
        }

        if (!url) throw new Error('アップロードURLの取得に失敗しました。');

        await authApi.v1UpdateAuthUserThumbnail({ body: { thumbnailUrl: url } });
        await this.fetchMe();
      }
      catch (e) {
        const message = e instanceof Error ? e.message : 'Failed to update thumbnail';
        this.error = message;
        console.error('User updateThumbnail failed:', e);
        throw e;
      }
      finally {
        this.isLoading = false;
      }
    },

    // 退会
    async deleteMe(): Promise<void> {
      this.isLoading = true;
      this.error = null;
      try {
        const token = useAuthStore().token?.accessToken;
        if (!token) throw new Error('Not authenticated');
        const api = this.authUserApiClient(token);
        await api.v1DeleteAuthUser();
        this.user = null;
      }
      catch (e) {
        const message = e instanceof Error ? e.message : 'Failed to delete user';
        this.error = message;
        console.error('User deleteMe failed:', e);
        throw e;
      }
      finally {
        this.isLoading = false;
      }
    },
  },
});

