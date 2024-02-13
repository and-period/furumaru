import type { ValidationArgs } from '@vuelidate/core'
import { maxLength, required } from '~/lib/validations'

export const CreateScheduleValidationRules: ValidationArgs = {
  title: { required, maxLength: maxLength(64) },
  description: { required, maxLength: maxLength(2000) },
  thumbnailUrl: { required }
}

export const UpdateScheduleValidationRules: ValidationArgs = {
  title: { required, maxLength: maxLength(64) },
  description: { required, maxLength: maxLength(2000) }
}
