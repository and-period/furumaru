<script setup lang="ts">
import { storeToRefs } from 'pinia'
import { useProductStore } from '~/store/product'

const router = useRouter()

const productStore = useProductStore()
const { fetchProducts } = productStore
const { productsFetchState, products } = storeToRefs(productStore)

const handleClick = (id: string) => {
  router.push(`/items/${id}`)
}

fetchProducts()
</script>

<template>
  <div class="flex flex-col bg-white px-[36px] py-[48px] text-main">
    <div class="w-full">
      <p class="text-center text-[20px] font-bold tracking-[2px]">
        すべての商品
      </p>
      <div class="mt-[38px] w-full">
        <div class="relative mx-auto md:w-[648px]">
          <the-search-icon class="absolute left-[24px] top-[12px]" />
          <input
            class="block w-full rounded-[28px] border border-typography py-3 pl-[56px] placeholder:text-center focus:border-2 focus:border-main focus:outline-none"
            type="text"
            placeholder="すべての商品からさがす"
          />
        </div>
      </div>
    </div>
    <hr class="mt-[40px]" />
    <div class="mt-[24px] w-full">
      <div class="text-right">
        <div class="inline-flex text-[13px] tracking-[1.3px] text-typography">
          <div class="mr-[16px]">並び替え：</div>
          <div class="inline-flex gap-[22px]">
            <button class="border-b border-main pb-4 text-main">新着順</button>
            <button class="pb-4">値段の安い順</button>
            <button class="pb-4">値段の高い順</button>
          </div>
        </div>
      </div>

      <div
        class="mx-auto mt-[24px] grid max-w-[1440px] grid-cols-2 gap-x-8 gap-y-6 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5"
      >
        <template v-if="productsFetchState.isLoading">
          <div
            v-for="i in [1, 2, 3, 4, 5]"
            :key="i"
            class="w-full animate-pulse"
          >
            <div class="aspect-square w-full bg-slate-200"></div>
            <div class="mt-2 h-[24px] w-[80%] rounded-lg bg-slate-200"></div>
            <div class="mt-2 h-[24px] w-[60%] rounded-lg bg-slate-200"></div>
          </div>
        </template>

        <template v-else>
          <the-product-list-item
            v-for="product in products"
            :key="product.id"
            class="cursor-pointer"
            :name="product.name"
            :price="product.price"
            :inventory="product.inventory"
            :has-stock="product.hasStock"
            :thumbnail="product.thumbnail"
            :coordinator="product.coordinator"
            :origin-city="product.originCity"
            @click="handleClick(product.id)"
          />
        </template>
      </div>
    </div>
  </div>
</template>
