<script setup lang="ts">
import { MOCK_PURCHASE_ITEMS } from '~/constants/mock'
import { CreateAddressRequest } from '~/types/api'
import { useAdressStore } from '~/store/address'

const router = useRouter()

const cartItem = MOCK_PURCHASE_ITEMS[0]

const discount = 0

const addressStore = useAdressStore()
const { registerAddress, fetchAddresses } = addressStore

const formData = ref<CreateAddressRequest>({
  lastname: '',
  firstname: '',
  lastnameKana: '',
  firstnameKana: '',
  postalCode: '',
  prefecture: 0,
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

const handleClickPreviousStepButton = () => {
  router.back()
}

const handleClickNextStepButton = async () => {
  await registerAddress(formData.value)
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
            <div class="mt-4 grid grid-cols-2 gap-4">
              <the-text-input
                v-model="formData.lastname"
                placeholder="性"
                :with-label="false"
                type="text"
                required
              />
              <the-text-input
                v-model="formData.firstname"
                placeholder="名"
                :with-label="false"
                type="text"
                required
              />
            </div>
            <div class="mt-4 grid grid-cols-2 gap-4">
              <the-text-input
                v-model="formData.lastnameKana"
                placeholder="フリガナ(セイ)"
                :with-label="false"
                type="text"
                required
              />
              <the-text-input
                v-model="formData.firstnameKana"
                placeholder="フリガナ(メイ)"
                :with-label="false"
                type="text"
                required
              />
            </div>
            <the-text-input
              v-model="formData.phoneNumber"
              placeholder="電話番号"
              :with-label="false"
              type="text"
              class="pt-4"
              required
            />
            <the-text-input
              placeholder="メールアドレス"
              :with-label="false"
              type="text"
              class="pt-4"
              required
            />
            <the-text-input
              placeholder="メールアドレス(確認用)"
              :with-label="false"
              type="text"
              class="pt-4"
              required
            />
            <div
              class="mt-4 flex flex-col gap-4 md:grid md:grid-cols-2 md:gap-4"
            >
              <the-text-input
                v-model="formData.postalCode"
                placeholder="郵便番号（ハイフンなし）"
                :with-label="false"
                type="text"
                required
              />
              <the-text-input
                v-model="formData.prefecture"
                placeholder="都道府県"
                :with-label="false"
                type="text"
                required
              />
            </div>
            <the-text-input
              v-model="formData.city"
              placeholder="住所（市区町村)"
              :with-label="false"
              type="text"
              class="pt-4"
              required
            />
            <the-text-input
              v-model="formData.addressLine1"
              placeholder="住所（それ以降）"
              :with-label="false"
              type="text"
              class="pt-4"
              required
            />
            <the-text-input
              v-model="formData.addressLine2"
              placeholder="住所（マンション名、部屋番号）"
              :with-label="false"
              type="text"
              class="pt-4"
              required
            />
            <div class="pt-8">
              <input
                v-model="formData.isDefault"
                class="check:before:border-main relative float-left mr-1 mt-0.5 h-5 w-5 appearance-none rounded-full border-2 border-solid border-neutral-300 before:pointer-events-none before:absolute before:h-4 before:w-4 before:scale-0 before:rounded-full before:bg-transparent before:opacity-0 before:shadow-[0px_0px_0px_13px_transparent] before:content-[''] after:absolute after:z-[1] after:block after:h-4 after:w-4 after:rounded-full after:content-[''] checked:border-main checked:before:opacity-[0.16] checked:after:absolute checked:after:left-1/2 checked:after:top-1/2 checked:after:h-[0.625rem] checked:after:w-[0.625rem] checked:after:rounded-full checked:after:bg-main checked:after:content-[''] checked:after:[transform:translate(-50%,-50%)] hover:cursor-pointer hover:before:opacity-[0.04] hover:before:shadow-[0px_0px_0px_13px_rgba(0,0,0,0.6)] focus:shadow-none focus:outline-none focus:ring-0 focus:before:scale-100 focus:before:opacity-[0.12] focus:before:shadow-[0px_0px_0px_13px_rgba(0,0,0,0.6)] focus:before:transition-[box-shadow_0.2s,transform_0.2s] checked:focus:border-main checked:focus:before:scale-100 checked:focus:before:shadow-[0px_0px_0px_13px_#3b71ca] checked:focus:before:transition-[box-shadow_0.2s,transform_0.2s] dark:border-neutral-600 dark:focus:before:shadow-[0px_0px_0px_13px_rgba(255,255,255,0.4)] dark:checked:focus:before:shadow-[0px_0px_0px_13px_#3b71ca]"
                type="radio"
                checked
              />
              <label class="pl-2 text-[14px] text-main md:text-[16px]">
                上記の住所にお届け
              </label>
            </div>

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
                  @click="handleClickNextStepButton"
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

              <div class="items-center border-b" />
              <div class="grid grid-cols-5 gap-2 py-6 md:grid-cols-4">
                <div class="col-span-3">
                  <input
                    type="text"
                    class="w-full border border-gray-300 bg-gray-50 p-2.5 text-[14px] md:text-[16px]"
                    placeholder="クーポンコード"
                  />
                </div>

                <button
                  class="col-span-2 w-full justify-self-end bg-main p-2 text-[14px] text-white md:col-span-1 md:text-[16px]"
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
