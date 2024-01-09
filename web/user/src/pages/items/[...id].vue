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

const selectedMediaIndex = ref<number>(-1)

const handleClickMediaItem = (index: number) => {
  selectedMediaIndex.value = index
}

useAsyncData(`product-${id.value}`, () => {
  return fetchProduct(id.value)
})

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
      class="gap-4 bg-white px-4 pb-6 pt-[40px] text-main md:grid md:grid-cols-2 md:px-[112px]"
    >
      <div class="mx-auto aspect-square w-full max-w-[500px]">
        <div class="flex aspect-square w-full justify-center">
          <img
            class="block w-full object-contain"
            :src="
              selectedMediaIndex === -1
                ? product.thumbnail.url
                : product.media[selectedMediaIndex].url
            "
            :alt="`${product.name}のサムネイル画像`"
          />
        </div>
        <div
          class="hidden-scrollbar mt-2 grid w-full grid-flow-col justify-start gap-2 overflow-x-scroll"
        >
          <img
            v-for="(m, i) in product.media"
            :key="i"
            :src="m.url"
            :alt="`${product.name}の画像_${i}`"
            class="aspect-square w-[72px] cursor-pointer"
            @click="handleClickMediaItem(i)"
          />
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

  <template v-if="product.producer">
    <div class="mx-auto mt-[40px] w-full px-4 xl:px-28">
      <div
        class="flex w-full flex-col rounded-3xl bg-white px-8 py-10 text-main xl:px-16"
      >
        <p
          class="mx-auto w-full rounded-full bg-base py-2 text-center text-[14px] font-bold text-main md:text-[16px]"
        >
          この商品の生産者
        </p>

        <div class="mt-[64px] flex w-full flex-col gap-4 md:flex-row lg:gap-10">
          <div
            class="flex min-w-max flex-col items-center justify-center gap-4 md:flex-row"
          >
            <img
              :src="product.producer.thumbnailUrl"
              :alt="`${product.producer.username}`"
              class="mx-auto block aspect-square w-[96px] rounded-full md:w-[120px]"
            />
            <div
              class="flex min-w-max grow flex-col items-center gap-2 md:items-start md:gap-2 md:whitespace-nowrap"
            >
              <p class="text-sm font-[500] tracking-[1.4px]">
                {{ `${product.originPrefecture} ${product.originCity}` }}
              </p>
              <p
                class="text-[16px] tracking-[1.4px] before:mr-1 before:text-[14px] before:font-medium before:content-['生産者'] md:text-[24px]"
              >
                {{ product.producer.username }}
              </p>
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
  </template>
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
