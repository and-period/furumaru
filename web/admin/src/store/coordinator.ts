import { defineStore } from 'pinia'

import { useProductTypeStore } from './product-type'
import { apiClient } from '~/plugins/api-client'
import type {
  Coordinator,
  CreateCoordinatorRequest,
  Producer,
  UpdateCoordinatorRequest,
  UploadImageResponse,
  UploadVideoResponse
} from '~/types/api'
import { uploadTimeout } from '~/plugins/axios'

export const useCoordinatorStore = defineStore('coordinator', {
  state: () => ({
    coordinator: {} as Coordinator,
    coordinators: [] as Coordinator[],
    producers: [] as Producer[],
    totalItems: 0,
    producerTotalItems: 0
  }),

  actions: {
    /**
     * コーディネーターの一覧を取得する非同期関数
     * @param limit 最大取得件数
     * @param offset 取得開始位置
     */
    async fetchCoordinators (limit = 20, offset = 0): Promise<void> {
      try {
        const res = await apiClient.coordinatorApi().v1ListCoordinators(limit, offset)

        const productTypeStore = useProductTypeStore()
        this.coordinators = res.data.coordinators
        this.totalItems = res.data.total
        productTypeStore.productTypes = res.data.productTypes
      } catch (err) {
        return this.errorHandler(err)
      }
    },

    /**
     * コーディネーターを検索する非同期関数
     * @param name コーディネーター名(あいまい検索)
     * @param coordinatorIds stateの更新時に残しておく必要があるコーディネーター情報
     */
    async searchCoordinators (name = '', coordinatorIds: string[] = []): Promise<void> {
      try {
        const res = await apiClient.coordinatorApi().v1ListCoordinators(undefined, undefined, name)
        const coordinators: Coordinator[] = []
        this.coordinators.forEach((coordinator: Coordinator): void => {
          if (!coordinatorIds.includes(coordinator.id)) {
            return
          }
          coordinators.push(coordinator)
        })
        res.data.coordinators.forEach((coordinator: Coordinator): void => {
          if (coordinators.find((v): boolean => v.id === coordinator.id)) {
            return
          }
          coordinators.push(coordinator)
        })
        this.coordinators = coordinators
        this.totalItems = res.data.total
      } catch (err) {
        return this.errorHandler(err)
      }
    },

    /**
     * コーディネーターの詳細情報を取得する非同期関数
     * @param coordinatorId 対象のコーディネーターのID
     */
    async getCoordinator (coordinatorId: string): Promise<void> {
      try {
        const res = await apiClient.coordinatorApi().v1GetCoordinator(coordinatorId)

        const productTypeStore = useProductTypeStore()
        this.coordinator = res.data.coordinator
        productTypeStore.productTypes = res.data.productTypes
      } catch (err) {
        return this.errorHandler(err)
      }
    },

    /**
     * コーディネーターを登録する非同期関数
     * @param payload
     */
    async createCoordinator (payload: CreateCoordinatorRequest) {
      try {
        const res = await apiClient.coordinatorApi().v1CreateCoordinator(payload)
        return res.data
      } catch (err) {
        return this.errorHandler(err)
      }
    },

    /**
     * コーディネーターの情報を更新する非同期関数
     * @param payload
     * @param coordinatorId 更新するコーディネーターのID
     */
    async updateCoordinator (coordinatorId: string, payload: UpdateCoordinatorRequest): Promise<void> {
      try {
        await apiClient.coordinatorApi().v1UpdateCoordinator(coordinatorId, payload)
      } catch (err) {
        return this.errorHandler(err)
      }
    },

    /**
     * コーディネーターのサムネイル画像をアップロードする非同期関数
     * @param payload サムネイル画像
     * @returns アップロードされた画像のURI
     */
    async uploadCoordinatorThumbnail (payload: File): Promise<UploadImageResponse> {
      try {
        const res = await apiClient.coordinatorApi().v1UploadCoordinatorThumbnail(
          payload,
          {
            headers: {
              'Content-Type': 'multipart/form-data'
            },
            timeout: uploadTimeout
          }
        )
        return res.data
      } catch (err) {
        return this.errorHandler(err)
      }
    },

    /**
     * コーディネーターのヘッダー画像をアップロードする非同期関数
     * @param payload ヘッダー画像
     * @returns アップロードされた画像のURI
     */
    async uploadCoordinatorHeader (payload: File): Promise<UploadImageResponse> {
      try {
        const res = await apiClient.coordinatorApi().v1UploadCoordinatorHeader(
          payload,
          {
            headers: {
              'Content-Type': 'multipart/form-data'
            },
            timeout: uploadTimeout
          }
        )
        return res.data
      } catch (err) {
        return this.errorHandler(err)
      }
    },

    /**
     * コーディネーターの紹介画像をアップロードする非同期関数
     * @param payload 紹介画像
     * @returns アップロードされた動画のURI
     */
    async uploadCoordinatorPromotionVideo (payload: File): Promise<UploadVideoResponse> {
      try {
        const res = await apiClient.coordinatorApi().v1UploadCoordinatorPromotionVideo(
          payload,
          {
            headers: {
              'Content-Type': 'multipart/form-data'
            },
            timeout: uploadTimeout
          }
        )
        return res.data
      } catch (err) {
        return this.errorHandler(err)
      }
    },

    /**
     * コーディネーターのサンキュー画像をアップロードする非同期関数
     * @param payload サンキュー画像
     * @returns アップロードされた動画のURI
     */
    async uploadCoordinatorBonusVideo (payload: File): Promise<UploadVideoResponse> {
      try {
        const res = await apiClient.coordinatorApi().v1UploadCoordinatorBonusVideo(
          payload,
          {
            headers: {
              'Content-Type': 'multipart/form-data'
            },
            timeout: uploadTimeout
          }
        )
        return res.data
      } catch (err) {
        return this.errorHandler(err)
      }
    },

    /**
     * コーディーネータを削除する非同期関数
     * @param id 削除するコーディネーターのID
     * @returns
     */
    async deleteCoordinator (id: string) {
      try {
        await apiClient.coordinatorApi().v1DeleteCoordinator(id)
      } catch (err) {
        return this.errorHandler(err)
      }
      this.fetchCoordinators()
    }
  }
})
