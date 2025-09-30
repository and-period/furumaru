import { defineStore } from 'pinia';
import ApiClientFactory from '@/plugins/helper/factory';
import { AuthUserApi } from '@/types/api/facility';
import { useAuthStore } from '@/stores/auth';
import type { AuthUserResponse, UpdateAuthUserRequest } from '@/types/api/facility';
import { buildApiErrorMessage } from '@/plugins/helper/error';

interface UserState {
  isLoading: boolean;
  error: string | null;
  profile: AuthUserResponse | null;
}

export const useUserStore = defineStore('user', {
  state: (): UserState => ({
    isLoading: false,
    error: null,
    profile: null,
  }),

  getters: {
    phoneNumber(state): string {
      if (state.profile) {
        return state.profile.phoneNumber.replace('+81', '0');
      }
      return '';
    },
    lastCheckInAt(state): string {
      if (state.profile) {
        if (!state.profile.lastCheckInAt) {
          return '—';
        }
        const datetime = new Date(state.profile.lastCheckInAt * 1000);
        return datetime.toLocaleString('ja-JP', {
          year: 'numeric',
          month: '2-digit',
          day: '2-digit',
          hour: '2-digit',
          minute: '2-digit',
        });
      }
      return '—';
    },
  },

  actions: {
    reset() {
      this.isLoading = false;
      this.error = null;
      this.profile = null;
    },

    async fetchMe(facilityId: string, accessToken: string) {
      this.isLoading = true;
      this.error = null;
      try {
        const factory = new ApiClientFactory();
        const api = factory.createFacility<AuthUserApi>(AuthUserApi, accessToken);
        const res = await api.facilitiesFacilityIdUsersMeGet({ facilityId });
        this.profile = res;
        return res;
      }
      catch (e) {
        console.error('User fetchMe failed:', e);
        this.error = await buildApiErrorMessage(e);
        this.profile = null;
        throw e;
      }
      finally {
        this.isLoading = false;
      }
    },

    async updateMe(facilityId: string, userData: UpdateAuthUserRequest) {
      this.isLoading = true;
      this.error = null;
      try {
        // authStoreからアクセストークンを取得してAPIに付与
        const authStore = useAuthStore();
        const accessToken = authStore.token!.accessToken;
        const factory = new ApiClientFactory();
        const api = factory.createFacility<AuthUserApi>(AuthUserApi, accessToken);

        await api.facilitiesFacilityIdUsersMePut({
          facilityId,
          updateAuthUserRequest: userData,
        });

        // 更新成功後に最新データを取得
        await this.fetchMe(facilityId, accessToken);

        return true;
      }
      catch (e) {
        console.error('User updateMe failed:', e);
        this.error = await buildApiErrorMessage(e);
        throw e;
      }
      finally {
        this.isLoading = false;
      }
    },
  },
});
