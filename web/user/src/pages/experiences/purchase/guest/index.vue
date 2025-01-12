<script setup lang="ts">
import { useAddressForm } from '~/hooks'
import { useAddressStore } from '~/store/address'
import { useExperienceCheckoutStore } from '~/store/experienceCheckout'
import { useShoppingCartStore } from '~/store/shopping'
import {
  type GuestCheckoutExperienceRequest,
  PaymentMethodType,
} from '~/types/api'
import { ApiBaseError } from '~/types/exception'
import type { I18n } from '~/types/locales'

const i18n = useI18n()

const dt = (str: keyof I18n['experiences']['purchase']) => {
  return i18n.t(`experiences.purchase.${str}`)
}

const gt = (str: keyof I18n['purchase']['guest']) => {
  return i18n.t(`purchase.guest.${str}`)
}

const route = useRoute()

const shoppingCartStore = useShoppingCartStore()
const { availablePaymentSystem } = storeToRefs(shoppingCartStore)
const { fetchAvailablePaymentOptions } = shoppingCartStore

const addressStore = useAddressStore()
const { searchAddressByPostalCode } = addressStore

const { fetchCheckoutTarget, checkoutByGuest } = useExperienceCheckoutStore()

/**
 * 体験ID（クエリパラメータから算出）
 */
const experienceId = computed<string>(() => {
  const id = route.query.id
  if (id) {
    return String(id)
  }
  return ''
})

/**
 * 大人の人数（クエリパラメータから算出）
 */
const adultCount = computed<number>(() => {
  return route.query.adultCount ? Number(route.query.adultCount) : 0
})

/**
 * 中学生の人数（クエリパラメータから算出）
 */
const juniorHighSchoolCount = computed<number>(() => {
  return route.query.juniorHighSchoolCount ? Number(route.query.juniorHighSchoolCount) : 0
})

/**
 * 小学生の人数（クエリパラメータから算出）
 */
const elementarySchoolCount = computed<number>(() => {
  return route.query.elementarySchoolCount ? Number(route.query.elementarySchoolCount) : 0
})

/**
 * 幼児の人数（クエリパラメータから算出）
 */
const preschoolCount = computed<number>(() => {
  return route.query.preschoolCount ? Number(route.query.preschoolCount) : 0
})

/**
 * シニアの人数（クエリパラメータから算出）
 */
const seniorCount = computed<number>(() => {
  return route.query.seniorCount ? Number(route.query.seniorCount) : 0
})

/**
 * 体験情報取得処理
 */
const {
  data: targetExperience,
  status: targetExperienceFetchStatus,
  error: targetExperienceFetchError,
} = useAsyncData('target-experience', async () => {
  if (experienceId.value) {
    return await fetchCheckoutTarget({
      experienceId: experienceId.value,
      adult: adultCount.value,
      juniorHighSchool: juniorHighSchoolCount.value,
      elementarySchool: elementarySchoolCount.value,
      preschool: preschoolCount.value,
      senior: seniorCount.value,
    })
  }
})

/**
 * フォームデータ
 */
