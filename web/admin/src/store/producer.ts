import { defineStore } from 'pinia'

import ApiClientFactory from '../plugins/factory'

import { useAuthStore } from './auth'
import { useCommonStore } from './common'

import {
  CreateProducerRequest,
  ProducerApi,
  ProducerResponse,
  ProducersResponse,
  UploadImageResponse,
} from '~/types/api'

export const useProducerStore = defineStore('Producer', {
  state: () => ({
    producers: [] as ProducersResponse['producers'],
  }),
  actions: {
    /**
     * 登録済みの生産者一覧を取得する非同期関数
     * @param limit 取得上限数
     * @param offset 取得開始位置
     */
    async fetchProducers(
      limit: number = 20,
      offset: number = 0
    ): Promise<void> {
      try {
        const authStore = useAuthStore()
        const accessToken = authStore.accessToken
        if (!accessToken) {
          return Promise.reject(new Error('認証エラー'))
        }

        const factory = new ApiClientFactory()
        const producersApiClient = factory.create(ProducerApi, accessToken)
        const res = await producersApiClient.v1ListProducers(limit, offset)
        this.producers = res.data.producers
      } catch (error) {
        // TODO: エラーハンドリング
        throw new Error('Internal Server Error')
      }
    },

    /**
     * 生産者を新規登録する非同期関数
     * @param payload
     */
    async createProducer(payload: CreateProducerRequest): Promise<void> {
      try {
        const authStore = useAuthStore()
        const accessToken = authStore.accessToken
        if (!accessToken) {
          return Promise.reject(new Error('認証エラー'))
        }

        const factory = new ApiClientFactory()
        const producersApiClient = factory.create(ProducerApi, accessToken)
        await producersApiClient.v1CreateProducer(payload)
        const commonStore = useCommonStore()
        commonStore.addSnackbar({
          message: `${payload.storeName}を作成しました。`,
          color: 'info',
        })
      } catch (error) {
        // TODO: エラーハンドリング
        console.log(error)
        throw new Error('Internal Server Error')
      }
    },

    /**
     * 生産者のサムネイル画像をアップロードする関数
     * @param payload サムネイル画像のファイルオブジェクト
     * @returns アップロード後のサムネイル画像のパスを含んだオブジェクト
     */
    async uploadProducerThumbnail(payload: File): Promise<UploadImageResponse> {
      try {
        const authStore = useAuthStore()
        const accessToken = authStore.accessToken
        if (!accessToken) {
          return Promise.reject(new Error('認証エラー'))
        }

        const factory = new ApiClientFactory()
        const producersApiClient = factory.create(ProducerApi, accessToken)
        const res = await producersApiClient.v1UploadProducerThumbnail(
          payload,
          {
            headers: {
              'Content-Type': 'multipart/form-data',
            },
          }
        )
        return res.data
      } catch (error) {
        throw new Error('Internal Server Error')
      }
    },

    /**
     * 生産者のヘッダー画像をアップロードする関数
     * @param payload ヘッダー画像のファイルオブジェクト
     * @returns アップロード後のヘッダー画像のパスを含んだオブジェクト
     */
    async uploadProducerHeader(payload: File): Promise<UploadImageResponse> {
      try {
        const authStore = useAuthStore()
        const accessToken = authStore.accessToken
        if (!accessToken) {
          return Promise.reject(new Error('認証エラー'))
        }

        const factory = new ApiClientFactory()
        const producersApiClient = factory.create(ProducerApi, accessToken)
        const res = await producersApiClient.v1UploadProducerHeader(payload, {
          headers: {
            'Content-Type': 'multipart/form-data',
          },
        })
        return res.data
      } catch (error) {
        throw new Error('Internal Server Error')
      }
    },

    /**
     * 生産者IDから生産者の情報を取得する非同期関数
     * @param id 生産者ID
     * @returns 生産者の情報
     */
    async getProducer(id: string): Promise<ProducerResponse> {
      try {
        const authStore = useAuthStore()
        const accessToken = authStore.accessToken
        if (!accessToken) throw new Error('認証エラー')

        const factory = new ApiClientFactory()
        const producersApiClient = factory.create(ProducerApi, accessToken)
        const res = await producersApiClient.v1GetProducer(id)
        return res.data
      } catch (error) {
        return Promise.reject(new Error('不明なエラーが発生しました。'))
      }
    },
  },
})
