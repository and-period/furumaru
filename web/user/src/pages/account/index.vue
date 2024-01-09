<script setup lang="ts">
import { storeToRefs } from 'pinia'
import { useAdressStore } from '~/store/address'
import { convertI18nToJapanesePhoneNumber } from '~/lib/phone-number'
import { useAuthStore } from '~/store/auth'

const addressStore = useAdressStore()
const { addresses } = storeToRefs(addressStore)
const { fetchAddresses } = addressStore

const authStore = useAuthStore()
const { user } = storeToRefs(authStore)
const { fetchUserInfo } = authStore

await useAsyncData('account', () => {
  return Promise.all([fetchUserInfo(), fetchAddresses()])
})

useSeoMeta({
  title: 'アカウント',
})
</script>

<template>
  <div class="container mx-auto flex flex-col gap-6 px-4 text-main xl:px-0">
    <!-- ユーザー情報表示エリア -->
    <div class="mt-4">
      <div class="mb-2 text-[16px] font-semibold">アカウント情報</div>
      <template v-if="user">
        <div class="rounded bg-white p-4">
          <dl>
            <div class="py-2 sm:grid sm:grid-cols-3 sm:gap-4">
              <dt>氏名（ふりがな）</dt>
              <dd class="sm:col-span-2">
                {{ `${user.lastname} ${user.firstname}` }}
                {{ `（${user.lastnameKana} ${user.firstnameKana}）` }}
              </dd>
            </div>
            <div class="py-2 sm:grid sm:grid-cols-3 sm:gap-4">
              <dt>ユーザーID</dt>
              <dd class="sm:col-span-2">
                {{ user.username }}
              </dd>
            </div>
            <div class="py-2 sm:grid sm:grid-cols-3 sm:gap-4">
              <dt>メールアドレス</dt>
              <dd class="sm:col-span-2">
                {{ user.email }}
              </dd>
            </div>
          </dl>
        </div>
      </template>
    </div>

    <!-- アドレス帳情報表示エリア -->
    <div>
      <div class="mb-2 text-[16px] font-semibold">アドレス帳</div>
      <div class="flex flex-col gap-4 sm:flex-row">
        <div
          v-for="(address, i) in addresses"
          :key="address.id"
          class="rounded bg-white p-4"
        >
          <div class="mb-2 inline-flex items-center gap-4">
            <div class="font-semibold">アドレス #{{ i + 1 }}</div>
            <div
              v-if="address.isDefault"
              class="inline-block rounded-2xl bg-orange px-2 py-1 text-[12px] text-white"
            >
              規定の住所
            </div>
          </div>
          <dl>
            <div class="py-2 sm:grid sm:grid-cols-3 sm:gap-4">
              <dt>氏名（ふりがな）</dt>
              <dd class="sm:col-span-2">
                {{ `${address.lastname} ${address.firstname}` }}
                {{ `（${address.lastnameKana} ${address.firstnameKana}）` }}
              </dd>
            </div>
            <div class="py-2 sm:grid sm:grid-cols-3 sm:gap-4">
              <dt>住所</dt>
              <dd class="sm:col-span-2">
                <div>〒{{ address.postalCode }}</div>
                {{
                  `${address.prefecture}${address.city}${address.addressLine1}${address.addressLine2}`
                }}
              </dd>
            </div>
            <div class="py-2 sm:grid sm:grid-cols-3 sm:gap-4">
              <dt>電話番号</dt>
              <dd class="sm:col-span-2">
                {{ convertI18nToJapanesePhoneNumber(address.phoneNumber) }}
              </dd>
            </div>
          </dl>
        </div>
      </div>
    </div>
  </div>
</template>
