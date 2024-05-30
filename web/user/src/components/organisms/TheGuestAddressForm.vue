<script setup lang="ts">
import type { GuestCheckoutAddress } from '~/types/api'
import { prefecturesList } from '~/constants/prefectures'

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

const formDataValue = computed({
  get: () => props.formData,
  set: (val: GuestCheckoutAddress) => emits('update:formData', val),
})

const email = defineModel<string>('email')

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
        placeholder="姓"
        :with-label="false"
        :error-message="nameErrorMessage"
        type="text"
        name="lastname"
        required
      />
      <the-text-input
        v-model="formDataValue.firstname"
        placeholder="名"
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
        placeholder="ふりがな(姓)"
        :with-label="false"
        :error-message="nameKanaErrorMessage"
        type="text"
        required
      />
      <the-text-input
        v-model="formDataValue.firstnameKana"
        placeholder="ふりがな(名)"
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
      placeholder="メールアドレス"
      :with-label="false"
      :error-message="emailErrorMessage"
      type="text"
      required
    />
    <div class="flex items-center gap-4">
      <the-text-input
        v-model="formDataValue.postalCode"
        placeholder="郵便番号（ハイフンなし）"
        :with-label="false"
        :error-message="postalCodeErrorMessage"
        type="text"
        name="postal-code"
        required
      />
      <button
        type="button"
        class="whitespace-nowrap bg-main px-4 py-1 text-white"
        @click="handleClickSearchAddressButton"
      >
        検索
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
        都道府県
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
      placeholder="住所（市区町村)"
      :with-label="false"
      :error-message="cityErrorMessage"
      name="address-line1"
      type="text"
      required
    />
    <the-text-input
      id="address-line2"
      v-model="formDataValue.addressLine1"
      placeholder="住所（それ以降）"
      :with-label="false"
      :error-message="addressErrorMessage"
      name="address-line2"
      type="text"
    />
    <the-text-input
      id="address-line3"
      v-model="formDataValue.addressLine2"
      placeholder="住所（マンション名、部屋番号）"
      :with-label="false"
      name="address-line3"
      type="text"
    />
  </form>
</template>
