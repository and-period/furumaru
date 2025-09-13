import { useAdminStore } from './admin'
import type {
  CreateNotificationRequest,
  Notification,
  UpdateNotificationRequest,
  V1NotificationsGetRequest,
  V1NotificationsNotificationIdDeleteRequest,
  V1NotificationsNotificationIdGetRequest,
  V1NotificationsNotificationIdPatchRequest,
  V1NotificationsPostRequest,
} from '~/types/api/v1'

export const useNotificationStore = defineStore('notification', {
  state: () => ({
    notification: {} as Notification,
    notifications: [] as Notification[],
    totalItems: 0,
  }),
  actions: {
    /**
     * 登録済みのお知らせ一覧を取得する非同期関数
     * @param limit 取得上限数
     * @param offset 取得開始位置
     * @param orders ソートキー
     * @returns
     */
    async fetchNotifications(limit = 20, offset = 0, orders = []): Promise<void> {
      try {
        const params: V1NotificationsGetRequest = {
          limit,
          offset,
          orders: orders.join(','),
        }
        const res = await this.notificationApi().v1NotificationsGet(params)

        const adminStore = useAdminStore()
        this.notifications = res.notifications
        this.totalItems = res.total
        adminStore.admins = res.admins
      }
      catch (err) {
        return this.errorHandler(err)
      }
    },

    /**
     * お知らせIDからお知らせ情報情報を取得する非同期関数
     * @param id お知らせID
     * @returns お知らせ情報
     */
    async getNotification(id: string): Promise<Notification> {
      try {
        const params: V1NotificationsNotificationIdGetRequest = {
          notificationId: id,
        }
        const res = await this.notificationApi().v1NotificationsNotificationIdGet(params)

        const adminStore = useAdminStore()
        this.notification = res.notification
        adminStore.admin = res.admin
        return res.notification
      }
      catch (err) {
        return this.errorHandler(err, { 404: '対象のお知らせが存在しません' })
      }
    },

    /**
     * お知らせを登録する非同期関数
     * @param payload
     */
    async createNotification(
      payload: CreateNotificationRequest,
    ): Promise<void> {
      try {
        const params: V1NotificationsPostRequest = {
          createNotificationRequest: payload,
        }
        await this.notificationApi().v1NotificationsPost(params)
      }
      catch (err) {
        return this.errorHandler(err, { 400: '必須項目が不足しているか、内容に誤りがあります' })
      }
    },

    /**
     * お知らせを削除する非同期関数
     * @param id お知らせID
     */
    async deleteNotification(id: string): Promise<void> {
      try {
        const params: V1NotificationsNotificationIdDeleteRequest = {
          notificationId: id,
        }
        await this.notificationApi().v1NotificationsNotificationIdDelete(params)
      }
      catch (err) {
        return this.errorHandler(err, { 404: '対象のお知らせが存在しません' })
      }
      this.fetchNotifications()
    },

    /**
     * お知らせ情報を編集する非同期関数
     * @param id セールID
     * @param payload
     */
    async updateNotification(
      id: string,
      payload: UpdateNotificationRequest,
    ): Promise<void> {
      try {
        const params: V1NotificationsNotificationIdPatchRequest = {
          notificationId: id,
          updateNotificationRequest: payload,
        }
        await this.notificationApi().v1NotificationsNotificationIdPatch(params)
      }
      catch (err) {
        return this.errorHandler(err, {
          400: '必須項目が不足しているか、内容に誤りがあります',
          404: '対象のお知らせが存在しません',
        })
      }
    },
  },
})
