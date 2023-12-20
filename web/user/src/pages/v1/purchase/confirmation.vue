<script setup lang="ts">
import { storeToRefs } from 'pinia'
import { useAdressStore } from '~/store/address'
import { useCheckoutStore } from '~/store/checkout'
import { useShoppingCartStore } from '~/store/shopping'
import type { CheckoutRequest } from '~/types/api'
import { ApiBaseError } from '~/types/exception'

const addressStore = useAdressStore()
const { address, addressFetchState } = storeToRefs(addressStore)
const { fetchAddress } = addressStore

const shoppingCartStore = useShoppingCartStore()
const { calcCartResponseItem, availablePaymentSystem } =
  storeToRefs(shoppingCartStore)
const { calcCartItemByCoordinatorId, fetchAvailablePaymentOptions } =
  shoppingCartStore

const checkoutStore = useCheckoutStore()
const { checkout } = checkoutStore

const route = useRoute()
const router = useRouter()

const addressId = computed<string>(() => {
  const id = route.query.id
  if (id) {
    return String(id)
  } else {
    return ''
  }
})

const coordinatorId = computed<string>(() => {
  const id = route.query.coordinatorId
  if (id) {
    return String(id)
  } else {
    return ''
  }
})

const checkoutFormData = ref<CheckoutRequest>({
  requestId: '',
  coordinatorId: '',
  boxNumber: 0,
  billingAddressId: '',
  shippingAddressId: '',
  promotionCode: '',
  paymentMethod: 0,
  callbackUrl: '',
  total: 0,
  creditCard: {
    number: '',
    month: 0,
    year: 0,
    verificationValue: '',
  },
})

const checkoutError = ref<string>('')

const priceFormatter = (price: number) => {
  return new Intl.NumberFormat('ja-JP', {
    style: 'currency',
    currency: 'JPY',
  }).format(price)
}

const handleClickPreviousStepButton = () => {
  router.back()
}

const handleClickNextStepButton = async () => {
  try {
    const url = await checkout(checkoutFormData.value)
    window.location.href = url
  } catch (error) {
    if (error instanceof ApiBaseError) {
      checkoutError.value = error.message
    }
  }
}

onMounted(async () => {
  fetchAvailablePaymentOptions().then(() => {
    if (availablePaymentSystem.value.length > 0) {
      checkoutFormData.value.paymentMethod =
        availablePaymentSystem.value[0].methodType
    }
  })

  if (addressId.value) {
    checkoutFormData.value.billingAddressId = addressId.value
    checkoutFormData.value.shippingAddressId = addressId.value
    await fetchAddress(addressId.value)
  }

  await calcCartItemByCoordinatorId(
    coordinatorId.value,
    address.value?.prefectureCode,
  )

  checkoutFormData.value.requestId = calcCartResponseItem.value?.requestId ?? ''
  checkoutFormData.value.coordinatorId = coordinatorId.value
  checkoutFormData.value.total = calcCartResponseItem.value?.total ?? 0
  checkoutFormData.value.callbackUrl = `${window.location.origin}/v1/purchase/complete`
})

useSeoMeta({
  title: 'ご購入手続き',
})
</script>

