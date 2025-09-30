import { defineStore } from 'pinia';
import { CheckoutApi, Configuration as FacilityConfiguration } from '@/types/api/facility';
import type {
  CheckoutRequest,
  CheckoutCreditCard,
  CheckoutResponse,
  CheckoutStateResponse,
  PaymentMethodType,
} from '@/types/api/facility';
import { useAuthStore } from '~/stores/auth';
import { useShoppingCartStore } from '~/stores/shopping';

interface CheckoutState {
  isLoading: boolean;
  error: string | null;
  redirectUrl: string | null;
  lastResponse: CheckoutResponse | null;
  lastStatus: CheckoutStateResponse | null;
}

export interface StartCheckoutPayload {
  callbackUrl: string;
  paymentMethod: PaymentMethodType;
  pickupAt: number; // Unix timestamp
  requestId?: string;
  coordinatorId?: string; // 未指定時はカートから推定
  boxNumber?: number;
  creditCard: CheckoutCreditCard;
  promotionCode?: string;
  total?: number;
}

export const useCheckoutStore = defineStore('checkout', {
  state: (): CheckoutState => ({
    isLoading: false,
    error: null,
    redirectUrl: null,
    lastResponse: null,
    lastStatus: null,
  }),

  actions: {
    // Checkout API クライアント
    checkoutApiClient() {
      const runtimeConfig = useRuntimeConfig();
      const authStore = useAuthStore();
      const accessToken = authStore.token?.accessToken;
      const headers = accessToken ? { Authorization: `Bearer ${accessToken}` } : undefined;
      const config = new FacilityConfiguration({
        headers,
        basePath: runtimeConfig.public.API_BASE_URL,
        credentials: 'include',
      });
      return new CheckoutApi(config);
    },

    // requestId を生成（未指定時）
    generateRequestId(): string {
      // シンプルな一意キー（環境に依存しない実装）
      const rand = Math.random().toString(36).slice(2, 10);
      const ts = Date.now().toString(36);
      return `req_${ts}_${rand}`;
    },

    // チェックアウト開始
    async startCheckout(payload: StartCheckoutPayload): Promise<CheckoutResponse> {
      this.isLoading = true;
      this.error = null;
      this.redirectUrl = null;
      this.lastResponse = null;

      try {
        const route = useRoute();
        const facilityId = String(route.params.facilityId ?? '');
        if (!facilityId) throw new Error('facilityId is not specified in params.');

        const cartStore = useShoppingCartStore();

        // coordinatorId 未指定の場合は、カートから推定（今回の要件: コーディネーターは1つ）
        const coordinatorId = payload.coordinatorId
          ?? cartStore.shoppingCart.carts[0]?.coordinatorId
          ?? '';
        if (!coordinatorId) {
          throw new Error('coordinatorId is required but not found in cart.');
        }

        const requestId = payload.requestId || this.generateRequestId();

        const body: CheckoutRequest = {
          boxNumber: payload.boxNumber || 0,
          callbackUrl: payload.callbackUrl,
          coordinatorId,
          creditCard: payload.creditCard,
          paymentMethod: payload.paymentMethod,
          promotionCode: payload.promotionCode || '',
          requestId,
          orderRequest: '',
          pickupLocation: '',
          pickupAt: payload.pickupAt,
          total: payload.total || 0,
        };

        const api = this.checkoutApiClient();
        const res = await api.facilitiesFacilityIdCheckoutsPost({
          facilityId,
          checkoutRequest: body,
        });

        this.lastResponse = res;
        this.redirectUrl = res.url ?? null;
        return res;
      }
      catch (e) {
        const message = e instanceof Error ? e.message : 'Failed to start checkout';
        this.error = message;
        console.error('Checkout start failed:', e);
        throw e;
      }
      finally {
        this.isLoading = false;
      }
    },

    // 支払い状態の取得
    async fetchCheckoutState(transactionId: string): Promise<CheckoutStateResponse> {
      this.isLoading = true;
      this.error = null;
      try {
        const route = useRoute();
        const facilityId = String(route.params.facilityId ?? '');
        if (!facilityId) throw new Error('facilityId is not specified in params.');

        const api = this.checkoutApiClient();
        const res = await api.facilitiesFacilityIdCheckoutsTransactionIdGet({
          facilityId,
          transactionId,
        });
        this.lastStatus = res;
        return res;
      }
      catch (e) {
        const message = e instanceof Error ? e.message : 'Failed to fetch checkout state';
        this.error = message;
        console.error('Checkout state fetch failed:', e);
        throw e;
      }
      finally {
        this.isLoading = false;
      }
    },
  },
});
