import { DiscountType } from '~/types/api/v1'
import type { CreatePromotionRequest, UpdatePromotionRequest } from '~/types/api/v1'

export const usePromotionValidation = (
  formData: Ref<CreatePromotionRequest | UpdatePromotionRequest>,
) => {
  const getDiscountErrorMessage = (): string => {
    switch (formData.value.discountType) {
      case DiscountType.DiscountTypeAmount:
        if (formData.value.discountRate >= 0) {
          return ''
        }
        return '0以上の値を指定してください'
      case DiscountType.DiscountTypeRate:
        if (formData.value.discountRate >= 0 && formData.value.discountRate <= 100) {
          return ''
        }
        return '0~100の値を指定してください'
      default:
        return ''
    }
  }

  return {
    getDiscountErrorMessage,
  }
}
