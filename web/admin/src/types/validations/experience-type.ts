import type { ValidationArgs } from '@vuelidate/core'
import { maxLength, required } from '~/lib/validations'

export const CreateExperienceTypeValidationRules: ValidationArgs = {
  name: { required, maxlength: maxLength(32) },
}

export const UpdateExperienceTypeValidationRules: ValidationArgs = {
  name: { required, maxlength: maxLength(32) },
}
