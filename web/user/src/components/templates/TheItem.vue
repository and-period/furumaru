<script setup lang="ts">
import { storeToRefs } from 'pinia'
import { useEventBus } from '@vueuse/core'
import { useProductStore } from '~/store/product'
import { useShoppingCartStore } from '~/store/shopping'
import type { Snackbar } from '~/types/props'
import type { I18n } from '~/types/locales'

const i18n = useI18n()

const router = useRouter()
const route = useRoute()

const productStore = useProductStore()
const shoppingCartStore = useShoppingCartStore()

const { fetchProducts } = productStore
const { addCart } = shoppingCartStore
const { products, totalProductsCount } = storeToRefs(productStore)

const { emit } = useEventBus('add-to-cart')

const lt = (str: keyof I18n['items']['list']) => {
  return i18n.t(`items.list.${str}`)
}

const handleClick = (id: string) => {
  router.push(`/items/${id}`)
}

const snackbarItems = ref<Snackbar[]>([])

const handleClickAddCartButton = async (
  name: string,
  id: string,
  quantity: number,
) => {
  await addCart({
    productId: id,
    quantity,
  })
  emit('add-to-cart')
}

// 1ページ当たりに表示する商品数
const pagePerItems = ref<number>(20)

// 現在のページ番号
const currentPage = computed<number>(() => {
  return route.query.page ? Number(route.query.page) : 1
})

// ページネーション情報
const pagination = computed<{
  limit: number
  offset: number
  pageArray: number[]
}>(() => {
  const totalPage = Math.ceil(totalProductsCount.value / pagePerItems.value)
  const pageArray = Array.from({ length: totalPage }, (_, i) => i + 1)

  return {
    limit: pagePerItems.value,
    offset: pagePerItems.value * (currentPage.value - 1),
    pageArray,
  }
})

const handleClickPage = (page: number) => {
  router.push({
    query: {
      ...route.query,
      page,
    },
  })
}

watch(currentPage, () => {
  fetchProducts(pagePerItems.value, pagination.value.offset)
})

const { status } = useAsyncData('products', () => {
  return fetchProducts(pagePerItems.value, pagination.value.offset)
})

const hideV1App = false
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

  <div
    class="flex flex-col bg-white px-[15px] py-[48px] text-main md:px-[36px]"
  >
    <div class="w-full">
      <p
        class="text-center text-[14px] font-bold tracking-[2px] md:text-[20px]"
      >
        {{ lt("allItemsTitle") }}
      </p>
      <div
        v-if="hideV1App"
        class="mt-[24px] w-full md:mt-[38px]"
      >
        <div class="relative mx-auto md:w-[648px]">
          <the-search-icon class="absolute left-[24px] top-[12px]" />
          <input
            class="block w-full rounded-[28px] border border-typography py-3 pl-[56px] text-[12px] placeholder:text-center focus:border-2 focus:border-main focus:outline-none md:text-[16px]"
            type="text"
            placeholder="すべての商品からさがす"
          >
        </div>
      </div>
    </div>
    <hr class="mt-[40px]">
    <div class="mt-[24px] w-full">
      <div
        v-if="hideV1App"
        class="text-right"
      >
        <div
          class="inline-flex text-[12px] tracking-[1.3px] text-typography md:text-[13px]"
        >
          <div class="mr-[16px]">
            並び替え：
          </div>
          <div class="inline-flex gap-[22px]">
            <button class="border-b border-main pb-2 text-main">
              新着順
            </button>
            <button class="pb-2">
              値段の安い順
            </button>
            <button class="pb-2">
              値段の高い順
            </button>
          </div>
        </div>
      </div>

      <div
        class="mx-auto mt-[24px] grid max-w-[1440px] grid-cols-2 gap-x-[19px] gap-y-6 md:grid-cols-3 md:gap-x-8 lg:grid-cols-4 xl:grid-cols-5"
      >
        <template v-if="status === 'pending'">
          <div
            v-for="i in [1, 2, 3, 4, 5]"
            :key="i"
            class="w-full animate-pulse"
          >
            <div class="aspect-square w-full bg-slate-200" />
            <div class="mt-2 h-[24px] w-[80%] rounded-lg bg-slate-200" />
            <div class="mt-2 h-[24px] w-[60%] rounded-lg bg-slate-200" />
          </div>
        </template>

        <template v-else>
          <the-product-list-item
            v-for="product in products"
            :id="product.id"
            :key="product.id"
            :status="product.status"
            :name="product.name"
            :price="product.price"
            :inventory="product.inventory"
            :has-stock="product.hasStock"
            :thumbnail="product.thumbnail"
            :coordinator="product.coordinator"
            :origin-city="product.originCity"
            :thumbnail-is-video="product.thumbnailIsVideo"
            @click:item="handleClick"
            @click:add-cart="handleClickAddCartButton"
          />
        </template>
      </div>
      <the-pagination
        class="mt-8"
        :current-page="currentPage"
        :page-array="pagination.pageArray"
        @change-page="handleClickPage"
      />
    </div>
  </div>
</template>