const formData = ref<GuestCheckoutExperienceRequest>({
  requestId: targetExperience.value?.requestId || '',
  billingAddressId: '',
  promotionCode: '',
  adultCount: adultCount.value,
  juniorHighSchoolCount: juniorHighSchoolCount.value,
  elementarySchoolCount: elementarySchoolCount.value,
  preschoolCount: preschoolCount.value,
  seniorCount: seniorCount.value,
  transportation: '',
  requestedDate: '',
  requestedTime: '',
  paymentMethod: 0,
  callbackUrl: '',
  total: targetExperience.value?.total || 0,
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

const isValidQueryParams = computed<boolean>(() => {
  // 必須パラメータが揃っているか
  // 体験IDがあるか
  if (!experienceId.value) {
    return false
  }

  // 大人などの人数が全て0の場合はエラー
  if (
    adultCount.value === 0
    && juniorHighSchoolCount.value === 0
    && elementarySchoolCount.value === 0
    && preschoolCount.value === 0
    && seniorCount.value === 0
  ) {
    return false
  }

  // 大人などの人数で10人以上の場合はエラー
  if (
    adultCount.value > 10
    || juniorHighSchoolCount.value > 10
    || elementarySchoolCount.value > 10
    || preschoolCount.value > 10
    || seniorCount.value > 10
  ) {
    return false
  }

  // 大人などの人数で0人以下のパラメータがある場合はエラー
  if (
    adultCount.value < 0
    || juniorHighSchoolCount.value < 0
    || elementarySchoolCount.value < 0
    || preschoolCount.value < 0
    || seniorCount.value < 0
  ) {
    return false
  }

  return true
})

const addressFormData = computed(() => formData.value.billingAddress)
const emailFormData = computed(() => formData.value.email)

const {
  // hasError,
  nameErrorMessage,
  nameKanaErrorMessage,
  postalCodeErrorMessage,
  phoneErrorMessage,
  cityErrorMessage,
  addressErrorMessage,
  emailErrorMessage,
  validate,
} = useAddressForm(addressFormData, emailFormData)

/**
 * 住所検索ボタンクリック時の処理
 */
const handleClickSearchAddressButton = async () => {
  try {
    const res = await searchAddressByPostalCode(formData.value.billingAddress.postalCode)
    formData.value.billingAddress.prefectureCode = res.prefectureCode
    formData.value.billingAddress.city = res.city
    formData.value.billingAddress.addressLine1 = res.town
  }
  catch (_) {
    postalCodeErrorMessage.value = gt('addressNotFoundErrorMessage')
  }
}

/**
 * 支払い方法取得処理
 */
useAsyncData('payment-options', async () => {
  await fetchAvailablePaymentOptions()
  if (availablePaymentSystem.value.length > 0) {
    formData.value.paymentMethod = availablePaymentSystem.value[0].methodType
  }
  // useAsyncDataの戻り値は何か返さないといけないのでtrueを返す
  return true
})

const submitErrorMessage = ref<string>('')

/**
 * 体験購入フォーム送信時の処理
 */
const handleSubmit = async () => {
  if (validate()) {
    return
  }

  try {
    const url = await checkoutByGuest(experienceId.value, formData.value)
    window.location.href = url
  }
  catch (e) {
    if (e instanceof ApiBaseError) {
      submitErrorMessage.value = e.message
    }
    else {
      submitErrorMessage.value = '不明なエラーが発生しました。'
    }
  }
}

onMounted(() => {
  formData.value.callbackUrl = `${window.location.origin}/v1/purchase/guest/complete`
  if (availablePaymentSystem.value.length > 0) {
    formData.value.paymentMethod = availablePaymentSystem.value[0].methodType
  }
})

useSeoMeta(
  {
    title: gt('seoTitle'),
  },
)
</script>

<template>
  <div class="container mx-auto">
    <!-- タイトル -->
    <div
      class="my-[32px] text-center text-[20px] font-bold tracking-[2px] text-main"
    >
      {{ dt("title") }}
    </div>

    <!-- エラー -->
    <div
      v-if="!isValidQueryParams || targetExperienceFetchStatus == 'error'"
      class="px-4 md:px-20"
    >
      <the-alert
        v-if="!isValidQueryParams"
        class="bg-white"
        type="error"
      >
        不正なパラメータが含まれています
      </the-alert>

      <the-alert
        v-if="targetExperienceFetchStatus == 'error'"
        class="bg-white"
        type="error"
      >
        {{ targetExperienceFetchError?.message }}
      </the-alert>
    </div>

    <template v-else>
      <!-- フォーム送信エラー -->
      <div
        v-if="submitErrorMessage"
        class="px-4 md:px-20 mb-6"
      >
        <the-alert
          v-if="submitErrorMessage"
          class="bg-white"
          type="error"
        >
          {{ submitErrorMessage }}
        </the-alert>
      </div>

      <div
        class="bg-white py-10 md:px-20 flex flex-col gap-8 px-4"
      >
        <div
          class="lg:grid grid-cols-2 gap-x-20 auto-rows-auto flex flex-col gap-y-4"
        >
          <!-- 顧客情報入力フォーム -->
          <form
            id="checkout-form"
            class="order-2 lg:order-1"
            @submit.prevent="handleSubmit"
          >
            <div
              class="mb-6 text-left text-[16px] font-bold tracking-[1.6px] text-main"
            >
              {{ dt("customerInformationTitle") }}
            </div>
            <the-guest-address-form
              v-model:form-data="formData.billingAddress"
              v-model:email="formData.email"
              form-id=""
              :name-error-message="nameErrorMessage"
              :name-kana-error-message="nameKanaErrorMessage"
              :postal-code-error-message="postalCodeErrorMessage"
              :phone-error-message="phoneErrorMessage"
              :city-error-message="cityErrorMessage"
              :address-error-message="addressErrorMessage"
              :email-error-message="emailErrorMessage"
              @click:search-address-button="handleClickSearchAddressButton"
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
          <template v-if="targetExperience?.experience">
            <the-experience-summary
              class="max-h-max order-1 lg:order-2"
              :experience="targetExperience.experience"
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
            {{ dt("submitButtonText") }}
          </button>
        </div>
      </div>
    </template>
  </div>
</template>
