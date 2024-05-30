<script setup lang="ts">
import { storeToRefs } from 'pinia'
import dayjs from 'dayjs'
import { useAddressStore } from '~/store/address'
import { convertI18nToJapanesePhoneNumber } from '~/lib/phone-number'
import { useAuthStore } from '~/store/auth'
import { useOrderStore } from '~/store/order'
import { getOrderStatusString } from '~/lib/order'
import { priceFormatter } from '~/lib/price'

const router = useRouter()
const route = useRoute()

const addressStore = useAddressStore()
const { addresses, defaultAddress } = storeToRefs(addressStore)
const { fetchAddresses } = addressStore

const authStore = useAuthStore()
const { user } = storeToRefs(authStore)
const { fetchUserInfo, logout } = authStore

const orderStore = useOrderStore()
const { fetchOrderHistoryList } = orderStore
const { orderHistories, total, fetchState } = storeToRefs(orderStore)

// 1ページ当たりに表示する注文履歴数
const orderPagePerItems = ref<number>(20)

// 注文履歴の現在のページ番号
const currentOrderPage = computed<number>(() => {
  return route.query.orderPage ? Number(route.query.orderPage) : 1
})

// 注文履歴のページネーション情報
const orderPagination = computed<{
  limit: number
  offset: number
  pageArray: number[]
}>(() => {
  const totalPage = Math.ceil(total.value / orderPagePerItems.value)
  const pageArray = Array.from({ length: totalPage }, (_, i) => i + 1)

  return {
    limit: orderPagePerItems.value,
    offset: orderPagePerItems.value * (currentOrderPage.value - 1),
    pageArray,
  }
})

// 指定した注文履歴ページへ遷移する
const handleChangeOrderPage = (page: number) => {
  router.push({
    query: {
      ...route.query,
      orderPage: page,
    },
  })
}

// 1ページ当たりに表示するアドレス帳数
const addressPagePerItems = ref<number>(20)

// アドレス帳の現在のページ番号
const currentAddressPage = computed<number>(() => {
  return route.query.addressPage ? Number(route.query.addressPage) : 1
})

// アドレス帳のページネーション情報
const addressPagination = computed<{
  limit: number
  offset: number
  pageArray: number[]
}>(() => {
  const totalPage = Math.ceil(
    addresses.value.length / addressPagePerItems.value,
  )
  const pageArray = Array.from({ length: totalPage }, (_, i) => i + 1)

  return {
    limit: addressPagePerItems.value,
    offset: addressPagePerItems.value * (currentAddressPage.value - 1),
    pageArray,
  }
})

const handleChangeAddressPage = (page: number) => {
  router.push({
    query: {
      ...route.query,
      addressPage: page,
    },
  })
}

watch(currentOrderPage, () => {
  fetchOrderHistoryList(
    orderPagination.value.limit,
    orderPagination.value.offset,
  )
})

watch(currentAddressPage, () => {
  fetchAddresses(addressPagination.value.limit, addressPagination.value.offset)
})

await useAsyncData('account', () => {
  return Promise.all([
    fetchUserInfo(),
    fetchAddresses(
      addressPagination.value.limit,
      addressPagination.value.offset,
    ),
    fetchOrderHistoryList(
      orderPagination.value.limit,
      orderPagination.value.offset,
    ),
  ])
})

const handleClickLogout = async () => {
  await logout()
  router.push('/')
}

useSeoMeta({
  title: 'アカウント',
})

definePageMeta({
  middleware: 'auth',
})
</script>

