import { OrderStatus } from '~/types/api'

/**
 * ステータスを文字列に変換する関数
 * https://github.com/and-period/furumaru/blob/main/docs/swagger/user/components/schemas/codes/order.yaml
 * @param status
 * @returns
 */
export function getOrderStatusString(status: OrderStatus): string {
  switch (status) {
    case OrderStatus.UNPAID:
      return '支払い待ち'
    case OrderStatus.PREPARING:
      return '配送対応中'
    case OrderStatus.COMPLETED:
      return '完了'
    case OrderStatus.CANCELED:
      return 'キャンセル済み'
    case OrderStatus.REFUNDED:
      return '返金済み'
    case OrderStatus.FAILED:
      return '失敗'
    default:
      return '不明'
  }
}

/** ステータスを3値に変換する関数 */
export function getOperationResultFromOrderStatus(status: OrderStatus): string {
  switch (status) {
    case OrderStatus.UNPAID:
    case OrderStatus.PREPARING:
    case OrderStatus.COMPLETED:
      return 'success'
    case OrderStatus.CANCELED:
    case OrderStatus.REFUNDED:
      return 'canceled'
    case OrderStatus.FAILED:
      return 'failed'
    default:
      return 'unknown'
  }
}
