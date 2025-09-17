import { OrderStatus, OrderType, ShippingCarrier, ShippingType } from "~/types/api/v1"

export const orderTypes = [
  { title: '商品', value: OrderType.OrderTypeProduct },
  { title: '体験', value: OrderType.OrderTypeExperience },
]

export const orderStatuses = [
  { title: '支払い待ち', value: OrderStatus.OrderStatusUnpaid },
  { title: '受注待ち', value: OrderStatus.OrderStatusWaiting },
  { title: '発送準備中', value: OrderStatus.OrderStatusPreparing },
  { title: '発送済み', value: OrderStatus.OrderStatusShipped },
  { title: '完了', value: OrderStatus.OrderStatusCompleted },
  { title: 'キャンセル', value: OrderStatus.OrderStatusCanceled },
  { title: '返金', value: OrderStatus.OrderStatusRefunded },
  { title: '失敗', value: OrderStatus.OrderStatusFailed },
]

export const orderShippingTypes = [
  { title: '常温・冷蔵便', value: ShippingType.ShippingTypeNormal },
  { title: '冷凍便', value: ShippingType.ShippingTypeFrozen },
  { title: '店舗受け取り', value: ShippingType.ShippingTypePickup },
]

export const fulfillmentCompanies = [
  { title: '指定なし', value: ShippingCarrier.ShippingCarrierUnknown },
  { title: '佐川急便', value: ShippingCarrier.ShippingCarrierSagawa },
  { title: 'ヤマト運輸', value: ShippingCarrier.ShippingCarrierYamato },
]