<template>
  <div
    class="container mx-auto grid gap-6 px-4 text-main md:grid-cols-2 xl:px-0"
  >
    <!-- ユーザー情報表示エリア -->
    <div class="flex flex-col gap-4">
      <div>
        <div class="rounded bg-white p-4 md:p-16">
          <div class="mb-6 text-[16px] font-semibold">
            アカウント情報
          </div>
          <template v-if="user">
            <dl
              class="divide-y divide-dashed divide-typography border-y border-dashed border-typography [&>div]:py-4 [&_dt]:text-typography"
            >
              <div class="flex items-center">
                <div class="grow items-center sm:grid sm:grid-cols-3 sm:gap-4">
                  <dt>プロフィール写真</dt>
                  <dd class="sm:col-span-2">
                    <template v-if="user.thumbnailUrl">
                      <nuxt-img
                        provider="cloudFront"
                        :src="user.thumbnailUrl"
                        :alt="`${user.username}のプロフィール写真`"
                        class="h-14 w-14 rounded-full"
                        width="64px"
                        height="64px"
                        format="webp"
                      />
                    </template>
                    <template v-else>
                      <the-account-icon
                        class="h-14 w-14"
                        color="white"
                      />
                    </template>
                  </dd>
                </div>
                <nuxt-link
                  to="/account/edit/thumbnail"
                  class="whitespace-nowrap bg-main px-4 py-1 text-white"
                >
                  変更
                </nuxt-link>
              </div>

              <div class="flex items-center">
                <div class="grow sm:grid sm:grid-cols-3 sm:gap-4">
                  <dt>ユーザー名</dt>
                  <dd class="sm:col-span-2">
                    {{ user.username }}
                  </dd>
                </div>
                <nuxt-link
                  to="/account/edit/username"
                  class="whitespace-nowrap bg-main px-4 py-1 text-white"
                >
                  変更
                </nuxt-link>
              </div>

              <div class="flex items-center">
                <div class="grow sm:grid sm:grid-cols-3 sm:gap-4">
                  <dt>ユーザーID</dt>
                  <dd class="sm:col-span-2">
                    {{ user.accountId }}
                  </dd>
                </div>
                <nuxt-link
                  to="/account/edit/id"
                  class="whitespace-nowrap bg-main px-4 py-1 text-white"
                >
                  変更
                </nuxt-link>
              </div>

              <div class="flex items-center">
                <div class="grow sm:grid sm:grid-cols-3 sm:gap-4">
                  <dt>パスワード</dt>
                  <dd class="sm:col-span-2">
                    ********
                  </dd>
                </div>
                <nuxt-link
                  to="/account/edit/password"
                  class="whitespace-nowrap bg-main px-4 py-1 text-white"
                >
                  変更
                </nuxt-link>
              </div>

              <div class="flex items-center">
                <div class="grow sm:grid sm:grid-cols-3 sm:gap-4">
                  <dt>メールアドレス</dt>
                  <dd class="sm:col-span-2">
                    {{ user.email }}
                  </dd>
                </div>
                <!-- 一時的に選択不可能にする
                <button
                  class="cursor-not-allowed bg-main/70 px-2 py-1 text-white"
                >
                  準備中
                </button>
                <nuxt-link
                  v-if="false"
                  to="/account/edit/email"
                  class="whitespace-nowrap bg-main px-4 py-1 text-white"
                >
                  変更
                </nuxt-link>
              -->
              </div>

              <div class="flex items-center">
                <div class="grow sm:grid sm:grid-cols-3 sm:gap-4">
                  <dt>メール配信</dt>
                  <dd class="sm:col-span-2">
                    {{
                      user.notificationEnabled
                        ? 'お知らせメールを受け取る'
                        : 'お知らせメールを受け取らない'
                    }}
                  </dd>
                </div>
                <nuxt-link
                  to="/account/edit/notification"
                  class="whitespace-nowrap bg-main px-4 py-1 text-white"
                >
                  変更
                </nuxt-link>
              </div>
            </dl>
          </template>
        </div>
      </div>

      <!-- 基本お届け先情報表示エリア -->
      <div
        v-if="defaultAddress"
        class="rounded bg-white p-4 md:p-16"
      >
        <div class="mb-6 text-[16px] font-semibold">
          基本お届け先情報
        </div>
        <dl class="[&_dt]:text-typography">
          <div class="py-2 sm:grid sm:grid-cols-3 sm:gap-4">
            <dt>氏名（ふりがな）</dt>
            <dd class="sm:col-span-2">
              {{ `${defaultAddress.lastname} ${defaultAddress.firstname}` }}
              {{
                `（${defaultAddress.lastnameKana} ${defaultAddress.firstnameKana}）`
              }}
            </dd>
          </div>
          <div class="py-2 sm:grid sm:grid-cols-3 sm:gap-4">
            <dt>住所</dt>
            <dd class="sm:col-span-2">
              <div>〒{{ defaultAddress.postalCode }}</div>
              {{
                `${defaultAddress.prefecture}${defaultAddress.city}${defaultAddress.addressLine1}${defaultAddress.addressLine2}`
              }}
            </dd>
          </div>
          <div class="py-2 sm:grid sm:grid-cols-3 sm:gap-4">
            <dt>電話番号</dt>
            <dd class="sm:col-span-2">
              {{ convertI18nToJapanesePhoneNumber(defaultAddress.phoneNumber) }}
            </dd>
          </div>
        </dl>
      </div>

      <!-- お届け先情報一覧表示エリア -->
      <div class="rounded bg-white p-4 md:p-16">
        <div class="mb-6 text-[16px] font-semibold">
          お届け先情報
        </div>
        <div
          class="flex flex-col gap-6 divide-y divide-main border-y border-main"
        >
          <div
            v-for="(address, i) in addresses"
            :key="address.id"
            class="py-4"
          >
            <div class="mb-2 inline-flex items-center gap-4">
              <div class="font-semibold">
                お届け先 #{{ i + 1 }}
              </div>
            </div>
            <dl class="[&_dt]:text-typography">
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
        <template v-if="addresses.length > 0">
          <the-pagination
            class="mt-8"
            :current-page="currentAddressPage"
            :page-array="addressPagination.pageArray"
            @change-page="handleChangeAddressPage"
          />
        </template>
        <template v-else>
          <div class="mt-4 text-center">
            登録されている住所がありません
          </div>
        </template>
      </div>
    </div>

    <!-- 注文履歴一覧表示エリア -->
    <div class="bg-white p-4 md:p-16">
      <div class="mb-6 text-[16px] font-semibold">
        注文履歴（{{ total }}件）
      </div>
      <div class="flex flex-col">
        <template v-if="fetchState.isLoading" />
        <div class="divide-y divide-main border-y border-main">
          <div
            v-for="order in orderHistories"
            :key="order.id"
            class="flex flex-col py-4 [&_dt]:text-typography"
          >
            <div class="flex items-center justify-between tracking-[10%]">
              <div>
                <div class="py-2 sm:grid sm:grid-cols-3 sm:gap-4">
                  <dt>注文日時</dt>
                  <dd class="text-[18px] font-semibold sm:col-span-2">
                    {{
                      dayjs
                        .unix(order.payment.orderedAt)
                        .format('YYYY/MM/DD HH:mm')
                    }}
                  </dd>
                </div>

                <div class="py-2 sm:grid sm:grid-cols-3 sm:gap-4">
                  <dt>ステータス</dt>
                  <dd class="text-[18px] font-semibold sm:col-span-2">
                    {{ getOrderStatusString(order.status) }}
                  </dd>
                </div>
              </div>
              <div>
                <nuxt-link
                  :to="`/account/orders/${order.id}`"
                  class="block bg-main px-4 py-2 text-white"
                >
                  詳細を見る
                </nuxt-link>
              </div>
            </div>
            <hr class="my-4 border-b border-dashed">
            <div />

            <div
              class="flex flex-col gap-4 sm:grid sm:grid-cols-3 md:grid-cols-4"
            >
              <div class="flex items-center justify-center sm:aspect-square">
                <nuxt-img
                  provider="cloudFront"
                  class="block max-w-[80px] rounded-full"
                  width="80px"
                  format="webp"
                  :src="order.coordinator?.thumbnailUrl"
                  :alt="`${order.coordinator?.username}のサムネイル`"
                />
              </div>
              <dl class="col-span-2 md:col-span-3">
                <div class="py-2 sm:grid sm:grid-cols-3 sm:gap-4">
                  <dt>注文ID</dt>
                  <dd class="line-clamp-1 sm:col-span-2">
                    {{ order.id }}
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
              </dl>
            </div>
            <div class="mt-2 text-right text-[14px]" />
          </div>
        </div>
      </div>
      <template v-if="orderHistories.length > 0">
        <!-- ページネーション -->
        <the-pagination
          class="mt-8"
          :current-page="currentOrderPage"
          :page-array="orderPagination.pageArray"
          @change-page="handleChangeOrderPage"
        />
      </template>
      <template v-else>
        <div class="mt-4 text-center">
          注文履歴がありません
        </div>
      </template>
    </div>
  </div>

  <div class="container mx-auto mt-16 text-right xl:px-0">
    <button
      class="underline"
      @click="handleClickLogout"
    >
      ログアウト
    </button>
  </div>
</template>
