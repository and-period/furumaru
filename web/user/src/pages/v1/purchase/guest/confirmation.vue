<script setup lang="ts">
import { storeToRefs } from 'pinia'
import dayjs from 'dayjs'
import { useAddressStore } from '~/store/address'
import { useProductCheckoutStore } from '~/store/productCheckout'
import { useShoppingCartStore } from '~/store/shopping'
import type { GuestCheckoutProductRequest } from '~/types/api'
import type { I18n } from '~/types/locales'
import { ApiBaseError } from '~/types/exception'

const i18n = useI18n()
const addressStore = useAddressStore()
const { addressFetchState } = storeToRefs(addressStore)

const shoppingCartStore = useShoppingCartStore()
const { calcCartResponseItem, availablePaymentSystem }
  = storeToRefs(shoppingCartStore)
const {
  calcCartItemByCoordinatorId,
  fetchAvailablePaymentOptions,
  verifyPromotionCode,
} = shoppingCartStore

const checkoutStore = useProductCheckoutStore()
const { guestCheckout } = checkoutStore

const route = useRoute()
const router = useRouter()

const ct = (str: keyof I18n['purchase']['confirmation']) => {
  return i18n.t(`purchase.confirmation.${str}`)
}

/**
 * コーディネーターID（クエリパラメータから算出）
 */
const coordinatorId = computed<string>(() => {
  const id = route.query.coordinatorId
  if (id) {
    return String(id)
  }
  else {
    return ''
  }
})

/**
 * 都道府県コード（クエリパラメータから算出）
 */
const prefectureCode = computed<number | null>(() => {
  const code = route.query.prefectureCode
  if (code) {
    return code
  }
  else {
    return null
  }
})

/**
 * プロモーションコード（クエリパラメータから算出）
 */
const promotionCode = computed<string | undefined>(() => {
  const code = route.query.promotionCode
  if (typeof code === 'string') {
    return code
  }
  else {
    return undefined
  }
})

/**
 * カート番号（クエリパラメータから算出）
 */
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

const { email, guestAddress } = storeToRefs(addressStore)

