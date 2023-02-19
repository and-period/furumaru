import { defineStore } from 'pinia'

/**
 * 買い物かごを管理するグローバルステート
 */
export const useShoppingStore = defineStore('shopping', {
  state: () => {
    return {
      cartItems: []
    }
  },

  getters: {
    cartIsEmpty: (state) => {
      return state.cartItems.length === 0
    }
  }
})
