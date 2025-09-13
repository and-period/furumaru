import type { Administrator, CreateAdministratorRequest, UpdateAdministratorRequest, V1AdministratorsAdminIdDeleteRequest, V1AdministratorsAdminIdGetRequest, V1AdministratorsAdminIdPatchRequest, V1AdministratorsGetRequest, V1AdministratorsPostRequest } from '~/types/api/v1'

export const useAdministratorStore = defineStore('administrator', {
  state: () => ({
    administrator: {} as Administrator,
    administrators: [] as Administrator[],
    total: 0,
  }),

  actions: {
    /**
     * 管理者一覧を取得する非同期関数
     * @param limit 取得上限数
     * @param offset 取得開始位置
     */
    async fetchAdministrators(limit = 20, offset = 0): Promise<void> {
      try {
        const params: V1AdministratorsGetRequest = {
          limit,
          offset,
        }
        const res = await this.administratorApi().v1AdministratorsGet(params)
        this.administrators = res.administrators
        this.total = res.total
      }
      catch (err) {
        return this.errorHandler(err)
      }
    },

    /**
     * 管理者を取得する非同期関数
     * @param administratorId 管理者ID
     */
    async getAdministrator(administratorId: string): Promise<void> {
      try {
        const params: V1AdministratorsAdminIdGetRequest = {
          adminId: administratorId,
        }
        const res = await this.administratorApi().v1AdministratorsAdminIdGet(params)
        this.administrator = res.administrator
      }
      catch (err) {
        return this.errorHandler(err, { 404: '対象のシステム管理者が存在しません' })
      }
    },

    /**
     * 管理者を登録する非同期関数
     * @param payload 登録リクエスト
     */
    async createAdministrator(payload: CreateAdministratorRequest): Promise<void> {
      try {
        const params: V1AdministratorsPostRequest = {
          createAdministratorRequest: payload,
        }
        await this.administratorApi().v1AdministratorsPost(params)
      }
      catch (err) {
        return this.errorHandler(err, {
          400: '必須項目が不足しているか、内容に誤りがあります',
          409: 'このメールアドレスはすでに登録されているため、登録できません。',
        })
      }
    },

    /**
     * 管理者を更新する非同期関数
     * @param administratorId 管理者ID
     * @param payload 更新リクエスト
     */
    async updateAdministrator(administratorId: string, payload: UpdateAdministratorRequest): Promise<void> {
      try {
        const params: V1AdministratorsAdminIdPatchRequest = {
          adminId: administratorId,
          updateAdministratorRequest: payload,
        }
        await this.administratorApi().v1AdministratorsAdminIdPatch(params)
      }
      catch (err) {
        return this.errorHandler(err, {
          400: '必須項目が不足しているか、内容に誤りがあります',
          404: '対象のシステム管理者が存在しません',
          409: 'このメールアドレスはすでに登録されています。',
        })
      }
    },

    /**
     * 管理者を削除する非同期関数
     * @param administratorId 管理者ID
     */
    async deleteAdministrator(administratorId: string): Promise<void> {
      try {
        const params: V1AdministratorsAdminIdDeleteRequest = {
          adminId: administratorId,
        }
        await this.administratorApi().v1AdministratorsAdminIdDelete(params)
      }
      catch (err: any) {
        return this.errorHandler(err, {
          400: '必須項目が不足しているか、内容に誤りがあります',
          404: '対象のシステム管理者が存在しません',
        })
      }
    },
  },
})
