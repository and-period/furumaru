import { defineStore } from 'pinia';
import type {
  CartResponse,
  CartItem as ApiCartItem,
  Product,
  Coordinator,
  Weekday,
  ProductStatus,
  DeliveryType,
  StorageMethodType,
} from '~/types/api/models';

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
    // カート情報を取得（モックデータで実装）
    async getCart() {
      // TODO: 実際のAPI呼び出しに置き換える
      // 現在はモックデータを使用
      this._shoppingCart = {
        carts: [
          {
            number: 1,
            type: 1,
            size: 1,
            coordinatorId: 'coordinator-1',
            items: [
              {
                productId: 'vNpfY5NoPfa9iVcyAgS8qU',
                quantity: 2,
              },
            ],
            rate: 0.5,
          },
        ],
        coordinators: [
          {
            id: 'coordinator-1',
            marcheName: '田中マルシェ',
            username: '田中農園',
            profile: 'こだわりの野菜を作っています',
            productTypeIds: ['type-1'],
            businessDays: [1, 2, 3, 4, 5], // Monday to Friday
            thumbnailUrl: 'https://example.com/coordinator.jpg',
            headerUrl: '',
            promotionVideoUrl: '',
            instagramId: '',
            facebookId: '',
            prefecture: '広島県',
            city: '東広島市',
          },
        ],
        products: [
          {
            id: 'vNpfY5NoPfa9iVcyAgS8qU',
            coordinatorId: 'coordinator-1',
            producerId: 'producer-1',
            categoryId: 'category-1',
            productTypeId: 'type-1',
            productTagIds: [],
            name: '【瀬戸内の名産】赤土じゃがいも5キロ(40〜50個)',
            description: '瀬戸内海を臨む安芸津町の赤土じゃがいも',
            status: 2, // FOR_SALE
            inventory: 93,
            weight: 5,
            itemUnit: '個',
            itemDescription: '5キロ',
            thumbnailUrl: 'https://assets.furumaru.and-period.co.jp/products/media/image/tic8TSBKJGWqGbdpi3h5z7.jpg',
            media: [
              {
                url: 'https://assets.furumaru.and-period.co.jp/products/media/image/tic8TSBKJGWqGbdpi3h5z7.jpg',
                isThumbnail: true,
              },
            ],
            price: 2700,
            expirationDate: 30,
            recommendedPoint1: '地元の名産・安芸津の赤土じゃがいも',
            recommendedPoint2: 'ホクホク＆しっとりで甘い',
            recommendedPoint3: '煮くずれしにくく調理しやすい',
            storageMethodType: 1, // NORMAL
            deliveryType: 1, // NORMAL
            box60Rate: 50,
            box80Rate: 40,
            box100Rate: 30,
            originPrefecture: '広島県',
            originCity: '東広島市',
            rate: {
              count: 0,
              average: 0,
              detail: {},
            },
            startAt: Math.floor(Date.now() / 1000),
            endAt: Math.floor((Date.now() + 86400000 * 30) / 1000),
          },
        ],
      };
    },
  },
});
