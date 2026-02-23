import { useShopStore } from './shop'
import { useApiClient } from '~/composables/useApiClient'
import { PromotionApi } from '~/types/api/v1'
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

export const usePromotionStore = defineStore('promotion', () => {
  const { create, errorHandler } = useApiClient()
  const promotionApi = () => create(PromotionApi)

  const promotion = ref<Promotion>({} as Promotion)
  const promotions = ref<Promotion[]>([])
  const total = ref<number>(0)

  /**
   * 登録済みのセール情報一覧を取得する非同期関数
   * @param limit 取得上限数
   * @param offset 取得開始位置
   * @param orders ソートキー
   */
  async function fetchPromotions(limit = 20, offset = 0, orders: string[] = []): Promise<void> {
    try {
      const params: V1PromotionsGetRequest = {
        limit,
        offset,
        orders: orders.join(','),
      }
      const res = await promotionApi().v1PromotionsGet(params)

      const shopStore = useShopStore()
      promotions.value = res.promotions
      total.value = res.total
      shopStore.shops = res.shops
    }
    catch (err) {
      return errorHandler(err)
    }
  }

  /**
   * セール情報を検索する非同期関数
   * @param name タイトル(あいまい検索)
   * @param promotionIds stateの更新時に残しておく必要があるセール情報
   */
  async function searchPromotions(name = '', promotionIds: string[] = []): Promise<void> {
    try {
      const params: V1PromotionsGetRequest = {
        title: name,
      }
      const res = await promotionApi().v1PromotionsGet(params)
      const merged: Promotion[] = []
      promotions.value.forEach((promotion: Promotion): void => {
        if (!promotionIds.includes(promotion.id)) {
          return
        }
        merged.push(promotion)
      })
      res.promotions.forEach((promotion: Promotion): void => {
        if (merged.find((v): boolean => v.id === promotion.id)) {
          return
        }
        merged.push(promotion)
      })

      const shopStore = useShopStore()
      promotions.value = merged
      total.value = res.total
      shopStore.shops = res.shops
    }
    catch (err) {
      return errorHandler(err)
    }
  }

  /**
   * セールIDからセール情報情報を取得する非同期関数
   * @param promotionId セールID
   * @returns セールの情報
   */
  async function getPromotion(promotionId: string): Promise<void> {
    try {
      const params: V1PromotionsPromotionIdGetRequest = {
        promotionId,
      }
      const res = await promotionApi().v1PromotionsPromotionIdGet(params)
      promotion.value = res.promotion
      if (!res.shop) {
        return
      }

      const shopStore = useShopStore()
      shopStore.shop = res.shop
    }
    catch (err) {
      return errorHandler(err, { 404: '対象のセール情報が存在しません。' })
    }
  }

  /**
   * セール情報を登録する非同期関数
   * @param payload
   */
  async function createPromotion(payload: CreatePromotionRequest): Promise<void> {
    try {
      const params: V1PromotionsPostRequest = {
        createPromotionRequest: payload,
      }
      await promotionApi().v1PromotionsPost(params)
    }
    catch (err) {
      return errorHandler(err, {
        400: '必須項目が不足しているか、内容に誤りがあります。',
        409: 'このクーポンコードはすでに登録されています。',
      })
    }
  }

  /**
   * セール情報を削除する非同期関数
   * @param promotionId お知らせID
   */
  async function deletePromotion(promotionId: string): Promise<void> {
    try {
      const params: V1PromotionsPromotionIdDeleteRequest = {
        promotionId,
      }
      await promotionApi().v1PromotionsPromotionIdDelete(params)
      fetchPromotions()
    }
    catch (err) {
      return errorHandler(err, { 404: '対象のセール情報が存在しません。' })
    }
  }

  /**
   * セール情報を編集する非同期関数
   * @param promotionId セールID
   * @param payload
   */
  async function updatePromotion(promotionId: string, payload: UpdatePromotionRequest): Promise<void> {
    try {
      const params: V1PromotionsPromotionIdPatchRequest = {
        promotionId,
        updatePromotionRequest: payload,
      }
      await promotionApi().v1PromotionsPromotionIdPatch(params)
    }
    catch (err) {
      return errorHandler(err, {
        400: '必須項目が不足しているか、内容に誤りがあります。',
        404: '対象のセール情報が存在しません。',
        409: 'このクーポンコードはすでに登録されています。',
      })
    }
  }

  return {
    // state
    promotion,
    promotions,
    total,
    // actions
    fetchPromotions,
    searchPromotions,
    getPromotion,
    createPromotion,
    deletePromotion,
    updatePromotion,
  }
})
