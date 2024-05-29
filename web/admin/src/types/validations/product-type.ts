import type { ValidationArgs } from '@vuelidate/core'
import { maxLength, required } from '~/lib/validations'

export const CreateProductTypeValidationRules: ValidationArgs = {
  name: { required, maxlength: maxLength(32) },
  iconUrl: { required },
}

export const UpdateProductTypeValidationRules: ValidationArgs = {
  name: { required, maxlength: maxLength(32) },
  iconUrl: { required },
}
