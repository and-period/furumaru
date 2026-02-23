import { useAdminStore } from './admin'
import { useApiClient } from '~/composables/useApiClient'
import { NotificationApi } from '~/types/api/v1'
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

export const useNotificationStore = defineStore('notification', () => {
  const { create, errorHandler } = useApiClient()
  const notificationApi = () => create(NotificationApi)

  // state
  const notification = ref<Notification>({} as Notification)
  const notifications = ref<Notification[]>([])
  const totalItems = ref<number>(0)

  /**
   * 登録済みのお知らせ一覧を取得する非同期関数
   * @param limit 取得上限数
   * @param offset 取得開始位置
   * @param orders ソートキー
   * @returns
   */
  async function fetchNotifications(limit = 20, offset = 0, orders: string[] = []): Promise<void> {
    try {
      const params: V1NotificationsGetRequest = {
        limit,
        offset,
        orders: orders.join(','),
      }
      const res = await notificationApi().v1NotificationsGet(params)

      const adminStore = useAdminStore()
      notifications.value = res.notifications
      totalItems.value = res.total
      adminStore.admins = res.admins
    }
    catch (err) {
      return errorHandler(err)
    }
  }

  /**
   * お知らせIDからお知らせ情報情報を取得する非同期関数
   * @param id お知らせID
   * @returns お知らせ情報
   */
  async function getNotification(id: string): Promise<Notification> {
    try {
      const params: V1NotificationsNotificationIdGetRequest = {
        notificationId: id,
      }
      const res = await notificationApi().v1NotificationsNotificationIdGet(params)

      const adminStore = useAdminStore()
      notification.value = res.notification
      adminStore.admin = res.admin
      return res.notification
    }
    catch (err) {
      return errorHandler(err, { 404: '対象のお知らせが存在しません' })
    }
  }

  /**
   * お知らせを登録する非同期関数
   * @param payload
   */
  async function createNotification(
    payload: CreateNotificationRequest,
  ): Promise<void> {
    try {
      const params: V1NotificationsPostRequest = {
        createNotificationRequest: payload,
      }
      await notificationApi().v1NotificationsPost(params)
    }
    catch (err) {
      return errorHandler(err, { 400: '必須項目が不足しているか、内容に誤りがあります' })
    }
  }

  /**
   * お知らせを削除する非同期関数
   * @param id お知らせID
   */
  async function deleteNotification(id: string): Promise<void> {
    try {
      const params: V1NotificationsNotificationIdDeleteRequest = {
        notificationId: id,
      }
      await notificationApi().v1NotificationsNotificationIdDelete(params)
    }
    catch (err) {
      return errorHandler(err, { 404: '対象のお知らせが存在しません' })
    }
    fetchNotifications()
  }

  /**
   * お知らせ情報を編集する非同期関数
   * @param id セールID
   * @param payload
   */
  async function updateNotification(
    id: string,
    payload: UpdateNotificationRequest,
  ): Promise<void> {
    try {
      const params: V1NotificationsNotificationIdPatchRequest = {
        notificationId: id,
        updateNotificationRequest: payload,
      }
      await notificationApi().v1NotificationsNotificationIdPatch(params)
    }
    catch (err) {
      return errorHandler(err, {
        400: '必須項目が不足しているか、内容に誤りがあります',
        404: '対象のお知らせが存在しません',
      })
    }
  }

  return {
    // state
    notification,
    notifications,
    totalItems,
    // actions
    fetchNotifications,
    getNotification,
    createNotification,
    deleteNotification,
    updateNotification,
  }
})
