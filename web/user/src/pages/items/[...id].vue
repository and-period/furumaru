<script setup lang="ts">
import { storeToRefs } from 'pinia'
import { useProductStore } from '~/store/product'
import { useShoppingCartStore } from '~/store/shopping'
import { Snackbar } from '~/types/props'

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
      class="md:grid animate-pulse md:grid-cols-2 bg-white px-[112px] pb-6 pt-[40px] text-main"
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
    <div class="md:grid md:grid-cols-2 bg-white px-4 md:px-[112px] pb-6 pt-[40px] text-main">
      <div class="w-full">
        <div class="mx-auto aspect-square md:h-[500px] md:w-[500px] w-full">
          <img
            class="w-full"
            :src="product.thumbnail.url"
            :alt="`${product.name}のサムネイル画像`"
          />
        </div>
      </div>

      <div class="flex w-full flex-col gap-4">
        <div class="break-words text-[24px] tracking-[2.4px]">
          {{ product.name }}
        </div>

        <div v-if="product.producer" class="mt-4 flex flex-col leading-[32px]">
          <div class="text-[16px] tracking-[1.6px]">
            生産者:
            <a href="#" class="font-bold underline">
              {{ product.producer.username }}
            </a>
          </div>
          <div class="text-[14px] tracking-[1.4px]">
            {{ product.originPrefecture }} {{ product.originCity }}
          </div>
        </div>

        <div
          v-if="product && product.recommendedPoint1"
          class="mt-8 w-full rounded-2xl bg-base px-[20px] py-[28px] text-main"
        >
          <p class="mb-[12px] text-[14px] tracking-[1.4px]">おすすめポイント</p>
          <ol
            class="recommend-list flex flex-col divide-y divide-dashed divide-main px-[4px] pl-[24px]"
          >
            <li v-if="product.recommendedPoint1" class="py-3">
              {{ product.recommendedPoint1 }}
            </li>
            <li v-if="product.recommendedPoint2" class="py-3">
              {{ product.recommendedPoint2 }}
            </li>
            <li v-if="product.recommendedPoint3" class="py-3">
              {{ product.recommendedPoint3 }}
            </li>
          </ol>
        </div>

        <div>
          <div
            class="mt-[60px] text-[32px] after:ml-2 after:text-[16px] after:content-['(税込)']"
          >
            {{ priceString }}
          </div>

          <div v-if="product" class="mt-8 inline-flex items-center">
            <label class="mr-2 block text-[16px]">数量</label>
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
          class="mt-8 w-full bg-main py-4 text-center text-white"
          @click="handleClickAddCartButton"
        >
          買い物カゴに入れる
        </button>

        <div class="mt-4 inline-flex gap-4">
          <span
            v-for="productTag in product.productTags"
            :key="productTag?.id"
            class="rounded-2xl border border-main px-4 py-1"
          >
            {{ productTag?.name }}
          </span>
        </div>
      </div>

      <div class="col-span-2 pb-16">
        <article class="leading-[32px]" v-html="product.description" />
      </div>

      <div class="flex flex-col border-main divide-main divide-y divide-dashed border-y border-dashed col-span-2">
        <div class="grid grid-cols-5 py-4">
          <p>・賞味期限</p>
          <p class="col-span-4">発送日より{{ product.expirationDate }}日</p>
        </div>
        <div class="grid grid-cols-5 py-4">
          <p>・内容量</p>
          <p class="col-span-4">{{ product.weight }}kg</p>
        </div>
        <div class="grid grid-cols-5 py-4">
          <p>・配送方法</p>
          <p class="col-span-4">{{ getDeliveryType(product.deliveryType) }}</p>
        </div>
        <div class="grid grid-cols-5 py-4">
          <p>・保存方法</p>
          <p class="col-span-4">{{ getStorageMethodType(product.storageMethodType) }}</p>
        </div>
      </div>
    </div>
  </template>

  <div class="mx-auto mt-[80px] w-[1216px] bg-white rounded-3xl">
    <div class="pt-10 px-[64px]">
      <p class="bg-base text-main mx-auto rounded-full py-2 text-center text-[16px] font-bold">この商品の生産者</p>
    </div>
    <div class="grid grid-cols-8 pt-16 px-[64px] mx-auto" v-if="product.producer">
      <img
        :src="product.producer.thumbnailUrl"
        :alt="`${product.producer.username}`"
        class="aspect-square w-[120px] rounded-full"
      />
      <div class="text-main col-span-2 ml-4">
        <div class="flex gap-4 test-sm pt-4">
          <p class="text-sm font-[500] tracking-[1.4px]">{{ product.originPrefecture }}</p>
          <p class="text-sm font-[500] tracking-[1.4px]">{{ product.originCity }}</p>
        </div>
        <div class="flex items-end pt-4 text-sm">
          <p class="font-[500] tracking-[1.4px]">生産者</p>
          <p class="pl-6 text-2xl">{{ product.producer.username }}</p>
        </div>
      </div>
      <div class="col-span-5 break-words">
        {{ product.producer.profile }}
      </div>
    </div>
    <div class="mx-auto flex px-[64px] flex-row-reverse pt-[45px] pb-12">
      <button class="text-main flex items-center">
        詳しく見る
        <div class="pl-4">
          <svg xmlns="http://www.w3.org/2000/svg" width="7" height="12" viewBox="0 0 7 12" fill="none">
            <path fill-rule="evenodd" clip-rule="evenodd" d="M1 11.3125L0.0302535 10.3428L4.71736 5.65565L0.0302528 0.968538L0.999999 -0.00120831L6.65685 5.65565L1 11.3125Z" fill="#604C3F"/>
          </svg>
        </div>
      </button>
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
