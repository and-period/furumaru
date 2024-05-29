import type { ValidationArgs } from '@vuelidate/core'
import { maxLength, maxLengthArray, required } from '~/lib/validations'

export const CreateNotificationValidationRules: ValidationArgs = {
  type: { required },
  promotionId: { required },
  title: { maxLength: maxLength(128) },
  body: { required, maxLength: maxLength(2000) },
  note: { maxLength: maxLength(2000) },
  targets: { maxLengthArray: maxLengthArray(4) },
}

export const UpdateNotificationValidationRules: ValidationArgs = {
  title: { maxLength: maxLength(128) },
  body: { required, maxLength: maxLength(2000) },
  note: { maxLength: maxLength(2000) },
  targets: { maxLengthArray: maxLengthArray(4) },
}
