import { defineStore } from 'pinia'
import { useAuthStore } from './auth'
import {
  type AddCartItemRequest,
  type CalcCartResponse,
  type CartResponse,
  type PaymentSystem,
} from '~/types/api'
import type {
  CalcCart,
  PaymentSystemStatus,
  ProductItem,
  ShoppingCart,
} from '~/types/store'
import { getPaymentMethodNameByPaymentMethodType } from '~/lib/order'

/**
 * 買い物かごを管理するグローバルステート
 */
export const useShoppingCartStore = defineStore('shopping-cart', {
  state: () => {
    return {
      cartItems: [],
      recommendProducts: [] as ProductItem[],

      _calcCartResponseItem: undefined as CalcCartResponse | undefined,

      _shoppingCart: {
        carts: [],
        coordinators: [],
        products: [],
      } as CartResponse,

      _paymentSystemStatus: [] as PaymentSystem[],
    }
  },

  getters: {
    shoppingCart(state) {
      return {
        carts: state._shoppingCart.carts.map((cart) => {
          const boxType = (type: number) => {
            switch (type) {
              case 1:
                return this.i18n.t('items.details.deliveryTypeStandard')
              case 2:
                return this.i18n.t('items.details.deliveryTypeRefrigerated')
              case 3:
                return this.i18n.t('items.details.deliveryTypeFrozen')
              default:
                return ''
            }
          }
          const boxSize = (size: number) => {
            switch (size) {
              case 1:
                return 60
              case 2:
                return 80
              case 3:
                return 100
              default:
                return 0
            }
          }

          return {
            ...cart,
            // コーディネーターのマッピング
            coordinator: state._shoppingCart.coordinators.find(
              coordinator => coordinator.id === cart.coordinatorId,
            ),
            // 箱タイプ
            boxType: boxType(cart.type),
            // 箱サイズ
            boxSize: boxSize(cart.size),
            // カート内の商品のマッピング
            items: cart.items.map((item) => {
              // マッピング用の商品オブジェクトを事前計算
              const product = state._shoppingCart.products.find(
                product => product.id === item.productId,
              )
              return {
                ...item,
                product: {
                  ...product,
                  // サムネイル画像のマッピング
                  thumbnail: product?.media.find(m => m.isThumbnail),
                },
              }
            }),
            // 占有率
            useRate: cart.rate,
          }
        }),
      }
    },

    calcCartResponseItem: (state): CalcCart | undefined => {
      if (state._calcCartResponseItem) {
        // 商品情報をマッピング
        const products: ProductItem[]
          = state._calcCartResponseItem.products.map((product) => {
            return {
              ...product,
              thumbnail: product.media.find(m => m.isThumbnail),
            }
          })

        return {
          ...state._calcCartResponseItem,
          // 買い物カゴの中身の商品をマッピング
          items: state._calcCartResponseItem.items.map((item) => {
            return {
              ...item,
              product: products.find(
                product => product.id === item.productId,
              ),
            }
          }),
        }
      }
      else {
        return undefined
      }
    },

    cartIsEmpty: (state) => {
      return state._shoppingCart.carts.length === 0
    },

    totalPrice() {
      const carts = this.shoppingCart.carts as ShoppingCart[]
      if (carts.length === 0) {
        return 0
      }
      const totalPrice = carts
        .map(cart =>
          cart.items
            .map(item => item.product.price)
            .filter(price => typeof price === 'number')
            .reduce((sum, price) => sum + price),
        )
        .reduce((sum, price) => sum + price)
      return totalPrice
    },

    availablePaymentSystem: (state): PaymentSystemStatus[] => {
      return state._paymentSystemStatus
        .map((item) => {
          return {
            ...item,
            methodName: getPaymentMethodNameByPaymentMethodType(
              item.methodType,
            ),
          }
        })
        .filter(item => item.status === 1)
    },
  },

  actions: {
    async getCart() {
      const response = await this.cartApiClient().v1GetCart()
      this._shoppingCart = response
    },

    async addCart(payload: AddCartItemRequest) {
      await this.cartApiClient().v1AddCartItem({ body: payload })
      this.getCart()
    },

    async removeProductFromCart(cartNumber: number, productId: string) {
      await this.cartApiClient().v1RemoveCartItem({
        number: cartNumber,
        productId,
      })
      this.getCart()
    },

    async calcCartItemByCoordinatorId(
      coordinatorId: string,
      cartNumber?: number,
      prefecture?: number,
      promotion?: string,
    ): Promise<string | undefined> {
      try {
        const authStore = useAuthStore()
        const res = await this.cartApiClient(authStore.accessToken).v1CalcCart({
          coordinatorId,
          number: cartNumber,
          prefecture,
          promotion,
        })
        this._calcCartResponseItem = res
        const requestId = res.requestId
        return requestId
      }
      catch (error) {
        return this.errorHandler(error)
      }
    },

    async fetchAvailablePaymentOptions() {
      const res = await this.statusApiClient().v1ListPaymentSystems()
      this._paymentSystemStatus = res.systems
    },

    /**
     * 有効なプロモーションコードかを検証する
     */
    async verifyPromotionCode(promotionCode: string): Promise<boolean> {
      try {
        const authStore = useAuthStore()
        await this.promotionApiClient(authStore.accessToken).v1GetPromotion({
          code: promotionCode,
        })
        return true
      }
      catch (_error) {
        return false
      }
    },

    /**
     * 有効なプロモーションコードかを検証する(ゲスト用)
     */
    async verifyGuestPromotionCode(promotionCode: string): Promise<boolean> {
      try {
        await this.promotionApiClient().v1GetPromotion({
          code: promotionCode,
        })
        return true
      }
      catch (_error) {
        return false
      }
    },
  },
})
