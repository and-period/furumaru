<script setup lang="ts">
import { storeToRefs } from 'pinia'
import { useEventBus } from '@vueuse/core'
import { useProductStore } from '~/store/product'
import { useShoppingCartStore } from '~/store/shopping'
import { ProductStatus } from '~/types/api'
import type { Snackbar } from '~/types/props'
import type { I18n } from '~/types/locales'
import { useProductReviewStore } from '~/store/productReview'

const i18n = useI18n()

const route = useRoute()

const productStore = useProductStore()
const shoppingCartStore = useShoppingCartStore()

const { fetchProduct } = productStore
const { addCart } = shoppingCartStore

const { product, productFetchState } = storeToRefs(productStore)

const productReviewStore = useProductReviewStore()
const { fetchReviews } = productReviewStore
const { reviews } = storeToRefs(productReviewStore)

const { emit } = useEventBus('add-to-cart')

const dt = (str: keyof I18n['items']['details']) => {
  return i18n.t(`items.details.${str}`)
}

const itemThumbnailAlt = computed<string>(() => {
  return i18n.t('items.list.itemThumbnailAlt', {
    itemName: product.value.name,
  })
})

const expirationDateText = computed<string>(() => {
  return i18n.t('items.details.expirationDateText', {
    expirationDate: product.value.expirationDate,
  })
})

const id = computed<string>(() => {
  const ids = route.params.id
  if (Array.isArray(ids)) {
    return ids[0]
  }
  else {
    return ids
  }
})

const snackbarItems = ref<Snackbar[]>([])

const quantity = ref<number>(1)

const priceString = computed<string>(() => {
  if (product.value) {
    return new Intl.NumberFormat('ja-JP', {
      style: 'currency',
      currency: 'JPY',
    }).format(product.value.price)
  }
  else {
    return ''
  }
})

const canAddCart = computed<boolean>(() => {
  // データ取得に失敗していたらfalseを返す
  if (status.value === 'success') {
    if (product.value) {
      return (
        product.value.status === ProductStatus.FOR_SALE
        && product.value.inventory > 0
      )
    }
    else {
      return false
    }
  }
  else {
    return false
  }
})

const handleClickAddCartButton = () => {
  addCart({
    productId: id.value,
    quantity: quantity.value,
  })
  emit('add-to-cart')
}

const getDeliveryType = (type: number) => {
  switch (type) {
    case 1:
      return dt('deliveryTypeStandard')
    case 2:
      return dt('deliveryTypeRefrigerated')
    case 3:
      return dt('deliveryTypeFrozen')
    default:
      return ''
  }
}

const getStorageMethodType = (type: number) => {
  switch (type) {
    case 0:
      return dt('storageTypeUnknown')
    case 1:
      return dt('storageTypeRoomTemperature')
    case 2:
      return dt('storageTypeCoolAndDark')
    case 3:
      return dt('storageTypeRefrigerated')
    case 4:
      return dt('storageTypeFrozen')
  }
}

const title = computed<string>(() => product.value.name)

const selectedMediaIndex = ref<number>(-1)

const selectMediaSrcUrl = computed<string>(() => {
  return selectedMediaIndex.value === -1
    ? product.value.thumbnailUrl
    : product.value.media[selectedMediaIndex.value].url
})

const handleClickMediaItem = (index: number) => {
  selectedMediaIndex.value = index
}

const { status, error } = useAsyncData(`product-${id.value}`, async () => {
  await fetchProduct(id.value)
  return true
})

useAsyncData(`reviews-${id.value}`, async () => {
  await fetchReviews(id.value)
  return true
})

useSeoMeta({
  title,
})
</script>

