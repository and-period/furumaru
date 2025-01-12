<script setup lang="ts">
import { useAddressForm } from '~/hooks'
import { useAddressStore } from '~/store/address'
import { useAuthStore } from '~/store/auth'
import { useExperienceCheckoutStore } from '~/store/experienceCheckout'
import { useShoppingCartStore } from '~/store/shopping'
import { type GuestCheckoutExperienceRequest, PaymentMethodType } from '~/types/api'
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

const authStore = useAuthStore()
const { user } = storeToRefs(authStore)

const { fetchCheckoutTarget, checkout } = useExperienceCheckoutStore()

const addressStore = useAddressStore()
const { defaultAddress } = storeToRefs(addressStore)
const { searchAddressByPostalCode, fetchAddresses, registerAddress } = addressStore
const targetAddress = ref<'default' | 'other'>('default')

const shoppingCartStore = useShoppingCartStore()
const { availablePaymentSystem } = storeToRefs(shoppingCartStore)
const { fetchAvailablePaymentOptions } = shoppingCartStore

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
 * 住所取得処理
 */
useAsyncData('address', async () => {
  await fetchAddresses()
  // useAsyncDataの戻り値は何か返さないといけないのでtrueを返す
  return true
})

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
  email: user.value?.email || '',
  billingAddress: {
    lastname: user.value?.lastname || '',
    firstname: user.value?.firstname || '',
    lastnameKana: user.value?.lastnameKana || '',
    firstnameKana: user.value?.firstnameKana || '',
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
  postalCodeErrorMessage.value = ''
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

const submitErrorMessage = ref<string>('')
const isSubmitting = ref<boolean>(false)

/**
 * フォーム送信処理
 */
const handleSubmit = async () => {
  isSubmitting.value = true
  if (defaultAddress.value && targetAddress.value === 'default') {
    // デフォルトの住所を使用する場合
    formData.value.billingAddressId = defaultAddress.value.id
  }
  else {
    if (validate()) {
      isSubmitting.value = false
      return
    }

    // 住所登録
    try {
      const address = await registerAddress({ ...formData.value.billingAddress, isDefault: true })
      if (address) {
        formData.value.billingAddressId = address.id
      }
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

  // チェックアウト処理
  try {
    const url = await checkout(experienceId.value, formData.value)
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
  isSubmitting.value = false
}

onMounted(() => {
  formData.value.callbackUrl = `${window.location.origin}/v1/purchase/complete`
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

    <!-- エラー表示 -->
    <div
      v-if="!isValidQueryParams || targetExperienceFetchStatus == 'error'"
      class="px-4 md:px-20"
    >
      <!-- クエリパラメータが不正な場合 -->
      <the-alert
        v-if="!isValidQueryParams"
        class="bg-white"
        type="error"
      >
        不正なパラメータが含まれています
      </the-alert>

      <!-- 購入対象の体験の取得に失敗した場合 -->
      <the-alert
        v-if="targetExperienceFetchStatus == 'error'"
        class="bg-white"
        type="error"
      >
        {{ targetExperienceFetchError?.message }}
      </the-alert>
    </div>

    <template
      v-else
    >
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
            <!-- 住所選択 -->
            <template v-if="defaultAddress">
              <the-address-info :address="defaultAddress" />
              <!-- デフォルトの住所を使用するかのセレクタ -->
              <hr class="my-[20px]">
              <div class="flex flex-col gap-4 mb-6">
                <div class="flex items-center gap-2">
                  <input
                    id="default-radio"
                    v-model="targetAddress"
                    type="radio"
                    class="h-4 w-4 accent-main"
                    value="default"
                  >
                  <label for="default-radio">
                    {{ dt("useDefaultAddressLabel") }}
                  </label>
                </div>
                <div class="flex items-center gap-2">
                  <input
                    id="other-radio"
                    v-model="targetAddress"
                    type="radio"
                    class="h-4 w-4 accent-main"
                    value="other"
                  >
                  <label for="other-radio">
                    {{ dt("useOtherAddressLabel") }}
                  </label>
                </div>
              </div>

              <template v-if="targetAddress === 'other'">
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
              </template>
            </template>

            <!-- 住所入力欄 -->
            <template v-else>
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
            </template>

            <div class="flex flex-col gap-3 my-4">
              <div
                v-for="p in availablePaymentSystem"
                :key="p.methodType"
                class="flex w-full items-center justify-between"
              >
                <div class="inline-flex items-center gap-2">
                  <input
                    :id="String(p.methodType)"
                    v-model="formData.paymentMethod"
                    class="h-4 w-4 accent-main"
                    type="radio"
                    :value="p.methodType"
                  >
                  <label
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
            class="bg-main text-white py-2 w-60 disabled:cursor-wait"
            type="submit"
            form="checkout-form"
            :disabled="isSubmitting"
          >
            <template v-if="isSubmitting">
              <div class="w-full flex justify-center items-center">
                <the-loading-icon />
              </div>
            </template>
            <template v-else>
              {{ dt("submitButtonText") }}
            </template>
          </button>
        </div>
      </div>
    </template>
  </div>
</template>
