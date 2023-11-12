import { defineStore } from 'pinia'
import { AddCartItemRequest, CartResponse } from '~/types/api'
import { ProductItem, ShoppingCart } from '~/types/store'

/**
 * 買い物かごを管理するグローバルステート
 */
export const useShoppingCartStore = defineStore('shopping-cart', {
  state: () => {
    return {
      cartItems: [],
      recommendProducts: [] as ProductItem[],

      _shoppingCart: {
        carts: [],
        coordinators: [],
        products: [],
      } as CartResponse,
    }
  },

  getters: {
    shoppingCart(state) {
      return {
        carts: state._shoppingCart.carts.map((cart) => {
          return {
            ...cart,
            // コーディネーターのマッピング
            coordinator: state._shoppingCart.coordinators.find(
              (coordinator) => coordinator.id === cart.coordinatorId,
            ),
            // カート内の商品のマッピング
            items: cart.items.map((item) => {
              // マッピング用の商品オブジェクトを事前計算
              const product = state._shoppingCart.products.find(
                (product) => product.id === item.productId,
              )
              return {
                ...item,
                product: {
                  ...product,
                  // サムネイル画像のマッピング
                  thumbnail: product?.media.find((m) => m.isThumbnail),
                },
              }
            }),
          }
        }),
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
        .map((cart) =>
          cart.items
            .map((item) => item.product.price)
            .filter((price) => typeof price === 'number')
            .reduce((sum, price) => sum + price),
        )
        .reduce((sum, price) => sum + price)
      return totalPrice
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
  },
})