<template>
  <template
    v-for="(snackbarItem, i) in snackbarItems"
    :key="i"
  >
    <the-snackbar
      v-model:is-show="snackbarItem.isShow"
      :text="snackbarItem.text"
    />
  </template>

  <!-- ロード状態 -->
  <template v-if="status === 'pending'">
    <div
      class="animate-pulse bg-white px-[112px] pb-6 pt-[40px] text-main md:grid md:grid-cols-2"
    >
      <div class="w-full">
        <div class="mx-auto aspect-square h-[500px] w-[500px] bg-slate-100" />
      </div>
      <div class="flex w-full flex-col gap-4">
        <div class="h-[24px] w-[80%] rounded-md bg-slate-100" />
        <div class="h-[24px] w-[60%] rounded-md bg-slate-100" />
      </div>
    </div>
  </template>

  <!-- エラーの場合の表示 -->
  <template v-if="status === 'error'">
    <the-alert>
      {{ error?.message }}
    </the-alert>
  </template>

  <!-- 商品情報の表示 -->
  <template v-if="status == 'success' && !productFetchState.isLoading && product.thumbnail">
    <div class="bg-white w-full">
      <div
        class="gap-10 px-4 pb-6 pt-[40px] text-main md:grid md:grid-cols-2 lg:px-[112px] w-full max-w-[1440px] mx-auto"
      >
        <div class="mx-auto w-full max-w-[100%]">
          <div class="flex aspect-square h-full w-full justify-center">
            <template v-if="selectMediaSrcUrl.endsWith('.mp4')">
              <the-item-video-player :src="selectMediaSrcUrl" />
            </template>
            <template v-else>
              <nuxt-img
                provider="cloudFront"
                fill="contain"
                class="block h-full w-full object-contain border"
                :src="selectMediaSrcUrl"
                :alt="itemThumbnailAlt"
              />
            </template>
          </div>
          <div
            class="hidden-scrollbar mt-2 grid w-full grid-flow-col justify-start gap-2 overflow-x-scroll"
          >
            <template
              v-for="(m, i) in product.media"
              :key="i"
            >
              <template v-if="m.url.endsWith('.mp4')">
                <div
                  class="aspect-square w-[72px] h-[72px] cursor-pointer border relative"
                  @click="handleClickMediaItem(i)"
                >
                  <video
                    :src="`${m.url}#t=0.1`"
                    class="aspect-square w-full object-contain h-full"
                  />
                  <div
                    class="absolute h-6 w-6 bottom-0 right-0 p-1 bg-main/80 rounded-full text-white"
                  >
                    <svg
                      xmlns="http://www.w3.org/2000/svg"
                      fill="none"
                      viewBox="0 0 24 24"
                      stroke-width="1.5"
                      stroke="currentColor"
                      class="size-6"
                    >
                      <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        d="m15.75 10.5 4.72-4.72a.75.75 0 0 1 1.28.53v11.38a.75.75 0 0 1-1.28.53l-4.72-4.72M4.5 18.75h9a2.25 2.25 0 0 0 2.25-2.25v-9a2.25 2.25 0 0 0-2.25-2.25h-9A2.25 2.25 0 0 0 2.25 7.5v9a2.25 2.25 0 0 0 2.25 2.25Z"
                      />
                    </svg>
                  </div>
                </div>
              </template>
              <template v-else>
                <nuxt-img
                  width="72px"
                  fill="contain"
                  provider="cloudFront"
                  :src="m.url"
                  :alt="`${itemThumbnailAlt}_${i}`"
                  class="aspect-square w-[72px] h-[72px] cursor-pointer object-contain border block"
                  @click="handleClickMediaItem(i)"
                />
              </template>
            </template>
          </div>
        </div>

        <div class="mt-4 flex w-full flex-col gap-4">
          <div
            class="break-words text-[16px] tracking-[1.6px] md:text-[24px] md:tracking-[2.4px]"
          >
            {{ product.name }}
          </div>

          <div
            v-if="product.producer"
            class="flex flex-col leading-[32px] md:mt-4"
          >
            <div
              class="text-[14px] tracking-[1.4px] md:text-[16px] md:tracking-[1.6px]"
            >
              {{ dt("producerLabel") }}:
              <a
                href="#"
                class="font-bold underline"
              >
                {{ product.producer.username }}
              </a>
            </div>
            <div class="text-[12px] tracking-[1.4px] md:text-[14px]">
              {{ product.originPrefecture }} {{ product.originCity }}
            </div>
          </div>

          <!-- 評価情報 -->
          <div>
            <div class="flex items-center gap-3">
              <div class="inline-flex items-center">
                <the-rating-star :rate="product.rate.average" />
                <p class="ms-2 text-sm font-bold text-main">
                  {{ product.rate.average }}
                </p>
              </div>
              <a
                href="#reviews"
                class="text-sm font-medium text-main underline hover:no-underline "
              >
                {{ product.rate.count }} {{ dt('reviewCountLabel') }}
              </a>
            </div>
          </div>

          <div
            v-if="product && product.recommendedPoint1"
            class="w-full rounded-2xl bg-base px-[20px] py-[28px] text-main md:mt-8"
          >
            <p
              class="mb-[12px] text-[12px] font-medium tracking-[1.4px] md:text-[14px]"
            >
              {{ dt("highlightsLabel") }}
            </p>
            <ol
              class="recommend-list flex flex-col divide-y divide-dashed divide-main px-[4px] pl-[24px]"
            >
              <li
                v-if="product.recommendedPoint1"
                class="py-3 text-[14px] font-medium md:text-[16px]"
              >
                {{ product.recommendedPoint1 }}
              </li>
              <li
                v-if="product.recommendedPoint2"
                class="py-3 text-[14px] font-medium md:text-[16px]"
              >
                {{ product.recommendedPoint2 }}
              </li>
              <li
                v-if="product.recommendedPoint3"
                class="py-3 text-[14px] font-medium md:text-[16px]"
              >
                {{ product.recommendedPoint3 }}
              </li>
            </ol>
          </div>

          <div>
            <div class="flex items-end justify-end md:justify-start">
              <div
                class="text-[24px] font-bold md:mt-[60px] md:flex md:flex-row md:text-[32px]"
              >
                {{ priceString }}
              </div>
              <p class="pb-1 pl-2 text-[12px] md:text-[16px]">
                {{ dt("itemPriceTaxIncludedText") }}
              </p>
            </div>

            <div
              v-if="product"
              class="mt-4 inline-flex items-center md:mt-8"
            >
              <label class="mr-2 block text-[14px] md:text-[16px]">{{
                dt("quantityLabel")
              }}</label>
              <select
                v-model="quantity"
                class="h-full border-[1px] border-main px-2"
                :disabled="!product.hasStock"
              >
                <option
                  v-for="(_, i) in Array.from({
                    length: product.inventory < 10 ? product.inventory : 10,
                  })"
                  :key="i + 1"
                  :value="i + 1"
                >
                  {{ i + 1 }}
                </option>
              </select>
            </div>
          </div>

          <button
            class="mt-2 w-full bg-orange py-4 text-center text-white transition-all duration-200 ease-in-out hover:shadow-lg active:scale-95 disabled:cursor-not-allowed disabled:bg-main/60 md:mt-8"
            :disabled="!canAddCart"
            @click="handleClickAddCartButton"
          >
            {{ dt("addToCartText") }}
          </button>

          <div class="mt-4 inline-flex gap-4">
            <span
              v-for="productTag in product.productTags"
              :key="productTag?.id"
              class="rounded-2xl border border-main px-4 py-1 text-[14px] md:text-[16px]"
            >
              {{ productTag?.name }}
            </span>
          </div>
        </div>

        <div class="col-span-2 mt-[40px] pb-10 md:mt-[80px] md:pb-16">
          <article
            class="text-[14px] leading-[32px] tracking-[1.4px] md:text-[16px] md:tracking-[1.6px] whitespace-pre-wrap"
            v-text="product.description"
          />
        </div>

        <div
          class="col-span-2 flex flex-col divide-y divide-dashed divide-main border-y border-dashed border-main text-[14px] md:text-[16px]"
        >
          <div class="grid grid-cols-5 py-4">
            <p class="col-span-2 md:col-span-1">
              {{ dt("expirationDateLabel") }}
            </p>
            <p class="col-span-3 md:col-span-4">
              {{ expirationDateText }}
            </p>
          </div>
          <div class="grid grid-cols-5 py-4">
            <p class="col-span-2 md:col-span-1">
              {{ dt("weightLabel") }}
            </p>
            <p class="col-span-3 md:col-span-4">
              {{ product.weight }}kg
            </p>
          </div>
          <div class="grid grid-cols-5 py-4">
            <p class="col-span-2 md:col-span-1">
              {{ dt("deliveryTypeLabel") }}
            </p>
            <p class="col-span-3 md:col-span-4">
              {{ getDeliveryType(product.deliveryType) }}
            </p>
          </div>
          <div class="grid grid-cols-5 py-4">
            <p class="col-span-2 md:col-span-1">
              {{ dt("storageTypeLabel") }}
            </p>
            <p class="col-span-3 md:col-span-4">
              {{ getStorageMethodType(product.storageMethodType) }}
            </p>
          </div>
        </div>
      </div>
    </div>
  </template>

  <template v-if="product.producer">
    <div class="w-full">
      <div class="mx-auto mt-[40px] w-full px-4 xl:px-28 max-w-[1440px]">
        <div
          class="flex w-full flex-col rounded-3xl bg-white px-8 py-10 text-main xl:px-16"
        >
          <p
            class="mx-auto w-full rounded-full bg-base py-2 text-center text-[14px] font-bold text-main md:text-[16px]"
          >
            {{ dt("producerInformationTitle") }}
          </p>

          <div
            class="mt-[64px] flex w-full flex-col gap-4 md:flex-row lg:gap-10"
          >
            <div
              class="flex min-w-max flex-col items-center justify-center gap-4 md:flex-row"
            >
              <template v-if="product.producer.thumbnailUrl">
                <nuxt-img
                  provider="cloudFront"
                  sizes="96px md:120px"
                  fit="cover"
                  :src="product.producer.thumbnailUrl"
                  :alt="`${product.producer.username}`"
                  class="mx-auto block aspect-square w-[96px] rounded-full md:w-[120px] object-cover"
                />
              </template>
              <template v-else>
                <img
                  class="mx-auto block aspect-square w-[96px] rounded-full md:w-[120px] object-cover"
                  src="/img/account.png"
                >
              </template>
              <div
                class="flex min-w-max grow flex-col items-center gap-2 md:items-start md:gap-2 md:whitespace-nowrap"
              >
                <p class="text-sm font-[500] tracking-[1.4px]">
                  {{ `${product.originPrefecture} ${product.originCity}` }}
                </p>
                <div
                  class="flex flex-row items-baseline grow text-[16px] tracking-[1.4px] md:text-[24px]"
                >
                  <p class="mr-1 text-[14px] font-medium">
                    {{ dt("producerLabel") }}
                  </p>
                  <p>{{ product.producer.username }}</p>
                </div>
              </div>
            </div>
            <div
              class="pt-2 text-[14px] tracking-[1.4px] md:pt-0 md:text-[16px] md:tracking-[1.6px]"
            >
              {{ product.producer.profile }}
            </div>
          </div>

          <!--
        <div class="mt-4 w-full text-right">
          <button class="inline-flex items-center">
            詳しく見る
            <the-right-arrow-icon class="ml-2 h-[12px] w-[12px]" />
          </button>
        </div>
        -->
        </div>
      </div>
    </div>
  </template>

  <div
    id="reviews"
    class="w-full"
  >
    <div class="mx-auto mt-[40px] w-full px-4 xl:px-28 max-w-[1440px]">
      <div
        class="flex w-full flex-col rounded-3xl bg-white px-8 py-10 text-main xl:px-16"
      >
        <p
          class="mx-auto w-full rounded-full bg-base py-2 text-center text-[14px] font-bold text-main md:text-[16px]"
        >
          {{ dt("reviewLabel") }}
        </p>

        <div class="mt-[32px] flex flex-col divide-y">
          <div
            v-if="reviews.length === 0"
            class="text-center py-4"
          >
            {{ dt("noReviewText") }}
          </div>
          <div
            v-for="review, i in reviews"
            :key="review.id"
            class="flex flex-col gap-2 py-4"
          >
            <div class="flex gap-2">
              <div v-if="review.thumbnailUrl">
                <nuxt-img
                  provider="cloudFront"
                  sizes="48px"
                  fit="cover"
                  :src="review.thumbnailUrl"
                  :alt="`${review.username}`"
                  class="block aspect-square w-[48px] rounded-full object-cover"
                />
              </div>
              <div v-else>
                <the-account-icon />
              </div>
              {{ review.username }}
            </div>
            <div>
              <div class="flex gap-2 items-center mb-2">
                <the-rating-star
                  :id="`review-${i}`"
                  :rate="review.rate"
                />
                <div class="font-semibold">
                  {{ review.title }}
                </div>
              </div>
              <div>
                {{ review.comment }}
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.recommend-list {
  list-style: none;
  counter-reset: li;
  position: relative;
}

.recommend-list li {
  padding-left: 16px;
}

.recommend-list li::before {
  content: counter(li);
  counter-increment: li;
  position: absolute;
  left: 0;
  background-color: #604c3f;
  color: #f9f6ea;
  border-radius: 100%;
  width: 24px;
  height: 24px;
  text-align: center;
}
</style>
