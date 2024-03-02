import { defineStore } from 'pinia'

import { apiClient } from '~/plugins/api-client'
import type { Administrator, CreateAdministratorRequest, UpdateAdministratorRequest } from '~/types/api'

export const useAdministratorStore = defineStore('administrator', {
  state: () => ({
    administrator: {} as Administrator,
    administrators: [] as Administrator[],
    total: 0
  }),

  actions: {
    /**
     * 管理者一覧を取得する非同期関数
     * @param limit 取得上限数
     * @param offset 取得開始位置
     */
    async fetchAdministrators (limit = 20, offset = 0): Promise<void> {
      try {
        const res = await apiClient.administratorApi().v1ListAdministrators(
          limit,
          offset
        )
        this.administrators = res.data.administrators
        this.total = res.data.total
      } catch (err) {
        return this.errorHandler(err)
      }
    },

    /**
     * 管理者を取得する非同期関数
     * @param administratorId 管理者ID
     */
    async getAdministrator (administratorId: string): Promise<void> {
      try {
        const res = await apiClient.administratorApi().v1GetAdministrator(administratorId)
        this.administrator = res.data.administrator
      } catch (err: any) {
        if (err?.response?.status === 404) {
          return this.errorHandler(err, { 404: '対象のシステム管理者が存在しません' })
        }
        return this.errorHandler(err)
      }
    },

    /**
     * 管理者を登録する非同期関数
     * @param payload 登録リクエスト
     */
    async createAdministrator (payload: CreateAdministratorRequest): Promise<void> {
      try {
        await apiClient.administratorApi().v1CreateAdministrator(payload)
      } catch (err: any) {
        if (err?.response?.status === 400) {
          return this.errorHandler(err, { 400: '必須項目が不足しているか、内容に誤りがあります' })
        }
        return this.errorHandler(err, { 409: 'このメールアドレスはすでに登録されているため、登録できません。' })
      }
    },

    /**
     * 管理者を更新する非同期関数
     * @param administratorId 管理者ID
     * @param payload 更新リクエスト
     */
    async updateAdministrator (administratorId: string, payload: UpdateAdministratorRequest): Promise<void> {
      try {
        await apiClient.administratorApi().v1UpdateAdministrator(administratorId, payload)
      } catch (err: any) {
        if (err?.response?.status === 400) {
          return this.errorHandler(err, { 400: '必須項目が不足しているか、内容に誤りがあります' })
        }
        if (err?.response?.status === 404) {
          return this.errorHandler(err, { 404: '対象のシステム管理者が存在しません' })
        }
        return this.errorHandler(err, { 409: 'このメールアドレスはすでに登録されています。' })
      }
    },

    /**
     * 管理者を削除する非同期関数
     * @param administratorId 管理者ID
     */
    async deleteAdministrator (administratorId: string): Promise<void> {
      try {
        await apiClient.administratorApi().v1DeleteAdministrator(administratorId)
      } catch (err: any) {
        if (err?.response?.status === 400) {
          return this.errorHandler(err, { 400: '必須項目が不足しているか、内容に誤りがあります' })
        }
        if (err?.response?.status === 404) {
          return this.errorHandler(err, { 404: '対象のシステム管理者が存在しません' })
        }
        return this.errorHandler(err)
      }
    }
  }
})
