import { defineStore } from 'pinia'

import { fileUpload } from './helper'
import { useCoordinatorStore } from './coordinator'
import { apiClient } from '~/plugins/api-client'
import type { ApproveScheduleRequest, BroadcastViewerLog, CreateScheduleRequest, GetUploadUrlRequest, PublishScheduleRequest, Schedule, UpdateScheduleRequest } from '~/types/api'

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
        const res = await apiClient.scheduleApi().v1ListSchedules(limit, offset)

        const coordinatorStore = useCoordinatorStore()
        this.schedules = res.data.schedules
        this.total = res.data.total
        coordinatorStore.coordinators = res.data.coordinators
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
        const res = await apiClient.scheduleApi().v1GetSchedule(scheduleId)

        const coordinatorStore = useCoordinatorStore()
        this.schedule = res.data.schedule
        coordinatorStore.coordinators.push(res.data.coordinator)
      }
      catch (err) {
        return this.errorHandler(err, { 404: '対象の開催スケジュールが見つかりません。' })
      }
    },

    async analyzeSchedule(scheduleId: string): Promise<void> {
      try {
        const res = await apiClient.scheduleApi().v1AnalyzeSchedule(scheduleId)

        this.viewerLogs = res.data.viewerLogs
      }
      catch (err) {
        return this.errorHandler(err, { 404: '対象の開催スケジュールが見つかりません。' })
      }
    },

    /**
     * マルシェ開催スケジュールを登録する非同期関数
     * @param payload
     */
    async createSchedule(payload: CreateScheduleRequest): Promise<Schedule> {
      try {
        const res = await apiClient.scheduleApi().v1CreateSchedule(payload)
        return res.data.schedule
      }
      catch (err) {
        return this.errorHandler(err, { 400: '必須項目が不足しているか、入力内容に誤りがあります。' })
      }
    },

    /**
     * マルシェ開催スケジュールを更新する非同期関数
     * @param scheduleId スケジュールID
     * @param payload
     */
    async updateSchedule(scheduleId: string, payload: UpdateScheduleRequest): Promise<void> {
      try {
        await apiClient.scheduleApi().v1UpdateSchedule(scheduleId, payload)
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
        await apiClient.scheduleApi().v1DeleteSchedule(scheduleId)
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
    async approveSchedule(scheduleId: string, approved: boolean): Promise<void> {
      try {
        const req: ApproveScheduleRequest = { approved }
        await apiClient.scheduleApi().v1ApproveSchedule(scheduleId, req)

        // データの更新
        const index = this.schedules.findIndex(schedule => schedule.id === scheduleId)
        if (index === -1) {
          return
        }
        const res = await apiClient.scheduleApi().v1GetSchedule(scheduleId)
        this.schedules.splice(index, 1, res.data.schedule)
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
    async publishSchedule(scheduleId: string, published: boolean): Promise<void> {
      try {
        const req: PublishScheduleRequest = { public: published }
        await apiClient.scheduleApi().v1PublishSchedule(scheduleId, req)

        // データの更新
        const index = this.schedules.findIndex(schedule => schedule.id === scheduleId)
        if (index === -1) {
          return
        }
        const res = await apiClient.scheduleApi().v1GetSchedule(scheduleId)
        this.schedules.splice(index, 1, res.data.schedule)
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
        const body: GetUploadUrlRequest = {
          fileType: contentType,
        }
        const res = await apiClient.scheduleApi().v1GetScheduleThumbnailUploadUrl(body)

        return await fileUpload(payload, res.data.key, res.data.url)
      }
      catch (err) {
        return this.errorHandler(err, { 400: 'このファイルはアップロードできません。' })
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
        const body: GetUploadUrlRequest = {
          fileType: contentType,
        }
        const res = await apiClient.scheduleApi().v1GetScheduleImageUploadUrl(body)

        return await fileUpload(payload, res.data.key, res.data.url)
      }
      catch (err) {
        return this.errorHandler(err, { 400: 'このファイルはアップロードできません。' })
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
        const body: GetUploadUrlRequest = {
          fileType: contentType,
        }
        const res = await apiClient.scheduleApi().v1GetScheduleOpeningVideoUploadUrl(body)

        return await fileUpload(payload, res.data.key, res.data.url)
      }
      catch (err) {
        return this.errorHandler(err, { 400: 'このファイルはアップロードできません。' })
      }
    },
  },
})
