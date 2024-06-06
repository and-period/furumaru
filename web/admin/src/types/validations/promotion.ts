import type { ValidationArgs } from '@vuelidate/core'
import { maxLength, minLength, minValue, required } from '~/lib/validations'

export const CreatePromotionValidationRules: ValidationArgs = {
  title: { required, maxLength: maxLength(64) },
  description: { required, maxLength: maxLength(2000) },
  discountType: {},
  discountRate: { minValue: minValue(0) },
  code: { required, minLength: minLength(8), maxLength: maxLength(8) },
}

export const UpdatePromotionValidationRules: ValidationArgs = {
  title: { required, maxLength: maxLength(64) },
  description: { required, maxLength: maxLength(2000) },
  discountType: {},
  discountRate: { minValue: minValue(0) },
  code: { required, minLength: minLength(8), maxLength: maxLength(8) },
}
