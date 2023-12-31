import type { ValidationArgs } from '@vuelidate/core'
import { required } from '~/lib/validations'

export const TimeDataValidationRules: ValidationArgs = {
  startDate: { required },
  startTime: { required },
  endDate: { required },
  endTime: { required }
}
