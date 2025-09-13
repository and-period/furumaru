import type { Shop, UpdateShopRequest, V1ShopsShopIdGetRequest, V1ShopsShopIdPatchRequest } from '~/types/api/v1'

export const useShopStore = defineStore('shop', {
  state: () => ({
    shop: {} as Shop,
    shops: [] as Shop[],
  }),

  actions: {
    /**
     * 店舗情報を取得する非同期関数
     * @param shopId 店舗ID
     */
    async fetchShop(shopId: string): Promise<void> {
      try {
        const params: V1ShopsShopIdGetRequest = {
          shopId,
        }
        const res = await this.shopApi().v1ShopsShopIdGet(params)
        this.shop = res.shop
      }
      catch (err) {
        return this.errorHandler(err, {
          404: '店舗が見つかりません。',
        })
      }
    },

    /**
     * 店舗情報を取得する非同期関数
     * @param shopId 店舗ID
     */
    async updateShop(shopId: string, payload: UpdateShopRequest): Promise<void> {
      try {
        const params: V1ShopsShopIdPatchRequest = {
          shopId,
          updateShopRequest: payload,
        }
        const res = await this.shopApi().v1ShopsShopIdPatch(params)
      }
      catch (err) {
        return this.errorHandler(err, {
          400: '必須項目が不足しているか、入力内容に誤りがあります。',
        })
      }
    },
  },
})
