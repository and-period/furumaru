import { DiscountType } from '~/types/api/v1'

export const DISCOUNT_METHODS = [
  { method: '円', value: DiscountType.DiscountTypeAmount },
  { method: '%', value: DiscountType.DiscountTypeRate },
  { method: '送料無料', value: DiscountType.DiscountTypeFreeShipping },
]

export const PROMOTION_STATUS_OPTIONS = [
  { status: '有効', value: true },
  { status: '無効', value: false },
]
