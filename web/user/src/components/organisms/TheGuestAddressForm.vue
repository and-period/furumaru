<script setup lang="ts">
import type { GuestCheckoutAddress } from '~/types/api'
import { prefecturesList } from '~/constants/prefectures'
import type { I18n } from '~/types/locales'

interface Props {
  formData: GuestCheckoutAddress
  formId: string
  nameErrorMessage: string
  nameKanaErrorMessage: string
  postalCodeErrorMessage: string
  phoneErrorMessage: string
  cityErrorMessage: string
  addressErrorMessage: string
  emailErrorMessage: string
}

const props = defineProps<Props>()

interface Emits {
  (e: 'update:formData', val: GuestCheckoutAddress): void
  (e: 'click:searchAddressButton', postalCode: string): void
  (e: 'submit'): void
}

const emits = defineEmits<Emits>()

const i18n = useI18n()

const gt = (str: keyof I18n['purchase']['guest']) => {
  return i18n.t(`purchase.guest.${str}`)
}

const formDataValue = computed({
  get: () => props.formData,
  set: (val: GuestCheckoutAddress) => emits('update:formData', val),
})

const postalCodeErrorMessageValue = computed<string>(() => {
  // エラーメッセージがあればそれを返す
  if (props.postalCodeErrorMessage) {
    return props.postalCodeErrorMessage
  }

  // 郵便番号が未入力の場合はエラーメッセージを表示しない
  if (props.formData.postalCode === '') {
    return ''
  }

  // ハイフンが含まれている場合はエラーメッセージを表示
  if (props.formData.postalCode.includes('-')) {
    return gt('postalCodeHyphenNotAllowedErrorMessage')
  }

  // 数字以外が含まれている場合はエラーメッセージを表示
  const postalCodeRegex = /^\d+$/
  if (!postalCodeRegex.test(props.formData.postalCode)) {
    return gt('postalCodeInvalidErrorMessage')
  }

  return ''
})

const email = defineModel<string>('email', { required: true })

const handleClickSearchAddressButton = () => {
  emits('click:searchAddressButton', props.formData.postalCode)
}

const handleSubmit = () => {
  emits('submit')
}
</script>

<template>
  <component
    :is="formId ? 'form' : 'div'"
    :id="formId"
    class="flex w-full flex-col gap-4"
    @submit.prevent="handleSubmit"
  >
    <div class="grid grid-cols-2 gap-4">
      <the-text-input
        v-model="formDataValue.lastname"
        :placeholder="gt('lastNamePlaceholder')"
        :with-label="false"
        :error-message="nameErrorMessage"
        type="text"
        name="lastname"
        required
      />
      <the-text-input
        v-model="formDataValue.firstname"
        :placeholder="gt('firstNamePlaceholder')"
        :with-label="false"
        :error-message="nameErrorMessage"
        name="firstName"
        type="text"
        required
      />
    </div>
    <div class="grid grid-cols-2 gap-4">
      <the-text-input
        v-model="formDataValue.lastnameKana"
        :placeholder="gt('lastNameKanaPlaceholder')"
        :with-label="false"
        :error-message="nameKanaErrorMessage"
        type="text"
        required
      />
      <the-text-input
        v-model="formDataValue.firstnameKana"
        :placeholder="gt('firstNameKanaPlaceholder')"
        :with-label="false"
        :error-message="nameKanaErrorMessage"
        type="text"
        required
      />
    </div>
    <the-phone-number-input
      v-model="formDataValue.phoneNumber"
      :error-message="phoneErrorMessage"
      required
    />
    <the-text-input
      v-model="email"
      :placeholder="gt('emailPlaceholder')"
      :with-label="false"
      :error-message="emailErrorMessage"
      type="text"
      required
    />
    <div class="flex items-center gap-4">
      <the-text-input
        v-model="formDataValue.postalCode"
        :placeholder="gt('postalCodeLabel')"
        :with-label="false"
        :error-message="postalCodeErrorMessageValue"
        type="text"
        name="postal-code"
        required
      />
      <button
        type="button"
        class="whitespace-nowrap bg-main px-4 py-1 text-white"
        @click="handleClickSearchAddressButton"
      >
        {{ gt('searchButtonText') }}
      </button>
    </div>
    <select
      v-model="formDataValue.prefectureCode"
      :class="{
        'mb-1 block w-full appearance-none rounded-none border-b border-main bg-transparent px-1 py-2 text-inherit focus:outline-none': true,
      }"
      required
    >
      <option
        disabled
        value="0"
      >
        {{ gt('prefectureLabel') }}
      </option>
      <option
        v-for="prefecture in prefecturesList"
        :key="prefecture.id"
        :value="prefecture.value"
      >
        {{ prefecture.text }}
      </option>
    </select>
    <the-text-input
      id="address-line1"
      v-model="formDataValue.city"
      :placeholder="gt('cityPlaceholder')"
      :with-label="false"
      :error-message="cityErrorMessage"
      name="address-line1"
      type="text"
      required
    />
    <the-text-input
      id="address-line2"
      v-model="formDataValue.addressLine1"
      :placeholder="gt('streetPlaceholder')"
      :with-label="false"
      :error-message="addressErrorMessage"
      name="address-line2"
      type="text"
    />
    <the-text-input
      id="address-line3"
      v-model="formDataValue.addressLine2"
      :placeholder="gt('apartmentPlaceholder')"
      :with-label="false"
      name="address-line3"
      type="text"
    />
  </component>
</template>
