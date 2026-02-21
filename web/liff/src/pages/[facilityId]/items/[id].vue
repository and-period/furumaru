<script setup lang="ts">
import { FmProductDetail } from '@furumaru/shared';
import { useProductStore } from '~/stores/product';
import { useShoppingCartStore } from '~/stores/shopping';
import { ResponseError } from '~/types/api/facility';

const route = useRoute();
const router = useRouter();

const productStore = useProductStore();
const shoppingCartStore = useShoppingCartStore();
const facilityId = computed<string>(() => String(route.params.facilityId || ''));

// Get product ID from route parameter
const productId = route.params.id as string;

// Find the product in the store
const product = computed(() => {
  return productStore.getProductById(productId);
});

const isProductLoading = ref(true);
const productFetchError = ref<string | null>(null);

const loadProductDetail = async () => {
  if (!facilityId.value || !productId) {
    throw createError({
      statusCode: 404,
      statusMessage: '商品が見つかりません',
    });
  }

  isProductLoading.value = true;
  productFetchError.value = null;

  try {
    await productStore.fetchFacilityProductDetail(facilityId.value, productId);

    if (!product.value) {
      throw createError({
        statusCode: 404,
        statusMessage: '商品が見つかりません',
      });
    }
  }
  catch (e) {
    console.error('Failed to fetch product detail', e);
    if (e instanceof ResponseError && e.response.status === 404) {
      throw createError({
        statusCode: 404,
        statusMessage: '商品が見つかりません',
      });
    }
    productFetchError.value = '商品情報の取得に失敗しました。時間をおいて再度お試しください。';
  }
  finally {
    isProductLoading.value = false;
  }
};

await loadProductDetail();

// Transform product data for FmProductDetail component
const mediaFiles = computed(() => {
  if (!product.value?.media) return [];
  return product.value.media.map(media => ({
    url: media.url,
    isThumbnail: media.isThumbnail,
  }));
});

// Add to cart handler
const isAdding = ref(false);
const addToCartError = ref<string | null>(null);

const addToCart = async () => {
  if (!product.value) return;
  if (product.value.inventory === 0) return;
  try {
    isAdding.value = true;
    addToCartError.value = null;
    await shoppingCartStore.addCartItem(facilityId.value, product.value.id, 1);
  }
  catch (e) {
    console.error('Failed to add to cart', e);
    let statusLabel = '';
    if (e instanceof ResponseError) {
      statusLabel = `（ステータスコード: ${e.response.status}）`;
    }
    addToCartError.value = `カゴへの追加に失敗しました${statusLabel}。\n時間をおいて再度お試しください。`;
  }
  finally {
    isAdding.value = false;
  }
};

const rating = computed(() => {
  if (!product.value?.rate) return { average: 0, count: 0 };
  return {
    average: product.value.rate.average,
    count: product.value.rate.count,
    detail: product.value.rate.detail,
  };
});
</script>

<template>
  <div class="min-h-screen bg-white relative mb-[186px] overflow-auto">
    <!-- Navigation header for mobile -->
    <div class="top-0 z-10 bg-white border-b border-gray-200 px-4 py-3 fixed w-full">
      <div class="flex items-center">
        <button
          class="flex items-center text-gray-600 hover:text-gray-800"
          @click="router.back()"
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
      <div
        v-if="isProductLoading"
        class="text-center text-gray-500 py-10"
      >
        商品情報を読み込み中です...
      </div>
      <div
        v-else-if="productFetchError"
        class="text-center text-red-600 bg-red-50 border border-red-200 rounded px-4 py-6"
      >
        {{ productFetchError }}
      </div>
      <fm-product-detail
        v-else-if="product"
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
      v-if="product && !isProductLoading && !productFetchError"
      class="fixed bottom-[56px] bg-white border-t border-gray-200 p-4 w-full"
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
        :disabled="product.inventory === 0 || isAdding"
        class="w-full bg-orange text-white py-3 px-4 rounded-lg font-semibold disabled:bg-gray-300 disabled:cursor-not-allowed cursor-pointer"
        @click="addToCart"
      >
        {{ product.inventory > 0 ? (isAdding ? '追加中...' : 'カゴに追加') : '在庫なし' }}
      </button>
      <div
        v-if="addToCartError"
        class="mt-3 px-4 py-3 bg-red-50 border border-red-200 text-red-700 text-sm rounded"
      >
        {{ addToCartError }}
      </div>
    </div>
  </div>
</template>