const checkoutFormData = ref<GuestCheckoutProductRequest>({
  requestId: '',
  coordinatorId: '',
  boxNumber: 0,
  promotionCode: '',
  paymentMethod: 0,
  callbackUrl: '',
  total: 0,
  email: '',
  isSameAddress: true,
  shippingAddress: {
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

const creditCardMonthValue = computed({
  get: () => {
    if (checkoutFormData.value.creditCard.month === 0) {
      return '0'
    }
    if (checkoutFormData.value.creditCard.month < 10) {
      return `0${checkoutFormData.value.creditCard.month}`
    }
    else {
      return String(checkoutFormData.value.creditCard.month)
    }
  },
  set: (val: string) => {
    const month = Number(val)
    if (!isNaN(month)) {
      checkoutFormData.value.creditCard.month = month
    }
  },
})

const checkoutError = ref<string>('')
const validPromotionCode = ref<boolean>(false)

const priceFormatter = (price: number) => {
  return new Intl.NumberFormat('ja-JP', {
    style: 'currency',
    currency: 'JPY',
  }).format(price)
}

const handleClickPreviousStepButton = () => {
  router.back()
}

/**
 * チェックアウト処理を実行するメソッド
 */
const doCheckout = async () => {
  try {
    const url = await guestCheckout({
      ...checkoutFormData.value,
      boxNumber: cartNumber.value ?? 0,
    })
    console.log('debug', 'doCheckout', url)
    addressStore.$reset()
    window.location.href = url
  }
  catch (error) {
    if (error instanceof ApiBaseError) {
      checkoutError.value = error.message
    }
  }
}

const handleClickNextStepButton = () => {
  if (checkoutFormData.value.paymentMethod === 2) {
    // クレジットカードなら何もしない
    return
  }
  doCheckout()
}

const handleSubmitCreditCardForm = () => {
  doCheckout()
}

onMounted(async () => {
  // 利用可能な支払い方法を取得
  fetchAvailablePaymentOptions().then(() => {
    if (availablePaymentSystem.value.length > 0) {
      checkoutFormData.value.paymentMethod
        = availablePaymentSystem.value[0].methodType
    }
  })

  // プロモーションコードの有効性を確認
  if (promotionCode.value) {
    const result = await verifyPromotionCode(promotionCode.value, coordinatorId.value)
    if (result) {
      validPromotionCode.value = result
      checkoutFormData.value.promotionCode = promotionCode.value
    }
    else {
      console.log('無効なプロモーションコードが設定されています。')
    }
  }

  if (prefectureCode.value !== null) {
    await calcCartItemByCoordinatorId(
      coordinatorId.value,
      cartNumber.value,
      prefectureCode.value,
      validPromotionCode.value ? promotionCode.value : undefined,
    )
  }

  checkoutFormData.value.requestId
    = calcCartResponseItem.value?.requestId ?? ''
  checkoutFormData.value.coordinatorId = coordinatorId.value
  checkoutFormData.value.total = calcCartResponseItem.value?.total ?? 0
  checkoutFormData.value.callbackUrl = `${window.location.origin}/v1/purchase/guest/complete`

  if (email.value !== null) {
    checkoutFormData.value.email = email.value
  }
  else {
    checkoutFormData.value.email = ''
  }

  if (guestAddress.value !== undefined) {
    checkoutFormData.value.shippingAddress = guestAddress.value
    checkoutFormData.value.billingAddress = guestAddress.value
  }
  else {
    checkoutFormData.value.email = ''
  }

  if (calcCartResponseItem.value?.promotion) {
    checkoutFormData.value.promotionCode
      = calcCartResponseItem.value.promotion.code
  }
})

useSeoMeta({
  title: 'ご購入手続き',
})
</script>

<template>
  <div class="container mx-auto">
    <div class="text-center text-[20px] font-bold tracking-[2px] text-main">
      {{ ct("checkoutTitle") }}
    </div>

    <the-alert
      v-if="checkoutError"
      class="mt-4 bg-white"
      type="error"
    >
      {{ checkoutError }}
    </the-alert>

    <the-alert
      v-if="prefectureCode === null"
      class="mt-4 bg-white"
      type="error"
    >
      {{ "都道府県が指定されていません。住所を再度入力してください。" }}
    </the-alert>

    <div
      class="relative my-10 gap-x-[80px] bg-white px-6 py-10 md:mx-0 md:grid md:grid-cols-2 md:grid-rows-[auto_auto] md:px-[80px]"
    >
      <template v-if="addressFetchState.isLoading">
        <div class="absolute h-2 w-full animate-pulse bg-main" />
      </template>

      <template v-else>
        <!-- 左側 -->
        <div class="row-span-1 self-start py-[24px] md:w-full md:py-10">
          <div
            class="mb-6 text-left text-[16px] font-bold tracking-[1.6px] text-main"
          >
            {{ ct("customerInformationTitle") }}
          </div>
          <the-guest-address-info
            v-if="guestAddress"
            :address="guestAddress"
          />

          <div class="items-center border-b py-2" />

          <div>
            <div
              class="pt-6 text-left text-[16px] font-bold tracking-[1.6px] text-main"
            >
              {{ ct("shippingInformationLabel") }}
            </div>
            <div class="pt-[27px] text-[14px] tracking-[1.4px] text-main">
              {{ ct("shippingAvobeAdderssLabel") }}
            </div>
          </div>

          <div class="items-center border-b py-2" />

          <div>
            <div
              class="pt-6 text-left text-[16px] font-bold tracking-[1.6px] text-main"
            >
              {{ ct("paymentInformationTitle") }}
            </div>

            <div class="mt-4 flex items-center justify-between">
              <div class="flex w-full flex-col items-center gap-4">
                <div
                  v-for="p in availablePaymentSystem"
                  :key="p.methodType"
                  class="flex w-full items-center justify-between"
                >
                  <div>
                    <input
                      :id="String(p.methodType)"
                      v-model="checkoutFormData.paymentMethod"
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

                  <template v-if="p.methodType === 2">
                    <div
                      class="flex h-[18px] min-w-max items-center gap-2 md:hidden"
                    >
                      <img
                        src="~/assets/img/cc/visa.png"
                        alt="visa icon"
                        class="h-full"
                      >
                      <img
                        src="~/assets/img/cc/jcb.png"
                        alt="jcb icon"
                        class="h-full"
                      >
                      <img
                        src="~/assets/img/cc/amex.png"
                        alt="amex icon"
                        class="h-full"
                      >
                      <img
                        src="~/assets/img/cc/master.png"
                        alt="master icon"
                        class="h-full"
                      >
                    </div>
                  </template>
                </div>
              </div>
            </div>

            <form
              v-if="checkoutFormData.paymentMethod === 2"
              id="credit-card-form"
              @submit.prevent="handleSubmitCreditCardForm"
            >
              <div class="mt-4 flex w-full items-center gap-4">
                <the-text-input
                  v-model="checkoutFormData.creditCard.number"
                  :placeholder="ct('creditCardNumberPlaceholder')"
                  :with-label="false"
                  name="cc-number"
                  type="text"
                  class="w-full"
                  pattern="[0-9]*"
                  required
                />
                <div
                  class="hidden h-[24px] min-w-max items-center gap-2 md:flex"
                >
                  <img
                    src="~/assets/img/cc/visa.png"
                    alt="visa icon"
                    class="h-full"
                  >
                  <img
                    src="~/assets/img/cc/jcb.png"
                    alt="jcb icon"
                    class="h-full"
                  >
                  <img
                    src="~/assets/img/cc/amex.png"
                    alt="amex icon"
                    class="h-full"
                  >
                  <img
                    src="~/assets/img/cc/master.png"
                    alt="master icon"
                    class="h-full"
                  >
                </div>
              </div>
              <the-text-input
                v-model="checkoutFormData.creditCard.name"
                :placeholder="ct('cardholderNamePlaceholder')"
                :with-label="false"
                name="cc-name"
                type="text"
                class="mt-2 w-full"
                required
              />
              <div class="mt-2 flex gap-4">
                <select
                  v-model="creditCardMonthValue"
                  class="mb-1 block w-full appearance-none rounded-none border-b border-main bg-transparent px-1 py-2 text-inherit focus:outline-none"
                >
                  <option
                    :value="0"
                    disabled
                  >
                    {{ ct("expirationMonthPlaceholder") }}
                  </option>
                  <option
                    v-for="i in 12"
                    :key="i"
                    :value="i < 10 ? `0${i}` : String(i)"
                  >
                    {{ i }}
                  </option>
                </select>

                <select
                  v-model="checkoutFormData.creditCard.year"
                  class="mb-1 block w-full appearance-none rounded-none border-b border-main bg-transparent px-1 py-2 text-inherit focus:outline-none"
                >
                  <option
                    value="0"
                    disabled
                  >
                    {{ ct("expirationYearPlaceholder") }}
                  </option>
                  <option
                    v-for="i in 11"
                    :key="i"
                    :value="dayjs().year() + i - 1"
                  >
                    {{ dayjs().year() + i - 1 }}
                  </option>
                </select>
              </div>
              <the-text-input
                v-model="checkoutFormData.creditCard.verificationValue"
                :placeholder="ct('securityCodePlaceholder')"
                :with-label="false"
                name="cc-csc"
                type="password"
                pattern="[0-9]*"
                class="mt-4 w-1/2"
                required
              />
            </form>
          </div>
        </div>

        <!-- 右側 -->
        <div
          class="row-span-2 self-start bg-base px-[16px] py-[24px] text-main md:w-full md:p-10"
        >
          <div class="text-[14px] font-bold tracking-[1.6px] md:text-[16px]">
            {{ ct("orderDetailsTitle") }}
          </div>
          <template v-if="calcCartResponseItem">
            <div class="my-[16px] text-[12px] tracking-[1.2px] md:my-6">
              <p>
                {{ calcCartResponseItem.coordinator.marcheName }}
              </p>
              <p>
                {{ ct("shipFromLabel") }}
                {{
                  `${calcCartResponseItem.coordinator.prefecture}${calcCartResponseItem.coordinator.city}`
                }}
              </p>
              <p>
                {{ ct("coordinatorLabel") }}
                {{ calcCartResponseItem.coordinator.username }}
              </p>
              <p>
                {{ ct("boxCountLabel") }}{{ calcCartResponseItem.carts.length }}
              </p>
            </div>
            <div>
              <div>
                <div
                  v-for="(item, i) in calcCartResponseItem.items"
                  :key="i"
                  class="grid grid-cols-5 border-t py-2 text-[12px] tracking-[1.2px]"
                >
                  <template v-if="item.product">
                    <template
                      v-if="item.product.thumbnail.url.endsWith('.mp4')"
                    >
                      <video
                        width="56px"
                        height="56px"
                        :src="item.product.thumbnail.url"
                        class="block aspect-square h-[56px] w-[56px]"
                      />
                    </template>
                    <template v-else>
                      <img
                        :src="item.product.thumbnailUrl"
                        :alt="`${item.product.name}の画像`"
                        class="block aspect-square h-[56px] w-[56px]"
                      >
                    </template>
                    <div class="col-span-3 pl-[24px] md:pl-0">
                      <div>{{ item.product?.name }}</div>
                      <div
                        class="mt-4 md:mt-0 md:items-center md:justify-self-end md:text-right"
                      >
                        {{ ct("quantityLabel") }}{{ item.quantity }}
                      </div>
                    </div>

                    <div class="flex items-center justify-self-end text-right">
                      {{ priceFormatter(item.product.price) }}
                    </div>
                  </template>
                </div>
              </div>

              <div
                class="grid grid-cols-5 gap-y-4 border-y border-main py-6 text-[12px] tracking-[1.4px] md:grid-cols-2 md:text-[14px]"
              >
                <div class="col-span-3 md:col-span-1">
                  {{ ct("itemTotalPriceLabel") }}
                </div>
                <div class="col-span-2 text-right md:col-span-1">
                  {{ priceFormatter(calcCartResponseItem.subtotal) }}
                </div>
                <div class="col-span-3 md:col-span-1">
                  {{ ct("applyCouponLabel") }}
                </div>
                <div class="col-span-2 text-right md:col-span-1">
                  {{ priceFormatter(calcCartResponseItem.discount) }}
                </div>
                <div class="col-span-3 md:col-span-1">
                  {{ ct("shippingFeeLabel") }}
                </div>
                <div class="col-span-2 text-right md:col-span-1">
                  {{ priceFormatter(calcCartResponseItem.shippingFee) }}
                </div>
              </div>

              <div
                class="mt-6 grid grid-cols-2 text-[14px] font-bold tracking-[1.4px]"
              >
                <div>{{ ct("totalPriceLabel") }}</div>
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
            @click="handleClickPreviousStepButton"
          >
            <the-left-arrow-icon class="h-4 w-4" />
            {{ ct("backToPreviousPageButtonText") }}
          </button>
          <button
            class="w-full bg-orange p-[14px] text-[16px] text-white md:order-1 md:w-[240px]"
            :type="checkoutFormData.paymentMethod === 2 ? 'submit' : 'button'"
            :form="
              checkoutFormData.paymentMethod === 2 ? 'credit-card-form' : ''
            "
            @click="handleClickNextStepButton"
          >
            {{ ct("proceedToPaymentButtonText") }}
          </button>
        </div>
      </template>
    </div>
  </div>
</template>
