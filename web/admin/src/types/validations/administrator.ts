import type { ValidationArgs } from '@vuelidate/core'
import { email, kana, maxLength, required, tel } from '~/lib/validations'

export const CreateAdministratorValidationRules: ValidationArgs = {
  lastname: { required, maxLength: maxLength(16) },
  firstname: { required, maxLength: maxLength(16) },
  lastnameKana: { required, maxLength: maxLength(32), kana },
  firstnameKana: { required, maxLength: maxLength(32), kana },
  email: { required, email },
  phoneNumber: { required, tel },
}

export const UpdateAdministratorValidationRules: ValidationArgs = {
  lastname: { required, maxLength: maxLength(16) },
  firstname: { required, maxLength: maxLength(16) },
  lastnameKana: { required, maxLength: maxLength(32), kana },
  firstnameKana: { required, maxLength: maxLength(32), kana },
  phoneNumber: { required, tel },
}
