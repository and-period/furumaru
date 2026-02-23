import { useApiClient } from '~/composables/useApiClient'
import { AdministratorApi } from '~/types/api/v1'
import type {
  Administrator,
  CreateAdministratorRequest,
  UpdateAdministratorRequest,
  V1AdministratorsAdminIdDeleteRequest,
  V1AdministratorsAdminIdGetRequest,
  V1AdministratorsAdminIdPatchRequest,
  V1AdministratorsGetRequest,
  V1AdministratorsPostRequest,
} from '~/types/api/v1'

export const useAdministratorStore = defineStore('administrator', () => {
  const { create, errorHandler } = useApiClient()
  const administratorApi = () => create(AdministratorApi)

  const administrator = ref<Administrator>({} as Administrator)
  const administrators = ref<Administrator[]>([])
  const total = ref<number>(0)

  async function fetchAdministrators(limit = 20, offset = 0): Promise<void> {
    try {
      const params: V1AdministratorsGetRequest = { limit, offset }
      const res = await administratorApi().v1AdministratorsGet(params)
      administrators.value = res.administrators
      total.value = res.total
    }
    catch (err) {
      return errorHandler(err)
    }
  }

  async function getAdministrator(administratorId: string): Promise<void> {
    try {
      const params: V1AdministratorsAdminIdGetRequest = { adminId: administratorId }
      const res = await administratorApi().v1AdministratorsAdminIdGet(params)
      administrator.value = res.administrator
    }
    catch (err) {
      return errorHandler(err, { 404: '対象のシステム管理者が存在しません' })
    }
  }

  async function createAdministrator(payload: CreateAdministratorRequest): Promise<void> {
    try {
      const params: V1AdministratorsPostRequest = { createAdministratorRequest: payload }
      await administratorApi().v1AdministratorsPost(params)
    }
    catch (err) {
      return errorHandler(err, {
        400: '必須項目が不足しているか、内容に誤りがあります',
        409: 'このメールアドレスはすでに登録されているため、登録できません。',
      })
    }
  }

  async function updateAdministrator(administratorId: string, payload: UpdateAdministratorRequest): Promise<void> {
    try {
      const params: V1AdministratorsAdminIdPatchRequest = {
        adminId: administratorId,
        updateAdministratorRequest: payload,
      }
      await administratorApi().v1AdministratorsAdminIdPatch(params)
    }
    catch (err) {
      return errorHandler(err, {
        400: '必須項目が不足しているか、内容に誤りがあります',
        404: '対象のシステム管理者が存在しません',
        409: 'このメールアドレスはすでに登録されています。',
      })
    }
  }

  async function deleteAdministrator(administratorId: string): Promise<void> {
    try {
      const params: V1AdministratorsAdminIdDeleteRequest = { adminId: administratorId }
      await administratorApi().v1AdministratorsAdminIdDelete(params)
    }
    catch (err) {
      return errorHandler(err, {
        400: '必須項目が不足しているか、内容に誤りがあります',
        404: '対象のシステム管理者が存在しません',
      })
    }
  }

  return {
    administrator,
    administrators,
    total,
    fetchAdministrators,
    getAdministrator,
    createAdministrator,
    updateAdministrator,
    deleteAdministrator,
  }
})
