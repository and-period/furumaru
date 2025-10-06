import { defineStore } from 'pinia';
import ApiClientFactory from '@/plugins/helper/factory';
import type { AuthApi, AuthResponse, SignInRequest, CreateAuthUserRequest } from '@/types/api/facility';
import { AuthUserApi } from '@/types/api/facility';
import { buildApiErrorMessage } from '@/plugins/helper/error';

interface AuthState {
  isLoading: boolean;
  error: string | null;
  isAuthenticated: boolean;
  token: Pick<AuthResponse, 'accessToken' | 'refreshToken' | 'expiresIn' | 'tokenType' | 'userId'> | null;
  // アクセストークンの有効期限（UNIXエポックミリ秒）
  expiresAt: number | null;
  // チェックインの登録が必要かどうかのフラグ
  checkInRequired: boolean;
}

/**
 * 認証を管理するグローバルステート (LIFF用)
 * - src/types/api/facility/apis/AuthApi.ts を利用
 */
export const useAuthStore = defineStore('auth', {
  state: (): AuthState => ({
    isLoading: false,
    error: null,
    isAuthenticated: false,
    token: null,
    expiresAt: null,
    checkInRequired: true,
  }),

  getters: {
    // トークンの有効期限切れ判定
    isTokenExpired: (state): boolean => {
      if (!state.expiresAt) {
        return true;
      }
      return Date.now() >= state.expiresAt;
    },
  },

  actions: {
    async signIn(facilityId: string, liffIDToken: string) {
      this.isLoading = true;
      this.error = null;

      try {
        const api = this.facilityAuthApiClient();

        const payload: SignInRequest = {
          authToken: liffIDToken,
        };

        const res = await api.facilitiesFacilityIdAuthPost({
          facilityId,
          signInRequest: payload,
        });

        // ステート更新（サーバー側Cookieでも管理される前提だが、UI用に保持）
        this.token = {
          accessToken: res.accessToken,
          refreshToken: res.refreshToken,
          expiresIn: res.expiresIn,
          tokenType: res.tokenType,
          userId: res.userId,
        };
        // サーバー応答の expiresIn（秒）から期限時刻（ミリ秒）を算出
        this.expiresAt = Date.now() + (res.expiresIn * 1000);
        this.isAuthenticated = Boolean(res.accessToken || res.userId);
        return res;
      }
      catch (e) {
        const message = e instanceof Error ? e.message : 'Failed to sign in';
        this.error = message;
        this.isAuthenticated = false;
        this.token = null;
        this.expiresAt = null;
        console.error('Auth signIn failed:', e);
        throw e;
      }
      finally {
        this.isLoading = false;
      }
    },

    // チェックイン（ユーザー登録）
    async registerUser(facilityId: string, payload: CreateAuthUserRequest) {
      this.isLoading = true;
      this.error = null;

      try {
        const factory = new ApiClientFactory();
        const api = factory.createFacility<AuthUserApi>(AuthUserApi);

        await api.facilitiesFacilityIdUsersPost({
          facilityId,
          createAuthUserRequest: payload,
        });

        // チェックイン完了後はフラグをオフ
        this.setCheckInRequired(false);
        return true;
      }
      catch (e) {
        console.error('Auth registerUser failed:', e);
        this.error = await buildApiErrorMessage(e);
        throw e;
      }
      finally {
        this.isLoading = false;
      }
    },

    async signOut(facilityId: string) {
      this.isLoading = true;
      this.error = null;

      try {
        const api = this.facilityAuthApiClient();
        await api.facilitiesFacilityIdAuthDelete({ facilityId });

        // サインアウト後はステートをクリア
        this.isAuthenticated = false;
        this.token = null;
        this.expiresAt = null;
      }
      catch (e) {
        const message = e instanceof Error ? e.message : 'Failed to sign out';
        this.error = message;
        console.error('Auth signOut failed:', e);
        throw e;
      }
      finally {
        this.isLoading = false;
      }
    },

    setCheckInRequired(required: boolean) {
      this.checkInRequired = required;
    },
  },
});

declare module 'pinia' {
  export interface PiniaCustomProperties {
    facilityAuthApiClient: (token?: string) => AuthApi;
  }
}
