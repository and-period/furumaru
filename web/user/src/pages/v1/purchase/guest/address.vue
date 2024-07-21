<script setup lang="ts">
import { storeToRefs } from 'pinia'
import type { GuestCheckoutAddress } from '~/types/api'
import { useAddressStore } from '~/store/address'
import { useShoppingCartStore } from '~/store/shopping'
import { ApiBaseError } from '~/types/exception'
import type { I18n } from '~/types/locales'

const i18n = useI18n()

const route = useRoute()
const router = useRouter()

const gt = (str: keyof I18n['purchase']['guest']) => {
  return i18n.t(`purchase.guest.${str}`)
}

const addressStore = useAddressStore()
const { searchAddressByPostalCode } = addressStore
const shoppingCartStore = useShoppingCartStore()
const { calcCartResponseItem } = storeToRefs(shoppingCartStore)
const { calcCartItemByCoordinatorId, verifyGuestPromotionCode }
  = shoppingCartStore

const calcCartResponseItemState = ref<{
  isLoading: boolean
  hasError: boolean
  errorMessage: string
}>({ isLoading: true, hasError: false, errorMessage: '' })

const formData = ref<GuestCheckoutAddress>({
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
})

const formEmailData = ref({
  email: '',
})

const promotionCodeFormValue = ref<string>('')
const invalidPromotion = ref<boolean>(false)
const validPromotion = ref<boolean>(false)

const coordinatorId = computed<string>(() => {
  const id = route.query.coordinatorId
  if (id) {
    return String(id)
  }
  else {
    return ''
  }
})

const cartNumber = computed<number | undefined>(() => {
  const id = route.query.cartNumber
  const idNumber = Number(id)
  if (idNumber === 0) {
    return undefined
  }
  if (isNaN(idNumber)) {
    return undefined
  }
  return idNumber
})

const promotionCode = computed<string | undefined>(() => {
  const code = route.query.promotionCode
  if (typeof code === 'string') {
    return code
  }
  else {
    return undefined
  }
})

const priceFormatter = (price: number) => {
  return new Intl.NumberFormat('ja-JP', {
    style: 'currency',
    currency: 'JPY',
  }).format(price)
}

const handleClickSearchAddressButton = async () => {
  const res = await searchAddressByPostalCode(formData.value.postalCode)
  formData.value.prefectureCode = res.prefectureCode
  formData.value.city = res.city
  formData.value.addressLine1 = res.town
}

const handleClickBackCartButton = () => {
  router.push('/purchase')
}

const hasError = ref<boolean>(false)
const nameErrorMessage = ref<string>('')
const nameKanaErrorMessage = ref<string>('')
const phoneErrorMessage = ref<string>('')
const postalCodeErrorMessage = ref<string>('')
const cityErrorMessage = ref<string>('')
const addressErrorMessage = ref<string>('')
const emailErrorMessage = ref<string>('')

// 次のページに行く上で入力されていないと困るのでvalidationを追加
const validate = () => {
  hasError.value = false
  nameErrorMessage.value = ''
  phoneErrorMessage.value = ''
  postalCodeErrorMessage.value = ''
  cityErrorMessage.value = ''
  addressErrorMessage.value = ''
  emailErrorMessage.value = ''

  if (formData.value.firstname === '' || formData.value.lastname === '') {
    nameErrorMessage.value = '氏名を入力してください'
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
    nameKanaErrorMessage.value = '氏名(かな)を入力してください'
    hasError.value = true
  }
  else if (
    !isKana(formData.value.firstnameKana)
    || !isKana(formData.value.lastnameKana)
  ) {
    nameKanaErrorMessage.value = '氏名(かな)を入力してください'
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
    phoneErrorMessage.value = '電話番号を入力してください'
    hasError.value = true
  }

  if (formData.value.postalCode === '') {
    postalCodeErrorMessage.value = '郵便番号を入力してください'
    hasError.value = true
  }

  if (formData.value.city === '') {
    cityErrorMessage.value = '市区町村を入力してください'
    hasError.value = true
  }

  if (formData.value.addressLine1 === '') {
    addressErrorMessage.value = '住所を入力してください'
    hasError.value = true
  }

  const validateEmail = (email: string): boolean => {
    // 正規表現を使用してメールアドレスをチェックする
    const regex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
    return regex.test(email)
  }

  if (formEmailData.value.email === '') {
    emailErrorMessage.value = 'メールアドレスを入力してください'
    hasError.value = true
  }
  else if (!validateEmail(formEmailData.value.email)) {
    emailErrorMessage.value = '正しいメールアドレスを入力してください'
    hasError.value = true
  }

  return hasError.value
}

const handleClickNextStepButton = () => {
  if (validate()) {
    return
  }

  addressStore.guestAddress = formData.value
  addressStore.email = formEmailData.value.email
  router.push({
    path: '/v1/purchase/guest/confirmation',
    query: {
      coordinatorId: coordinatorId.value,
      cartNumber: cartNumber.value,
      promotionCode: promotionCode.value,
      prefectureCode: formData.value.prefectureCode,
    },
  })
}

