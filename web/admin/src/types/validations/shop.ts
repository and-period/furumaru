import type { ValidationArgs } from '@vuelidate/core'
import { maxLength, maxLengthArray, required } from '~/lib/validations'

export const UpdateShopValidationRules: ValidationArgs = {
  name: { required, maxLength: maxLength(64) },
  businessDays: { maxLengthArray: maxLengthArray(7) },
}
