import { defineStore } from 'pinia'

import { useCoordinatorStore } from './coordinator'
import { apiClient } from '~/plugins/api-client'
import type { ApproveScheduleRequest, CreateScheduleRequest, PublishScheduleRequest, Schedule, UpdateScheduleRequest, UploadImageResponse, UploadVideoResponse } from '~/types/api'
import { uploadTimeout } from '~/plugins/axios'

export const useScheduleStore = defineStore('schedule', {
  state: () => ({
    schedule: {} as Schedule,
    schedules: [] as Schedule[],
    total: 0
  }),

  actions: {
    /**
     * マルシェ開催スケジュール一覧を取得する非同期関数
     * @param limit
     * @param offset
     */
    async fetchSchedules (limit = 20, offset = 0): Promise<void> {
      try {
        const res = await apiClient.scheduleApi().v1ListSchedules(limit, offset)

        const coordinatorStore = useCoordinatorStore()
        this.schedules = res.data.schedules
        this.total = res.data.total
        coordinatorStore.coordinators = res.data.coordinators
      } catch (err) {
        return this.errorHandler(err)
      }
    },

    /**
     * マルシェ開催スケジュールを取得する非同期関数
     * @param scheduleId スケジュールID
     * @returns
     */
    async getSchedule (scheduleId: string): Promise<void> {
      try {
        const res = await apiClient.scheduleApi().v1GetSchedule(scheduleId)

        const coordinatorStore = useCoordinatorStore()
        this.schedule = res.data.schedule
        coordinatorStore.coordinators.push(res.data.coordinator)
      } catch (err) {
        return this.errorHandler(err)
      }
    },

    /**
     * マルシェ開催スケジュールを登録する非同期関数
     * @param payload
     */
    async createSchedule (payload: CreateScheduleRequest): Promise<Schedule> {
      try {
        const res = await apiClient.scheduleApi().v1CreateSchedule(payload)
        return res.data.schedule
      } catch (err) {
        return this.errorHandler(err)
      }
    },

    /**
     * マルシェ開催スケジュールを更新する非同期関数
     * @param scheduleId スケジュールID
     * @param payload
     */
    async updateSchedule (scheduleId: string, payload: UpdateScheduleRequest): Promise<void> {
      try {
        await apiClient.scheduleApi().v1UpdateSchedule(scheduleId, payload)
      } catch (err) {
        return this.errorHandler(err)
      }
    },

    /**
     * マルシェ開催スケジュールの承認/却下をする非同期関数
     * @param schedule スケジュール
     * @returns
     */
    async approveSchedule (schedule: Schedule): Promise<void> {
      try {
        const req: ApproveScheduleRequest = { approved: !schedule.approved }
        await apiClient.scheduleApi().v1ApproveSchedule(schedule.id, req)
      } catch (err) {
        return this.errorHandler(err)
      }
    },

    /**
     * マルシェ開催スケジュールの公開/非公開をする非同期関数
     * @param scheduleId スケジュールID
     * @param public 公開フラグ
     * @returns
     */
    async publishSchedule (scheduleId: string, published: boolean): Promise<void> {
      try {
        const req: PublishScheduleRequest = { public: published }
        await apiClient.scheduleApi().v1PublishSchedule(scheduleId, req)
      } catch (err) {
        return this.errorHandler(err)
      }
    },

    /**
     * サムネイル画像をアップロードする非同期関数
     * @param payload
     * @returns アップロード先URL
     */
    async uploadScheduleThumbnail (payload: File): Promise<UploadImageResponse> {
      try {
        const res = await apiClient.scheduleApi().v1UploadScheduleThumbnail(
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
        return this.errorHandler(err, { 400: 'このファイルはアップロードできません。' })
      }
    },

    /**
     * 蓋絵画像をアップロードする非同期関数
     * @param payload
     * @returns アップロード先URL
     */
    async uploadScheduleImage (payload: File): Promise<UploadImageResponse> {
      try {
        const res = await apiClient.scheduleApi().v1UploadScheduleImage(
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
        return this.errorHandler(err, { 400: 'このファイルはアップロードできません。' })
      }
    },

    /**
     * オープニング動画をアップロードする非同期関数
     * @param payload
     * @returns アップロード先URL
     */
    async uploadScheduleOpeningVideo (payload: File): Promise<UploadVideoResponse> {
      try {
        const res = await apiClient.scheduleApi().v1UploadScheduleOpeningVideo(
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
        return this.errorHandler(err, { 400: 'このファイルはアップロードできません。' })
      }
    }
  }
})
