<script setup lang="ts">
import { storeToRefs } from 'pinia'
import { useProductStore } from '~/store/product'
import { useShoppingCartStore } from '~/store/shopping'

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
}
</script>

<template>
  <template v-if="productFetchState.isLoading">
    <div
      class="grid animate-pulse grid-cols-2 bg-white px-[112px] pb-6 pt-[40px] text-main"
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
    <div class="grid grid-cols-2 bg-white px-[112px] pb-6 pt-[40px] text-main">
      <div class="w-full">
        <div class="mx-auto aspect-square h-[500px] w-[500px]">
          <img
            class="w-full"
            :src="product.thumbnail.url"
            :alt="`${product.name}のサムネイル画像`"
          />
        </div>
      </div>

      <div class="flex w-full flex-col gap-4">
        <div class="break-words text-[24px] tracking-[2.4px]">
          {{ selectItem?.name }}
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

      <div class="col-span-2">
        <article class="leading-[32px]" v-html="product.description" />
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
