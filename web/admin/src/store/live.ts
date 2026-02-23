import { useApiClient } from '~/composables/useApiClient'
import { useProducerStore } from './producer'
import { useProductStore } from './product'
import { LiveApi } from '~/types/api/v1'
import type { CreateLiveRequest, Live, UpdateLiveRequest, V1SchedulesScheduleIdLivesGetRequest, V1SchedulesScheduleIdLivesLiveIdDeleteRequest, V1SchedulesScheduleIdLivesLiveIdPatchRequest, V1SchedulesScheduleIdLivesPostRequest } from '~/types/api/v1'

export const useLiveStore = defineStore('live', () => {
  const { create, errorHandler } = useApiClient()
  const liveApi = () => create(LiveApi)

  const lives = ref<Live[]>([])
  const total = ref<number>(0)

  async function fetchLives(scheduleId: string): Promise<void> {
    try {
      const params: V1SchedulesScheduleIdLivesGetRequest = { scheduleId }
      const res = await liveApi().v1SchedulesScheduleIdLivesGet(params)

      const producerStore = useProducerStore()
      const productStore = useProductStore()
      lives.value = res.lives
      total.value = res.total
      producerStore.producers = res.producers
      productStore.products = res.products
    }
    catch (err) {
      return errorHandler(err, { 404: 'マルシェタイムテーブルが存在しません' })
    }
  }

  async function fetchLiveComments(scheduleId: string): Promise<void> {}

  async function createLive(scheduleId: string, payload: CreateLiveRequest): Promise<void> {
    try {
      const params: V1SchedulesScheduleIdLivesPostRequest = {
        scheduleId,
        createLiveRequest: payload,
      }
      const res = await liveApi().v1SchedulesScheduleIdLivesPost(params)
      lives.value.push(res.live)
    }
    catch (err) {
      return errorHandler(err, {
        400: '必須項目が不足しているか、内容に誤りがあります',
        404: 'マルシェタイムテーブルが存在しません',
        412: '開催時間, 生産者が重複しています',
      })
    }
  }

  async function updateLive(scheduleId: string, liveId: string, payload: UpdateLiveRequest): Promise<void> {
    try {
      const params: V1SchedulesScheduleIdLivesLiveIdPatchRequest = {
        scheduleId,
        liveId,
        updateLiveRequest: payload,
      }
      await liveApi().v1SchedulesScheduleIdLivesLiveIdPatch(params)

      const index = lives.value.findIndex((l: Live): boolean => l.id === liveId)
      const live = lives.value[index]
      lives.value.splice(index, 1, { ...live, ...payload })
    }
    catch (err) {
      return errorHandler(err, {
        400: '必須項目が不足しているか、内容に誤りがあります',
        404: 'マルシェ開催スケジュール, マルシェタイムテーブルが存在しません',
        412: '開催時間, 生産者が重複しています',
      })
    }
  }

  async function deleteLive(scheduleId: string, liveId: string): Promise<void> {
    try {
      const params: V1SchedulesScheduleIdLivesLiveIdDeleteRequest = { scheduleId, liveId }
      await liveApi().v1SchedulesScheduleIdLivesLiveIdDelete(params)

      const index = lives.value.findIndex((l: Live): boolean => l.id === liveId)
      lives.value.splice(index, 1)
    }
    catch (err) {
      return errorHandler(err, {
        400: '必須項目が不足しているか、内容に誤りがあります',
        404: 'マルシェ開催スケジュール, マルシェタイムテーブルが存在しません',
      })
    }
  }

  return {
    lives,
    total,
    fetchLives,
    fetchLiveComments,
    createLive,
    updateLive,
    deleteLive,
  }
})
