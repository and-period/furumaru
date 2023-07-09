import { defineStore } from 'pinia'
import { useCommonStore } from './common'
import { useCoordinatorStore } from './coordinator'
import { useShippingStore } from './shipping'
import { apiClient } from '~/plugins/api-client'
import { CreateScheduleRequest, ScheduleResponse, SchedulesResponse, UploadImageResponse, UploadVideoResponse } from '~/types/api'

export const useScheduleStore = defineStore('schedule', {
  state: () => ({
    schedules: [] as SchedulesResponse['schedules'],
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
        const res = await apiClient.scheduleApi().v1ListSchedules()

        const coordinatorStore = useCoordinatorStore()
        const shippingStore = useShippingStore()
        this.schedules = res.data.schedules
        this.total = res.data.total
        coordinatorStore.coordinators = res.data.coordinators
        shippingStore.shippings = res.data.shippings
      } catch (err) {
        return this.errorHandler(err)
      }
    },

    /**
     * マルシェ開催スケジュールを登録する非同期関数
     * @param payload
     */
    async createSchedule (payload: CreateScheduleRequest): Promise<ScheduleResponse> {
      try {
        const res = await apiClient.scheduleApi().v1CreateSchedule(payload)
        const commonStore = useCommonStore()
        commonStore.addSnackbar({
          message: `${payload.title}を作成しました。`,
          color: 'info'
        })
        return res.data
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
