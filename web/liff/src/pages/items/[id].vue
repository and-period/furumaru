<script setup lang="ts">
import { FmProductDetail } from '@furumaru/shared';
import { useProductStore } from '~/stores/product';

const route = useRoute();

const productStore = useProductStore();

// Get product ID from route parameter
const productId = route.params.id as string;

// Find the product in the store
const product = computed(() => {
  return productStore.getProductById(productId);
});

// Transform product data for FmProductDetail component
const mediaFiles = computed(() => {
  if (!product.value?.media) return [];
  return product.value.media.map(media => ({
    url: media.url,
    isThumbnail: media.isThumbnail,
  }));
});

const rating = computed(() => {
  if (!product.value?.rate) return { average: 0, count: 0 };
  return {
    average: product.value.rate.average,
    count: product.value.rate.count,
    detail: product.value.rate.detail,
  };
});

// Handle case where product is not found
if (!product.value) {
  throw createError({
    statusCode: 404,
    statusMessage: '商品が見つかりません',
  });
}
</script>

<template>
  <div class="min-h-screen bg-white">
    <!-- Navigation header for mobile -->
    <div class="sticky top-0 z-10 bg-white border-b border-gray-200 px-4 py-3">
      <div class="flex items-center">
        <button
          class="flex items-center text-gray-600 hover:text-gray-800"
          @click="$router.back()"
        >
          <svg
            class="w-6 h-6 mr-2"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M15 19l-7-7 7-7"
            />
          </svg>
          戻る
        </button>
      </div>
    </div>

    <!-- Product detail content -->
    <div class="px-4 py-6">
      <FmProductDetail
        v-if="product"
        :media-files="mediaFiles"
        :name="product.name"
        :description="product.description"
        :origin-prefecture="product.originPrefecture"
        :origin-city="product.originCity"
        :rating="rating"
        :recommended-point1="product.recommendedPoint1"
        :recommended-point2="product.recommendedPoint2"
        :recommended-point3="product.recommendedPoint3"
        :expiration-date="product.expirationDate"
        :weight="product.weight"
        :delivery-type="product.deliveryType"
        :storage-method-type="product.storageMethodType"
      />
    </div>

    <!-- Price and purchase section for mobile -->
    <div
      v-if="product"
      class="sticky bottom-0 bg-white border-t border-gray-200 p-4"
    >
      <div class="flex items-center justify-between mb-4">
        <div class="text-2xl font-bold text-gray-900">
          {{ new Intl.NumberFormat('ja-JP', { style: 'currency', currency: 'JPY' }).format(product.price) }}
        </div>
        <div class="text-sm text-gray-500">
          在庫: {{ product.inventory }}個
        </div>
      </div>

      <button
        :disabled="product.inventory === 0"
        class="w-full bg-orange-500 text-white py-3 px-4 rounded-lg font-semibold disabled:bg-gray-300 disabled:cursor-not-allowed"
        @click="() => {}"
      >
        {{ product.inventory > 0 ? 'カートに追加' : '在庫なし' }}
      </button>
    </div>
  </div>
</template>

<style scoped>
/* Mobile-first responsive design for LINE app context */
@media (max-width: 640px) {
  .sticky {
    position: -webkit-sticky;
    position: sticky;
  }
}
</style>
