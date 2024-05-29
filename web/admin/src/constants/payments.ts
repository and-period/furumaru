import { PaymentMethodType } from '~/types/api'

export interface PaymentListItem {
  name: string
  value: PaymentMethodType
}

/**
 * 決済システム一覧
 */
export const paymentsList: PaymentListItem[] = [
  { name: '代引支払い', value: PaymentMethodType.CASH },
  { name: 'クレジットカード決済', value: PaymentMethodType.CREDIT_CARD },
  { name: 'コンビニ決済', value: PaymentMethodType.KONBINI },
  { name: '銀行振込決済', value: PaymentMethodType.BANK_TRANSFER },
  { name: 'QR決済（PayPay）', value: PaymentMethodType.PAYPAY },
  { name: 'QR決済（Line Pay）', value: PaymentMethodType.LINE_PAY },
  { name: 'QR決済（メルペイ）', value: PaymentMethodType.MERPAY },
  { name: 'QR決済（楽天ペイ）', value: PaymentMethodType.RAKUTEN_PAY },
  { name: 'QR決済（au PAY）', value: PaymentMethodType.AU_PAY },
]
