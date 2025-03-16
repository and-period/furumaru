import { defineStore } from 'pinia'

import { useProducerStore } from './producer'
import { useProductStore } from './product'
import { apiClient } from '~/plugins/api-client'
import type { CreateLiveRequest, Live, UpdateLiveRequest } from '~/types/api'

export const useLiveStore = defineStore('live', {
  state: () => ({
    lives: [] as Live[],
    total: 0,
  }),

  actions: {
    /**
     * ライブ配信スケジュール一覧を取得する非同期関数
     * @param scheduleId 開催スケジュールID
     * @returns
     */
    async fetchLives(scheduleId: string): Promise<void> {
      try {
        const res = await apiClient.liveApi().v1ListLives(scheduleId)

        const producerStore = useProducerStore()
        const productStore = useProductStore()
        this.lives = res.data.lives
        this.total = res.data.total
        producerStore.producers = res.data.producers
        productStore.products = res.data.products
      }
      catch (err) {
        return this.errorHandler(err, { 404: 'マルシェタイムテーブルが存在しません' })
      }
    },

    /**
     * ライブ配信コメント一覧を取得する非同期関数
     * @param scheduleId 開催スケジュールID
     */
    async fetchLiveComments(scheduleId: string): Promise<void> {},

    /**
     * ライブ配信スケジュールを登録する非同期関数
     * @param scheduleId 開催スケジュールID
     * @param payload
     * @returns
     */
    async createLive(scheduleId: string, payload: CreateLiveRequest): Promise<void> {
      try {
        const res = await apiClient.liveApi().v1CreateLive(scheduleId, payload)

        this.lives.push(res.data.live)
      }
      catch (err) {
        return this.errorHandler(err, {
          400: '必須項目が不足しているか、内容に誤りがあります',
          404: 'マルシェタイムテーブルが存在しません',
          412: '開催時間, 生産者が重複しています',
        })
      }
    },

    /**
     * ライブ配信スケジュールを更新する非同期関数
     * @param scheduleId 開催スケジュールID
     * @param liveId
     * @param payload
     * @returns
     */
    async updateLive(scheduleId: string, liveId: string, payload: UpdateLiveRequest): Promise<void> {
      try {
        await apiClient.liveApi().v1UpdateLive(scheduleId, liveId, payload)

        const index = this.lives.findIndex((live: Live): boolean => live.id === liveId)
        const live = this.lives[index]
        this.lives.splice(index, 1, { ...live, ...payload })
      }
      catch (err) {
        return this.errorHandler(err, {
          400: '必須項目が不足しているか、内容に誤りがあります',
          404: 'マルシェ開催スケジュール, マルシェタイムテーブルが存在しません',
          412: '開催時間, 生産者が重複しています',
        })
      }
    },

    /**
     * ライブ配信スケジュールを削除する非同期関数
     * @param scheduleId 開催スケジュールID
     * @param liveId
     * @returns
     */
    async deleteLive(scheduleId: string, liveId: string): Promise<void> {
      try {
        await apiClient.liveApi().v1DeleteLive(scheduleId, liveId)

        const index = this.lives.findIndex((live: Live): boolean => live.id === liveId)
        this.lives.splice(index, 1)
      }
      catch (err) {
        return this.errorHandler(err, {
          400: '必須項目が不足しているか、内容に誤りがあります',
          404: 'マルシェ開催スケジュール, マルシェタイムテーブルが存在しません',
        })
      }
    },
  },
})
