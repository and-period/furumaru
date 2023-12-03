<script setup lang="ts">
import { storeToRefs } from 'pinia'
import { MOCK_PURCHASE_ITEMS } from '~/constants/mock'
import { CreateAddressRequest } from '~/types/api'
import { useAdressStore } from '~/store/address'
import { convertJapaneseToI18nPhoneNumber } from '~/lib/phone-number'

const router = useRouter()

const cartItem = MOCK_PURCHASE_ITEMS[0]

const discount = 0

const addressStore = useAdressStore()
const { addressesFetchState } = storeToRefs(addressStore)
const { fetchAddresses, searchAddressByPostalCode, registerAddress } =
  addressStore

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

const itemsTotalPrice = computed(() => {
  return cartItem.cartItems[0].items
    .map((item) => item.price)
    .reduce((sum, price) => sum + price)
})

const totalPrice = computed(() => {
  return itemsTotalPrice.value + discount
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

const handleClickPreviousStepButton = () => {
  router.back()
}

const handleClickNextStepButton = async () => {
  await registerAddress({
    ...formData.value,
    phoneNumber: convertJapaneseToI18nPhoneNumber(formData.value.phoneNumber),
  })
  router.push('/v1/purchase/confirmation')
}

onMounted(() => {
  fetchAddresses()
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
      ご購入手続き
    </div>

    <div class="mx-[15px] my-10 bg-white px-6 pb-10 md:mx-0">
      <div class="gap-[80px] md:grid md:grid-cols-2">
        <div class="md:pl-10">
          <div>
            <div
              class="pt-[24px] text-left text-[16px] font-bold tracking-[1.6px] text-main md:pt-[80px]"
            >
              お客様情報
            </div>
            <the-new-address-form
              v-model:form-data="formData"
              form-id="new-address-form"
              @click:search-address-button="handleClickSearchAddressButton"
              @submit="handleClickNextStepButton"
            />
            <div class="hidden md:block">
              <div class="mt-12 md:grid md:grid-cols-2">
                <div class="flex items-center">
                  <button
                    class="inline-flex"
                    @click="handleClickPreviousStepButton"
                  >
                    <the-left-arrow-icon class="h-4 w-4" />
                    <p class="pl-2 text-[12px] tracking-[1.2px] text-main">
                      前のページへ戻る
                    </p>
                  </button>
                </div>

                <button
                  class="w-full bg-main p-[14px] text-[16px] text-white md:w-[240px] md:justify-self-end"
                  type="submit"
                  form="new-address-form"
                >
                  お支払方法の選択へ
                </button>
              </div>
            </div>
          </div>
        </div>

        <div class="mt-[24px] md:mr-10 md:mt-10">
          <div class="w-full bg-base px-[16px] py-[24px] text-main md:p-10">
            <div class="text-[14px] font-bold tracking-[1.6px] md:text-[16px]">
              注文内容
            </div>
            <div class="my-[16px] text-[12px] tracking-[1.2px] md:my-6">
              <p>
                {{ cartItem.marche }}
              </p>
              <p>発想地：{{ cartItem.address }}</p>
              <p>
                取扱元：
                {{ cartItem.sender }}
              </p>
              <p>箱の数：2（常温・冷蔵 ✕ 2）</p>
            </div>
            <div>
              <div>
                <div
                  v-for="(item, i) in cartItem.cartItems[0].items"
                  :key="i"
                  class="grid grid-cols-5 border-t py-2 text-[12px] tracking-[1.2px]"
                >
                  <img
                    :src="item.imgSrc"
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

              <div class="mb-4 flex gap-2">
                <div class="grow">
                  <input
                    type="text"
                    class="w-full border border-gray-300 bg-gray-50 p-2.5 text-[14px] md:text-[16px]"
                    placeholder="クーポンコード"
                  />
                </div>
                <button
                  class="whitespace-nowrap bg-main p-2 text-[14px] text-white md:text-[16px]"
                >
                  適用する
                </button>
              </div>

              <div
                class="grid grid-cols-5 gap-y-4 border-y border-main py-6 text-[12px] tracking-[1.4px] md:grid-cols-2 md:text-[14px]"
              >
                <div class="col-span-2 md:col-span-1">商品合計（税込）</div>
                <div class="col-span-3 text-right md:col-span-1">
                  {{ priceFormatter(itemsTotalPrice) }}
                </div>
                <div class="col-span-2 md:col-span-1">クーポン利用</div>
                <div class="col-span-3 text-right md:col-span-1">
                  {{ priceFormatter(discount) }}
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
                  {{ priceFormatter(totalPrice) }}
                </div>
              </div>
            </div>
          </div>
          <div class="block md:hidden">
            <div class="mt-12">
              <button
                class="w-full bg-main p-[14px] text-[14px] tracking-[1.4px] text-white md:w-[240px] md:justify-self-end md:text-[16px] md:tracking-[1.6px]"
                @click="handleClickNextStepButton"
              >
                お支払方法の選択へ
              </button>
              <div class="mt-[40px] items-center">
                <button
                  class="inline-flex"
                  @click="handleClickPreviousStepButton"
                >
                  <the-left-arrow-icon class="h-4 w-4" />
                  <p class="pl-2 text-[12px] tracking-[1.2px] text-main">
                    前のページへ戻る
                  </p>
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
