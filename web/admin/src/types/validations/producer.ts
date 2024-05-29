import type { ValidationArgs } from '@vuelidate/core'
import { email, kana, maxLength, tel, required } from '~/lib/validations'

export const CreateProducerValidationRules: ValidationArgs = {
  lastname: { required, maxLength: maxLength(16) },
  firstname: { required, maxLength: maxLength(16) },
  lastnameKana: { required, maxLength: maxLength(32), kana },
  firstnameKana: { required, maxLength: maxLength(32), kana },
  username: { required, maxLength: maxLength(64) },
  profile: { maxLength: maxLength(2000) },
  email: { email },
  phoneNumber: { tel },
  instagramId: { maxLength: maxLength(30) },
  facebookId: { maxLength: maxLength(50) },
}

export const UpdateProducerValidationRules: ValidationArgs = {
  lastname: { required, maxLength: maxLength(16) },
  firstname: { required, maxLength: maxLength(16) },
  lastnameKana: { required, maxLength: maxLength(32), kana },
  firstnameKana: { required, maxLength: maxLength(32), kana },
  username: { required, maxLength: maxLength(64) },
  profile: { maxLength: maxLength(2000) },
  email: { email },
  phoneNumber: { tel },
  instagramId: { maxLength: maxLength(30) },
  facebookId: { maxLength: maxLength(50) },
}
