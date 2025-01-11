import type { Ref } from 'vue'
import type { GuestCheckoutAddress } from '~/types/api'
import type { I18n } from '~/types/locales'

export function useAddressForm(formData: Ref<GuestCheckoutAddress>, email: Ref<string>) {
  const i18n = useI18n()

  const gt = (str: keyof I18n['purchase']['guest']) => {
    return i18n.t(`purchase.guest.${str}`)
  }

  const hasError = ref<boolean>(false)
  const nameErrorMessage = ref<string>('')
  const nameKanaErrorMessage = ref<string>('')
  const postalCodeErrorMessage = ref<string>('')
  const phoneErrorMessage = ref<string>('')
  const cityErrorMessage = ref<string>('')
  const addressErrorMessage = ref<string>('')
  const emailErrorMessage = ref<string>('')

  // バリデーションを実施する関数
  const validate = () => {
    hasError.value = false
    nameErrorMessage.value = ''
    phoneErrorMessage.value = ''
    postalCodeErrorMessage.value = ''
    cityErrorMessage.value = ''
    addressErrorMessage.value = ''
    emailErrorMessage.value = ''

    if (
      formData.value.firstname === ''
      || formData.value.lastname === ''
    ) {
      nameErrorMessage.value = gt('nameErrorMessage')
      hasError.value = true
    }
    else {
      nameErrorMessage.value = ''
    }

    const isKana = (input: string): boolean => {
    // ひらがなの正規表現
      const kanaRegex = /^[\u3040-\u309F]+$/
      return kanaRegex.test(input)
    }

    if (
      formData.value.firstnameKana === ''
      || formData.value.lastnameKana === ''
    ) {
      nameKanaErrorMessage.value = gt('nameKanaErrorMessage')
      hasError.value = true
    }
    else if (
      !isKana(formData.value.firstnameKana)
      || !isKana(formData.value.lastnameKana)
    ) {
      nameKanaErrorMessage.value = gt('nameKanaErrorMessage')
      hasError.value = true
    }
    else {
      nameKanaErrorMessage.value = ''
    }

    const isValidJapanesePhoneNumber = (phoneNumber: string): boolean => {
      const regex = /^0\d{1,4}-\d{1,4}-\d{3,4}$/
      return regex.test(phoneNumber)
    }

    if (
      formData.value.phoneNumber === ''
      || !isValidJapanesePhoneNumber(formData.value.phoneNumber)
    ) {
      phoneErrorMessage.value = gt('phoneErrorMessage')
      hasError.value = true
    }

    if (formData.value.postalCode === '') {
      postalCodeErrorMessage.value = gt('postalCodeErrorMessage')
      hasError.value = true
    }

    if (formData.value.city === '') {
      cityErrorMessage.value = gt('cityErrorMessage')
      hasError.value = true
    }

    if (formData.value.addressLine1 === '') {
      addressErrorMessage.value = gt('addressErrorMessage')
      hasError.value = true
    }

    const validateEmail = (email: string): boolean => {
    // 正規表現を使用してメールアドレスをチェックする
      const regex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
      return regex.test(email)
    }

    if (email.value === '') {
      emailErrorMessage.value = gt('emailErrorMessage')
      hasError.value = true
    }
    else if (!validateEmail(email.value)) {
      emailErrorMessage.value = gt('emailInvalidErrorMessage')
      hasError.value = true
    }

    return hasError.value
  }

  return {
    hasError,
    nameErrorMessage,
    nameKanaErrorMessage,
    postalCodeErrorMessage,
    phoneErrorMessage,
    cityErrorMessage,
    addressErrorMessage,
    emailErrorMessage,
    validate,
  }
}