/**
 * クーポンコードを適用するボタンをクリックしたときの処理
 */
const handleClickUsePromotionCodeButton = async () => {
  const result = await verifyGuestPromotionCode(promotionCodeFormValue.value)
  if (result) {
    invalidPromotion.value = false
    validPromotion.value = true
  }
  else {
    invalidPromotion.value = true
    validPromotion.value = false
  }
}

/**
 * 適用されているクーポンコードを取り消すボタンをクリックしたときの処理
 */
const handleClickCancelPromotionCodeButton = async () => {
  promotionCodeFormValue.value = ''
  invalidPromotion.value = false
  validPromotion.value = false
  await calcCartItemByCoordinatorId(
    coordinatorId.value,
    cartNumber.value,
    undefined, // 都道府県
    undefined,
  )
}

watch(validPromotion, (newValue, oldValue) => {
  if (newValue && promotionCodeFormValue.value) {
    // フォーム入力でvalidPromotionがtrueになったときに、クーポンコードを適用して再計算する
    calcCartItemByCoordinatorId(
      coordinatorId.value,
      cartNumber.value,
      undefined, // 都道府県
      promotionCodeFormValue.value,
    )
    // リロード対策で、クエリパラメータにクーポンコードを追加する
    router.push({
      path: '/v1/purchase/guest/address',
      query: {
        coordinatorId: coordinatorId.value,
        cartNumber: cartNumber.value,
        promotionCode: promotionCodeFormValue.value,
      },
    })
  }
  if (oldValue && !newValue) {
    calcCartItemByCoordinatorId(
      coordinatorId.value,
      cartNumber.value,
      undefined, // 都道府県
      undefined,
    )
    // クーポンコードが適用されている状態から、適用されていない状態に変わったときに、クエリパラメータからクーポンコードを削除する
    router.push({
      path: '/v1/purchase/guest/address',
      query: {
        coordinatorId: coordinatorId.value,
        cartNumber: cartNumber.value,
        promotionCode: undefined,
      },
    })
  }
})

onMounted(async () => {
  try {
    calcCartResponseItemState.value.isLoading = true
    await calcCartItemByCoordinatorId(
      coordinatorId.value,
      cartNumber.value,
      undefined, // 都道府県
      promotionCode.value,
    )
    // クーポンコードが指定されている場合は、クーポンコードを適用する
    if (promotionCode.value) {
      validPromotion.value = true
    }
  }
  catch (error) {
    calcCartResponseItemState.value.hasError = true
    if (error instanceof ApiBaseError) {
      calcCartResponseItemState.value.errorMessage = error.message
    }
    else {
      calcCartResponseItemState.value.errorMessage
        = '不明なエラーが発生しました。'
    }
  }
  finally {
    calcCartResponseItemState.value.isLoading = false
  }
})

useSeoMeta({
  title: 'ご購入手続き',
})
</script>

