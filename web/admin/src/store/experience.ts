import { useApiClient } from '~/composables/useApiClient'
import { useExperienceTypeStore } from './experience-type'
import { fileUpload } from './helper'
import { useProducerStore } from './producer'
import { ExperienceApi, UploadApi } from '~/types/api/v1'
import type {
  CreateExperienceRequest,
  Experience,
  GetUploadURLRequest,
  UpdateExperienceRequest,
  UploadURLResponse,
  V1ExperiencesExperienceIdDeleteRequest,
  V1ExperiencesExperienceIdGetRequest,
  V1ExperiencesExperienceIdPatchRequest,
  V1ExperiencesGetRequest,
  V1ExperiencesPostRequest,
  V1UploadProductsImagePostRequest,
  V1UploadProductsVideoPostRequest,
} from '~/types/api/v1'

export const useExperienceStore = defineStore('experience', () => {
  const { create, errorHandler } = useApiClient()
  const experienceApi = () => create(ExperienceApi)
  const uploadApi = () => create(UploadApi)

  const experience = ref<Experience>({} as Experience)
  const experiences = ref<Experience[]>([])
  const totalItems = ref<number>(0)

  async function searchExperiences(name: string = '', producerId: string = '', experienceId: string[] = []) {
    try {
      const params: V1ExperiencesGetRequest = { name, producerId }
      const res = await experienceApi().v1ExperiencesGet(params)
      const merged: Experience[] = []
      experiences.value.forEach((e) => {
        if (experienceId.includes(e.id)) {
          merged.push(e)
        }
      })
      res.experiences.forEach((e) => {
        if (merged.find(v => v.id === e.id)) {
          return
        }
        merged.push(e)
      })
      experiences.value = merged
      totalItems.value = res.total
    }
    catch (err) {
      return errorHandler(err)
    }
  }

  async function fetchExperiences(limit = 20, offset = 0): Promise<void> {
    try {
      const params: V1ExperiencesGetRequest = { limit, offset }
      const res = await experienceApi().v1ExperiencesGet(params)

      totalItems.value = res.total
      experiences.value = res.experiences

      const producerStore = useProducerStore()
      producerStore.producers = res.producers

      const experienceTypeStore = useExperienceTypeStore()
      experienceTypeStore.experienceTypes = res.experienceTypes
    }
    catch (err) {
      return errorHandler(err)
    }
  }

  async function fetchExperience(experienceId: string) {
    try {
      const params: V1ExperiencesExperienceIdGetRequest = { experienceId }
      const res = await experienceApi().v1ExperiencesExperienceIdGet(params)
      return res
    }
    catch (err) {
      return errorHandler(err)
    }
  }

  async function createExperience(payload: CreateExperienceRequest) {
    try {
      const params: V1ExperiencesPostRequest = { createExperienceRequest: payload }
      const res = await experienceApi().v1ExperiencesPost(params)
      return res
    }
    catch (err) {
      return errorHandler(err, {
        400: '必須項目が不足しているか、入力内容に誤りがあります。',
        403: '登録する権限がありません。管理者にお問い合わせください',
      })
    }
  }

  async function updateExperience(experienceId: string, payload: UpdateExperienceRequest) {
    try {
      const params: V1ExperiencesExperienceIdPatchRequest = {
        experienceId,
        updateExperienceRequest: payload,
      }
      const res = await experienceApi().v1ExperiencesExperienceIdPatch(params)
      return res
    }
    catch (err) {
      return errorHandler(err)
    }
  }

  async function getExperienceMediaUploadUrl(contentType: string): Promise<UploadURLResponse> {
    const body: GetUploadURLRequest = { fileType: contentType }
    if (contentType.includes('image/')) {
      try {
        const params: V1UploadProductsImagePostRequest = { getUploadURLRequest: body }
        return await uploadApi().v1UploadProductsImagePost(params)
      }
      catch (err) {
        return errorHandler(err, { 400: 'このファイルはアップロードできません。' })
      }
    }
    if (contentType.includes('video/')) {
      try {
        const params: V1UploadProductsVideoPostRequest = { getUploadURLRequest: body }
        return await uploadApi().v1UploadProductsVideoPost(params)
      }
      catch (err) {
        return errorHandler(err, { 400: 'このファイルはアップロードできません。' })
      }
    }
    throw new Error('不明なMINEタイプです。')
  }

  async function uploadExperienceMedia(payload: File): Promise<string> {
    const contentType = payload.type
    try {
      const res = await getExperienceMediaUploadUrl(contentType)
      return await fileUpload(uploadApi(), payload, res.key, res.url)
    }
    catch (err) {
      return errorHandler(err, { 400: 'このファイルはアップロードできません。' })
    }
  }

  async function deleteExperience(experienceId: string) {
    try {
      const params: V1ExperiencesExperienceIdDeleteRequest = { experienceId }
      await experienceApi().v1ExperiencesExperienceIdDelete(params)
      const index = experiences.value.findIndex(e => e.id === experienceId)
      experiences.value.splice(index, 1)
      totalItems.value--
    }
    catch (err) {
      return errorHandler(err, {
        403: '体験を削除する権限がありません',
        404: '対象の商品が存在しません',
      })
    }
  }

  return {
    experience,
    experiences,
    totalItems,
    searchExperiences,
    fetchExperiences,
    fetchExperience,
    createExperience,
    updateExperience,
    uploadExperienceMedia,
    getExperienceMediaUploadUrl,
    deleteExperience,
  }
})
