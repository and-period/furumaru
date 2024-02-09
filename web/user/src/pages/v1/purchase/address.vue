<script setup lang="ts">
import { storeToRefs } from 'pinia'
import type { CreateAddressRequest } from '~/types/api'
import { useAdressStore } from '~/store/address'
import { convertJapaneseToI18nPhoneNumber } from '~/lib/phone-number'
import { useShoppingCartStore } from '~/store/shopping'
import { ApiBaseError } from '~/types/exception'

const route = useRoute()
const router = useRouter()

const addressStore = useAdressStore()
const { addressesFetchState, defaultAddress } = storeToRefs(addressStore)
const { fetchAddresses, searchAddressByPostalCode, registerAddress } =
  addressStore

const shoppingCartStore = useShoppingCartStore()
const { calcCartResponseItem } = storeToRefs(shoppingCartStore)
const { calcCartItemByCoordinatorId, verifyPromotionCode } = shoppingCartStore

const targetAddress = ref<string>('default')
const calcCartResponseItemState = ref<{
  isLoading: boolean
  hasError: boolean
  errorMessage: string
}>({ isLoading: true, hasError: false, errorMessage: '' })

const formData = ref<CreateAddressRequest>({
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
  isDefault: true,
})

const promotionCodeFormValue = ref<string>('')
const invalidPromotion = ref<boolean>(false)
const validPromotion = ref<boolean>(false)

const coordinatorId = computed<string>(() => {
  const id = route.query.coordinatorId
  if (id) {
    return String(id)
  } else {
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
  } else {
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

const handleSubmitNewAddressForm = async () => {
  const registeredAddress = await registerAddress({
    ...formData.value,
    phoneNumber: convertJapaneseToI18nPhoneNumber(formData.value.phoneNumber),
  })
  router.push({
    path: '/v1/purchase/confirmation',
    query: {
      id: registeredAddress.id,
      coordinatorId: coordinatorId.value,
      cartNumber: cartNumber.value,
      promotionCode: promotionCode.value,
    },
  })
}

const handleClickNextStepButton = (id: string) => {
  router.push({
    path: '/v1/purchase/confirmation',
    query: {
      id,
      coordinatorId: coordinatorId.value,
      cartNumber: cartNumber.value,
      promotionCode: promotionCode.value,
    },
  })
}

/**
 * クーポンコードを適用するボタンをクリックしたときの処理
 */
const handleClickUsePromotionCodeButton = async () => {
  const result = await verifyPromotionCode(promotionCodeFormValue.value)
  if (result) {
    invalidPromotion.value = false
    validPromotion.value = true
  } else {
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
      path: '/v1/purchase/address',
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
      path: '/v1/purchase/address',
      query: {
        coordinatorId: coordinatorId.value,
        cartNumber: cartNumber.value,
        promotionCode: undefined,
      },
    })
  }
})

onMounted(() => {
  fetchAddresses()
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
  } catch (error) {
    calcCartResponseItemState.value.hasError = true
    if (error instanceof ApiBaseError) {
      calcCartResponseItemState.value.errorMessage = error.message
    } else {
      calcCartResponseItemState.value.errorMessage =
        '不明なエラーが発生しました。'
    }
  } finally {
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
      class="text-main mt-[32px] text-center text-[20px] font-bold tracking-[2px]"
    >
      ご購入手続き
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
      <template v-if="addressesFetchState.isLoading">
        <div class="bg-main absolute h-2 w-full animate-pulse"></div>
      </template>

      <template v-else>
        <!-- 左側 -->
        <div class="row-span-1 self-start py-[24px] md:w-full md:py-10">
          <div
            class="text-main mb-6 text-left text-[16px] font-bold tracking-[1.6px]"
          >
            お客様情報
          </div>

          <!-- デフォルトの住所が登録されている場合 -->
          <template v-if="defaultAddress">
            <the-address-info :address="defaultAddress" />
            <hr class="my-[20px]" />
            <div class="flex flex-col gap-4">
              <div class="flex items-center gap-2">
                <input
                  id="default-radio"
                  v-model="targetAddress"
                  type="radio"
                  class="accent-main h-4 w-4"
                  value="default"
                />
                <label for="default-radio">上記の住所にお届け</label>
              </div>
              <div class="flex items-center gap-2">
                <input
                  id="other-radio"
                  v-model="targetAddress"
                  type="radio"
                  class="accent-main h-4 w-4"
                  value="other"
                />
                <label for="other-radio">その他の住所にお届け</label>
              </div>
            </div>
            <template v-if="targetAddress === 'other'">
              <div
                class="text-main my-6 text-[16px] font-bold tracking-[1.6px]"
              >
                お届け先情報
              </div>
              <the-new-address-form
                v-model:form-data="formData"
                form-id="new-address-form"
                @click:search-address-button="handleClickSearchAddressButton"
                @submit="handleSubmitNewAddressForm"
              />
            </template>
          </template>

          <!-- デフォルトの住所が登録されていない場合 -->
          <template v-else>
            <the-new-address-form
              v-model:form-data="formData"
              form-id="new-address-form"
              @click:search-address-button="handleClickSearchAddressButton"
              @submit="handleSubmitNewAddressForm"
            />
          </template>
        </div>

        <!-- 右側 -->
        <div
          class="bg-base text-main row-span-2 self-start px-[16px] py-[24px] md:w-full md:p-10"
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
                発想地：{{
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
                    <img
                      v-if="item.product.thumbnail"
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
                  class="leading-[1.4px border-orange text-orange mt-4 flex justify-between rounded-lg border p-2 text-[12px]"
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
                    />
                  </div>
                  <button
                    class="bg-main whitespace-nowrap p-2 text-[14px] text-white md:text-[16px]"
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
                class="border-main mt-4 grid grid-cols-5 gap-y-4 border-y py-6 text-[12px] tracking-[1.4px] md:grid-cols-2 md:text-[14px]"
              >
                <div class="col-span-2 md:col-span-1">商品合計（税込）</div>
                <div class="col-span-3 text-right md:col-span-1">
                  {{ priceFormatter(calcCartResponseItem.subtotal) }}
                </div>
                <div class="col-span-2 md:col-span-1">クーポン利用</div>
                <div class="col-span-3 text-right md:col-span-1">
                  {{ priceFormatter(calcCartResponseItem.discount) }}
                </div>
                <div class="col-span-2 md:col-span-1">送料（税込）</div>
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
            class="text-main order-2 inline-flex w-full gap-2 text-left text-[12px] tracking-[1.2px] md:order-1 md:max-w-max"
            @click="handleClickBackCartButton"
          >
            <the-left-arrow-icon class="h-4 w-4" />
            買い物カゴへ戻る
          </button>

          <template v-if="defaultAddress && targetAddress === 'default'">
            <!-- 通常のボタンの場合 -->
            <button
              class="bg-main w-full p-[14px] text-[16px] text-white md:order-1 md:w-[240px]"
              @click="handleClickNextStepButton(defaultAddress.id)"
            >
              お支払方法の選択へ
            </button>
          </template>
          <template v-else>
            <!-- フォーム要素の場合 -->
            <button
              class="bg-main w-full p-[14px] text-[16px] text-white md:order-1 md:w-[240px]"
              type="submit"
              form="new-address-form"
            >
              お支払方法の選択へ
            </button>
          </template>
        </div>
      </template>
    </div>
  </div>
</template>
