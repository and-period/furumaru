import { useShopStore } from './shop'
import type {
  CreatePromotionRequest,
  Promotion,
  UpdatePromotionRequest,
  V1PromotionsGetRequest,
  V1PromotionsPostRequest,
  V1PromotionsPromotionIdDeleteRequest,
  V1PromotionsPromotionIdGetRequest,
  V1PromotionsPromotionIdPatchRequest,
} from '~/types/api/v1'

export const usePromotionStore = defineStore('promotion', {
  state: () => ({
    promotion: {} as Promotion,
    promotions: [] as Promotion[],
    total: 0,
  }),
  actions: {
    /**
     * 登録済みのセール情報一覧を取得する非同期関数
     * @param limit 取得上限数
     * @param offset 取得開始位置
     * @param orders ソートキー
     */
    async fetchPromotions(limit = 20, offset = 0, orders: string[] = []): Promise<void> {
      try {
        const params: V1PromotionsGetRequest = {
          limit,
          offset,
          orders: orders.join(','),
        }
        const res = await this.promotionApi().v1PromotionsGet(params)

        const shopStore = useShopStore()
        this.promotions = res.promotions
        this.total = res.total
        shopStore.shops = res.shops
      }
      catch (err) {
        return this.errorHandler(err)
      }
    },

    /**
     * セール情報を検索する非同期関数
     * @param name タイトル(あいまい検索)
     * @param promotionIds stateの更新時に残しておく必要があるセール情報
     */
    async searchPromotions(name = '', promotionIds: string[] = []): Promise<void> {
      try {
        const params: V1PromotionsGetRequest = {
          title: name,
        }
        const res = await this.promotionApi().v1PromotionsGet(params)
        const promotions: Promotion[] = []
        this.promotions.forEach((promotion: Promotion): void => {
          if (!promotionIds.includes(promotion.id)) {
            return
          }
          promotions.push(promotion)
        })
        res.promotions.forEach((promotion: Promotion): void => {
          if (promotions.find((v): boolean => v.id === promotion.id)) {
            return
          }
          promotions.push(promotion)
        })

        const shopStore = useShopStore()
        this.promotions = promotions
        this.total = res.total
        shopStore.shops = res.shops
      }
      catch (err) {
        return this.errorHandler(err)
      }
    },

    /**
     * セールIDからセール情報情報を取得する非同期関数
     * @param promotionId セールID
     * @returns セールの情報
     */
    async getPromotion(promotionId: string): Promise<void> {
      try {
        const params: V1PromotionsPromotionIdGetRequest = {
          promotionId,
        }
        const res = await this.promotionApi().v1PromotionsPromotionIdGet(params)
        this.promotion = res.promotion
        if (!res.shop) {
          return
        }

        const shopStore = useShopStore()
        shopStore.shop = res.shop
      }
      catch (err) {
        return this.errorHandler(err, { 404: '対象のセール情報が存在しません。' })
      }
    },

    /**
     * セール情報を登録する非同期関数
     * @param payload
     */
    async createPromotion(payload: CreatePromotionRequest): Promise<void> {
      try {
        const params: V1PromotionsPostRequest = {
          createPromotionRequest: payload,
        }
        await this.promotionApi().v1PromotionsPost(params)
      }
      catch (err) {
        return this.errorHandler(err, {
          400: '必須項目が不足しているか、内容に誤りがあります。',
          409: 'このクーポンコードはすでに登録されています。',
        })
      }
    },

    /**
     * セール情報を削除する非同期関数
     * @param promotionId お知らせID
     */
    async deletePromotion(promotionId: string): Promise<void> {
      try {
        const params: V1PromotionsPromotionIdDeleteRequest = {
          promotionId,
        }
        await this.promotionApi().v1PromotionsPromotionIdDelete(params)
        this.fetchPromotions()
      }
      catch (err) {
        return this.errorHandler(err, { 404: '対象のセール情報が存在しません。' })
      }
    },

    /**
     * セール情報を編集する非同期関数
     * @param promotionId セールID
     * @param payload
     */
    async updatePromotion(promotionId: string, payload: UpdatePromotionRequest): Promise<void> {
      try {
        const params: V1PromotionsPromotionIdPatchRequest = {
          promotionId,
          updatePromotionRequest: payload,
        }
        await this.promotionApi().v1PromotionsPromotionIdPatch(params)
      }
      catch (err) {
        return this.errorHandler(err, {
          400: '必須項目が不足しているか、内容に誤りがあります。',
          404: '対象のセール情報が存在しません。',
          409: 'このクーポンコードはすでに登録されています。',
        })
      }
    },
  },
})
