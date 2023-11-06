import { defineStore } from 'pinia'
import { useCommonStore } from './common'
import { useCoordinatorStore } from './coordinator'
import { apiClient } from '~/plugins/api-client'
import { ApproveScheduleRequest, CreateScheduleRequest, Schedule, UpdateScheduleRequest, UploadImageResponse, UploadVideoResponse } from '~/types/api'

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
        const commonStore = useCommonStore()
        commonStore.addSnackbar({
          message: `${payload.title}を作成しました。`,
          color: 'info'
        })
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
        const commonStore = useCommonStore()
        commonStore.addSnackbar({
          message: `${payload.title}を更新しました。`,
          color: 'info'
        })
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
        const commonStore = useCommonStore()
        commonStore.addSnackbar({
          message: `${schedule.title}を更新しました。`,
          color: 'info'
        })
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
            }
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
            }
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
            }
          }
        )
        return res.data
      } catch (err) {
        return this.errorHandler(err, { 400: 'このファイルはアップロードできません。' })
      }
    }
  }
})
