import { useProducerStore } from './producer'
import { useProductStore } from './product'
import type { CreateLiveRequest, Live, UpdateLiveRequest, V1SchedulesScheduleIdLivesGetRequest, V1SchedulesScheduleIdLivesLiveIdDeleteRequest, V1SchedulesScheduleIdLivesLiveIdPatchRequest, V1SchedulesScheduleIdLivesPostRequest } from '~/types/api/v1'

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
        const params: V1SchedulesScheduleIdLivesGetRequest = {
          scheduleId,
        }
        const res = await this.liveApi().v1SchedulesScheduleIdLivesGet(params)

        const producerStore = useProducerStore()
        const productStore = useProductStore()
        this.lives = res.lives
        this.total = res.total
        producerStore.producers = res.producers
        productStore.products = res.products
      }
      catch (err) {
        return this.errorHandler(err, {
          404: 'マルシェタイムテーブルが存在しません',
        })
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
    async createLive(
      scheduleId: string,
      payload: CreateLiveRequest,
    ): Promise<void> {
      try {
        const params: V1SchedulesScheduleIdLivesPostRequest = {
          scheduleId,
          createLiveRequest: payload,
        }
        const res = await this.liveApi().v1SchedulesScheduleIdLivesPost(params)

        this.lives.push(res.live)
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
    async updateLive(
      scheduleId: string,
      liveId: string,
      payload: UpdateLiveRequest,
    ): Promise<void> {
      try {
        const params: V1SchedulesScheduleIdLivesLiveIdPatchRequest = {
          scheduleId,
          liveId,
          updateLiveRequest: payload,
        }
        await this.liveApi().v1SchedulesScheduleIdLivesLiveIdPatch(params)

        const index = this.lives.findIndex(
          (live: Live): boolean => live.id === liveId,
        )
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
        const params: V1SchedulesScheduleIdLivesLiveIdDeleteRequest = {
          scheduleId,
          liveId,
        }
        await this.liveApi().v1SchedulesScheduleIdLivesLiveIdDelete(params)

        const index = this.lives.findIndex(
          (live: Live): boolean => live.id === liveId,
        )
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
