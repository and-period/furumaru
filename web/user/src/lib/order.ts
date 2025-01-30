import { OrderStatus, PaymentMethodType } from '~/types/api'

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

/**
 * 決済手段を文字列に変換する関数
 * @param methodType
 * @returns
 */
export function getPaymentMethodNameByPaymentMethodType(
  methodType: PaymentMethodType,
): string {
  switch (methodType) {
    case PaymentMethodType.CASH:
      return '現金支払い'
    case PaymentMethodType.CREDIT_CARD:
      return 'クレジットカード決済'
    case PaymentMethodType.KONBINI:
      return 'コンビニ決済'
    case PaymentMethodType.BANK_TRANSFER:
      return '銀行振込決済'
    case PaymentMethodType.PAYPAY:
      return 'QR決済（PayPay）'
    case PaymentMethodType.LINE_PAY:
      return 'QR決済（Line Pay）'
    case PaymentMethodType.MERPAY:
      return 'QR決済（メルペイ）'
    case PaymentMethodType.RAKUTEN_PAY:
      return 'QR決済（楽天ペイ）'
    case PaymentMethodType.AU_PAY:
      return 'QR決済（au PAY）'
    case PaymentMethodType.PAIDY:
      return 'ペイディ（Paidy）'
    case PaymentMethodType.PAY_EASY:
      return 'ペイジー（Pay-easy）'
    case PaymentMethodType.UNKNOWN:
    default:
      return ''
  }
}
