<script lang="ts" setup>
import { useShoppingStore } from '~/store/shopping'
import { MOCK_LIVE_ITEMS, MOCK_ARCHIVES_ITEMS, MOCK_RECOMMEND_ITEMS } from '~/constants/mock'

const shoppingStore = useShoppingStore()

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
      >
        <div class="px-20 grid grid-cols-3 gap-x-10 gap-y-8">
          <the-live-item
            v-for="liveItem in MOCK_LIVE_ITEMS"
            :id="liveItem.id"
            :key="liveItem.id"
            :title="liveItem.title"
            :img-src="liveItem.imgSrc"
            :start-at="liveItem.startAt"
            :published="liveItem.published"
          />
        </div>
        <div class="w-full text-center mt-10">
          <button class="bg-main text-white py-2 w-60">
            もっと見る
          </button>
        </div>
      </the-content-box>

      <the-content-box
        class="px-8"
        title="archive"
        sub-title="過去のマルシェ"
      >
        <div class="flex flex-nowrap gap-x-8 overflow-x-scroll">
          <the-archive-item
            v-for="archiveItem in MOCK_ARCHIVES_ITEMS"
            :id="archiveItem.id"
            :key="archiveItem.id"
            :title="archiveItem.title"
            :img-src="archiveItem.imgSrc"
          />
        </div>

        <div class="w-full text-center mt-10">
          <button class="bg-main text-white py-2 w-60">
            一覧を見る
          </button>
        </div>
      </the-content-box>

      <the-content-box
        title="recommend"
        sub-title="おすすめの商品"
      >
        <div class="grid grid-cols-5 gap-x-4 gap-y-6">
          <the-product-list-item
            v-for="productItem in MOCK_RECOMMEND_ITEMS"
            :id="productItem.id"
            :key="productItem.id"
            :name="productItem.name"
            :price="productItem.price"
            :img-src="productItem.imgSrc"
            :inventory="productItem.inventory"
            :address="productItem.address"
            :cn-name="productItem.cnName"
            :cn-img-src="productItem.cnImgSrc"
          />
        </div>
      </the-content-box>
    </div>
  </div>
</template>
