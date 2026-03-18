<script setup lang="ts">
import { useRecentlyViewed } from '~/hooks/useRecentlyViewed'
import type {
  Coordinator,
  ExperienceMediaInner,
  Product,
  ProductStatus } from '~/types/api'
import { ProductResponseFromJSON } from '~/types/api'

interface RecentlyViewedProduct {
  id: string
  name: string
  price: number
  status: ProductStatus
  inventory: number
  hasStock: boolean
  originCity: string
  coordinator: Coordinator | undefined
  thumbnail: ExperienceMediaInner | undefined
  thumbnailIsVideo: boolean
}

const { getItems } = useRecentlyViewed()

const recentProducts = ref<RecentlyViewedProduct[]>([])
const isLoading = ref(false)

async function fetchRecentProducts() {
  if (import.meta.server) {
    return
  }

  const ids = getItems()
  if (ids.length === 0) {
    return
  }

  isLoading.value = true
  try {
    const runtimeConfig = useRuntimeConfig()
    const results: RecentlyViewedProduct[] = []

    const responses = await Promise.allSettled(
      ids.map(async (id) => {
        const url = `${runtimeConfig.public.API_BASE_URL}/v1/products/${id}`
        const res = await fetch(url, { credentials: 'include' })
        if (!res.ok) {
          return null
        }
        const json = await res.json()
        return { id, data: ProductResponseFromJSON(json) }
      }),
    )

    for (const response of responses) {
      if (response.status === 'fulfilled' && response.value) {
        const { data } = response.value
        const product: Product = data.product
        const thumbnail = product.media.find(m => m.isThumbnail)
        results.push({
          id: product.id,
          name: product.name,
          price: product.price,
          status: product.status,
          inventory: product.inventory,
          hasStock: product.inventory > 0,
          originCity: product.originCity,
          coordinator: data.coordinator,
          thumbnail,
          thumbnailIsVideo: thumbnail ? thumbnail.url.endsWith('.mp4') : false,
        })
      }
    }

    recentProducts.value = results
  }
  finally {
    isLoading.value = false
  }
}

onMounted(() => {
  fetchRecentProducts()
})
</script>

<template>
  <div
    v-if="recentProducts.length > 0"
    class="mx-auto mt-[40px] w-full max-w-[1440px] px-4 xl:px-28"
  >
    <div
      class="flex w-full flex-col rounded-3xl bg-white px-8 py-10 text-main xl:px-16"
    >
      <p
        class="mx-auto w-full rounded-full bg-base py-2 text-center text-[14px] font-bold text-main md:text-[16px]"
      >
        最近チェックした商品
      </p>

      <div
        class="hidden-scrollbar mt-[24px] flex gap-4 overflow-x-auto pb-4"
      >
        <div
          v-for="product in recentProducts"
          :key="product.id"
          class="w-[180px] flex-shrink-0 md:w-[220px]"
        >
          <the-product-list-item
            :id="product.id"
            :status="product.status"
            :name="product.name"
            :price="product.price"
            :inventory="product.inventory"
            :has-stock="product.hasStock"
            :thumbnail="product.thumbnail"
            :coordinator="product.coordinator"
            :origin-city="product.originCity"
            :thumbnail-is-video="product.thumbnailIsVideo"
          />
        </div>
      </div>
    </div>
  </div>
</template>
