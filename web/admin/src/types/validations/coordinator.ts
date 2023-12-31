import type { ValidationArgs } from '@vuelidate/core'
import { email, kana, maxLength, maxLengthArray, required, tel } from '~/lib/validations'

export const CreateCoordinatorValidationRules: ValidationArgs = {
  lastname: { required, maxLength: maxLength(16) },
  firstname: { required, maxLength: maxLength(16) },
  lastnameKana: { required, maxLength: maxLength(32), kana },
  firstnameKana: { required, maxLength: maxLength(32), kana },
  marcheName: { required, maxLength: maxLength(64) },
  username: { required, maxLength: maxLength(64) },
  profile: { maxLength: maxLength(2000) },
  instagramId: { maxLength: maxLength(30) },
  facebookId: { maxLength: maxLength(50) },
  email: { required, email },
  phoneNumber: { required, tel },
  businessDays: { maxLengthArray: maxLengthArray(7) }
}

export const UpdateCoordinatorValidationRules: ValidationArgs = {
  lastname: { required, maxLength: maxLength(16) },
  firstname: { required, maxLength: maxLength(16) },
  lastnameKana: { required, maxLength: maxLength(32), kana },
  firstnameKana: { required, maxLength: maxLength(32), kana },
  marcheName: { required, maxLength: maxLength(64) },
  username: { required, maxLength: maxLength(64) },
  profile: { maxLength: maxLength(2000) },
  instagramId: { maxLength: maxLength(30) },
  facebookId: { maxLength: maxLength(50) },
  phoneNumber: { required, tel },
  businessDays: { maxLengthArray: maxLengthArray(7) }
}
