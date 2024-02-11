import type { ValidationArgs } from '@vuelidate/core'
import { required } from '~/lib/validations'

export const TimeDataValidationRules: ValidationArgs = {
  date: { required },
  time: { required }
}
