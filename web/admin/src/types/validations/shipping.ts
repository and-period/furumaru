import type { ValidationArgs } from '@vuelidate/core'
import { maxLengthArray, maxValue, minLengthArray, minValue, required } from '~/lib/validations'

export const UpsertShippingValidationRules: ValidationArgs = {
  hasFreeShipping: {},
  box60Frozen: { required, minValue: minValue(0), maxValue: maxValue(9999999999) },
  box80Frozen: { required, minValue: minValue(0), maxValue: maxValue(9999999999) },
  box100Frozen: { required, minValue: minValue(0), maxValue: maxValue(9999999999) }
}

export const UpsertShippingRateValidationRules: ValidationArgs = {
  name: { required, maxValue: maxValue(64) },
  price: { required, minValue: minValue(1), maxValue: maxValue(9999999999) },
  prefectureCodes: { minLengthArray: minLengthArray(1), maxLengthArray: maxLengthArray(47) }
}
