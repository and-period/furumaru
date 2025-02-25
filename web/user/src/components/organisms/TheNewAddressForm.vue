<script setup lang="ts">
import type { CreateAddressRequest } from '~/types/api'
import { prefecturesList } from '~/constants/prefectures'
import type { I18n } from '~/types/locales'

interface Props {
  formData: CreateAddressRequest
  formId: string
  postalCodeErrorMessage: string
}

const i18n = useI18n()
const props = defineProps<Props>()

interface Emits {
  (e: 'update:formData', val: CreateAddressRequest): void
  (e: 'click:searchAddressButton', postalCode: string): void
  (e: 'submit'): void
}

const emits = defineEmits<Emits>()

const at = (str: keyof I18n['purchase']['address']) => {
  return i18n.t(`purchase.address.${str}`)
}

const formDataValue = computed({
  get: () => props.formData,
  set: (val: CreateAddressRequest) => emits('update:formData', val),
})

const handleClickSearchAddressButton = () => {
  emits('click:searchAddressButton', props.formData.postalCode)
}

const handleSubmit = () => {
  emits('submit')
}
</script>

<template>
  <form
    :id="formId"
    class="flex w-full flex-col gap-4"
    @submit.prevent="handleSubmit"
  >
    <div class="grid grid-cols-2 gap-4">
      <the-text-input
        v-model="formDataValue.lastname"
        :placeholder="at('lastNamePlaceholder')"
        :with-label="false"
        type="text"
        name="lastname"
        required
      />
      <the-text-input
        v-model="formDataValue.firstname"
        :placeholder="at('firstNamePlaceholder')"
        :with-label="false"
        name="firstName"
        type="text"
        required
      />
    </div>
    <div class="grid grid-cols-2 gap-4">
      <the-text-input
        v-model="formDataValue.lastnameKana"
        :placeholder="at('lastNameKanaPlaceholder')"
        :with-label="false"
        type="text"
        required
      />
      <the-text-input
        v-model="formDataValue.firstnameKana"
        :placeholder="at('firstNameKanaPlaceholder')"
        :with-label="false"
        type="text"
        required
      />
    </div>
    <the-phone-number-input
      v-model="formDataValue.phoneNumber"
      required
    />
    <div class="flex items-center gap-4">
      <the-text-input
        v-model="formDataValue.postalCode"
        :placeholder="at('postalCodePlaceholder')"
        :with-label="false"
        type="text"
        name="postal-code"
        required
        :error-message="postalCodeErrorMessage"
      />
      <button
        type="button"
        class="whitespace-nowrap bg-main px-4 py-1 text-white"
        @click="handleClickSearchAddressButton"
      >
        {{ at("searchButtonText") }}
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
        {{ at("prefectureLabel") }}
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
      :placeholder="at('cityPlaceholder')"
      :with-label="false"
      name="address-line1"
      type="text"
      required
    />
    <the-text-input
      id="address-line2"
      v-model="formDataValue.addressLine1"
      :placeholder="at('streetPlaceholder')"
      :with-label="false"
      name="address-line2"
      type="text"
    />
    <the-text-input
      id="address-line3"
      v-model="formDataValue.addressLine2"
      :placeholder="at('apartmentPlaceholder')"
      :with-label="false"
      name="address-line3"
      type="text"
    />
    <div class="flex items-center gap-2">
      <input
        id="isDefault"
        v-model="formDataValue.isDefault"
        type="checkbox"
        class="h-4 w-4 rounded accent-main"
      >
      <label for="isDefault">{{ at("setDefaultAddressLabel") }}</label>
    </div>
  </form>
</template>
