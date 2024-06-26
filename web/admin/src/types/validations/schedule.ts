import type { ValidationArgs } from '@vuelidate/core'
import { maxLength, required } from '~/lib/validations'

export const CreateScheduleValidationRules: ValidationArgs = {
  title: { required, maxLength: maxLength(64) },
  description: { required, maxLength: maxLength(2000) },
  thumbnailUrl: { required },
  openingVideoUrl: { required },
  imageUrl: { required },
}

export const UpdateScheduleValidationRules: ValidationArgs = {
  title: { required, maxLength: maxLength(64) },
  description: { required, maxLength: maxLength(2000) },
}
