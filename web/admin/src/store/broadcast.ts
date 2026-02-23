import { useApiClient } from '~/composables/useApiClient'
import { fileUpload } from './helper'
import { BroadcastApi, GuestApi, UploadApi } from '~/types/api/v1'
import type { AuthYoutubeBroadcastRequest, Broadcast, CallbackAuthYoutubeBroadcastRequest, CreateYoutubeBroadcastRequest, GuestBroadcast, V1GuestsSchedulesBroadcastsYoutubeAuthCompletePostRequest, V1GuestsSchedulesBroadcastsYoutubePostRequest, V1SchedulesScheduleIdBroadcastsArchiveVideoPostRequest, V1SchedulesScheduleIdBroadcastsDeleteRequest, V1SchedulesScheduleIdBroadcastsGetRequest, V1SchedulesScheduleIdBroadcastsMp4PostRequest, V1SchedulesScheduleIdBroadcastsPostRequest, V1SchedulesScheduleIdBroadcastsRtmpPostRequest, V1SchedulesScheduleIdBroadcastsStaticImageDeleteRequest, V1SchedulesScheduleIdBroadcastsStaticImagePostRequest, V1SchedulesScheduleIdBroadcastsYoutubeAuthPostRequest, V1UploadSchedulesBroadcastsLivePostRequest, V1UploadSchedulesScheduleIdBroadcastsArchivePostRequest } from '~/types/api/v1'

