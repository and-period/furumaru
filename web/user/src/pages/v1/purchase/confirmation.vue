<script setup lang="ts">
import { MOCK_USER_INFOMATION, MOCK_PURCHASE_ITEMS } from '~/constants/mock'

const userInformationItem = MOCK_USER_INFOMATION
const cartItem = MOCK_PURCHASE_ITEMS[0]

const discount = -500
const shipping = 2500

const itemsTotalPrice = computed(() => {
  return cartItem.cartItems[0].items
    .map((item) => item.price)
    .reduce((sum, price) => sum + price)
})

const totalPrice = computed(() => {
  return itemsTotalPrice.value + discount + shipping
})

const priceFormatter = (price: number) => {
  return new Intl.NumberFormat('ja-JP', {
    style: 'currency',
    currency: 'JPY',
  }).format(price)
}
</script>

<template>
  <div class="container mx-auto">
    <div class="text-center text-[20px] font-bold tracking-[2px] text-main">
      ご購入手続き
    </div>
    <div class="my-10 bg-white px-6 pb-10">
      <div class="grid grid-cols-2 gap-[80px]">
        <div class="pl-10">
          <div>
            <div
              class="pt-[80px] text-left text-[16px] font-bold tracking-[1.6px] text-main"
            >
              お客様情報
            </div>
            <dl
              class="mt-[28px] w-full text-[14px] leading-[24px] tracking-[1.4px]"
            >
              <div class="grid grid-cols-3">
                <dt>氏名</dt>
                <dd class="col-span-2">
                  {{ userInformationItem.name }}
                </dd>
              </div>
              <div class="grid grid-cols-3">
                <dt>電話番号</dt>
                <dd class="col-span-2">
                  {{ userInformationItem.phoneNumber }}
                </dd>
              </div>
              <div class="grid grid-cols-3">
                <dt>メールアドレス</dt>
                <dd class="col-span-2">
                  {{ userInformationItem.email }}
                </dd>
              </div>
              <div class="grid grid-cols-3">
                <dt>住所</dt>
                <dd class="col-span-2">
                  {{ userInformationItem.address }}
                </dd>
              </div>
            </dl>
          </div>
          <div class="pr-8 pt-4 text-right">
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
            <div class="pr-8 pt-4 text-right">
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
              <input
                class="check:before:border-main relative float-left mr-1 mt-0.5 h-5 w-5 appearance-none rounded-full border-2 border-solid border-neutral-300 before:pointer-events-none before:absolute before:h-4 before:w-4 before:scale-0 before:rounded-full before:bg-transparent before:opacity-0 before:shadow-[0px_0px_0px_13px_transparent] before:content-[''] after:absolute after:z-[1] after:block after:h-4 after:w-4 after:rounded-full after:content-[''] checked:border-main checked:before:opacity-[0.16] checked:after:absolute checked:after:left-1/2 checked:after:top-1/2 checked:after:h-[0.625rem] checked:after:w-[0.625rem] checked:after:rounded-full checked:after:bg-main checked:after:content-[''] checked:after:[transform:translate(-50%,-50%)] hover:cursor-pointer hover:before:opacity-[0.04] hover:before:shadow-[0px_0px_0px_13px_rgba(0,0,0,0.6)] focus:shadow-none focus:outline-none focus:ring-0 focus:before:scale-100 focus:before:opacity-[0.12] focus:before:shadow-[0px_0px_0px_13px_rgba(0,0,0,0.6)] focus:before:transition-[box-shadow_0.2s,transform_0.2s] checked:focus:border-main checked:focus:before:scale-100 checked:focus:before:shadow-[0px_0px_0px_13px_#3b71ca] checked:focus:before:transition-[box-shadow_0.2s,transform_0.2s] dark:border-neutral-600 dark:focus:before:shadow-[0px_0px_0px_13px_rgba(255,255,255,0.4)] dark:checked:focus:before:shadow-[0px_0px_0px_13px_#3b71ca]"
                type="radio"
                checked
              />
              <label class="pl-2 text-main"> クレジットカード支払い </label>
              <div class="flex w-full items-center gap-4">
                <the-text-input
                  placeholder="カード番号"
                  :with-label="false"
                  name="cc-number"
                  type="text"
                  class="mt-4 w-full"
                  required
                />
                <div class="flex h-[24px] min-w-max gap-2">
                  <img src="~/assets/img/cc/visa.png" alt="visa icon" />
                  <img src="~/assets/img/cc/jcb.png" alt="jcb icon" />
                  <img src="~/assets/img/cc/amex.png" alt="amex icon" />
                  <img src="~/assets/img/cc/master.png" alt="master icon" />
                </div>
              </div>
              <the-text-input
                placeholder="カード名義"
                name="cc-name"
                :with-label="false"
                type="text"
                class="mt-4"
                required
              />
              <div class="flex gap-4">
                <the-text-input
                  placeholder="有効期限 (月)"
                  :with-label="false"
                  name="cc-exp-month"
                  type="text"
                  pattern="[0-9]*"
                  class="mt-4 w-1/2"
                  required
                />
                <the-text-input
                  placeholder="有効期限 (年)"
                  :with-label="false"
                  name="cc-exp-year"
                  type="text"
                  pattern="[0-9]*"
                  class="mt-4 w-1/2"
                  required
                />
                <the-text-input
                  placeholder="セキュリティコード"
                  :with-label="false"
                  name="cc-csc"
                  type="text"
                  pattern="[0-9]*"
                  class="mt-4 w-1/2"
                  required
                />
              </div>
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

            <div class="mt-12 grid grid-cols-2">
              <div class="inline-flex items-center">
                <the-left-arrow-icon class="h-4 w-4" />
                <p class="pl-2 text-[12px] tracking-[1.2px] text-main">
                  前のページへ戻る
                </p>
              </div>

              <button
                class="w-[240px] justify-self-end bg-main p-[14px] text-[16px] text-white"
              >
                注文確定
              </button>
            </div>
          </div>
        </div>

        <div class="mr-10 mt-10">
          <div class="w-full bg-base p-10 text-main">
            <div class="text-[16px] font-bold tracking-[1.6px]">注文内容</div>
            <div class="my-6 text-[12px] tracking-[1.2px]">
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
                  <div class="col-span-2">{{ item.name }}</div>
                  <div
                    class="flex w-full items-center justify-self-end text-right"
                  >
                    数量：{{ 1 }}
                  </div>
                  <div class="flex items-center justify-self-end text-right">
                    {{ priceFormatter(item.price) }}
                  </div>
                </div>
              </div>

              <div
                class="grid grid-cols-2 gap-y-4 border-y border-main py-6 text-[14px] tracking-[1.4px]"
              >
                <div>商品合計（税込み）</div>
                <div class="text-right">
                  {{ priceFormatter(itemsTotalPrice) }}
                </div>
                <div>クーポン利用</div>
                <div class="text-right">
                  {{ priceFormatter(discount) }}
                </div>
                <div>送料（合計）</div>
                <div class="text-right">
                  {{ priceFormatter(shipping) }}
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
        </div>
      </div>
    </div>
  </div>
</template>