<template>
  <div class="container mx-auto">
    <div class="text-center text-[20px] font-bold tracking-[2px] text-main">
      ご購入手続き
    </div>

    <the-alert v-if="checkoutError" class="mt-4 bg-white" type="error">{{
      checkoutError
    }}</the-alert>

    <div
      class="relative my-10 gap-x-[80px] bg-white px-6 py-10 md:mx-0 md:grid md:grid-cols-2 md:grid-rows-[auto_auto] md:px-[80px]"
    >
      <template v-if="addressFetchState.isLoading">
        <div class="absolute h-2 w-full animate-pulse bg-main"></div>
      </template>

      <template v-else>
        <!-- 左側 -->
        <div class="row-span-1 self-start py-[24px] md:w-full md:py-10">
          <div
            class="mb-6 text-left text-[16px] font-bold tracking-[1.6px] text-main"
          >
            お客様情報
          </div>
          <the-address-info v-if="address" :address="address" />
          <div class="pt-4 text-right tracking-[1.4px]">
            <a href="#" class="underline">変更</a>
          </div>

          <div class="items-center border-b py-2" />

          <div>
            <div
              class="pt-6 text-left text-[16px] font-bold tracking-[1.6px] text-main"
            >
              お届け先情報
            </div>
            <div class="pt-[27px] text-[14px] tracking-[1.4px] text-main">
              上記の住所にお届け
            </div>
            <div class="pt-4 text-right tracking-[1.4px]">
              <a href="#" class="underline">変更</a>
            </div>
          </div>

          <div class="items-center border-b py-2" />

          <div>
            <div
              class="pt-6 text-left text-[16px] font-bold tracking-[1.6px] text-main"
            >
              お支払い情報
            </div>

            <div class="pt-4">
              <div class="flex items-center justify-between">
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
                      />
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
                        />
                        <img
                          src="~/assets/img/cc/jcb.png"
                          alt="jcb icon"
                          class="h-full"
                        />
                        <img
                          src="~/assets/img/cc/amex.png"
                          alt="amex icon"
                          class="h-full"
                        />
                        <img
                          src="~/assets/img/cc/master.png"
                          alt="master icon"
                          class="h-full"
                        />
                      </div>
                    </template>
                  </div>
                </div>
              </div>

              <template v-if="checkoutFormData.paymentMethod === 2">
                <div class="mt-4 flex w-full items-center gap-4">
                  <the-text-input
                    v-model="checkoutFormData.creditCard.number"
                    placeholder="カード番号"
                    :with-label="false"
                    name="cc-number"
                    type="text"
                    class="w-full"
                    required
                  />
                  <div
                    class="hidden h-[24px] min-w-max items-center gap-2 md:flex"
                  >
                    <img
                      src="~/assets/img/cc/visa.png"
                      alt="visa icon"
                      class="h-full"
                    />
                    <img
                      src="~/assets/img/cc/jcb.png"
                      alt="jcb icon"
                      class="h-full"
                    />
                    <img
                      src="~/assets/img/cc/amex.png"
                      alt="amex icon"
                      class="h-full"
                    />
                    <img
                      src="~/assets/img/cc/master.png"
                      alt="master icon"
                      class="h-full"
                    />
                  </div>
                </div>
                <div class="flex gap-4">
                  <the-text-input
                    v-model="checkoutFormData.creditCard.month"
                    placeholder="有効期限 (月)"
                    :with-label="false"
                    name="cc-exp-month"
                    type="number"
                    pattern="[0-9]*"
                    class="mt-4 w-1/2"
                    required
                  />
                  <the-text-input
                    v-model="checkoutFormData.creditCard.year"
                    placeholder="有効期限 (年)"
                    :with-label="false"
                    name="cc-exp-year"
                    type="number"
                    pattern="[0-9]*"
                    class="mt-4 w-1/2"
                    required
                  />
                </div>
                <the-text-input
                  v-model="checkoutFormData.creditCard.verificationValue"
                  placeholder="セキュリティコード"
                  :with-label="false"
                  name="cc-csc"
                  type="text"
                  pattern="[0-9]*"
                  class="mt-4 w-1/2"
                  required
                />
              </template>
            </div>

            <div
              class="mt-12 text-left text-[16px] font-bold tracking-[1.6px] text-main"
            >
              お届け日の指定
            </div>
            <div class="mt-4 grid grid-cols-2 gap-4">
              <the-text-input
                placeholder="お届け希望日"
                :with-label="false"
                type="date"
                pattern="[0-9]*"
                class="w-full"
                required
              />
              <the-text-input
                placeholder="お届け時間帯"
                :with-label="false"
                type="time"
                pattern="[0-9]*"
                class="w-full"
                required
              />
            </div>
          </div>
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
              <div>
                <div
                  v-for="(item, i) in calcCartResponseItem.products"
                  :key="i"
                  class="grid grid-cols-5 border-t py-2 text-[12px] tracking-[1.2px]"
                >
                  <img
                    :src="item.media[0].url"
                    :alt="`${item.name}の画像`"
                    class="block aspect-square h-[56px] w-[56px]"
                  />
                  <div class="col-span-3 pl-[24px] md:pl-0">
                    <div>{{ item.name }}</div>
                    <div
                      class="mt-4 md:mt-0 md:items-center md:justify-self-end md:text-right"
                    >
                      数量：{{ 1 }}
                    </div>
                  </div>

                  <div class="flex items-center justify-self-end text-right">
                    {{ priceFormatter(item.price) }}
                  </div>
                </div>
              </div>

              <div
                class="grid grid-cols-5 gap-y-4 border-y border-main py-6 text-[12px] tracking-[1.4px] md:grid-cols-2 md:text-[14px]"
              >
                <div class="col-span-3 md:col-span-1">商品合計（税込み）</div>
                <div class="col-span-2 text-right md:col-span-1">
                  {{ priceFormatter(calcCartResponseItem.subtotal) }}
                </div>
                <div class="col-span-3 md:col-span-1">クーポン利用</div>
                <div class="col-span-2 text-right md:col-span-1">
                  {{ priceFormatter(calcCartResponseItem.discount) }}
                </div>
                <div class="col-span-3 md:col-span-1">送料（合計）</div>
                <div class="col-span-2 text-right md:col-span-1">
                  {{ priceFormatter(calcCartResponseItem.shippingFee) }}
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
            @click="handleClickPreviousStepButton"
          >
            <the-left-arrow-icon class="h-4 w-4" />
            前のページへ戻る
          </button>
          <button
            class="w-full bg-main p-[14px] text-[16px] text-white md:order-1 md:w-[240px]"
            @click="handleClickNextStepButton"
          >
            支払い画面へ
          </button>
        </div>
      </template>
    </div>
  </div>
</template>
