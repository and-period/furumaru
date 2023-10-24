<script setup lang="ts">
import { MOCK_USER_INFOMATION } from '~/constants/mock'

const userInforamtionItem = MOCK_USER_INFOMATION
</script>

<template>
  <div class="container mx-auto">
    <div class="text-center text-[20px] font-bold tracking-[2px] text-main">
      ご購入手続き
    </div>
    <div class="my-10 bg-white px-6">
      <div class="grid grid-cols-2 gap-2">
        <div>
          <div class="pl-10">
            <div class="grid grid-cols-4">
              <div class="col-span-1">
                <div
                  class="pt-[80px] text-left text-[16px] font-bold tracking-[1.6px] text-main"
                >
                  お客様情報
                </div>
                <div class="pt-[27px] text-[14px] tracking-[1.4px]">
                  <p>氏名</p>
                  <p>電話番号</p>
                  <p>メールアドレス</p>
                  <p>住所</p>
                </div>
              </div>
              <div class="col-span-3 pl-8 pt-[125px] text-[14px] tracking-[1.4px]">
                <p>
                  {{ userInforamtionItem.name }}
                </p>
                <p class="">
                  {{ userInforamtionItem.phoneNumber }}
                </p>
                <p>
                  {{ userInforamtionItem.email }}
                </p>
                <p>
                  {{ userInforamtionItem.address }}
                </p>
              </div>
            </div>
          </div>
        </div>
        <div class="mr-10 mt-10 bg-base">

        </div>
      </div>
    </div>
  </div>
</template>
