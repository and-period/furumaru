<script setup lang="ts">
import { useExperienceCheckoutStore } from '~/store/experienceCheckout'
import type { GuestCheckoutExperienceRequest } from '~/types/api'
import type { I18n } from '~/types/locales'

const i18n = useI18n()

const dt = (str: keyof I18n['experiences']['purchase']) => {
  return i18n.t(`experiences.purchase.${str}`)
}

const route = useRoute()

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
const { data, status, error } = useAsyncData('target-experience', async () => {
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

const formData = ref<GuestCheckoutExperienceRequest>({
  requestId: '',
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
  total: data.value?.total || 0,
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
      v-if="!isValidQueryParams || status == 'error'"
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
        v-if="status == 'error'"
        class="bg-white"
        type="error"
      >
        {{ error?.message }}
      </the-alert>
    </div>

    <div
      v-else
      class="bg-white py-10 md:px-20 flex flex-col gap-8 px-4"
    >
      {{ data }}

      <!-- 購入内容確認 -->
      <template v-if="data?.experience">
        <the-experience-summary
          class="max-h-max order-1 md:order-2"
          :experience="data.experience"
          :adult-count="formData.adultCount"
          :junior-high-school-count="formData.juniorHighSchoolCount"
          :elementary-school-count="formData.elementarySchoolCount"
          :preschool-count="formData.preschoolCount"
          :senior-count="formData.seniorCount"
        />
      </template>
    </div>
  </div>
</template>
