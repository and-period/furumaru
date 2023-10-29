<script setup lang="ts">
import { MOCK_ALL_PRODUCT_ITEMS } from '~/constants/mock'

const router = useRouter()

const handleClick = (id: string) => {
  router.push(`/items/${id}`)
}
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
        <the-product-list-item
          class="cursor-pointer"
          v-for="productItem in MOCK_ALL_PRODUCT_ITEMS"
          :key="productItem.id"
          :id="productItem.id"
          :name="productItem.name"
          :price="productItem.price"
          :img-src="productItem.imgSrc"
          :inventory="productItem.inventory"
          :address="productItem.address"
          :cn-name="productItem.cnName"
          :cn-img-src="productItem.cnImgSrc"
          @click="handleClick(productItem.id)"
        />
      </div>
    </div>
  </div>
</template>
