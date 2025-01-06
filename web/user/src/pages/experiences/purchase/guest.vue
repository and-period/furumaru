<script setup lang="ts">
import { useExperienceStore } from '~/store/experience'
import { useShoppingCartStore } from '~/store/shopping'
import { ProductResponseFromJSON, type GuestCheckoutExperienceRequest, PaymentMethodType } from '~/types/api'
import type { I18n } from '~/types/locales'

const i18n = useI18n()

const dt = (str: keyof I18n['items']['experiences']) => {
  return i18n.t(`items.experiences.${str}`)
}

const route = useRoute()

const experienceStore = useExperienceStore()
const { fetchExperience } = experienceStore

const shoppingCartStore = useShoppingCartStore()
const { calcCartResponseItem, availablePaymentSystem }
  = storeToRefs(shoppingCartStore)
const { fetchAvailablePaymentOptions } = shoppingCartStore

const experienceId = computed<string>(() => {
  const ids = route.query.id
  if (Array.isArray(ids)) {
    return ids[0]
  }
  else if (typeof ids === 'string') {
    return ids
  }
  return ''
})

const formData = ref<GuestCheckoutExperienceRequest>({
  requestId: '',
  billingAddressId: '',
  promotionCode: '',
  adultCount: route.query.adultCount ? Number(route.query.adultCount) : 0,
  juniorHighSchoolCount: route.query.juniorHighSchoolCount ? Number(route.query.juniorHighSchoolCount) : 0,
  elementarySchoolCount: route.query.elementarySchoolCount ? Number(route.query.elementarySchoolCount) : 0,
  preschoolCount: route.query.preschoolCount ? Number(route.query.preschoolCount) : 0,
  seniorCount: route.query.seniorCount ? Number(route.query.seniorCount) : 0,
  transportation: '',
  requestedDate: '',
  requestedTime: '',
  paymentMethod: 0,
  callbackUrl: '',
  total: 0,
  email: '',
  billingAddress: {
    lastname: '',
    firstname: '',
    lastnameKana: '',
    firstnameKana: '',
    postalCode: '',
    prefectureCode: 0,
    city: '',
    addressLine1: '',
    addressLine2: '',
    phoneNumber: '',
  },
  creditCard: {
    name: '',
    number: '',
    month: 0,
    year: 0,
    verificationValue: '',
  },
})

const { data, status } = useAsyncData('experience', async () => {
  if (experienceId.value) {
    return await fetchExperience(experienceId.value)
  }
})

const postalCodeErrorMessage = computed(() => {
  if (formData.value.billingAddress.postalCode.includes('-')) {
    return 'ハイフンを含めずに入力してください'
  }
  else {
    return ''
  }
})

fetchAvailablePaymentOptions().then(() => {
  if (availablePaymentSystem.value.length > 0) {
    formData.value.paymentMethod
        = availablePaymentSystem.value[0].methodType
  }
})
</script>

<template>
  <div class="container mx-auto">
    <!-- タイトル -->
    <div
      class="my-[32px] text-center text-[20px] font-bold tracking-[2px] text-main"
    >
      ご購入手続き
    </div>
    <div class="bg-white py-10 md:px-20 flex flex-col gap-8 px-4">
      <div class="md:grid grid-cols-2 gap-x-20 auto-rows-auto flex flex-col gap-y-4">
        <!-- 顧客情報入力フォーム -->
        <form
          id="checkout-form"
          class="order-2 md:order-1"
        >
          <the-guest-address-form
            v-model:form-data="formData.billingAddress"
            form-id=""
            :postal-code-error-message="postalCodeErrorMessage"
          />
          <div class="flex flex-col gap-3 my-4">
            <div
              v-for="p in availablePaymentSystem"
              :key="p.methodType"
              class="flex w-full items-center justify-between"
            >
              <div class="inline-flex items-center">
                <input
                  :id="String(p.methodType)"
                  v-model="formData.paymentMethod"
                  class="check:before:border-main relative float-left mr-1 mt-0.5 h-5 w-5 appearance-none rounded-full border-2 border-solid border-neutral-300 before:pointer-events-none before:absolute before:h-4 before:w-4 before:scale-0 before:rounded-full before:bg-transparent before:opacity-0 before:shadow-[0px_0px_0px_13px_transparent] before:content-[''] after:absolute after:z-[1] after:block after:h-4 after:w-4 after:rounded-full after:content-[''] checked:border-main checked:before:opacity-[0.16] checked:after:absolute checked:after:left-1/2 checked:after:top-1/2 checked:after:h-[0.625rem] checked:after:w-[0.625rem] checked:after:rounded-full checked:after:bg-main checked:after:content-[''] checked:after:[transform:translate(-50%,-50%)] hover:cursor-pointer hover:before:opacity-[0.04] hover:before:shadow-[0px_0px_0px_13px_rgba(0,0,0,0.6)] focus:shadow-none focus:outline-none focus:ring-0 focus:before:scale-100 focus:before:opacity-[0.12] focus:before:shadow-[0px_0px_0px_13px_rgba(0,0,0,0.6)] focus:before:transition-[box-shadow_0.2s,transform_0.2s] checked:focus:border-main checked:focus:before:scale-100 checked:focus:before:shadow-[0px_0px_0px_13px_#3b71ca] checked:focus:before:transition-[box-shadow_0.2s,transform_0.2s] dark:border-neutral-600 dark:focus:before:shadow-[0px_0px_0px_13px_rgba(255,255,255,0.4)] dark:checked:focus:before:shadow-[0px_0px_0px_13px_#3b71ca]"
                  type="radio"
                  :value="p.methodType"
                >
                <label
                  class="pl-2 text-[14px] text-main"
                  :for="String(p.methodType)"
                >
                  {{ p.methodName }}
                </label>
              </div>
            </div>
          </div>

          <the-payment-form
            v-if="formData.paymentMethod === PaymentMethodType.CREDIT_CARD"
            v-model="formData.creditCard"
            form-id=""
          />
        </form>

        <!-- 購入内容確認 -->
        <template v-if="data?.experience">
          <the-experience-summary
            class=" max-h-max order-1 md:order-2"
            :experience="data.experience"
            :adult-count="formData.adultCount"
            :junior-high-school-count="formData.juniorHighSchoolCount"
            :elementary-school-count="formData.elementarySchoolCount"
            :preschool-count="formData.preschoolCount"
            :senior-count="formData.seniorCount"
          />
        </template>
      </div>

      <div class="text-center">
        <button
          class="bg-main text-white py-2 w-60"
          type="submit"
          form="checkout-form"
        >
          購入する
        </button>
      </div>
    </div>
  </div>
</template>
