import { defineStore } from 'pinia'
import { ProductItem } from '~/types/store'
/**
 * 買い物かごを管理するグローバルステート
 */
export const useShoppingStore = defineStore('shopping', {
  state: () => {
    return {
      cartItems: [],
      recommendProducts: [] as ProductItem[]
    }
  },

  getters: {
    cartIsEmpty: (state) => {
      return state.cartItems.length === 0
    }
  },

  actions: {
    /**
     * ダミーデータセットアップ用の関数
    */
    setupDummyData () {
      const baseItem: ProductItem = {
        id: '',
        name: 'たまねぎ',
        description: '',
        producerId: '',
        storeName: '',
        categoryId: '',
        categoryName: '',
        productTypeId: '',
        productTypeName: '',
        productTypeIconUrl: '',
        public: false,
        inventory: 0,
        weight: 0,
        itemUnit: '',
        itemDescription: '',
        media: [
          {
            url: '~/assets/img/sample.png',
            isThumbnail: true,
            images: [
              {
                url: 'https://and-period.jp/thumbnail_240.png',
                size: 1
              },
              {
                url: 'https://and-period.jp/thumbnail_675.png',
                size: 2
              },
              {
                url: 'https://and-period.jp/thumbnail_900.png',
                size: 3
              }
            ]
          }
        ],
        price: 3000,
        deliveryType: 0,
        box60Rate: 0,
        box80Rate: 0,
        box100Rate: 0,
        originPrefecture: '',
        originCity: '',
        createdAt: 0,
        updatedAt: 0
      }

      const items = Array.from(Array(5)).map((_) => {
        return { ...baseItem }
      })

      this.recommendProducts = items
    }
  }
})
