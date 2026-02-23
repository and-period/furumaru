import { useApiClient } from '~/composables/useApiClient'
import { ShippingApi } from '~/types/api/v1'
import type { CreateShippingRequest, Shipping, ShippingsResponse, UpdateDefaultShippingRequest, UpdateShippingRequest, V1CoordinatorsCoordinatorIdShippingsGetRequest, V1CoordinatorsCoordinatorIdShippingsPostRequest, V1CoordinatorsCoordinatorIdShippingsShippingIdActivationPatchRequest, V1CoordinatorsCoordinatorIdShippingsShippingIdDeleteRequest, V1CoordinatorsCoordinatorIdShippingsShippingIdGetRequest, V1CoordinatorsCoordinatorIdShippingsShippingIdPatchRequest, V1ShippingsDefaultPatchRequest } from '~/types/api/v1'

export const useShippingStore = defineStore('shipping', () => {
  const { create, errorHandler } = useApiClient()
  const shippingApi = () => create(ShippingApi)

  const shipping = ref<Shipping>({} as Shipping)
  const shippings = ref<Shipping[]>([])
  const total = ref<number>(0)

  async function fetchShippings(coordinatorId: string, limit?: number, offset?: number): Promise<ShippingsResponse> {
    try {
      const params: V1CoordinatorsCoordinatorIdShippingsGetRequest = { coordinatorId, limit, offset }
      const res = await shippingApi().v1CoordinatorsCoordinatorIdShippingsGet(params)
      shippings.value = res.shippings
      total.value = res.total
      return res
    }
    catch (err) {
      return errorHandler(err, { 404: '対象のコーディネーターが見つかりません。' })
    }
  }

  async function fetchShipping(coordinatorId: string, shippingId: string): Promise<Shipping> {
    try {
      const params: V1CoordinatorsCoordinatorIdShippingsShippingIdGetRequest = { coordinatorId, shippingId }
      const res = await shippingApi().v1CoordinatorsCoordinatorIdShippingsShippingIdGet(params)
      return res.shipping
    }
    catch (err) {
      return errorHandler(err, { 404: '配送設定が見つかりません。' })
    }
  }

  async function createShipping(coordinatorId: string, payload: CreateShippingRequest): Promise<void> {
    try {
      const params: V1CoordinatorsCoordinatorIdShippingsPostRequest = {
        coordinatorId,
        createShippingRequest: payload,
      }
      await shippingApi().v1CoordinatorsCoordinatorIdShippingsPost(params)
    }
    catch (err) {
      return errorHandler(err)
    }
  }

  async function updateShipping(coordinatorId: string, shippingId: string, payload: UpdateShippingRequest): Promise<void> {
    try {
      const params: V1CoordinatorsCoordinatorIdShippingsShippingIdPatchRequest = {
        coordinatorId,
        shippingId,
        updateShippingRequest: payload,
      }
      await shippingApi().v1CoordinatorsCoordinatorIdShippingsShippingIdPatch(params)
    }
    catch (err) {
      return errorHandler(err)
    }
  }

  async function deleteShipping(coordinatorId: string, shippingId: string): Promise<void> {
    try {
      const params: V1CoordinatorsCoordinatorIdShippingsShippingIdDeleteRequest = { coordinatorId, shippingId }
      await shippingApi().v1CoordinatorsCoordinatorIdShippingsShippingIdDelete(params)
    }
    catch (err) {
      return errorHandler(err, { 404: '対象の配送設定が見つかりません。' })
    }
  }

  async function activeShipping(coordinatorId: string, shippingId: string): Promise<void> {
    try {
      const params: V1CoordinatorsCoordinatorIdShippingsShippingIdActivationPatchRequest = { coordinatorId, shippingId }
      await shippingApi().v1CoordinatorsCoordinatorIdShippingsShippingIdActivationPatch(params)
    }
    catch (err) {
      return errorHandler(err, { 404: '対象の配送設定が見つかりません。' })
    }
  }

  async function fetchDefaultShipping(): Promise<void> {
    try {
      const res = await shippingApi().v1ShippingsDefaultGet()
      shipping.value = res.shipping
    }
    catch (err) {
      return errorHandler(err)
    }
  }

  async function updateDefaultShipping(payload: UpdateDefaultShippingRequest): Promise<void> {
    try {
      const params: V1ShippingsDefaultPatchRequest = { updateDefaultShippingRequest: payload }
      await shippingApi().v1ShippingsDefaultPatch(params)
    }
    catch (err) {
      return errorHandler(err, { 400: '必須項目が不足しているか、入力内容に誤りがあります。' })
    }
  }

  return {
    shipping,
    shippings,
    total,
    fetchShippings,
    fetchShipping,
    createShipping,
    updateShipping,
    deleteShipping,
    activeShipping,
    fetchDefaultShipping,
    updateDefaultShipping,
  }
})
