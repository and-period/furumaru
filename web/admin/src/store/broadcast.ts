import { defineStore } from 'pinia'

import { fileUpload } from './helper'
import { apiClient } from '~/plugins/api-client'
import type { ActivateBroadcastMP4Request, Broadcast, GetUploadUrlRequest, UpdateBroadcastArchiveRequest } from '~/types/api'

export const useBroadcastStore = defineStore('broadcast', {
  state: () => ({
    broadcast: {} as Broadcast
  }),

  actions: {
    /**
     * ライブ配信情報を取得する非同期関数
     * @param scheduleId マルシェ開催スケジュールID
     * @returns
     */
    async getBroadcastByScheduleId (scheduleId: string): Promise<void> {
      try {
        const res = await apiClient.broadcastApi().v1GetBroadcast(scheduleId)
        this.broadcast = res.data.broadcast
      } catch (err) {
        return this.errorHandler(err)
      }
    },

    /**
     * ライブ配信を一時停止する非同期関数
     * @param scheduleId マルシェ開催スケジュールID
     * @returns
     */
    async pause (scheduleId: string): Promise<void> {
      try {
        await apiClient.broadcastApi().v1PauseBroadcast(scheduleId)
      } catch (err) {
        return this.errorHandler(err)
      }
    },

    /**
     * ライブ配信の一時停止を解除する非同期関数
     * @param scheduleId マルシェ開催スケジュールID
     * @returns
     */
    async unpause (scheduleId: string): Promise<void> {
      try {
        await apiClient.broadcastApi().v1UnpauseBroadcast(scheduleId)
      } catch (err) {
        return this.errorHandler(err)
      }
    },

    /**
     * ライブ配信のふた絵を有効化する非同期関数
     * @param scheduleId マルシェ開催スケジュールID
     * @returns
     */
    async activateStaticImage (scheduleId: string): Promise<void> {
      try {
        await apiClient.broadcastApi().v1ActivateBroadcastStaticImage(scheduleId)
      } catch (err) {
        return this.errorHandler(err)
      }
    },

    /**
     * ライブ配信のふた絵を無効化する非同期関数
     * @param scheduleId マルシェ開催スケジュールID
     * @returns
     */
    async deactivateStaticImage (scheduleId: string): Promise<void> {
      try {
        await apiClient.broadcastApi().v1DeactivateBroadcastStaticImage(scheduleId)
      } catch (err) {
        return this.errorHandler(err)
      }
    },

    /**
     * ライブ配信の入力チャンネルをMP4に切り替え
     * @param scheduleId マルシェ開催スケジュールID
     * @param payload ライブ動画
     * @returns
     */
    async activateMp4Input (scheduleId: string, payload: File): Promise<void> {
      try {
        const body: GetUploadUrlRequest = {
          fileType: payload.type
        }
        const res = await apiClient.broadcastApi().v1GetBroadcastLiveUploadUrl(body)

        const inputUrl = await fileUpload(payload, res.data.url)

        const req: ActivateBroadcastMP4Request = {
          inputUrl
        }
        await apiClient.broadcastApi().v1ActivateBroadcastMP4(scheduleId, req)
      } catch (err) {
        return this.errorHandler(err)
      }
    },

    /**
     * ライブ配信の入力チャンネルをRTMPに切り替え
     * @param scheduleId マルシェ開催スケジュールID
     * @returns
     */
    async activateRtmpInput (scheduleId: string): Promise<void> {
      try {
        await apiClient.broadcastApi().v1ActivateBroadcastRTMP(scheduleId)
      } catch (err) {
        return this.errorHandler(err)
      }
    },

    /**
     * オンデマンド配信の動画を差し替え
     * @param scheduleId マルシェ開催スケジュールID
     * @param payload オンデマンド動画
     * @returns
     */
    async uploadArchiveMp4 (scheduleId: string, payload: File): Promise<void> {
      try {
        const body: GetUploadUrlRequest = {
          fileType: payload.type
        }
        const res = await apiClient.broadcastApi().v1GetBroadcastArchiveUploadUrl(scheduleId, body)

        const archiveUrl = await fileUpload(payload, res.data.url)

        const req: UpdateBroadcastArchiveRequest = {
          archiveUrl
        }
        await apiClient.broadcastApi().v1UpdateBroadcastArchive(scheduleId, req)
      } catch (err) {
        return this.errorHandler(err)
      }
    }
  }
})
