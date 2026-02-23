import { useApiClient } from '~/composables/useApiClient'
import { ExperienceTypeApi } from '~/types/api/v1'
import type { CreateExperienceTypeRequest, ExperienceType, V1ExperienceTypesExperienceTypeIdDeleteRequest, V1ExperienceTypesExperienceTypeIdPatchRequest, V1ExperienceTypesGetRequest, V1ExperienceTypesPostRequest } from '~/types/api/v1'

export const useExperienceTypeStore = defineStore('experienceType', () => {
  const { create, errorHandler } = useApiClient()
  const experienceTypeApi = () => create(ExperienceTypeApi)

  const experienceType = ref<ExperienceType>({} as ExperienceType)
  const experienceTypes = ref<ExperienceType[]>([])
  const totalItems = ref<number>(0)

  async function fetchExperienceTypes(limit = 20, offset = 0): Promise<void> {
    try {
      const params: V1ExperienceTypesGetRequest = { limit, offset }
      const res = await experienceTypeApi().v1ExperienceTypesGet(params)
      experienceTypes.value = res.experienceTypes
      totalItems.value = res.total
    }
    catch (err) {
      return errorHandler(err)
    }
  }

  async function searchExperienceTypes(name = '', experienceTypeIds: string[] = []): Promise<void> {
    try {
      const params: V1ExperienceTypesGetRequest = { name }
      const res = await experienceTypeApi().v1ExperienceTypesGet(params)
      const merged: ExperienceType[] = []
      experienceTypes.value.forEach((et: ExperienceType): void => {
        if (!experienceTypeIds.includes(et.id)) {
          return
        }
        merged.push(et)
      })
      res.experienceTypes.forEach((et: ExperienceType): void => {
        if (merged.find((v): boolean => v.id === et.id)) {
          return
        }
        merged.push(et)
      })
      experienceTypes.value = merged
      totalItems.value = res.total
    }
    catch (err) {
      return errorHandler(err)
    }
  }

  async function createExperienceType(payload: CreateExperienceTypeRequest): Promise<void> {
    try {
      const params: V1ExperienceTypesPostRequest = { createExperienceTypeRequest: payload }
      await experienceTypeApi().v1ExperienceTypesPost(params)
    }
    catch (err) {
      return errorHandler(err, {
        400: '必須項目が不足しているか、内容に誤りがあります。',
        409: 'この体験カテゴリ名はすでに登録されています。',
      })
    }
  }

  async function updateExperienceType(experienceTypeId: string, payload: CreateExperienceTypeRequest): Promise<void> {
    try {
      const params: V1ExperienceTypesExperienceTypeIdPatchRequest = {
        experienceTypeId,
        updateExperienceTypeRequest: payload,
      }
      await experienceTypeApi().v1ExperienceTypesExperienceTypeIdPatch(params)
    }
    catch (err) {
      return errorHandler(err, {
        400: '必須項目が不足しているか、内容に誤りがあります。',
        404: 'この体験カテゴリは存在しません。',
        409: 'この体験カテゴリ名はすでに登録されています。',
      })
    }
  }

  async function deleteExperienceType(experienceTypeId: string): Promise<void> {
    try {
      const params: V1ExperienceTypesExperienceTypeIdDeleteRequest = { experienceTypeId }
      await experienceTypeApi().v1ExperienceTypesExperienceTypeIdDelete(params)
    }
    catch (err) {
      return errorHandler(err, { 404: 'この体験カテゴリは存在しません。' })
    }
  }

  return {
    experienceType,
    experienceTypes,
    totalItems,
    fetchExperienceTypes,
    searchExperienceTypes,
    createExperienceType,
    updateExperienceType,
    deleteExperienceType,
  }
})
