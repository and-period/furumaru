<script setup lang="ts">
import { MOCK_USER_INFOMATION } from '~/constants/mock'

const userInformationItem = MOCK_USER_INFOMATION
</script>

<template>
  <div class="container mx-auto">
    <div class="text-center text-[20px] font-bold tracking-[2px] text-main">
      ご購入手続き
    </div>
    <div class="my-10 bg-white px-6">
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
              <the-text-input
                placeholder="カード番号"
                :with-label="false"
                name="cc-number"
                type="text"
                class="mt-4"
                required
              />
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
                  type="text"
                  pattern="[0-9]*"
                  class="mt-4 w-1/2"
                  required
                />
              </div>
            </div>
            <div
              class="pt-12 text-left text-[16px] font-bold tracking-[1.6px] text-main"
            >
              お届け日の指定
            </div>
            <div class="flex gap-4">
              <the-text-input
                placeholder="お届け希望日"
                :with-label="false"
                type="date"
                pattern="[0-9]*"
                class="mt-4 w-1/2"
                required
              />
              <the-text-input
                placeholder="お届け時間帯"
                :with-label="false"
                type="time"
                pattern="[0-9]*"
                class="mt-4 w-1/2"
                required
              />
            </div>
            <div class="flex pb-4 pt-[80px]">
              <div class="grid grid-cols-2">
                <div class="flex items-center">
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    width="7"
                    height="12"
                    viewBox="0 0 7 12"
                    fill="none"
                  >
                    <path
                      fill-rule="evenodd"
                      clip-rule="evenodd"
                      d="M5.65674 0L6.62648 0.969746L1.93938 5.65685L6.62648 10.344L5.65674 11.3137L-0.000115871 5.65685L5.65674 0Z"
                      fill="#604C3F"
                    />
                  </svg>
                  <p
                    class="pl-2 text-[16px] font-bold tracking-[1.6px] text-main"
                  >
                    前のページへ戻る
                  </p>
                </div>
              </div>
              <div class="pb-6">
                <button
                  class="w-[240px] bg-main p-[14px] text-[16px] text-white"
                >
                  注文確定
                </button>
              </div>
            </div>
          </div>
        </div>

        <div class="mr-10 mt-10">
          <div class="w-full bg-base p-10 text-main">
            <div class="text-[16px] font-bold tracking-[1.6px]">注文内容</div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
