import { defineStore } from 'pinia'

import { apiClient } from '~/plugins/api-client'
import type { CreateSpotTypeRequest, SpotType, UpdateSpotTypeRequest } from '~/types/api'

export const useSpotTypeStore = defineStore('spotType', {
  state: () => ({
    spotType: {} as SpotType,
    spotTypes: [] as SpotType[],
    total: 0,
  }),

  actions: {
    /**
     * スポット種別一覧を取得する非同期関数
     * @param limit 取得上限数
     * @param offset 取得開始位置
     */
    async fetchSpotTypes(limit = 20, offset = 0): Promise<void> {
      try {
        const res = await apiClient.spotTypeApi().v1ListSpotTypes(limit, offset, '')
        this.spotTypes = res.data.spotTypes
        this.total = res.data.total
      }
      catch (err) {
        return this.errorHandler(err)
      }
    },

    /**
     * スポット種別を検索する非同期関数
     * @param name スポット種別名(あいまい検索)
     * @param spotTypeIds stateの更新時に残しておく必要があるスポット種別情報
     */
    async searchSpotTypes(name = '', spotTypeIds: string[] = []): Promise<void> {
      try {
        const res = await apiClient.spotTypeApi().v1ListSpotTypes(undefined, undefined, name)
        const spotTypes: SpotType[] = []
        this.spotTypes.forEach((spotType: SpotType): void => {
          if (!spotTypeIds.includes(spotType.id)) {
            return
          }
          spotTypes.push(spotType)
        })
        res.data.spotTypes.forEach((spotType: SpotType): void => {
          if (spotTypes.find((v): boolean => v.id === spotType.id)) {
            return
          }
          spotTypes.push(spotType)
        })
        this.spotTypes = spotTypes
        this.total = res.data.total
      }
      catch (err) {
        return this.errorHandler(err)
      }
    },

    /**
     * スポット種別を新規登録する非同期関数
     * @param payload
     */
    async createSpotType(payload: CreateSpotTypeRequest): Promise<void> {
      try {
        const res = await apiClient.spotTypeApi().v1CreateSpotType(payload)
        this.spotTypes.unshift(res.data.spotType)
      }
      catch (err) {
        return this.errorHandler(err, {
          400: '必須項目が不足しているか、内容に誤りがあります。',
          409: 'このスポット名はすでに登録されています。',
        })
      }
    },

    /**
     * スポット種別を更新する非同期関数
     * @param spotTypeId スポット種別ID
     * @param payload
     */
    async updateSpotType(spotTypeId: string, payload: UpdateSpotTypeRequest): Promise<void> {
      try {
        await apiClient.spotTypeApi().v1UpdateSpotType(spotTypeId, payload)
      }
      catch (err) {
        return this.errorHandler(err, {
          400: '必須項目が不足しているか、内容に誤りがあります。',
          404: 'このスポット種別は存在しません。',
          409: 'このスポット種別名はすでに登録されています。',
        })
      }
    },

    /**
     * スポット種別を削除する非同期関数
     * @param spotTypeId スポット種別ID
     */
    async deleteSpotType(spotTypeId: string): Promise<void> {
      try {
        await apiClient.spotTypeApi().v1DeleteSpotType(spotTypeId)
      }
      catch (err) {
        return this.errorHandler(err, { 404: 'このスポット種別は存在しません。' })
      }
    },
  },
})
