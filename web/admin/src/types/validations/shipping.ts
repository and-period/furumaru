import type { ValidationArgs } from '@vuelidate/core'
import { maxLengthArray, maxLength, maxValue, minLengthArray, minValue, required } from '~/lib/validations'

export const ShippingValidationRules: ValidationArgs = {
  name: { required, maxLength: maxLength(64) },
  hasFreeShipping: {},
  box60Frozen: { required, minValue: minValue(0), maxValue: maxValue(9999999999) },
  box80Frozen: { required, minValue: minValue(0), maxValue: maxValue(9999999999) },
  box100Frozen: { required, minValue: minValue(0), maxValue: maxValue(9999999999) },
}

export const ShippingRateValidationRules: ValidationArgs = {
  name: { required, maxLength: maxLength(64) },
  price: { required, minValue: minValue(0), maxValue: maxValue(9999999999) },
  prefectureCodes: { minLengthArray: minLengthArray(1), maxLengthArray: maxLengthArray(47) },
}

export const CreateShippingValidationRules: ValidationArgs = {
  name: { required, maxLength: maxLength(64) },
  hasFreeShipping: {},
  freeShippingRates: { required, minValue: minValue(0), maxValue: maxValue(9999999999) },
  box60Frozen: { required, minValue: minValue(0), maxValue: maxValue(9999999999) },
  box80Frozen: { required, minValue: minValue(0), maxValue: maxValue(9999999999) },
  box100Frozen: { required, minValue: minValue(0), maxValue: maxValue(9999999999) },
}

export const UpdateShippingValidationRules: ValidationArgs = {
  name: { required, maxLength: maxLength(64) },
  hasFreeShipping: {},
  box60Frozen: { required, minValue: minValue(0), maxValue: maxValue(9999999999) },
  box80Frozen: { required, minValue: minValue(0), maxValue: maxValue(9999999999) },
  box100Frozen: { required, minValue: minValue(0), maxValue: maxValue(9999999999) },
}

export const UpdateDefaultShippingValidationRules: ValidationArgs = {
  hasFreeShipping: {},
  box60Frozen: { required, minValue: minValue(0), maxValue: maxValue(9999999999) },
  box80Frozen: { required, minValue: minValue(0), maxValue: maxValue(9999999999) },
  box100Frozen: { required, minValue: minValue(0), maxValue: maxValue(9999999999) },
}
