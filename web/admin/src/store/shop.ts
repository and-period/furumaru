import { useApiClient } from '~/composables/useApiClient'
import { useProducerStore } from './producer'
import { useProductTypeStore } from './product-type'
import { ShopApi } from '~/types/api/v1'
import type { Shop, UpdateShopRequest, V1ShopsShopIdGetRequest, V1ShopsShopIdPatchRequest } from '~/types/api/v1'

export const useShopStore = defineStore('shop', () => {
  const { create, errorHandler } = useApiClient()
  const shopApi = () => create(ShopApi)

  const shop = ref<Shop>({} as Shop)
  const shops = ref<Shop[]>([])

  async function fetchShop(shopId: string): Promise<void> {
    try {
      const params: V1ShopsShopIdGetRequest = { shopId }
      const res = await shopApi().v1ShopsShopIdGet(params)
      shop.value = res.shop

      const producerStore = useProducerStore()
      producerStore.producers = res.producers

      const productTypeStore = useProductTypeStore()
      productTypeStore.productTypes = res.productTypes
    }
    catch (err) {
      return errorHandler(err, { 404: '店舗が見つかりません。' })
    }
  }

  async function updateShop(shopId: string, payload: UpdateShopRequest): Promise<void> {
    try {
      const params: V1ShopsShopIdPatchRequest = {
        shopId,
        updateShopRequest: payload,
      }
      await shopApi().v1ShopsShopIdPatch(params)
    }
    catch (err) {
      return errorHandler(err, { 400: '必須項目が不足しているか、入力内容に誤りがあります。' })
    }
  }

  return {
    shop,
    shops,
    fetchShop,
    updateShop,
  }
})
