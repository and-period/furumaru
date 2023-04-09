import { defineStore } from 'pinia'

import { getAccessToken } from './auth'
import { useCommonStore } from './common'
import {
  CoordinatorResponse,
  CoordinatorsResponse,
  CreateCoordinatorRequest,
  ProducersResponse,
  RelateProducersRequest,
  UpdateCoordinatorRequest,
  UploadImageResponse
} from '~/types/api'

export const useCoordinatorStore = defineStore('Coordinator', {
  state: () => ({
    coordinators: [] as CoordinatorsResponse['coordinators'],
    producers: [] as ProducersResponse['producers'],
    totalItems: 0,
    producerTotalItems: 0,
  }),
  actions: {
    /**
     * コーディネータの一覧を取得する非同期関数
     * @param limit 最大取得件数
     * @param offset 取得開始位置
     * @returns
     */
    async fetchCoordinators (limit = 20, offset = 0): Promise<void> {
      try {
        const accessToken = getAccessToken()
        const res = await this.coordinatorApiClient(accessToken).v1ListCoordinators(
          limit,
          offset
        )
        this.coordinators = res.data.coordinators
        this.totalItems = res.data.total
      } catch (error) {
        return this.errorHandler(error)
      }
    },

    /**
     * コーディネータを登録する非同期関数
     * @param payload
     * @returns
     */
    async createCoordinator (payload: CreateCoordinatorRequest) {
      try {
        const accessToken = getAccessToken()
        const res = await this.coordinatorApiClient(accessToken).v1CreateCoordinator(payload)
        return res.data
      } catch (error) {
        console.log(error)
        return this.errorHandler(error)
      }
    },

    /**
     * コーディネータの詳細情報を取得する非同期関数
     * @param id 対象のコーディネータのID
     * @returns
     */
    async getCoordinator (id: string): Promise<CoordinatorResponse> {
      try {
        const accessToken = getAccessToken()
        const res = await this.coordinatorApiClient(accessToken).v1GetCoordinator(id)
        return res.data
      } catch (error) {
        console.log(error)
        return this.errorHandler(error)
      }
    },

    /**
     * コーディネータの情報を更新する非同期関数
     * @param payload
     * @param coordinatorId 更新するコーディネータのID
     * @returns
     */
    async updateCoordinator (
      payload: UpdateCoordinatorRequest,
      coordinatorId: string
    ): Promise<void> {
      try {
        const accessToken = getAccessToken()
        await this.coordinatorApiClient(accessToken).v1UpdateCoordinator(coordinatorId, payload)
        const commonStore = useCommonStore()
        commonStore.addSnackbar({
          message: 'コーディネータ情報が更新されました。',
          color: 'info'
        })
      } catch (error) {
        console.log(error)
        return this.errorHandler(error)
      }
    },

    /**
     * コーディネータのサムネイル画像をアップロードする非同期関数
     * @param payload サムネイル画像
     * @returns アップロードされた画像のURI
     */
    async uploadCoordinatorThumbnail (
      payload: File
    ): Promise<UploadImageResponse> {
      try {
        const accessToken = getAccessToken()
        const res = await this.coordinatorApiClient(accessToken).v1UploadCoordinatorThumbnail(
          payload,
          {
            headers: {
              'Content-Type': 'multipart/form-data'
            }
          }
        )
        return res.data
      } catch (error) {
        console.log(error)
        return this.errorHandler(error)
      }
    },

    /**
     * コーディネータのヘッダー画像をアップロードする非同期関数
     * @param payload ヘッダー画像
     * @returns アップロードされた画像のURI
     */
    async uploadCoordinatorHeader (payload: File): Promise<UploadImageResponse> {
      try {
        const accessToken = getAccessToken()
        const res = await this.coordinatorApiClient(accessToken).v1UploadCoordinatorHeader(
          payload,
          {
            headers: {
              'Content-Type': 'multipart/form-data'
            }
          }
        )
        return res.data
      } catch (error) {
        console.log(error)
        return this.errorHandler(error)
      }
    },

    /**
     * コーディーネータを削除する非同期関数
     * @param id 削除するコーディネータのID
     * @returns
     */
    async deleteCoordinator (id: string) {
      try {
        const accessToken = getAccessToken()
        await this.coordinatorApiClient(accessToken).v1DeleteCoordinator(id)
        const commonStore = useCommonStore()
        commonStore.addSnackbar({
          message: 'コーディネーターの削除が完了しました',
          color: 'info'
        })
      } catch (error) {
        console.log(error)
        return this.errorHandler(error)
      }
      this.fetchCoordinators()
    },

    /**
     * コーディーネータに生産者を紐づける非同期関数
     * @param id 生産者を紐づけるコーディネータのID
     * @param payload コーディネーターに紐づく生産者
     * @returns
     */
    async relateProducers (
      id: string,
      payload: RelateProducersRequest
    ): Promise<void> {
      try {
        const accessToken = getAccessToken()
        await this.coordinatorApiClient(accessToken).v1RelateProducers(id, payload)
        const commonStore = useCommonStore()
        commonStore.addSnackbar({
          message: 'コーディネーターと生産者の紐付けが完了しました',
          color: 'info'
        })
      } catch (error) {
        console.log(error)
        return this.errorHandler(error)
      }
    },

    /**
     * コーディーネータに紐づいている生産者を取得する非同期関数
     * @param id コーディネータのID
     * @returns
     */
    async fetchRelatedProducers (
      id: string,
      limit = 20,
      offset = 0
    ): Promise<void> {
      try {
        const accessToken = getAccessToken()
        const res = await this.coordinatorApiClient(accessToken).v1ListRelatedProducers(
          id,
          limit,
          offset
        )
        this.producers = res.data.producers
      } catch (error) {
        console.log(error)
        return this.errorHandler(error)
      }
    }
  }
})
