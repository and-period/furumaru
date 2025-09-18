import { defineStore } from 'pinia';
import type { OrderApi, OrdersResponse, OrderResponse } from '@/types/api/facility';
import { useAuthStore } from '@/stores/auth';

interface OrderState {
  isLoading: boolean;
  error: string | null;
  orders: OrdersResponse | null;
  orderDetail: OrderResponse | null;
}

/**
 * 注文を管理するグローバルステート (LIFF用)
 * - src/types/api/facility/apis/OrderApi.ts を利用
 */
export const useOrderStore = defineStore('order', {
  state: (): OrderState => ({
    isLoading: false,
    error: null,
    orders: null,
    orderDetail: null,
  }),

  actions: {
    /**
     * 注文一覧を取得する
     * @param facilityId 施設ID
     * @param limit 取得件数制限（省略可能）
     * @param offset オフセット（省略可能）
     * @param types 注文タイプ（省略可能）
     */
    async getOrders(facilityId: string, limit?: number, offset?: number, types?: Array<number>) {
      this.isLoading = true;
      this.error = null;

      try {
        const authStore = useAuthStore();

        if (!authStore.token?.accessToken) {
          throw new Error('Access token is required to fetch orders.');
        }

        const api = this.facilityOrderApiClient(authStore.token.accessToken);

        const response = await api.facilitiesFacilityIdOrdersGet({
          facilityId,
          limit,
          offset,
          types,
        });

        this.orders = response;
        return response;
      }
      catch (e) {
        const message = e instanceof Error ? e.message : 'Failed to fetch orders';
        this.error = message;
        console.error('Order fetch failed:', e);
        throw e;
      }
      finally {
        this.isLoading = false;
      }
    },

    /**
     * 注文詳細を取得する
     * @param facilityId 施設ID
     * @param orderId 注文ID
     */
    async getOrderDetail(facilityId: string, orderId: string) {
      this.isLoading = true;
      this.error = null;

      try {
        const authStore = useAuthStore();

        if (!authStore.token?.accessToken) {
          throw new Error('Access token is required to fetch order detail.');
        }

        const api = this.facilityOrderApiClient(authStore.token.accessToken);

        const response = await api.facilitiesFacilityIdOrdersOrderIdGet({
          facilityId,
          orderId,
        });

        this.orderDetail = response;
        return response;
      }
      catch (e) {
        const message = e instanceof Error ? e.message : 'Failed to fetch order detail';
        this.error = message;
        console.error('Order detail fetch failed:', e);
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
    facilityOrderApiClient: (token?: string) => OrderApi;
  }
}
