import { fileUpload } from './helper'
import { useCoordinatorStore } from './coordinator'
import type {
  BroadcastViewerLog,
  CreateScheduleRequest,
  Schedule,
  UpdateScheduleRequest,
  V1SchedulesGetRequest,
  V1SchedulesPostRequest,
  V1SchedulesScheduleIdAnalyticsGetRequest,
  V1SchedulesScheduleIdApprovalPatchRequest,
  V1SchedulesScheduleIdDeleteRequest,
  V1SchedulesScheduleIdGetRequest,
  V1SchedulesScheduleIdPatchRequest,
  V1SchedulesScheduleIdPublishPatchRequest,
  V1UploadSchedulesImagePostRequest,
  V1UploadSchedulesOpeningVideoPostRequest,
  V1UploadSchedulesThumbnailPostRequest,
} from '~/types/api/v1'

export const useScheduleStore = defineStore('schedule', {
  state: () => ({
    schedule: {} as Schedule,
    schedules: [] as Schedule[],
    viewerLogs: [] as BroadcastViewerLog[],
    total: 0,
  }),

  actions: {
    /**
     * マルシェ開催スケジュール一覧を取得する非同期関数
     * @param limit
     * @param offset
     */
    async fetchSchedules(limit = 20, offset = 0): Promise<void> {
      try {
        const params: V1SchedulesGetRequest = {
          limit,
          offset,
        }
        const res = await this.scheduleApi().v1SchedulesGet(params)

        const coordinatorStore = useCoordinatorStore()
        this.schedules = res.schedules
        this.total = res.total
        coordinatorStore.coordinators = res.coordinators
      }
      catch (err) {
        return this.errorHandler(err)
      }
    },

    /**
     * マルシェ開催スケジュールを取得する非同期関数
     * @param scheduleId スケジュールID
     * @returns
     */
    async getSchedule(scheduleId: string): Promise<void> {
      try {
        const params: V1SchedulesScheduleIdGetRequest = {
          scheduleId,
        }
        const res = await this.scheduleApi().v1SchedulesScheduleIdGet(params)

        const coordinatorStore = useCoordinatorStore()
        this.schedule = res.schedule
        coordinatorStore.coordinators.push(res.coordinator)
      }
      catch (err) {
        return this.errorHandler(err, {
          404: '対象の開催スケジュールが見つかりません。',
        })
      }
    },

    async analyzeSchedule(scheduleId: string): Promise<void> {
      try {
        const params: V1SchedulesScheduleIdAnalyticsGetRequest = {
          scheduleId,
        }
        const res = await this.scheduleApi().v1SchedulesScheduleIdAnalyticsGet(params)

        this.viewerLogs = res.viewerLogs
      }
      catch (err) {
        return this.errorHandler(err, {
          404: '対象の開催スケジュールが見つかりません。',
        })
      }
    },

    /**
     * マルシェ開催スケジュールを登録する非同期関数
     * @param payload
     */
    async createSchedule(payload: CreateScheduleRequest): Promise<Schedule> {
      try {
        const params: V1SchedulesPostRequest = {
          createScheduleRequest: payload,
        }
        const res = await this.scheduleApi().v1SchedulesPost(params)
        return res.schedule
      }
      catch (err) {
        return this.errorHandler(err, {
          400: '必須項目が不足しているか、入力内容に誤りがあります。',
        })
      }
    },

    /**
     * マルシェ開催スケジュールを更新する非同期関数
     * @param scheduleId スケジュールID
     * @param payload
     */
    async updateSchedule(
      scheduleId: string,
      payload: UpdateScheduleRequest,
    ): Promise<void> {
      try {
        const params: V1SchedulesScheduleIdPatchRequest = {
          scheduleId,
          updateScheduleRequest: payload,
        }
        await this.scheduleApi().v1SchedulesScheduleIdPatch(params)
      }
      catch (err) {
        return this.errorHandler(err, {
          400: '必須項目が不足しているか、内容に誤りがあります。',
          404: '対象の開催スケジュールが見つかりません。',
        })
      }
    },

    /**
     * マルシェ開催スケジュールを削除する非同期関数
     * @param scheduleId スケジュールID
     * @returns
     */
    async deleteSchedule(scheduleId: string): Promise<void> {
      try {
        const params: V1SchedulesScheduleIdDeleteRequest = {
          scheduleId,
        }
        await this.scheduleApi().v1SchedulesScheduleIdDelete(params)
      }
      catch (err) {
        return this.errorHandler(err, {
          404: '対象の開催スケジュールが見つかりません。',
          412: 'ライブ配信中のため削除できません。',
        })
      }
    },

    /**
     * マルシェ開催スケジュールの承認/却下をする非同期関数
     * @param scheduleId スケジュールID
     * @param approved 承認フラグ
     * @returns
     */
    async approveSchedule(
      scheduleId: string,
      approved: boolean,
    ): Promise<void> {
      try {
        const approveParams: V1SchedulesScheduleIdApprovalPatchRequest = {
          scheduleId,
          approveScheduleRequest: {
            approved,
          },
        }
        await this.scheduleApi().v1SchedulesScheduleIdApprovalPatch(approveParams)

        // データの更新
        const index = this.schedules.findIndex(
          schedule => schedule.id === scheduleId,
        )
        if (index === -1) {
          return
        }
        const getParams: V1SchedulesScheduleIdGetRequest = {
          scheduleId,
        }
        const res = await this.scheduleApi().v1SchedulesScheduleIdGet(getParams)
        this.schedules.splice(index, 1, res.schedule)
      }
      catch (err) {
        return this.errorHandler(err, {
          400: '必須項目が不足しているか、内容に誤りがあります。',
          404: '対象の開催スケジュールが見つかりません。',
        })
      }
    },

    /**
     * マルシェ開催スケジュールの公開/非公開をする非同期関数
     * @param scheduleId スケジュールID
     * @param public 公開フラグ
     * @returns
     */
    async publishSchedule(
      scheduleId: string,
      published: boolean,
    ): Promise<void> {
      try {
        const publishParams: V1SchedulesScheduleIdPublishPatchRequest = {
          scheduleId,
          publishScheduleRequest: {
            _public: published,
          },
        }
        await this.scheduleApi().v1SchedulesScheduleIdPublishPatch(publishParams)

        // データの更新
        const index = this.schedules.findIndex(
          schedule => schedule.id === scheduleId,
        )
        if (index === -1) {
          return
        }
        const getParams: V1SchedulesScheduleIdGetRequest = {
          scheduleId,
        }
        const res = await this.scheduleApi().v1SchedulesScheduleIdGet(getParams)
        this.schedules.splice(index, 1, res.schedule)
      }
      catch (err) {
        return this.errorHandler(err, {
          400: '必須項目が不足しているか、内容に誤りがあります。',
          404: '対象の開催スケジュールが見つかりません。',
        })
      }
    },

    /**
     * サムネイル画像をアップロードする非同期関数
     * @param payload
     * @returns アップロード先URL
     */
    async uploadScheduleThumbnail(payload: File): Promise<string> {
      const contentType = payload.type
      try {
        const params: V1UploadSchedulesThumbnailPostRequest = {
          getUploadURLRequest: {
            fileType: contentType,
          },
        }
        const res = await this.uploadApi().v1UploadSchedulesThumbnailPost(params)

        return await fileUpload(this.uploadApi(), payload, res.key, res.url)
      }
      catch (err) {
        return this.errorHandler(err, {
          400: 'このファイルはアップロードできません。',
        })
      }
    },

    /**
     * 蓋絵画像をアップロードする非同期関数
     * @param payload
     * @returns アップロード先URL
     */
    async uploadScheduleImage(payload: File): Promise<string> {
      const contentType = payload.type
      try {
        const params: V1UploadSchedulesImagePostRequest = {
          getUploadURLRequest: {
            fileType: contentType,
          },
        }
        const res = await this.uploadApi().v1UploadSchedulesImagePost(params)

        return await fileUpload(this.uploadApi(), payload, res.key, res.url)
      }
      catch (err) {
        return this.errorHandler(err, {
          400: 'このファイルはアップロードできません。',
        })
      }
    },

    /**
     * オープニング動画をアップロードする非同期関数
     * @param payload
     * @returns アップロード先URL
     */
    async uploadScheduleOpeningVideo(payload: File): Promise<string> {
      const contentType = payload.type
      try {
        const params: V1UploadSchedulesOpeningVideoPostRequest = {
          getUploadURLRequest: {
            fileType: contentType,
          },
        }
        const res = await this.uploadApi().v1UploadSchedulesOpeningVideoPost(params)

        return await fileUpload(this.uploadApi(), payload, res.key, res.url)
      }
      catch (err) {
        return this.errorHandler(err, {
          400: 'このファイルはアップロードできません。',
        })
      }
    },
  },
})
