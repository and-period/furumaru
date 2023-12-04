import { defineStore } from 'pinia'
import type { AddCartItemRequest, CartResponse } from '~/types/api'
import type { ProductItem, ShoppingCart } from '~/types/store'

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
          const boxType = (type: number) => {
            switch (type) {
              case 1:
                return '通常便'
              case 2:
                return '冷蔵便'
              case 3:
                return '冷凍便'
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
              (coordinator) => coordinator.id === cart.coordinatorId,
            ),
            // 箱タイプ
            boxType: boxType(cart.type),
            // 箱サイズ
            boxSize: boxSize(cart.size),
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
            // 占有率
            useRate: cart.items
              .map((item) => {
                const product = state._shoppingCart.products.find(
                  (product) => product.id === item.productId,
                )
                return product
              })
              .map((product) => {
                if (product) {
                  switch (cart.size) {
                    case 1:
                      return product.box60Rate
                    case 2:
                      return product.box80Rate
                    case 3:
                      return product.box100Rate
                    default:
                      return 0
                  }
                } else {
                  return 0
                }
              })
              .reduce((sum, rate) => sum + rate),
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

    async removeProductFromCart(cartNumber: number, productId: string) {
      await this.cartApiClient().v1RemoveCartItem({
        number: cartNumber,
        productId,
      })
      this.getCart()
    },
  },
})
