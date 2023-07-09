import { defineStore } from 'pinia'

import { useCommonStore } from './common'
import { useAdminStore } from './admin'
import {
  CreateNotificationRequest,
  Notification,
  UpdateNotificationRequest
} from '~/types/api'
import { apiClient } from '~/plugins/api-client'

export const useNotificationStore = defineStore('notification', {
  state: () => ({
    notification: {} as Notification,
    notifications: [] as Notification[],
    totalItems: 0
  }),
  actions: {
    /**
     * 登録済みのお知らせ一覧を取得する非同期関数
     * @param limit 取得上限数
     * @param offset 取得開始位置
     * @param orders ソートキー
     * @returns
     */
    async fetchNotifications (limit = 20, offset = 0, orders = []): Promise<void> {
      try {
        const res = await apiClient.notificationApi().v1ListNotifications(limit, offset, undefined, undefined, orders.join(''))

        const adminStore = useAdminStore()
        this.notifications = res.data.notifications
        this.totalItems = res.data.total
        adminStore.admins = res.data.admins
      } catch (err) {
        return this.errorHandler(err)
      }
    },

    /**
     * お知らせIDからお知らせ情報情報を取得する非同期関数
     * @param id お知らせID
     * @returns お知らせ情報
     */
    async getNotification (id: string): Promise<Notification> {
      try {
        const res = await apiClient.notificationApi().v1GetNotification(id)

        const adminStore = useAdminStore()
        this.notification = res.data.notification
        adminStore.admin = res.data.admin
        return res.data.notification
      } catch (err) {
        return this.errorHandler(err)
      }
    },

    /**
     * お知らせを登録する非同期関数
     * @param payload
     */
    async createNotification (
      payload: CreateNotificationRequest
    ): Promise<void> {
      try {
        await apiClient.notificationApi().v1CreateNotification(payload)
        const commonStore = useCommonStore()
        commonStore.addSnackbar({
          message: `${payload.title}を作成しました。`,
          color: 'info'
        })
      } catch (err) {
        return this.errorHandler(err)
      }
    },

    /**
     * お知らせを削除する非同期関数
     * @param id お知らせID
     */
    async deleteNotification (id: string): Promise<void> {
      const commonStore = useCommonStore()
      try {
        await apiClient.notificationApi().v1DeleteNotification(id)
        commonStore.addSnackbar({
          message: '品物削除が完了しました',
          color: 'info'
        })
      } catch (err) {
        return this.errorHandler(err)
      }
      this.fetchNotifications()
    },

    /**
     * お知らせ情報を編集する非同期関数
     * @param id セールID
     * @param payload
     */
    async updateNotification (
      id: string,
      payload: UpdateNotificationRequest
    ): Promise<void> {
      const commonStore = useCommonStore()
      try {
        await apiClient.notificationApi().v1UpdateNotification(id, payload)
        commonStore.addSnackbar({
          message: 'お知らせ情報の編集が完了しました',
          color: 'info'
        })
      } catch (err) {
        return this.errorHandler(err)
      }
    }
  }
})
