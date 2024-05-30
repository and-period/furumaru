import type { ValidationArgs } from '@vuelidate/core'
import { helpers } from '@vuelidate/validators'
import { email, maxLength, minLength, required, sameAs } from '~/lib/validations'

export const UpdateAuthEmailValidationRules: ValidationArgs = {
  email: { required, email },
}

export const VerifyAuthEmailValidationRules: ValidationArgs = {
  verifyCode: {
    required,
    minLength: helpers.withMessage('検証コードは6文字で入力してください。', minLength(6)),
    maxLength: helpers.withMessage('検証コードは6文字で入力してください。', maxLength(6)),
  },
}

export const UpdateAuthPasswordValidationRules = (newPassword: string): ValidationArgs => ({
  oldPassword: { required },
  newPassword: { required, minLength: minLength(8), maxLength: maxLength(32) },
  passwordConfirmation: { required, sameAs: sameAs(newPassword) },
})

export const ResetAuthPasswordValidationRules = (password: string): ValidationArgs => ({
  verifyCode: {
    required,
    minLength: helpers.withMessage('検証コードは6文字で入力してください。', minLength(6)),
    maxLength: helpers.withMessage('検証コードは6文字で入力してください。', maxLength(6)),
  },
  password: { required, minLength: minLength(8), maxLength: maxLength(32) },
  passwordConfirmation: { required, sameAs: sameAs(password) },
})
