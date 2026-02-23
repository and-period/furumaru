import { useApiClient } from '~/composables/useApiClient'
import { fileUpload } from './helper'
import { useCoordinatorStore } from './coordinator'
import { useShopStore } from './shop'
import { ProducerApi, UploadApi } from '~/types/api/v1'
import type {
  CreateProducerRequest,
  ProducerResponse,
  Producer,
  UpdateProducerRequest,
  V1ProducersGetRequest,
  V1ProducersProducerIdGetRequest,
  V1ProducersPostRequest,
  V1UploadProducersThumbnailPostRequest,
  V1UploadProducersHeaderPostRequest,
  V1UploadProducersPromotionVideoPostRequest,
  V1UploadProducersBonusVideoPostRequest,
  V1ProducersProducerIdPatchRequest,
  V1ProducersProducerIdDeleteRequest,
} from '~/types/api/v1'

export const useProducerStore = defineStore('producer', () => {
  const { create, errorHandler } = useApiClient()
  const producerApi = () => create(ProducerApi)
  const uploadApi = () => create(UploadApi)

  const producer = ref<Producer>({} as Producer)
  const producers = ref<Producer[]>([])
  const totalItems = ref<number>(0)

  async function fetchProducers(limit = 20, offset = 0, options = ''): Promise<void> {
    try {
      const params: V1ProducersGetRequest = { limit, offset }
      const res = await producerApi().v1ProducersGet(params)

      const coordinatorStore = useCoordinatorStore()
      const shopStore = useShopStore()
      producers.value = res.producers
      totalItems.value = res.total
      coordinatorStore.coordinators = res.coordinators
      shopStore.shops = res.shops
    }
    catch (err) {
      return errorHandler(err)
    }
  }

  async function searchProducers(name = '', producerIds: string[] = []): Promise<void> {
    try {
      const params: V1ProducersGetRequest = { username: name }
      const res = await producerApi().v1ProducersGet(params)
      const merged: Producer[] = []
      producers.value.forEach((p: Producer): void => {
        if (!producerIds.includes(p.id)) {
          return
        }
        merged.push(p)
      })
      res.producers.forEach((p: Producer): void => {
        if (merged.find((v): boolean => v.id === p.id)) {
          return
        }
        merged.push(p)
      })
      const shopStore = useShopStore()
      producers.value = merged
      totalItems.value = res.total
      shopStore.shops = res.shops
    }
    catch (err) {
      return errorHandler(err)
    }
  }

  async function getProducer(producerId: string): Promise<ProducerResponse> {
    try {
      const params: V1ProducersProducerIdGetRequest = { producerId }
      const res = await producerApi().v1ProducersProducerIdGet(params)

      const coordinatorStore = useCoordinatorStore()
      const shopStore = useShopStore()
      producer.value = res.producer
      coordinatorStore.coordinators = res.coordinators
      shopStore.shops = res.shops
      return res
    }
    catch (err) {
      return errorHandler(err, {
        403: '生産者の情報は閲覧権限がありません。',
        404: 'この生産者は存在しません。',
      })
    }
  }

  async function createProducer(payload: CreateProducerRequest): Promise<void> {
    try {
      const params: V1ProducersPostRequest = { createProducerRequest: payload }
      await producerApi().v1ProducersPost(params)
    }
    catch (err) {
      return errorHandler(err, {
        400: '必須項目が不足しているか、内容に誤りがあります',
        409: 'このメールアドレスはすでに登録されているため、登録できません。',
      })
    }
  }

  async function uploadProducerThumbnail(payload: File): Promise<string> {
    const contentType = payload.type
    try {
      const params: V1UploadProducersThumbnailPostRequest = {
        getUploadURLRequest: { fileType: contentType },
      }
      const res = await uploadApi().v1UploadProducersThumbnailPost(params)
      return await fileUpload(uploadApi(), payload, res.key, res.url)
    }
    catch (err) {
      return errorHandler(err, { 400: 'このファイルはアップロードできません。' })
    }
  }

  async function uploadProducerHeader(payload: File): Promise<string> {
    const contentType = payload.type
    try {
      const params: V1UploadProducersHeaderPostRequest = {
        getUploadURLRequest: { fileType: contentType },
      }
      const res = await uploadApi().v1UploadProducersHeaderPost(params)
      return await fileUpload(uploadApi(), payload, res.key, res.url)
    }
    catch (err) {
      return errorHandler(err, { 400: 'このファイルはアップロードできません。' })
    }
  }

  async function uploadProducerPromotionVideo(payload: File): Promise<string> {
    const contentType = payload.type
    try {
      const params: V1UploadProducersPromotionVideoPostRequest = {
        getUploadURLRequest: { fileType: contentType },
      }
      const res = await uploadApi().v1UploadProducersPromotionVideoPost(params)
      return await fileUpload(uploadApi(), payload, res.key, res.url)
    }
    catch (err) {
      return errorHandler(err, { 400: 'このファイルはアップロードできません。' })
    }
  }

  async function uploadProducerBonusVideo(payload: File): Promise<string> {
    const contentType = payload.type
    try {
      const params: V1UploadProducersBonusVideoPostRequest = {
        getUploadURLRequest: { fileType: contentType },
      }
      const res = await uploadApi().v1UploadProducersBonusVideoPost(params)
      return await fileUpload(uploadApi(), payload, res.key, res.url)
    }
    catch (err) {
      return errorHandler(err, { 400: 'このファイルはアップロードできません。' })
    }
  }

  async function updateProducer(producerId: string, payload: UpdateProducerRequest) {
    try {
      const params: V1ProducersProducerIdPatchRequest = {
        producerId,
        updateProducerRequest: payload,
      }
      await producerApi().v1ProducersProducerIdPatch(params)
    }
    catch (err) {
      return errorHandler(err, {
        403: '生産者の情報を更新する権限がありません。',
        404: 'この生産者は存在しません。',
      })
    }
  }

  async function deleteProducer(producerId: string) {
    try {
      const params: V1ProducersProducerIdDeleteRequest = { producerId }
      await producerApi().v1ProducersProducerIdDelete(params)
    }
    catch (err) {
      return errorHandler(err, {
        400: '必須項目が不足しているか、内容に誤りがあります',
        403: '生産者を削除する権限がありません。',
        404: 'この生産者は存在しません。',
      })
    }
  }

  return {
    producer,
    producers,
    totalItems,
    fetchProducers,
    searchProducers,
    getProducer,
    createProducer,
    uploadProducerThumbnail,
    uploadProducerHeader,
    uploadProducerPromotionVideo,
    uploadProducerBonusVideo,
    updateProducer,
    deleteProducer,
  }
})
