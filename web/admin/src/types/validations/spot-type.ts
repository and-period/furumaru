import type { ValidationArgs } from '@vuelidate/core'
import { maxLength, required } from '~/lib/validations'

export const CreateSpotTypeValidationRules: ValidationArgs = {
  name: { required, maxlength: maxLength(32) },
}

export const UpdateSpotTypeValidationRules: ValidationArgs = {
  name: { required, maxlength: maxLength(32) },
}
