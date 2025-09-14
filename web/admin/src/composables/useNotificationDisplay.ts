import { unix } from 'dayjs'
import { DiscountType } from '~/types/api/v1'
import type { Promotion } from '~/types/api/v1'

export const useNotificationDisplay = () => {
  const getDateTime = (unixTime: number): string => {
    if (unixTime === 0) {
      return ''
    }
    return unix(unixTime).format('YYYY/MM/DD HH:mm')
  }

  const getPromotionTerm = (promotion: Promotion | undefined): string => {
    if (!promotion) {
      return ''
    }

    const startAt = getDateTime(promotion.startAt)
    const endAt = getDateTime(promotion.endAt)
    return `${startAt} ~ ${endAt}`
  }

  const getPromotionDiscount = (promotion: Promotion | undefined): string => {
    if (!promotion) {
      return ''
    }

    switch (promotion.discountType) {
      case DiscountType.DiscountTypeAmount:
        return '￥' + promotion.discountRate.toLocaleString()
      case DiscountType.DiscountTypeRate:
        return promotion.discountRate + '%'
      case DiscountType.DiscountTypeFreeShipping:
        return '送料無料'
      default:
        return ''
    }
  }

  return {
    getDateTime,
    getPromotionTerm,
    getPromotionDiscount,
  }
}
