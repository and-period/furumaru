import { defineStore } from 'pinia'

import { useCommonStore } from './common'
import {
  CreateProducerRequest,
  ProducerResponse,
  ProducersResponse,
  UpdateProducerRequest,
  UploadImageResponse,
  UploadVideoResponse
} from '~/types/api'
import { apiClient } from '~/plugins/api-client'

export const useProducerStore = defineStore('producer', {
  state: () => ({
    producer: {} as ProducerResponse,
    producers: [] as ProducersResponse['producers'],
    totalItems: 0
  }),

  actions: {
    /**
     * 登録済みの生産者一覧を取得する非同期関数
     * @param limit 取得上限数
     * @param offset 取得開始位置
     */
    async fetchProducers (limit = 20, offset = 0, options = ''): Promise<void> {
      try {
        const res = await apiClient.producerApi().v1ListProducers(limit, offset, options)
        this.producers = res.data.producers
        this.totalItems = res.data.total
      } catch (err) {
        return this.errorHandler(err)
      }
    },

    /**
     * 生産者IDから生産者の情報を取得する非同期関数
     * @param producerId 生産者ID
     * @returns 生産者の情報
     */
    async getProducer (producerId: string): Promise<ProducerResponse> {
      try {
        const res = await apiClient.producerApi().v1GetProducer(producerId)
        this.producer = res.data
        return res.data
      } catch (err) {
        return this.errorHandler(err)
      }
    },

    /**
     * 生産者を新規登録する非同期関数
     * @param payload
     */
    async createProducer (payload: CreateProducerRequest): Promise<void> {
      try {
        await apiClient.producerApi().v1CreateProducer(payload)
        const commonStore = useCommonStore()
        commonStore.addSnackbar({
          message: `${payload.username}を作成しました。`,
          color: 'info'
        })
      } catch (err) {
        return this.errorHandler(err, { 409: 'このメールアドレスはすでに登録されているため、登録できません。' })
      }
    },

    /**
     * 生産者のサムネイル画像をアップロードする関数
     * @param payload サムネイル画像のファイルオブジェクト
     * @returns アップロード後のサムネイル画像のパスを含んだオブジェクト
     */
    async uploadProducerThumbnail (payload: File): Promise<UploadImageResponse> {
      try {
        const res = await apiClient.producerApi().v1UploadProducerThumbnail(
          payload,
          {
            headers: {
              'Content-Type': 'multipart/form-data'
            }
          }
        )
        return res.data
      } catch (err) {
        return this.errorHandler(err, { 400: 'このファイルはアップロードできません。' })
      }
    },

    /**
     * 生産者のヘッダー画像をアップロードする関数
     * @param payload ヘッダー画像のファイルオブジェクト
     * @returns アップロード後のヘッダー画像のパスを含んだオブジェクト
     */
    async uploadProducerHeader (payload: File): Promise<UploadImageResponse> {
      try {
        const res = await apiClient.producerApi().v1UploadProducerHeader(
          payload,
          {
            headers: {
              'Content-Type': 'multipart/form-data'
            }
          }
        )
        return res.data
      } catch (err) {
        return this.errorHandler(err, { 400: 'このファイルはアップロードできません。' })
      }
    },

    /**
     * 生産者の紹介画像をアップロードする非同期関数
     * @param payload 紹介画像
     * @returns アップロードされた動画のURI
     */
    async uploadProducerPromotionVideo (payload: File): Promise<UploadVideoResponse> {
      try {
        const res = await apiClient.producerApi().v1UploadProducerPromotionVideo(
          payload,
          {
            headers: {
              'Content-Type': 'multipart/form-data'
            }
          }
        )
        return res.data
      } catch (err) {
        return this.errorHandler(err)
      }
    },

    /**
     * 生産者のサンキュー画像をアップロードする非同期関数
     * @param payload サンキュー画像
     * @returns アップロードされた動画のURI
     */
    async uploadProducerBonusVideo (payload: File): Promise<UploadVideoResponse> {
      try {
        const res = await apiClient.producerApi().v1UploadProducerBonusVideo(
          payload,
          {
            headers: {
              'Content-Type': 'multipart/form-data'
            }
          }
        )
        return res.data
      } catch (err) {
        return this.errorHandler(err)
      }
    },

    /**
     * 生産者を更新する非同期関数
     * @param producerId 更新対象の生産者ID
     * @param payload
     * @returns
     */
    async updateProducer (producerId: string, payload: UpdateProducerRequest) {
      try {
        await apiClient.producerApi().v1UpdateProducer(producerId, payload)
      } catch (err) {
        return this.errorHandler(err)
      }
    },

    /**
     * 生産者を削除する非同期関数
     * @param producerId 削除する生産者のID
     * @returns
     */
    async deleteProducer (producerId: string) {
      try {
        await apiClient.producerApi().v1DeleteProducer(producerId)
      } catch (err) {
        return this.errorHandler(err)
      }
    }
  }
})
