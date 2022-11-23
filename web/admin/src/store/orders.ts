import { defineStore } from 'pinia'

import { OrderResponseItemsInner } from '~/types/api'

export const useOrderStore = defineStore('order', {
  state: () => {
    return {
      orders: [] as OrderResponseItemsInner[],
      totalItems: 0,
    }
  },

  actions: {},
})
