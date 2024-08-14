import type { Composer, UseI18nOptions } from 'vue-i18n'
import { OrderStatus, PaymentMethodType } from '~/types/api'
import type { I18n } from '~/types/locales/i18n'

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
  i18n: Composer<NonNullable<UseI18nOptions['messages']>, NonNullable<UseI18nOptions['datetimeFormats']>, NonNullable<UseI18nOptions['numberFormats']>, UseI18nOptions['locale'] extends unknown ? string : UseI18nOptions['locale']> ,
): string {
  const paymentMethodText = (str: keyof I18n['purchase']['confirmation']) => {
    return i18n.t(`purchase.confirmation.${str}`)
  }

  switch (methodType) {
    case PaymentMethodType.CASH:
      return paymentMethodText('paymentMethodCashText')
    case PaymentMethodType.CREDIT_CARD:
      return paymentMethodText('paymentMethodCreditCardText')
    case PaymentMethodType.KONBINI:
      return paymentMethodText('paymentMethodConvinienceStoreText')
    case PaymentMethodType.BANK_TRANSFER:
      return paymentMethodText('paymentMethodBankTransferText')
    case PaymentMethodType.PAYPAY:
      return paymentMethodText('paymentMethodPayPayText')
    case PaymentMethodType.LINE_PAY:
      return paymentMethodText('paymentMethodLinePayText')
    case PaymentMethodType.MERPAY:
      return paymentMethodText('paymentMethodMerPayText')
    case PaymentMethodType.RAKUTEN_PAY:
      return paymentMethodText('paymentMethodRakutenPayText')
    case PaymentMethodType.AU_PAY:
      return paymentMethodText('paymentMethodAUPayText')
    case PaymentMethodType.UNKNOWN:
    default:
      return ''
  }
}
