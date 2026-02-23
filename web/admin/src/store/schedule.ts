import { fileUpload } from './helper'
import { useCoordinatorStore } from './coordinator'
import { useApiClient } from '~/composables/useApiClient'
import { ScheduleApi, UploadApi } from '~/types/api/v1'
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

export const useScheduleStore = defineStore('schedule', () => {
  const { create, errorHandler } = useApiClient()
  const scheduleApi = () => create(ScheduleApi)
  const uploadApi = () => create(UploadApi)

  const schedule = ref<Schedule>({} as Schedule)
  const schedules = ref<Schedule[]>([])
  const viewerLogs = ref<BroadcastViewerLog[]>([])
  const totalViewers = ref<number>(0)
  const total = ref<number>(0)

  /**
   * マルシェ開催スケジュール一覧を取得する非同期関数
   * @param limit
   * @param offset
   */
  async function fetchSchedules(limit = 20, offset = 0): Promise<void> {
    try {
      const params: V1SchedulesGetRequest = {
        limit,
        offset,
      }
      const res = await scheduleApi().v1SchedulesGet(params)

      const coordinatorStore = useCoordinatorStore()
      schedules.value = res.schedules
      total.value = res.total
      coordinatorStore.coordinators = res.coordinators
    }
    catch (err) {
      return errorHandler(err)
    }
  }

  /**
   * マルシェ開催スケジュールを取得する非同期関数
   * @param scheduleId スケジュールID
   * @returns
   */
  async function getSchedule(scheduleId: string): Promise<void> {
    try {
      const params: V1SchedulesScheduleIdGetRequest = {
        scheduleId,
      }
      const res = await scheduleApi().v1SchedulesScheduleIdGet(params)

      const coordinatorStore = useCoordinatorStore()
      schedule.value = res.schedule
      coordinatorStore.coordinators.push(res.coordinator)
    }
    catch (err) {
      return errorHandler(err, {
        404: '対象の開催スケジュールが見つかりません。',
      })
    }
  }

  async function analyzeSchedule(scheduleId: string): Promise<void> {
    try {
      const params: V1SchedulesScheduleIdAnalyticsGetRequest = {
        scheduleId,
      }
      const res = await scheduleApi().v1SchedulesScheduleIdAnalyticsGet(params)

      viewerLogs.value = res.viewerLogs
      totalViewers.value = res.totalViewers
    }
    catch (err) {
      return errorHandler(err, {
        404: '対象の開催スケジュールが見つかりません。',
      })
    }
  }

  /**
   * マルシェ開催スケジュールを登録する非同期関数
   * @param payload
   */
  async function createSchedule(payload: CreateScheduleRequest): Promise<Schedule> {
    try {
      const params: V1SchedulesPostRequest = {
        createScheduleRequest: payload,
      }
      const res = await scheduleApi().v1SchedulesPost(params)
      return res.schedule
    }
    catch (err) {
      return errorHandler(err, {
        400: '必須項目が不足しているか、入力内容に誤りがあります。',
      })
    }
  }

  /**
   * マルシェ開催スケジュールを更新する非同期関数
   * @param scheduleId スケジュールID
   * @param payload
   */
  async function updateSchedule(
    scheduleId: string,
    payload: UpdateScheduleRequest,
  ): Promise<void> {
    try {
      const params: V1SchedulesScheduleIdPatchRequest = {
        scheduleId,
        updateScheduleRequest: payload,
      }
      await scheduleApi().v1SchedulesScheduleIdPatch(params)
    }
    catch (err) {
      return errorHandler(err, {
        400: '必須項目が不足しているか、内容に誤りがあります。',
        404: '対象の開催スケジュールが見つかりません。',
      })
    }
  }

  /**
   * マルシェ開催スケジュールを削除する非同期関数
   * @param scheduleId スケジュールID
   * @returns
   */
  async function deleteSchedule(scheduleId: string): Promise<void> {
    try {
      const params: V1SchedulesScheduleIdDeleteRequest = {
        scheduleId,
      }
      await scheduleApi().v1SchedulesScheduleIdDelete(params)
    }
    catch (err) {
      return errorHandler(err, {
        404: '対象の開催スケジュールが見つかりません。',
        412: 'ライブ配信中のため削除できません。',
      })
    }
  }

  /**
   * マルシェ開催スケジュールの承認/却下をする非同期関数
   * @param scheduleId スケジュールID
   * @param approved 承認フラグ
   * @returns
   */
  async function approveSchedule(
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
      await scheduleApi().v1SchedulesScheduleIdApprovalPatch(approveParams)

      // データの更新
      const index = schedules.value.findIndex(
        s => s.id === scheduleId,
      )
      if (index === -1) {
        return
      }
      const getParams: V1SchedulesScheduleIdGetRequest = {
        scheduleId,
      }
      const res = await scheduleApi().v1SchedulesScheduleIdGet(getParams)
      schedules.value.splice(index, 1, res.schedule)
    }
    catch (err) {
      return errorHandler(err, {
        400: '必須項目が不足しているか、内容に誤りがあります。',
        404: '対象の開催スケジュールが見つかりません。',
      })
    }
  }

  /**
   * マルシェ開催スケジュールの公開/非公開をする非同期関数
   * @param scheduleId スケジュールID
   * @param public 公開フラグ
   * @returns
   */
  async function publishSchedule(
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
      await scheduleApi().v1SchedulesScheduleIdPublishPatch(publishParams)

      // データの更新
      const index = schedules.value.findIndex(
        s => s.id === scheduleId,
      )
      if (index === -1) {
        return
      }
      const getParams: V1SchedulesScheduleIdGetRequest = {
        scheduleId,
      }
      const res = await scheduleApi().v1SchedulesScheduleIdGet(getParams)
      schedules.value.splice(index, 1, res.schedule)
    }
    catch (err) {
      return errorHandler(err, {
        400: '必須項目が不足しているか、内容に誤りがあります。',
        404: '対象の開催スケジュールが見つかりません。',
      })
    }
  }

  /**
   * サムネイル画像をアップロードする非同期関数
   * @param payload
   * @returns アップロード先URL
   */
  async function uploadScheduleThumbnail(payload: File): Promise<string> {
    const contentType = payload.type
    try {
      const params: V1UploadSchedulesThumbnailPostRequest = {
        getUploadURLRequest: {
          fileType: contentType,
        },
      }
      const res = await uploadApi().v1UploadSchedulesThumbnailPost(params)

      return await fileUpload(uploadApi(), payload, res.key, res.url)
    }
    catch (err) {
      return errorHandler(err, {
        400: 'このファイルはアップロードできません。',
      })
    }
  }

  /**
   * 蓋絵画像をアップロードする非同期関数
   * @param payload
   * @returns アップロード先URL
   */
  async function uploadScheduleImage(payload: File): Promise<string> {
    const contentType = payload.type
    try {
      const params: V1UploadSchedulesImagePostRequest = {
        getUploadURLRequest: {
          fileType: contentType,
        },
      }
      const res = await uploadApi().v1UploadSchedulesImagePost(params)

      return await fileUpload(uploadApi(), payload, res.key, res.url)
    }
    catch (err) {
      return errorHandler(err, {
        400: 'このファイルはアップロードできません。',
      })
    }
  }

  /**
   * オープニング動画をアップロードする非同期関数
   * @param payload
   * @returns アップロード先URL
   */
  async function uploadScheduleOpeningVideo(payload: File): Promise<string> {
    const contentType = payload.type
    try {
      const params: V1UploadSchedulesOpeningVideoPostRequest = {
        getUploadURLRequest: {
          fileType: contentType,
        },
      }
      const res = await uploadApi().v1UploadSchedulesOpeningVideoPost(params)

      return await fileUpload(uploadApi(), payload, res.key, res.url)
    }
    catch (err) {
      return errorHandler(err, {
        400: 'このファイルはアップロードできません。',
      })
    }
  }

  return {
    schedule,
    schedules,
    viewerLogs,
    totalViewers,
    total,
    fetchSchedules,
    getSchedule,
    analyzeSchedule,
    createSchedule,
    updateSchedule,
    deleteSchedule,
    approveSchedule,
    publishSchedule,
    uploadScheduleThumbnail,
    uploadScheduleImage,
    uploadScheduleOpeningVideo,
  }
})
