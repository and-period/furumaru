import { useApiClient } from '~/composables/useApiClient'
import { fileUpload } from './helper'
import { useProductTypeStore } from './product-type'
import { useShopStore } from './shop'
import { CoordinatorApi, UploadApi } from '~/types/api/v1'
import type {
  Coordinator,
  CreateCoordinatorRequest,
  Producer,
  UpdateCoordinatorRequest,
  V1CoordinatorsCoordinatorIdDeleteRequest,
  V1CoordinatorsCoordinatorIdGetRequest,
  V1CoordinatorsCoordinatorIdPatchRequest,
  V1CoordinatorsGetRequest,
  V1CoordinatorsPostRequest,
  V1UploadCoordinatorsBonusVideoPostRequest,
  V1UploadCoordinatorsHeaderPostRequest,
  V1UploadCoordinatorsPromotionVideoPostRequest,
  V1UploadCoordinatorsThumbnailPostRequest,
} from '~/types/api/v1'

export const useCoordinatorStore = defineStore('coordinator', () => {
  const { create, errorHandler } = useApiClient()
  const coordinatorApi = () => create(CoordinatorApi)
  const uploadApi = () => create(UploadApi)

  const coordinator = ref<Coordinator>({} as Coordinator)
  const coordinators = ref<Coordinator[]>([])
  const producers = ref<Producer[]>([])
  const totalItems = ref<number>(0)
  const producerTotalItems = ref<number>(0)

  async function fetchCoordinators(limit = 20, offset = 0): Promise<void> {
    try {
      const params: V1CoordinatorsGetRequest = { limit, offset }
      const res = await coordinatorApi().v1CoordinatorsGet(params)

      const productTypeStore = useProductTypeStore()
      const shopStore = useShopStore()
      coordinators.value = res.coordinators
      totalItems.value = res.total
      productTypeStore.productTypes = res.productTypes
      shopStore.shops = res.shops
    }
    catch (err) {
      return errorHandler(err)
    }
  }

  async function searchCoordinators(name = '', coordinatorIds: string[] = []): Promise<void> {
    try {
      const params: V1CoordinatorsGetRequest = { username: name }
      const res = await coordinatorApi().v1CoordinatorsGet(params)
      const merged: Coordinator[] = []
      coordinators.value.forEach((c: Coordinator): void => {
        if (!coordinatorIds.includes(c.id)) {
          return
        }
        merged.push(c)
      })
      res.coordinators.forEach((c: Coordinator): void => {
        if (merged.find((v): boolean => v.id === c.id)) {
          return
        }
        merged.push(c)
      })
      const shopStore = useShopStore()
      coordinators.value = merged
      totalItems.value = res.total
      shopStore.shops = res.shops
    }
    catch (err) {
      return errorHandler(err)
    }
  }

  async function getCoordinator(coordinatorId: string): Promise<void> {
    try {
      const params: V1CoordinatorsCoordinatorIdGetRequest = { coordinatorId }
      const res = await coordinatorApi().v1CoordinatorsCoordinatorIdGet(params)

      const productTypeStore = useProductTypeStore()
      const shopStore = useShopStore()
      coordinator.value = res.coordinator
      productTypeStore.productTypes = res.productTypes
      shopStore.shop = res.shop
    }
    catch (err) {
      return errorHandler(err, {
        404: 'コーディネーター情報が見つかりません。',
      })
    }
  }

  async function createCoordinator(payload: CreateCoordinatorRequest) {
    try {
      const params: V1CoordinatorsPostRequest = { createCoordinatorRequest: payload }
      const res = await coordinatorApi().v1CoordinatorsPost(params)
      return res
    }
    catch (err) {
      return errorHandler(err, {
        400: '必須項目が不足しているか、内容に誤りがあります',
        409: 'このメールアドレスはすでに登録されているため、登録できません。',
      })
    }
  }

  async function updateCoordinator(coordinatorId: string, payload: UpdateCoordinatorRequest): Promise<void> {
    try {
      const params: V1CoordinatorsCoordinatorIdPatchRequest = {
        coordinatorId,
        updateCoordinatorRequest: payload,
      }
      await coordinatorApi().v1CoordinatorsCoordinatorIdPatch(params)
    }
    catch (err) {
      return errorHandler(err, {
        400: '必須項目が不足しているか、入力内容に誤りがあります。',
        404: '対象のコーディネーターが存在しません',
      })
    }
  }

  async function uploadCoordinatorThumbnail(payload: File): Promise<string> {
    try {
      const params: V1UploadCoordinatorsThumbnailPostRequest = {
        getUploadURLRequest: { fileType: payload.type },
      }
      const res = await uploadApi().v1UploadCoordinatorsThumbnailPost(params)
      return await fileUpload(uploadApi(), payload, res.key, res.url)
    }
    catch (err) {
      return errorHandler(err, { 400: 'ファイルのアップロードに失敗しました' })
    }
  }

  async function uploadCoordinatorHeader(payload: File): Promise<string> {
    try {
      const params: V1UploadCoordinatorsHeaderPostRequest = {
        getUploadURLRequest: { fileType: payload.type },
      }
      const res = await uploadApi().v1UploadCoordinatorsHeaderPost(params)
      return await fileUpload(uploadApi(), payload, res.key, res.url)
    }
    catch (err) {
      return errorHandler(err, { 400: 'ファイルのアップロードに失敗しました' })
    }
  }

  async function uploadCoordinatorPromotionVideo(payload: File): Promise<string> {
    try {
      const params: V1UploadCoordinatorsPromotionVideoPostRequest = {
        getUploadURLRequest: { fileType: payload.type },
      }
      const res = await uploadApi().v1UploadCoordinatorsPromotionVideoPost(params)
      return await fileUpload(uploadApi(), payload, res.key, res.url)
    }
    catch (err) {
      return errorHandler(err, { 400: 'ファイルのアップロードに失敗しました' })
    }
  }

  async function uploadCoordinatorBonusVideo(payload: File): Promise<string> {
    try {
      const params: V1UploadCoordinatorsBonusVideoPostRequest = {
        getUploadURLRequest: { fileType: payload.type },
      }
      const res = await uploadApi().v1UploadCoordinatorsBonusVideoPost(params)
      return await fileUpload(uploadApi(), payload, res.key, res.url)
    }
    catch (err) {
      return errorHandler(err, { 400: 'ファイルのアップロードに失敗しました' })
    }
  }

  async function deleteCoordinator(id: string) {
    try {
      const params: V1CoordinatorsCoordinatorIdDeleteRequest = { coordinatorId: id }
      await coordinatorApi().v1CoordinatorsCoordinatorIdDelete(params)
    }
    catch (err) {
      return errorHandler(err, {
        400: '必須項目が不足しているか、内容に誤りがあります',
        404: '対象のコーディネーターが存在しません',
      })
    }
    fetchCoordinators()
  }

  return {
    coordinator,
    coordinators,
    producers,
    totalItems,
    producerTotalItems,
    fetchCoordinators,
    searchCoordinators,
    getCoordinator,
    createCoordinator,
    updateCoordinator,
    uploadCoordinatorThumbnail,
    uploadCoordinatorHeader,
    uploadCoordinatorPromotionVideo,
    uploadCoordinatorBonusVideo,
    deleteCoordinator,
  }
})
