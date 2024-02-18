import type { ValidationArgs } from '@vuelidate/core'
import { required, notSameAs } from '~/lib/validations'

export const TimeDataValidationRules: ValidationArgs = {
  date: { required },
  time: { required }
}

export const NotSameTimeDataValidationRules = (startAt: number, otherName?: string): ValidationArgs => ({
  startAt: { required },
  endAt: { required, notSameas: notSameAs(startAt, otherName) }
})
