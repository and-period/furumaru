import { useApiClient } from '~/composables/useApiClient'
import { SpotTypeApi } from '~/types/api/v1'
import type { CreateSpotTypeRequest, SpotType, UpdateSpotTypeRequest, V1SpotTypesGetRequest, V1SpotTypesPostRequest, V1SpotTypesSpotTypeIdDeleteRequest, V1SpotTypesSpotTypeIdPatchRequest } from '~/types/api/v1'

export const useSpotTypeStore = defineStore('spotType', () => {
  const { create, errorHandler } = useApiClient()
  const spotTypeApi = () => create(SpotTypeApi)

  const spotType = ref<SpotType>({} as SpotType)
  const spotTypes = ref<SpotType[]>([])
  const total = ref<number>(0)

  async function fetchSpotTypes(limit = 20, offset = 0): Promise<void> {
    try {
      const params: V1SpotTypesGetRequest = { limit, offset, name: '' }
      const res = await spotTypeApi().v1SpotTypesGet(params)
      spotTypes.value = res.spotTypes
      total.value = res.total
    }
    catch (err) {
      return errorHandler(err)
    }
  }

  async function searchSpotTypes(name = '', spotTypeIds: string[] = []): Promise<void> {
    try {
      const params: V1SpotTypesGetRequest = { name }
      const res = await spotTypeApi().v1SpotTypesGet(params)
      const merged: SpotType[] = []
      spotTypes.value.forEach((st: SpotType): void => {
        if (!spotTypeIds.includes(st.id)) {
          return
        }
        merged.push(st)
      })
      res.spotTypes.forEach((st: SpotType): void => {
        if (merged.find((v): boolean => v.id === st.id)) {
          return
        }
        merged.push(st)
      })
      spotTypes.value = merged
      total.value = res.total
    }
    catch (err) {
      return errorHandler(err)
    }
  }

  async function createSpotType(payload: CreateSpotTypeRequest): Promise<void> {
    try {
      const params: V1SpotTypesPostRequest = { createSpotTypeRequest: payload }
      const res = await spotTypeApi().v1SpotTypesPost(params)
      spotTypes.value.unshift(res.spotType)
    }
    catch (err) {
      return errorHandler(err, {
        400: '必須項目が不足しているか、内容に誤りがあります。',
        409: 'このスポット名はすでに登録されています。',
      })
    }
  }

  async function updateSpotType(spotTypeId: string, payload: UpdateSpotTypeRequest): Promise<void> {
    try {
      const params: V1SpotTypesSpotTypeIdPatchRequest = {
        spotTypeId,
        updateSpotTypeRequest: payload,
      }
      await spotTypeApi().v1SpotTypesSpotTypeIdPatch(params)
    }
    catch (err) {
      return errorHandler(err, {
        400: '必須項目が不足しているか、内容に誤りがあります。',
        404: 'このスポット種別は存在しません。',
        409: 'このスポット種別名はすでに登録されています。',
      })
    }
  }

  async function deleteSpotType(spotTypeId: string): Promise<void> {
    try {
      const params: V1SpotTypesSpotTypeIdDeleteRequest = { spotTypeId }
      await spotTypeApi().v1SpotTypesSpotTypeIdDelete(params)
    }
    catch (err) {
      return errorHandler(err, { 404: 'このスポット種別は存在しません。' })
    }
  }

  return {
    spotType,
    spotTypes,
    total,
    fetchSpotTypes,
    searchSpotTypes,
    createSpotType,
    updateSpotType,
    deleteSpotType,
  }
})
