import { NotificationTarget, NotificationType } from '~/types/api/v1'

export const NOTIFICATION_TYPES = [
  { title: 'システム関連', value: NotificationType.NotificationTypeSystem },
  { title: 'ライブ関連', value: NotificationType.NotificationTypeLive },
  { title: 'セール関連', value: NotificationType.NotificationTypePromotion },
  { title: 'その他', value: NotificationType.NotificationTypeOther },
]

export const NOTIFICATION_TARGETS = [
  { title: 'ユーザー', value: NotificationTarget.NotificationTargetUsers },
  { title: '生産者', value: NotificationTarget.NotificationTargetProducers },
  { title: 'コーディネーター', value: NotificationTarget.NotificationTargetCoordinators },
  { title: '管理者', value: NotificationTarget.NotificationTargetAdministrators },
]
