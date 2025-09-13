import { PaymentMethodType } from '~/types/api/v1'

export interface PaymentListItem {
  name: string
  value: PaymentMethodType
}

/**
 * 決済システム一覧
 */
export const paymentsList: PaymentListItem[] = [
  { name: '代引支払い', value: PaymentMethodType.PaymentMethodTypeCash },
  { name: 'クレジットカード決済', value: PaymentMethodType.PaymentMethodTypeCreditCard },
  { name: 'コンビニ決済', value: PaymentMethodType.PaymentMethodTypeKonbini },
  { name: '銀行振込決済', value: PaymentMethodType.PaymentMethodTypeBankTransfer },
  { name: 'QR決済（PayPay）', value: PaymentMethodType.PaymentMethodTypePayPay },
  { name: 'QR決済（Line Pay）', value: PaymentMethodType.PaymentMethodTypeLinePay },
  { name: 'QR決済（メルペイ）', value: PaymentMethodType.PaymentMethodTypeMerpay },
  { name: 'QR決済（楽天ペイ）', value: PaymentMethodType.PaymentMethodTypeRakutenPay },
  { name: 'QR決済（au PAY）', value: PaymentMethodType.PaymentMethodTypeAUPay },
  { name: 'ペイディ（Paidy）', value: PaymentMethodType.PaymentMethodTypePaidy },
  { name: 'ペイジー（Pay-easy）', value: PaymentMethodType.PaymentMethodTypePayEasy },
]
