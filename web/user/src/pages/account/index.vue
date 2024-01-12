<script setup lang="ts">
import { storeToRefs } from 'pinia'
import dayjs from 'dayjs'
import { useAdressStore } from '~/store/address'
import { convertI18nToJapanesePhoneNumber } from '~/lib/phone-number'
import { useAuthStore } from '~/store/auth'
import { useOrderStore } from '~/store/order'
import {
  getOrderStatusString,
  getOperationResultFromOrderStatus,
} from '~/lib/order'
import { priceFormatter } from '~/lib/price'
import type { OrderStatus } from '~/types/api'

const router = useRouter()

const addressStore = useAdressStore()
const { addresses } = storeToRefs(addressStore)
const { fetchAddresses } = addressStore

const authStore = useAuthStore()
const { user } = storeToRefs(authStore)
const { fetchUserInfo, logout } = authStore

const orderStore = useOrderStore()
const { fetchOrderHsitoryList } = orderStore
const { orderHistories, total } = storeToRefs(orderStore)

await useAsyncData('account', () => {
  return Promise.all([
    fetchUserInfo(),
    fetchAddresses(),
    fetchOrderHsitoryList(),
  ])
})

const handleClickLogout = async () => {
  await logout()
  router.push('/')
}

const getClassNameFromOrderStatus = (status: OrderStatus) => {
  const operationResult = getOperationResultFromOrderStatus(status)
  switch (operationResult) {
    case 'success':
      return 'bg-green'
    case 'failed':
      return 'bg-red-700'
    case 'canceled':
      return 'bg-yellow-600'
    default:
      return ''
  }
}

useSeoMeta({
  title: 'アカウント',
})

definePageMeta({
  middleware: 'auth',
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
              <dt>表示名</dt>
              <dd class="sm:col-span-2">
                {{ user.username }}
              </dd>
            </div>
            <div class="py-2 sm:grid sm:grid-cols-3 sm:gap-4">
              <dt>氏名（ふりがな）</dt>
              <dd class="sm:col-span-2">
                {{ `${user.lastname} ${user.firstname}` }}
                {{ `（${user.lastnameKana} ${user.firstnameKana}）` }}
              </dd>
            </div>
            <div class="py-2 sm:grid sm:grid-cols-3 sm:gap-4">
              <dt>アカウントID</dt>
              <dd class="sm:col-span-2">
                {{ user.accountId }}
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

    <!-- 注文履歴一覧表示エリア -->
    <div>
      <div class="mb-2 text-[16px] font-semibold">
        注文履歴（{{ total }}件）
      </div>
      <div class="grid gap-4 lg:grid-cols-2 xl:grid-cols-3">
        <div
          v-for="order in orderHistories"
          :key="order.id"
          class="flex flex-col bg-white p-4"
        >
          <div
            class="flex flex-col gap-4 sm:grid sm:grid-cols-3 md:grid-cols-4"
          >
            <div class="flex items-center justify-center sm:aspect-square">
              <img
                class="block max-w-[80px] rounded-full"
                :src="order.coordinator?.thumbnailUrl"
                :alt="`${order.coordinator?.username}のサムネイル`"
              />
            </div>
            <dl class="col-span-2 md:col-span-3">
              <div class="py-2 sm:grid sm:grid-cols-3 sm:gap-4">
                <dt>注文ID</dt>
                <dd class="sm:col-span-2">
                  {{ order.id }}
                </dd>
              </div>
              <div class="py-2 sm:grid sm:grid-cols-3 sm:gap-4">
                <dt>注文日時</dt>
                <dd class="sm:col-span-2">
                  {{
                    dayjs
                      .unix(order.payment.orderedAt)
                      .format('YYYY/MM/DD HH:mm')
                  }}
                </dd>
              </div>
              <div
                v-if="order.coordinator"
                class="py-2 sm:grid sm:grid-cols-3 sm:gap-4"
              >
                <dt>マルシェ名</dt>
                <dd class="sm:col-span-2">
                  {{ order.coordinator.marcheName }}
                </dd>
              </div>
              <div class="py-2 sm:grid sm:grid-cols-3 sm:gap-4">
                <dt>注文品数</dt>
                <dd>
                  {{ order.items.length }}
                </dd>
              </div>

              <div class="py-2 sm:grid sm:grid-cols-3 sm:gap-4">
                <dt>支払い金額</dt>
                <dd>
                  {{ priceFormatter(order.payment.total) }}
                </dd>
              </div>

              <div class="py-2 sm:grid sm:grid-cols-3 sm:gap-4">
                <dt>ステータス</dt>
                <dd class="sm:col-span-2">
                  <span
                    :class="getClassNameFromOrderStatus(order.status)"
                    class="inline-block rounded-lg px-2 py-1 text-[14px] text-white sm:mt-1"
                  >
                    {{ getOrderStatusString(order.status) }}
                  </span>
                </dd>
              </div>
            </dl>
          </div>
          <div class="mt-2 text-right text-[14px]">
            <nuxt-link :to="`/account/orders/${order.id}`" class="underline">
              詳細を見る
            </nuxt-link>
          </div>
        </div>
      </div>
    </div>

    <!-- アドレス帳情報表示エリア -->
    <div>
      <div class="mb-2 text-[16px] font-semibold">アドレス帳</div>
      <div class="grid gap-4 sm:grid-cols-2">
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

    <div class="text-right">
      <button class="underline" @click="handleClickLogout">ログアウト</button>
    </div>
  </div>
</template>
