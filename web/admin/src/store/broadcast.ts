import { fileUpload } from './helper'
import type { AuthYoutubeBroadcastRequest, Broadcast, CallbackAuthYoutubeBroadcastRequest, CreateYoutubeBroadcastRequest, GuestBroadcast, V1GuestsSchedulesBroadcastsYoutubeAuthCompletePostRequest, V1GuestsSchedulesBroadcastsYoutubePostRequest, V1SchedulesScheduleIdBroadcastsArchiveVideoPostRequest, V1SchedulesScheduleIdBroadcastsDeleteRequest, V1SchedulesScheduleIdBroadcastsGetRequest, V1SchedulesScheduleIdBroadcastsMp4PostRequest, V1SchedulesScheduleIdBroadcastsPostRequest, V1SchedulesScheduleIdBroadcastsRtmpPostRequest, V1SchedulesScheduleIdBroadcastsStaticImageDeleteRequest, V1SchedulesScheduleIdBroadcastsStaticImagePostRequest, V1SchedulesScheduleIdBroadcastsYoutubeAuthPostRequest, V1UploadSchedulesBroadcastsLivePostRequest, V1UploadSchedulesScheduleIdBroadcastsArchivePostRequest } from '~/types/api/v1'

export const useBroadcastStore = defineStore('broadcast', {
  state: () => ({
    broadcast: {} as Broadcast,
    guestBroadcast: {} as GuestBroadcast,
  }),

  actions: {
    /**
     * ライブ配信情報を取得する非同期関数
     * @param scheduleId マルシェ開催スケジュールID
     * @returns
     */
    async getBroadcastByScheduleId(scheduleId: string): Promise<void> {
      try {
        const params: V1SchedulesScheduleIdBroadcastsGetRequest = {
          scheduleId,
        }
        const res = await this.broadcastApi().v1SchedulesScheduleIdBroadcastsGet(params)
        this.broadcast = res.broadcast
      }
      catch (err) {
        return this.errorHandler(err)
      }
    },

    /**
     * ライブ配信を一時停止する非同期関数
     * @param scheduleId マルシェ開催スケジュールID
     * @returns
     */
    async pause(scheduleId: string): Promise<void> {
      try {
        const params: V1SchedulesScheduleIdBroadcastsDeleteRequest = {
          scheduleId,
        }
        await this.broadcastApi().v1SchedulesScheduleIdBroadcastsDelete(params)
      }
      catch (err) {
        return this.errorHandler(err)
      }
    },

    /**
     * ライブ配信の一時停止を解除する非同期関数
     * @param scheduleId マルシェ開催スケジュールID
     * @returns
     */
    async unpause(scheduleId: string): Promise<void> {
      try {
        const params: V1SchedulesScheduleIdBroadcastsPostRequest = {
          scheduleId,
        }
        await this.broadcastApi().v1SchedulesScheduleIdBroadcastsPost(params)
      }
      catch (err) {
        return this.errorHandler(err)
      }
    },

    /**
     * ライブ配信のふた絵を有効化する非同期関数
     * @param scheduleId マルシェ開催スケジュールID
     * @returns
     */
    async activateStaticImage(scheduleId: string): Promise<void> {
      try {
        const params: V1SchedulesScheduleIdBroadcastsStaticImagePostRequest = {
          scheduleId,
        }
        await this.broadcastApi().v1SchedulesScheduleIdBroadcastsStaticImagePost(params)
      }
      catch (err) {
        return this.errorHandler(err, {
          404: '指定したマルシェの配信が見つかりません。',
          412: 'マルシェの配信中ではないため、ふた絵を有効化できません。',
        })
      }
    },

    /**
     * ライブ配信のふた絵を無効化する非同期関数
     * @param scheduleId マルシェ開催スケジュールID
     * @returns
     */
    async deactivateStaticImage(scheduleId: string): Promise<void> {
      try {
        const params: V1SchedulesScheduleIdBroadcastsStaticImageDeleteRequest = {
          scheduleId,
        }
        await this.broadcastApi().v1SchedulesScheduleIdBroadcastsStaticImageDelete(params)
      }
      catch (err) {
        return this.errorHandler(err, {
          404: '指定したマルシェの配信が見つかりません。',
          412: 'マルシェの配信中ではないため、ふた絵を無効化できません。',
        })
      }
    },

    /**
     * ライブ配信の入力チャンネルをMP4に切り替え
     * @param scheduleId マルシェ開催スケジュールID
     * @param payload ライブ動画
     * @returns
     */
    async activateMp4Input(scheduleId: string, payload: File): Promise<void> {
      try {
        const uploadParams: V1UploadSchedulesBroadcastsLivePostRequest = {
          getUploadURLRequest: {
            fileType: payload.type,
          },
        }
        const res = await this.uploadApi().v1UploadSchedulesBroadcastsLivePost(uploadParams)

        const inputUrl = await fileUpload(this.uploadApi(), payload, res.key, res.url)

        const activateParams: V1SchedulesScheduleIdBroadcastsMp4PostRequest = {
          scheduleId,
          activateBroadcastMP4Request: {
            inputUrl,
          },
        }
        await this.broadcastApi().v1SchedulesScheduleIdBroadcastsMp4Post(activateParams)
      }
      catch (err) {
        return this.errorHandler(err, {
          404: '指定したマルシェの配信が見つかりません。',
          412: 'マルシェの配信中ではないため、MP4に切り替えできません。',
        })
      }
    },

    /**
     * ライブ配信の入力チャンネルをRTMPに切り替え
     * @param scheduleId マルシェ開催スケジュールID
     * @returns
     */
    async activateRtmpInput(scheduleId: string): Promise<void> {
      try {
        const params: V1SchedulesScheduleIdBroadcastsRtmpPostRequest = {
          scheduleId,
        }
        await this.broadcastApi().v1SchedulesScheduleIdBroadcastsRtmpPost(params)
      }
      catch (err) {
        return this.errorHandler(err, {
          404: '指定したマルシェの配信が見つかりません。',
          412: 'マルシェの配信中ではないため、RTMPに切り替えできません。',
        })
      }
    },

    /**
     * オンデマンド配信の動画を差し替え
     * @param scheduleId マルシェ開催スケジュールID
     * @param payload オンデマンド動画
     * @returns
     */
    async uploadArchiveMp4(scheduleId: string, payload: File): Promise<void> {
      try {
        const uploadParams: V1UploadSchedulesScheduleIdBroadcastsArchivePostRequest = {
          scheduleId,
          getUploadURLRequest: {
            fileType: payload.type,
          },
        }
        const res = await this.uploadApi().v1UploadSchedulesScheduleIdBroadcastsArchivePost(uploadParams)

        const archiveUrl = await fileUpload(this.uploadApi(), payload, res.key, res.url)

        const archiveParams: V1SchedulesScheduleIdBroadcastsArchiveVideoPostRequest = {
          scheduleId,
          updateBroadcastArchiveRequest: {
            archiveUrl,
          },
        }
        await this.broadcastApi().v1SchedulesScheduleIdBroadcastsArchiveVideoPost(archiveParams)
      }
      catch (err) {
        return this.errorHandler(err, {
          404: '指定したマルシェの配信が見つかりません。',
          412: 'マルシェの配信が終了していないため、オンデマンド動画を差し替えできません。',
        })
      }
    },

    /**
     * ゲスト用のライブ配信情報を取得する非同期関数
     * @returns
     */
    async getGuestBroadcast(): Promise<void> {
      try {
        const res = await this.broadcastApi().v1GuestsSchedulesBroadcastsGet()
        this.guestBroadcast = res.broadcast
      }
      catch (err) {
        return this.errorHandler(err, {
          401: '認証情報に誤りがあります。再度認証を行ってください。',
          403: '認証情報に誤りがあります。再度認証を行ってください。',
          404: '指定したマルシェの配信が見つかりません。',
        })
      }
    },

    /**
     * YouTube認証を行う非同期関数
     * @param scheduleId マルシェ開催スケジュールID
     * @param payload YouTube認証情報
     * @returns
     */
    async authYouTube(scheduleId: string, payload: AuthYoutubeBroadcastRequest): Promise<string> {
      try {
        const params: V1SchedulesScheduleIdBroadcastsYoutubeAuthPostRequest = {
          scheduleId,
          authYoutubeBroadcastRequest: payload,
        }
        const res = await this.broadcastApi().v1SchedulesScheduleIdBroadcastsYoutubeAuthPost(params)
        return res.url
      }
      catch (err) {
        return this.errorHandler(err, {
          404: '指定したマルシェの配信が見つかりません。',
          412: 'マルシェの配信が開始しているため、YouTube認証を行えません。',
        })
      }
    },

    /**
     * YouTube連携を行う非同期関数
     * @param payload YouTube連携情報
     * @returns
     */
    async connectYouTube(payload: CallbackAuthYoutubeBroadcastRequest): Promise<void> {
      try {
        const params: V1GuestsSchedulesBroadcastsYoutubeAuthCompletePostRequest = {
          callbackAuthYoutubeBroadcastRequest: payload,
        }
        const res = await this.guestApi().v1GuestsSchedulesBroadcastsYoutubeAuthCompletePost(params)
        this.guestBroadcast = res.broadcast
      }
      catch (err) {
        return this.errorHandler(err, {
          401: 'YouTubeの認証情報に誤りがあります。再度認証を行ってください。',
          403: 'YouTubeの認証情報に誤りがあります。再度認証を行ってください。',
          404: '指定したマルシェの配信が見つかりません。',
          412: 'マルシェの配信が開始しているため、YouTube認証を行えません。',
        })
      }
    },

    /**
     * YouTube配信の登録を行う非同期関数
     * @param payload YouTube配信情報
     * @returns
     */
    async createYoutubeLive(payload: CreateYoutubeBroadcastRequest): Promise<void> {
      try {
        const params: V1GuestsSchedulesBroadcastsYoutubePostRequest = {
          createYoutubeBroadcastRequest: payload,
        }
        await this.guestApi().v1GuestsSchedulesBroadcastsYoutubePost(params)
      }
      catch (err) {
        return this.errorHandler(err, {
          401: 'YouTubeの認証情報に誤りがあります。再度認証を行ってください。',
          403: 'YouTubeの認証情報に誤りがあります。再度認証を行ってください。',
          404: '指定したマルシェの配信が見つかりません。',
          412: 'マルシェの配信が開始しているため、YouTube認証を行えません。',
        })
      }
    },
  },
})
