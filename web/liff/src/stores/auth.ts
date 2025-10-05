import { defineStore } from 'pinia';
import type { AuthApi, AuthResponse, SignInRequest } from '@/types/api/facility';

interface AuthState {
  isLoading: boolean;
  error: string | null;
  isAuthenticated: boolean;
  token: Pick<AuthResponse, 'accessToken' | 'refreshToken' | 'expiresIn' | 'tokenType' | 'userId'> | null;
  // アクセストークンの有効期限（UNIXエポックミリ秒）
  expiresAt: number | null;
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
  },
});

declare module 'pinia' {
  export interface PiniaCustomProperties {
    facilityAuthApiClient: (token?: string) => AuthApi;
  }
}
