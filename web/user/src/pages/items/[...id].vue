<script setup lang="ts">
import { storeToRefs } from 'pinia'
import { useProductStore } from '~/store/product'
import { useShoppingCartStore } from '~/store/shopping'
import type { Snackbar } from '~/types/props'

const route = useRoute()

const productStore = useProductStore()
const shoppingCartStore = useShoppingCartStore()

const { fetchProduct } = productStore
const { addCart } = shoppingCartStore

const { product, productFetchState } = storeToRefs(productStore)

const id = computed<string>(() => {
  const ids = route.params.id
  if (Array.isArray(ids)) {
    return ids[0]
  } else {
    return ids
  }
})

fetchProduct(id.value)

const snackbarItems = ref<Snackbar[]>([])

const quantity = ref<number>(1)

const priceString = computed<string>(() => {
  if (product.value) {
    return new Intl.NumberFormat('ja-JP', {
      style: 'currency',
      currency: 'JPY',
    }).format(product.value.price)
  } else {
    return ''
  }
})

const handleClickAddCartButton = () => {
  addCart({
    productId: id.value,
    quantity: quantity.value,
  })
  snackbarItems.value.push({
    text: `買い物カゴに「${product.value.name}」を追加しました`,
    isShow: true,
  })
}

const getDeliveryType = (type: number) => {
  switch (type) {
    case 1:
      return '通常便'
    case 2:
      return '冷蔵便'
    case 3:
      return '冷凍便'
    default:
      return ''
  }
}

const getStorageMethodType = (type: number) => {
  switch (type) {
    case 0:
      return '不明'
    case 1:
      return '常温保存'
    case 2:
      return '冷暗所保存'
    case 3:
      return '冷蔵保存'
    case 4:
      return '冷凍保存'
  }
}

const title = computed<string>(() => product.value.name)

useSeoMeta({
  title,
})
</script>

<template>
  <template v-for="(snackbarItem, i) in snackbarItems" :key="i">
    <the-snackbar
      v-model:is-show="snackbarItem.isShow"
      :text="snackbarItem.text"
    />
  </template>

  <template v-if="productFetchState.isLoading">
    <div
      class="animate-pulse bg-white px-[112px] pb-6 pt-[40px] text-main md:grid md:grid-cols-2"
    >
      <div class="w-full">
        <div
          class="mx-auto aspect-square h-[500px] w-[500px] bg-slate-100"
        ></div>
      </div>
      <div class="flex w-full flex-col gap-4">
        <div class="h-[24px] w-[80%] rounded-md bg-slate-100"></div>
        <div class="h-[24px] w-[60%] rounded-md bg-slate-100"></div>
      </div>
    </div>
  </template>

  <template v-if="!productFetchState.isLoading && product.thumbnail">
    <div
      class="bg-white px-4 pb-6 pt-[40px] text-main md:grid md:grid-cols-2 md:px-[112px]"
    >
      <div class="w-full">
        <div class="mx-auto aspect-square w-full md:h-[500px] md:w-[500px]">
          <img
            class="w-full"
            :src="product.thumbnail.url"
            :alt="`${product.name}のサムネイル画像`"
          />
        </div>
      </div>

      <div class="flex w-full flex-col gap-4">
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
            生産者:
            <a href="#" class="font-bold underline">
              {{ product.producer.username }}
            </a>
          </div>
          <div class="text-[12px] tracking-[1.4px] md:text-[14px]">
            {{ product.originPrefecture }} {{ product.originCity }}
          </div>
        </div>

        <div
          v-if="product && product.recommendedPoint1"
          class="w-full rounded-2xl bg-base px-[20px] py-[28px] text-main md:mt-8"
        >
          <p
            class="mb-[12px] text-[12px] font-medium tracking-[1.4px] md:text-[14px]"
          >
            おすすめポイント
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
            <p class="pb-1 pl-2 text-[12px] md:text-[16px]">(税込)</p>
          </div>

          <div v-if="product" class="mt-4 inline-flex items-center md:mt-8">
            <label class="mr-2 block text-[14px] md:text-[16px]">数量</label>
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
          class="mt-2 w-full bg-main py-4 text-center text-white md:mt-8"
          @click="handleClickAddCartButton"
        >
          買い物カゴに入れる
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
          class="text-[14px] leading-[32px] tracking-[1.4px] md:text-[16px] md:tracking-[1.6px]"
          v-html="product.description"
        />
      </div>

      <div
        class="col-span-2 flex flex-col divide-y divide-dashed divide-main border-y border-dashed border-main text-[14px] md:text-[16px]"
      >
        <div class="grid grid-cols-5 py-4">
          <p class="col-span-2 md:col-span-1">・賞味期限</p>
          <p class="col-span-3 md:col-span-4">
            発送日より{{ product.expirationDate }}日
          </p>
        </div>
        <div class="grid grid-cols-5 py-4">
          <p class="col-span-2 md:col-span-1">・内容量</p>
          <p class="col-span-3 md:col-span-4">{{ product.weight }}kg</p>
        </div>
        <div class="grid grid-cols-5 py-4">
          <p class="col-span-2 md:col-span-1">・配送方法</p>
          <p class="col-span-3 md:col-span-4">
            {{ getDeliveryType(product.deliveryType) }}
          </p>
        </div>
        <div class="grid grid-cols-5 py-4">
          <p class="col-span-2 md:col-span-1">・保存方法</p>
          <p class="col-span-3 md:col-span-4">
            {{ getStorageMethodType(product.storageMethodType) }}
          </p>
        </div>
      </div>
    </div>
  </template>

  <div class="mx-auto mt-[40px] w-full px-[16px] md:mt-[80px] md:w-[1216px]">
    <div class="w-full rounded-3xl bg-white">
      <div class="px-[16px] pt-10 md:px-[64px]">
        <p
          class="mx-auto rounded-full bg-base py-2 text-center text-[14px] font-bold text-main md:text-[16px]"
        >
          この商品の生産者
        </p>
      </div>
      <div
        v-if="product.producer"
        class="mx-auto px-[16px] pt-16 md:grid md:grid-cols-8 md:px-[64px]"
      >
        <img
          :src="product.producer.thumbnailUrl"
          :alt="`${product.producer.username}`"
          class="mx-auto aspect-square w-[120px] rounded-full"
        />
        <div class="text-main md:col-span-2 md:ml-4">
          <div class="test-sm flex justify-center gap-3 pt-4 md:gap-4">
            <p class="text-sm font-[500] tracking-[1.4px]">
              {{ product.originPrefecture }}
            </p>
            <p class="text-sm font-[500] tracking-[1.4px]">
              {{ product.originCity }}
            </p>
          </div>
          <div class="flex items-end justify-center pt-4">
            <p class="text-[14px] font-[500] tracking-[1.4px]">生産者</p>
            <p class="pl-6 text-[16px] tracking-[1.4px] md:text-[24px]">{{ product.producer.username }}</p>
          </div>
        </div>
        <div class="break-words pt-[24px] tracking-[1.4px] md:col-span-5 md:pt-0 md:tracking-[1.6px]">
          {{ product.producer.profile }}
        </div>
      </div>
      <div class="flex justify-center px-[16px] pb-12 pt-[24px] md:justify-end md:px-[64px] md:pt-[45px]">
        <button class="flex items-center text-main">
          詳しく見る
          <div class="pl-4">
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
                d="M1 11.3125L0.0302535 10.3428L4.71736 5.65565L0.0302528 0.968538L0.999999 -0.00120831L6.65685 5.65565L1 11.3125Z"
                fill="#604C3F"
              />
            </svg>
          </div>
        </button>
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
