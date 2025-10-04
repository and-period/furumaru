import { defineStore } from 'pinia';
import type {
  CartResponse,
  CartItem as ApiCartItem,
  Product,
  Coordinator,
} from '~/types/api/facility/models';
import { CartApi, Configuration as FacilityConfiguration } from '@/types/api/facility';
import { useAuthStore } from '~/stores/auth';
import type { AddCartItemRequest } from '@/types/api/facility';

export interface CartItem extends ApiCartItem {
  product?: Product & {
    thumbnail?: { url: string; isThumbnail: boolean };
  };
}

export interface ShoppingCart {
  carts: Array<{
    number: number;
    type: number;
    size: number;
    coordinatorId: string;
    items: CartItem[];
    coordinator?: Coordinator;
    useRate: number;
  }>;
}

/**
 * 買い物かごを管理するグローバルステート (LIFF用)
 */
export const useShoppingCartStore = defineStore('shopping-cart', {
  state: () => {
    return {
      _shoppingCart: {
        carts: [],
        coordinators: [],
        products: [],
      } as CartResponse,
    };
  },

  getters: {
    shoppingCart(state): ShoppingCart {
      return {
        carts: state._shoppingCart.carts.map((cart) => {
          return {
            ...cart,
            // コーディネーターのマッピング
            coordinator: state._shoppingCart.coordinators.find(
              coordinator => coordinator.id === cart.coordinatorId,
            ),
            // カート内の商品のマッピング
            items: cart.items.map((item) => {
              // マッピング用の商品オブジェクトを事前計算
              const product = state._shoppingCart.products.find(
                product => product.id === item.productId,
              );
              return {
                ...item,
                product: product
                  ? {
                      ...product,
                      // サムネイル画像のマッピング
                      thumbnail: product.media.find(m => m.isThumbnail),
                    }
                  : undefined,
              };
            }),
            // 占有率
            useRate: cart.rate,
          };
        }),
      };
    },

    cartIsEmpty: (state) => {
      return state._shoppingCart.carts.length === 0;
    },

    // 全カート内の商品アイテムを平坦化したリスト
    allCartItems(): CartItem[] {
      return this.shoppingCart.carts.flatMap(cart => cart.items);
    },

    // 総計金額
    totalPrice(): number {
      return this.allCartItems.reduce((total, item) => {
        const price = item.product?.price || 0;
        return total + price * item.quantity;
      }, 0);
    },

    // 総商品数
    totalQuantity(): number {
      return this.allCartItems.reduce((total, item) => total + item.quantity, 0);
    },
  },

  actions: {
    // カート情報を取得（実API呼び出し）
    async getCart(facilityId: string) {
      try {
        const runtimeConfig = useRuntimeConfig();
        const authStore = useAuthStore();

        if (!facilityId) {
          console.warn('facilityId is not specified in params. Skipping cart fetch.');
          this._shoppingCart = { carts: [], coordinators: [], products: [] } as CartResponse;
          return;
        }

        const accessToken = authStore.token?.accessToken;
        const headers = accessToken ? { Authorization: `Bearer ${accessToken}` } : undefined;
        const config = new FacilityConfiguration({
          headers,
          basePath: runtimeConfig.public.API_BASE_URL,
          credentials: 'include',
        });

        const api = new CartApi(config);
        const res = await api.facilitiesFacilityIdCartsGet({ facilityId });

        this._shoppingCart = {
          carts: res.carts ?? [],
          coordinators: res.coordinators ?? [],
          products: res.products ?? [],
        } as CartResponse;
      }
      catch (error) {
        console.error('Failed to fetch cart:', error);
        this._shoppingCart = { carts: [], coordinators: [], products: [] } as CartResponse;
      }
    },

    // カートにアイテムを追加
    async addCartItem(productId: string, quantity: number = 1) {
      try {
        const runtimeConfig = useRuntimeConfig();
        const route = useRoute();
        const authStore = useAuthStore();

        const facilityId = String(route.params.facilityId ?? '');

        if (!facilityId) {
          console.warn('facilityId is not specified in params. Skipping addCartItem.');
          return;
        }

        const accessToken = authStore.token?.accessToken;
        const headers = accessToken ? { Authorization: `Bearer ${accessToken}` } : undefined;

        const config = new FacilityConfiguration({
          headers,
          basePath: runtimeConfig.public.API_BASE_URL,
          credentials: 'include',
        });

        const api = new CartApi(config);

        const payload: AddCartItemRequest = {
          productId,
          quantity,
        };

        await api.facilitiesFacilityIdCartsItemsPost({
          facilityId,
          addCartItemRequest: payload,
        });

        // 追加後にカート情報を更新
        await this.getCart(facilityId);
      }
      catch (error) {
        console.error('Failed to add item to cart:', error);
        throw error;
      }
    },

    // カートからアイテムを削除
    async removeCartItem(productId: string) {
      try {
        const runtimeConfig = useRuntimeConfig();
        const route = useRoute();
        const authStore = useAuthStore();

        const facilityId = String(route.params.facilityId ?? '');

        if (!facilityId) {
          console.warn('facilityId is not specified in params. Skipping removeCartItem.');
          return;
        }

        const accessToken = authStore.token?.accessToken;
        const headers = accessToken ? { Authorization: `Bearer ${accessToken}` } : undefined;

        const config = new FacilityConfiguration({
          headers,
          basePath: runtimeConfig.public.API_BASE_URL,
          credentials: 'include',
        });

        const api = new CartApi(config);

        await api.facilitiesFacilityIdCartsItemsProductIdDelete({
          facilityId,
          productId,
        });

        // 削除後にカート情報を更新
        await this.getCart(facilityId);
      }
      catch (error) {
        console.error('Failed to remove item from cart:', error);
        throw error;
      }
    },
  },
});
