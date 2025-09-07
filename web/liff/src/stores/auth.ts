import { defineStore } from 'pinia';
import type { AuthApi, AuthResponse, RequestSignInRequest } from '@/types/api/facility';

interface AuthState {
  isLoading: boolean;
  error: string | null;
  isAuthenticated: boolean;
  token: Pick<AuthResponse, 'accessToken' | 'refreshToken' | 'expiresIn' | 'tokenType' | 'userId'> | null;
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
  }),

  actions: {
    async signIn(liffIDToken: string) {
      this.isLoading = true;
      this.error = null;

      try {
        const route = useRoute();

        const facilityId = String(route.params.facilityId ?? '');
        if (!facilityId) {
          throw new Error('facilityId is not specified in params.');
        }

        const api = this.facilityAuthApiClient();

        const payload: RequestSignInRequest = {
          authToken: liffIDToken,
        };

        const res = await api.facilitiesFacilityIdAuthPost({
          facilityId,
          requestSignInRequest: payload,
        });

        // ステート更新（サーバー側Cookieでも管理される前提だが、UI用に保持）
        this.token = {
          accessToken: res.accessToken,
          refreshToken: res.refreshToken,
          expiresIn: res.expiresIn,
          tokenType: res.tokenType,
          userId: res.userId,
        };
        this.isAuthenticated = Boolean(res.accessToken || res.userId);
        return res;
      }
      catch (e) {
        const message = e instanceof Error ? e.message : 'Failed to sign in';
        this.error = message;
        this.isAuthenticated = false;
        this.token = null;
        console.error('Auth signIn failed:', e);
        throw e;
      }
      finally {
        this.isLoading = false;
      }
    },

    async signOut() {
      this.isLoading = true;
      this.error = null;

      try {
        const route = useRoute();

        const facilityId = String(route.params.facilityId ?? '');
        if (!facilityId) {
          throw new Error('facilityId is not specified in params.');
        }
        const api = this.facilityAuthApiClient();
        await api.facilitiesFacilityIdAuthDelete({ facilityId });

        // サインアウト後はステートをクリア
        this.isAuthenticated = false;
        this.token = null;
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
