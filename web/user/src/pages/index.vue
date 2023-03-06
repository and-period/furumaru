<script lang="ts" setup>
import { storeToRefs } from 'pinia'
import { useShoppingStore } from '~/store/shopping'

const shoppingStore = useShoppingStore()
const { recommendProducts } = storeToRefs(shoppingStore)

shoppingStore.setupDummyData()

const banners: string[] = [
  '/img/banner.png',
  '/img/banner.png',
  '/img/banner.png'
]
</script>

<template>
  <div>
    <the-carousel :images="banners" />

    <div class="my-6 flex flex-col gap-y-16">
      <the-content-box
        title="live"
        sub-title="配信中・配信予定のマルシェ"
      />

      <the-content-box
        title="archive"
        sub-title="過去のマルシェ"
      />

      <the-content-box
        title="recommend"
        sub-title="おすすめの商品"
      >
        <div class="grid grid-cols-4 gap-x-4 gap-y-6">
          <the-product-list-item v-for="product in recommendProducts" :key="product.id" :item="product" />
        </div>
      </the-content-box>
    </div>
  </div>
</template>
