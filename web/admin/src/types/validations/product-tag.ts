import type { ValidationArgs } from '@vuelidate/core'
import { maxLength, required } from '~/lib/validations'

export const CreateProductTagValidationRules: ValidationArgs = {
  name: { required, maxlength: maxLength(32) }
}

export const UpdateProductTagValidationRules: ValidationArgs = {
  name: { required, maxlength: maxLength(32) }
}
