import type { ValidationArgs } from '@vuelidate/core'
import { maxLength, required } from '~/lib/validations'

export const CreateCategoryValidationRules: ValidationArgs = {
  name: { required, maxlength: maxLength(32) },
}

export const UpdateCategoryValidationRules: ValidationArgs = {
  name: { required, maxlength: maxLength(32) },
}