<template>
  <div class="container mx-auto">
    <div
      class="mt-[32px] text-center text-[20px] font-bold tracking-[2px] text-main"
    >
      {{ gt('purchaseTitle') }}
    </div>

    <the-alert
      v-if="calcCartResponseItemState.hasError"
      class="mx-auto my-4 max-w-4xl bg-white"
    >
      {{ calcCartResponseItemState.errorMessage }}
    </the-alert>

    <div
      class="relative my-10 gap-x-[80px] bg-white px-6 py-10 md:mx-0 md:grid md:grid-cols-2 md:grid-rows-[auto_auto] md:px-[80px]"
    >
      <!-- 左側 -->
      <div class="row-span-1 self-start py-[24px] md:w-full md:py-10">
        <div
          class="mb-6 text-left text-[16px] font-bold tracking-[1.6px] text-main"
        >
          {{ gt('customerInformationTitle') }}
        </div>

        <the-guest-address-form
          v-model:form-data="formData"
          v-model:email="formEmailData.email"
          form-id="new-address-form"
          :name-error-message="nameErrorMessage"
          :name-kana-error-message="nameKanaErrorMessage"
          :phone-error-message="phoneErrorMessage"
          :postal-code-error-message="postalCodeErrorMessage"
          :city-error-message="cityErrorMessage"
          :address-error-message="addressErrorMessage"
          :email-error-message="emailErrorMessage"
          @click:search-address-button="handleClickSearchAddressButton"
        />
      </div>

      <!-- 右側 -->
      <div
        class="row-span-2 self-start bg-base px-[16px] py-[24px] text-main md:w-full md:p-10"
      >
        <div class="text-[14px] font-bold tracking-[1.6px] md:text-[16px]">
          注文内容
        </div>

        <template v-if="calcCartResponseItem">
          <div class="my-[16px] text-[12px] tracking-[1.2px] md:my-6">
            <p>
              {{ calcCartResponseItem.coordinator.marcheName }}
            </p>
            <p>
              発送地：{{
                `${calcCartResponseItem.coordinator.prefecture}${calcCartResponseItem.coordinator.city}`
              }}
            </p>
            <p>
              取扱元：
              {{ calcCartResponseItem.coordinator.username }}
            </p>
            <p>箱の数：{{ calcCartResponseItem.carts.length }}</p>
          </div>
          <div>
            <div class="divide-y border-y">
              <div
                v-for="(item, i) in calcCartResponseItem.items"
                :key="i"
                class="grid grid-cols-5 py-2 text-[12px] tracking-[1.2px]"
              >
                <template v-if="item.product">
                  <nuxt-img
                    v-if="item.product.thumbnail"
                    width="56px"
                    height="56px"
                    provider="cloudFront"
                    :src="item.product.thumbnail.url"
                    :alt="`${item.product.name}の画像`"
                    class="block aspect-square h-[56px] w-[56px]"
                  />
                  <div class="col-span-3 pl-[24px] md:pl-0">
                    <div>{{ item.product.name }}</div>
                    <div
                      class="mt-4 md:mt-0 md:items-center md:justify-self-end md:text-right"
                    >
                      数量：{{ item.quantity }}
                    </div>
                  </div>

                  <div class="flex items-center justify-self-end text-right">
                    {{ priceFormatter(item.product.price) }}
                  </div>
                </template>
              </div>
            </div>

            <template v-if="validPromotion">
              <div
                class="leading-[1.4px mt-4 flex justify-between rounded-lg border border-orange p-2 text-[12px] text-orange"
              >
                <div class="flex items-center gap-1">
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke-width="1.5"
                    stroke="currentColor"
                    class="h-4 w-4"
                  >
                    <path
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      d="M9 12.75 11.25 15 15 9.75M21 12a9 9 0 1 1-18 0 9 9 0 0 1 18 0Z"
                    />
                  </svg>
                  クーポンコード適用済み
                </div>
                <button @click="handleClickCancelPromotionCodeButton">
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke-width="1.5"
                    stroke="currentColor"
                    class="h-4 w-4"
                  >
                    <path
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      d="M6 18 18 6M6 6l12 12"
                    />
                  </svg>
                </button>
              </div>
            </template>

            <template v-else>
              <div class="mt-4 flex gap-2">
                <div class="grow">
                  <input
                    v-model="promotionCodeFormValue"
                    type="text"
                    class="w-full border border-gray-300 bg-gray-50 p-2.5 text-[14px] md:text-[16px]"
                    placeholder="クーポンコード"
                  >
                </div>
                <button
                  class="whitespace-nowrap bg-main p-2 text-[14px] text-white md:text-[16px]"
                  @click="handleClickUsePromotionCodeButton"
                >
                  適用する
                </button>
              </div>
              <div
                v-if="invalidPromotion"
                class="mt-2 px-1 text-[12px] leading-[1.2px]"
              >
                指定したクーポンコードは無効です。
              </div>
            </template>

            <div
              class="mt-4 grid grid-cols-5 gap-y-4 border-y border-main py-6 text-[12px] tracking-[1.4px] md:grid-cols-2 md:text-[14px]"
            >
              <div class="col-span-2 md:col-span-1">
                商品合計（税込）
              </div>
              <div class="col-span-3 text-right md:col-span-1">
                {{ priceFormatter(calcCartResponseItem.subtotal) }}
              </div>
              <div class="col-span-2 md:col-span-1">
                クーポン利用
              </div>
              <div class="col-span-3 text-right md:col-span-1">
                {{ priceFormatter(calcCartResponseItem.discount) }}
              </div>
              <div class="col-span-2 md:col-span-1">
                送料（税込）
              </div>
              <div class="col-span-3 text-right md:col-span-1">
                次ページで計算されます
              </div>
            </div>

            <div
              class="mt-6 grid grid-cols-2 text-[14px] font-bold tracking-[1.4px]"
            >
              <div>合計（税込み）</div>
              <div class="text-right">
                {{ priceFormatter(calcCartResponseItem.total) }}
              </div>
            </div>
          </div>
        </template>
      </div>

      <div
        class="mt-[24px] flex w-full flex-col items-center gap-4 self-start md:flex-row md:justify-between"
      >
        <button
          class="order-2 inline-flex w-full gap-2 text-left text-[12px] tracking-[1.2px] text-main md:order-1 md:max-w-max"
          @click="handleClickBackCartButton"
        >
          <the-left-arrow-icon class="h-4 w-4" />
          買い物カゴへ戻る
        </button>

        <div>
          <button
            class="w-full bg-main p-[14px] text-[16px] text-white md:order-1 md:w-[240px]"
            @click="handleClickNextStepButton()"
          >
            お支払方法の選択へ
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