export const useBroadcastStore = defineStore('broadcast', () => {
  const { create, errorHandler } = useApiClient()
  const broadcastApi = () => create(BroadcastApi)
  const guestApi = () => create(GuestApi)
  const uploadApi = () => create(UploadApi)

  const broadcast = ref<Broadcast>({} as Broadcast)
  const guestBroadcast = ref<GuestBroadcast>({} as GuestBroadcast)

  async function getBroadcastByScheduleId(scheduleId: string): Promise<void> {
    try {
      const params: V1SchedulesScheduleIdBroadcastsGetRequest = { scheduleId }
      const res = await broadcastApi().v1SchedulesScheduleIdBroadcastsGet(params)
      broadcast.value = res.broadcast
    }
    catch (err) {
      return errorHandler(err)
    }
  }

  async function pause(scheduleId: string): Promise<void> {
    try {
      const params: V1SchedulesScheduleIdBroadcastsDeleteRequest = { scheduleId }
      await broadcastApi().v1SchedulesScheduleIdBroadcastsDelete(params)
    }
    catch (err) {
      return errorHandler(err)
    }
  }

  async function unpause(scheduleId: string): Promise<void> {
    try {
      const params: V1SchedulesScheduleIdBroadcastsPostRequest = { scheduleId }
      await broadcastApi().v1SchedulesScheduleIdBroadcastsPost(params)
    }
    catch (err) {
      return errorHandler(err)
    }
  }

  async function activateStaticImage(scheduleId: string): Promise<void> {
    try {
      const params: V1SchedulesScheduleIdBroadcastsStaticImagePostRequest = { scheduleId }
      await broadcastApi().v1SchedulesScheduleIdBroadcastsStaticImagePost(params)
    }
    catch (err) {
      return errorHandler(err, {
        404: '指定したマルシェの配信が見つかりません。',
        412: 'マルシェの配信中ではないため、ふた絵を有効化できません。',
      })
    }
  }

  async function deactivateStaticImage(scheduleId: string): Promise<void> {
    try {
      const params: V1SchedulesScheduleIdBroadcastsStaticImageDeleteRequest = { scheduleId }
      await broadcastApi().v1SchedulesScheduleIdBroadcastsStaticImageDelete(params)
    }
    catch (err) {
      return errorHandler(err, {
        404: '指定したマルシェの配信が見つかりません。',
        412: 'マルシェの配信中ではないため、ふた絵を無効化できません。',
      })
    }
  }

  async function activateMp4Input(scheduleId: string, payload: File): Promise<void> {
    try {
      const uploadParams: V1UploadSchedulesBroadcastsLivePostRequest = {
        getUploadURLRequest: { fileType: payload.type },
      }
      const res = await uploadApi().v1UploadSchedulesBroadcastsLivePost(uploadParams)
      const inputUrl = await fileUpload(uploadApi(), payload, res.key, res.url)

      const activateParams: V1SchedulesScheduleIdBroadcastsMp4PostRequest = {
        scheduleId,
        activateBroadcastMP4Request: { inputUrl },
      }
      await broadcastApi().v1SchedulesScheduleIdBroadcastsMp4Post(activateParams)
    }
    catch (err) {
      return errorHandler(err, {
        404: '指定したマルシェの配信が見つかりません。',
        412: 'マルシェの配信中ではないため、MP4に切り替えできません。',
      })
    }
  }

  async function activateRtmpInput(scheduleId: string): Promise<void> {
    try {
      const params: V1SchedulesScheduleIdBroadcastsRtmpPostRequest = { scheduleId }
      await broadcastApi().v1SchedulesScheduleIdBroadcastsRtmpPost(params)
    }
    catch (err) {
      return errorHandler(err, {
        404: '指定したマルシェの配信が見つかりません。',
        412: 'マルシェの配信中ではないため、RTMPに切り替えできません。',
      })
    }
  }

  async function uploadArchiveMp4(scheduleId: string, payload: File): Promise<void> {
    try {
      const uploadParams: V1UploadSchedulesScheduleIdBroadcastsArchivePostRequest = {
        scheduleId,
        getUploadURLRequest: { fileType: payload.type },
      }
      const res = await uploadApi().v1UploadSchedulesScheduleIdBroadcastsArchivePost(uploadParams)
      const archiveUrl = await fileUpload(uploadApi(), payload, res.key, res.url)

      const archiveParams: V1SchedulesScheduleIdBroadcastsArchiveVideoPostRequest = {
        scheduleId,
        updateBroadcastArchiveRequest: { archiveUrl },
      }
      await broadcastApi().v1SchedulesScheduleIdBroadcastsArchiveVideoPost(archiveParams)
    }
    catch (err) {
      return errorHandler(err, {
        404: '指定したマルシェの配信が見つかりません。',
        412: 'マルシェの配信が終了していないため、オンデマンド動画を差し替えできません。',
      })
    }
  }

  async function getGuestBroadcast(): Promise<void> {
    try {
      const res = await broadcastApi().v1GuestsSchedulesBroadcastsGet()
      guestBroadcast.value = res.broadcast
    }
    catch (err) {
      return errorHandler(err, {
        401: '認証情報に誤りがあります。再度認証を行ってください。',
        403: '認証情報に誤りがあります。再度認証を行ってください。',
        404: '指定したマルシェの配信が見つかりません。',
      })
    }
  }

  async function authYouTube(scheduleId: string, payload: AuthYoutubeBroadcastRequest): Promise<string> {
    try {
      const params: V1SchedulesScheduleIdBroadcastsYoutubeAuthPostRequest = {
        scheduleId,
        authYoutubeBroadcastRequest: payload,
      }
      const res = await broadcastApi().v1SchedulesScheduleIdBroadcastsYoutubeAuthPost(params)
      return res.url
    }
    catch (err) {
      return errorHandler(err, {
        404: '指定したマルシェの配信が見つかりません。',
        412: 'マルシェの配信が開始しているため、YouTube認証を行えません。',
      })
    }
  }

  async function connectYouTube(payload: CallbackAuthYoutubeBroadcastRequest): Promise<void> {
    try {
      const params: V1GuestsSchedulesBroadcastsYoutubeAuthCompletePostRequest = {
        callbackAuthYoutubeBroadcastRequest: payload,
      }
      const res = await guestApi().v1GuestsSchedulesBroadcastsYoutubeAuthCompletePost(params)
      guestBroadcast.value = res.broadcast
    }
    catch (err) {
      return errorHandler(err, {
        401: 'YouTubeの認証情報に誤りがあります。再度認証を行ってください。',
        403: 'YouTubeの認証情報に誤りがあります。再度認証を行ってください。',
        404: '指定したマルシェの配信が見つかりません。',
        412: 'マルシェの配信が開始しているため、YouTube認証を行えません。',
      })
    }
  }

  async function createYoutubeLive(payload: CreateYoutubeBroadcastRequest): Promise<void> {
    try {
      const params: V1GuestsSchedulesBroadcastsYoutubePostRequest = {
        createYoutubeBroadcastRequest: payload,
      }
      await guestApi().v1GuestsSchedulesBroadcastsYoutubePost(params)
    }
    catch (err) {
      return errorHandler(err, {
        401: 'YouTubeの認証情報に誤りがあります。再度認証を行ってください。',
        403: 'YouTubeの認証情報に誤りがあります。再度認証を行ってください。',
        404: '指定したマルシェの配信が見つかりません。',
        412: 'マルシェの配信が開始しているため、YouTube認証を行えません。',
      })
    }
  }

  return {
    broadcast,
    guestBroadcast,
    getBroadcastByScheduleId,
    pause,
    unpause,
    activateStaticImage,
    deactivateStaticImage,
    activateMp4Input,
    activateRtmpInput,
    uploadArchiveMp4,
    getGuestBroadcast,
    authYouTube,
    connectYouTube,
    createYoutubeLive,
  }
})
