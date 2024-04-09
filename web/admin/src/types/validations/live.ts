import type { ValidationArgs } from '@vuelidate/core'
import { maxLength, maxLengthArray, required } from '~/lib/validations'

export const CreateLiveValidationRules: ValidationArgs = {
  producerId: { required },
  productIds: { maxLengthArray: maxLengthArray(100) },
  comment: { required, maxLength: maxLength(2000) }
}

export const UpdateLiveValidationRules: ValidationArgs = {
  productIds: { maxLengthArray: maxLengthArray(100) },
  comment: { required, maxLength: maxLength(2000) }
}
