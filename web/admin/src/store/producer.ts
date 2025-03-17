import { defineStore } from 'pinia'
import { fileUpload } from './helper'
import { useCoordinatorStore } from './coordinator'
import { useShopStore } from './shop'
import { apiClient } from '~/plugins/api-client'
import type {
  CreateProducerRequest,
  GetUploadUrlRequest,
  ProducerResponse,
  Producer,
  UpdateProducerRequest,
  Shop,
} from '~/types/api'

export const useProducerStore = defineStore('producer', {
  state: () => ({
    producer: {} as Producer,
    producers: [] as Producer[],
    totalItems: 0,
  }),

  actions: {
    /**
     * 登録済みの生産者一覧を取得する非同期関数
     * @param limit 取得上限数
     * @param offset 取得開始位置
     */
    async fetchProducers(limit = 20, offset = 0, options = ''): Promise<void> {
      try {
        const res = await apiClient.producerApi().v1ListProducers(limit, offset, options)

        const coordinatorStore = useCoordinatorStore()
        const shopStore = useShopStore()
        this.producers = res.data.producers
        this.totalItems = res.data.total
        coordinatorStore.coordinators = res.data.coordinators
        shopStore.shops = res.data.shops
      }
      catch (err) {
        return this.errorHandler(err)
      }
    },

    /**
     * 生産者を検索する非同期関数
     * @param name 生産者名(あいまい検索)
     * @param producerIds stateの更新時に残しておく必要がある生産者情報
     */
    async searchProducers(name = '', producerIds: string[] = []): Promise<void> {
      try {
        const res = await apiClient.producerApi().v1ListProducers(undefined, undefined, name)
        const producers: Producer[] = []
        this.producers.forEach((producer: Producer): void => {
          if (!producerIds.includes(producer.id)) {
            return
          }
          producers.push(producer)
        })
        res.data.producers.forEach((producer: Producer): void => {
          if (producers.find((v): boolean => v.id === producer.id)) {
            return
          }
          producers.push(producer)
        })
        const shopStore = useShopStore()
        this.producers = producers
        this.totalItems = res.data.total
        shopStore.shops = res.data.shops
      }
      catch (err) {
        return this.errorHandler(err)
      }
    },

    /**
     * 生産者IDから生産者の情報を取得する非同期関数
     * @param producerId 生産者ID
     * @returns 生産者の情報
     */
    async getProducer(producerId: string): Promise<ProducerResponse> {
      try {
        const res = await apiClient.producerApi().v1GetProducer(producerId)

        const coordinatorStore = useCoordinatorStore()
        const shopStore = useShopStore()
        this.producer = res.data.producer
        coordinatorStore.coordinators = res.data.coordinators
        shopStore.shops = res.data.shops
        return res.data
      }
      catch (err) {
        return this.errorHandler(err, {
          403: '生産者の情報は閲覧権限がありません。',
          404: 'この生産者は存在しません。',
        })
      }
    },

    /**
     * 生産者を新規登録する非同期関数
     * @param payload
     */
    async createProducer(payload: CreateProducerRequest): Promise<void> {
      try {
        await apiClient.producerApi().v1CreateProducer(payload)
      }
      catch (err) {
        return this.errorHandler(err, {
          400: '必須項目が不足しているか、内容に誤りがあります',
          409: 'このメールアドレスはすでに登録されているため、登録できません。',
        })
      }
    },

    /**
     * 生産者のサムネイル画像をアップロードする関数
     * @param payload サムネイル画像のファイルオブジェクト
     * @returns アップロード後のサムネイル画像のパスを含んだオブジェクト
     */
    async uploadProducerThumbnail(payload: File): Promise<string> {
      const contentType = payload.type
      try {
        const body: GetUploadUrlRequest = {
          fileType: contentType,
        }
        const res = await apiClient.producerApi().v1GetProducerThumbnailUploadUrl(body)

        return await fileUpload(payload, res.data.key, res.data.url, res.data.headers)
      }
      catch (err) {
        return this.errorHandler(err, { 400: 'このファイルはアップロードできません。' })
      }
    },

    /**
     * 生産者のヘッダー画像をアップロードする関数
     * @param payload ヘッダー画像のファイルオブジェクト
     * @returns アップロード後のヘッダー画像のパスを含んだオブジェクト
     */
    async uploadProducerHeader(payload: File): Promise<string> {
      const contentType = payload.type
      try {
        const body: GetUploadUrlRequest = {
          fileType: contentType,
        }
        const res = await apiClient.producerApi().v1GetProducerHeaderUploadUrl(body)

        return await fileUpload(payload, res.data.key, res.data.url, res.data.headers)
      }
      catch (err) {
        return this.errorHandler(err, { 400: 'このファイルはアップロードできません。' })
      }
    },

    /**
     * 生産者の紹介画像をアップロードする非同期関数
     * @param payload 紹介画像
     * @returns アップロードされた動画のURI
     */
    async uploadProducerPromotionVideo(payload: File): Promise<string> {
      const contentType = payload.type
      try {
        const body: GetUploadUrlRequest = {
          fileType: contentType,
        }
        const res = await apiClient.producerApi().v1GetProducerPromotionVideoUploadUrl(body)

        return await fileUpload(payload, res.data.key, res.data.url, res.data.headers)
      }
      catch (err) {
        return this.errorHandler(err, { 400: 'このファイルはアップロードできません。' })
      }
    },

    /**
     * 生産者のサンキュー画像をアップロードする非同期関数
     * @param payload サンキュー画像
     * @returns アップロードされた動画のURI
     */
    async uploadProducerBonusVideo(payload: File): Promise<string> {
      const contentType = payload.type
      try {
        const body: GetUploadUrlRequest = {
          fileType: contentType,
        }
        const res = await apiClient.producerApi().v1GetProducerBonusVideoUploadUrl(body)

        return await fileUpload(payload, res.data.key, res.data.url, res.data.headers)
      }
      catch (err) {
        return this.errorHandler(err, { 400: 'このファイルはアップロードできません。' })
      }
    },

    /**
     * 生産者を更新する非同期関数
     * @param producerId 更新対象の生産者ID
     * @param payload
     * @returns
     */
    async updateProducer(producerId: string, payload: UpdateProducerRequest) {
      try {
        await apiClient.producerApi().v1UpdateProducer(producerId, payload)
      }
      catch (err) {
        return this.errorHandler(err, {
          403: '生産者の情報を更新する権限がありません。',
          404: 'この生産者は存在しません。',
        })
      }
    },

    /**
     * 生産者を削除する非同期関数
     * @param producerId 削除する生産者のID
     * @returns
     */
    async deleteProducer(producerId: string) {
      try {
        await apiClient.producerApi().v1DeleteProducer(producerId)
      }
      catch (err) {
        return this.errorHandler(err, {
          400: '必須項目が不足しているか、内容に誤りがあります',
          403: '生産者を削除する権限がありません。',
          404: 'この生産者は存在しません。',
        })
      }
    },
  },
})
