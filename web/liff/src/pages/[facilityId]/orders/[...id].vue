<script setup lang="ts">
import { storeToRefs } from 'pinia';
import { useOrderStore } from '~/stores/order';
import { useAuthStore } from '~/stores/auth';

const route = useRoute();
const router = useRouter();

const facilityId = computed<string>(() => String(route.params.facilityId || ''));
const orderId = computed<string>(() => {
  // [...id] は配列として取得されるため、最初の要素を使用
  const id = route.params.id;
  return Array.isArray(id) ? id[0] : String(id || '');
});

// ストア
const orderStore = useOrderStore();
const authStore = useAuthStore();
const { orderDetail, isLoading, error } = storeToRefs(orderStore);

// 認証チェック
const checkAuth = () => {
  if (!authStore.isAuthenticated) {
    router.push(`/${facilityId.value}`);
    return false;
  }
  return true;
};

// 注文詳細取得
const fetchOrderDetail = async () => {
  if (!checkAuth()) return;

  if (!orderId.value) {
    console.error('Order ID is required');
    return;
  }

  try {
    await orderStore.getOrderDetail(facilityId.value, orderId.value);
  }
  catch (e) {
    console.error('Failed to fetch order detail:', e);
  }
};

// 注文一覧に戻る
const goBackToOrders = () => {
  router.push(`/${facilityId.value}/orders`);
};

// 日付フォーマット関数
const formatDate = (timestamp: number): string => {
  if (!timestamp) return '-';
  const date = new Date(timestamp * 1000);
  return date.toLocaleDateString('ja-JP', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
  });
};

// 価格フォーマット関数
const formatPrice = (price: number): string => {
  return `¥${price.toLocaleString('ja-JP')}`;
};

// 注文ステータスの表示名
const getOrderStatusText = (status: number): string => {
  const statusMap: Record<number, string> = {
    1: '未確定',
    2: '確定',
    3: '配送中',
    4: '完了',
    5: 'キャンセル',
  };
  return statusMap[status] || '不明';
};

// 注文ステータスのスタイルクラス
const getOrderStatusClass = (status: number): string => {
  const classMap: Record<number, string> = {
    1: 'bg-yellow-100 text-yellow-800',
    2: 'bg-blue-100 text-blue-800',
    3: 'bg-purple-100 text-purple-800',
    4: 'bg-green-100 text-green-800',
    5: 'bg-red-100 text-red-800',
  };
  return classMap[status] || 'bg-gray-100 text-gray-800';
};

// 商品情報を取得する関数
const getProductInfo = (productId: string) => {
  return orderDetail.value?.products?.find(product => product.id === productId);
};

// マウント時に注文詳細を取得
onMounted(() => {
  fetchOrderDetail();
});
</script>

<template>
  <div class="container mx-auto px-4 py-6">
    <!-- 戻るボタン -->
    <div class="mb-4">
      <button
        class="flex items-center text-blue-600 hover:text-blue-800 transition-colors"
        @click="goBackToOrders"
      >
        <svg
          xmlns="http://www.w3.org/2000/svg"
          fill="none"
          viewBox="0 0 24 24"
          stroke-width="1.5"
          stroke="currentColor"
          class="w-5 h-5 mr-1"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            d="M15.75 19.5L8.25 12l7.5-7.5"
          />
        </svg>
        注文一覧に戻る
      </button>
    </div>

    <h1 class="text-2xl font-bold text-center mb-6">
      注文詳細
    </h1>

    <!-- Loading state -->
    <div
      v-if="isLoading"
      class="text-center py-8"
    >
      <p class="text-gray-600">
        注文詳細を読み込み中...
      </p>
    </div>

    <!-- Error state -->
    <div
      v-else-if="error"
      class="text-center py-8"
    >
      <p class="text-red-600 mb-4">
        エラーが発生しました: {{ error }}
      </p>
      <button
        class="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600 transition-colors"
        @click="fetchOrderDetail"
      >
        再試行
      </button>
    </div>

    <!-- 注文詳細 -->
    <div
      v-else-if="orderDetail?.order"
      class="space-y-6"
    >
      <!-- 注文基本情報 -->
      <div class="bg-white border border-gray-200 rounded-lg p-6">
        <div class="flex justify-between items-start mb-4">
          <div>
            <h2 class="text-xl font-semibold mb-2">
              注文 #{{ orderDetail.order.id.slice(-8) }}
            </h2>
            <p class="text-gray-600">
              注文日時: {{ formatDate(orderDetail.order.payment.orderedAt) }}
            </p>
          </div>
          <span
            class="px-3 py-1 rounded-full text-sm font-medium"
            :class="getOrderStatusClass(orderDetail.order.payment.status)"
          >
            {{ getOrderStatusText(orderDetail.order.payment.status) }}
          </span>
        </div>
      </div>

      <!-- コーディネーター情報 -->
      <div
        v-if="orderDetail.coordinator"
        class="bg-white border border-gray-200 rounded-lg p-6"
      >
        <h3 class="text-lg font-semibold mb-3">
          コーディネーター情報
        </h3>
        <div class="space-y-2">
          <p><span class="font-medium">名前:</span> {{ orderDetail.coordinator.username }}</p>
          <p><span class="font-medium">所在地:</span> {{ orderDetail.coordinator.prefecture }}{{ orderDetail.coordinator.city }}</p>
          <p
            v-if="orderDetail.coordinator.profile"
            class="text-sm text-gray-600"
          >
            {{ orderDetail.coordinator.profile }}
          </p>
        </div>
      </div>

      <!-- 注文商品一覧 -->
      <div class="bg-white border border-gray-200 rounded-lg p-6">
        <h3 class="text-lg font-semibold mb-4">
          注文商品
        </h3>
        <div
          v-if="orderDetail.order.items && orderDetail.order.items.length > 0"
          class="space-y-4"
        >
          <div
            v-for="item in orderDetail.order.items"
            :key="item.productId"
            class="border border-gray-100 rounded-lg p-4"
          >
            <div class="flex gap-4">
              <!-- 商品画像 -->
              <div class="flex-shrink-0">
                <img
                  v-if="getProductInfo(item.productId)?.thumbnailUrl"
                  :src="getProductInfo(item.productId)?.thumbnailUrl"
                  :alt="getProductInfo(item.productId)?.name || '商品画像'"
                  class="w-20 h-20 object-cover rounded-lg"
                >
                <div
                  v-else
                  class="w-20 h-20 bg-gray-200 rounded-lg flex items-center justify-center"
                >
                  <span class="text-gray-400 text-xs">画像なし</span>
                </div>
              </div>

              <!-- 商品情報 -->
              <div class="flex-grow">
                <h4 class="font-semibold mb-2">
                  {{ getProductInfo(item.productId)?.name || `商品ID: ${item.productId}` }}
                </h4>
                <div class="space-y-1 text-sm text-gray-600">
                  <p>単価: {{ formatPrice(item.price) }}</p>
                  <p>数量: {{ item.quantity }}{{ getProductInfo(item.productId)?.itemUnit || '個' }}</p>
                  <p class="font-semibold text-base text-gray-900">
                    小計: {{ formatPrice(item.price * item.quantity) }}
                  </p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 配送情報 -->
      <div class="bg-white border border-gray-200 rounded-lg p-6">
        <h3 class="text-lg font-semibold mb-4">
          配送情報
        </h3>
        <div class="space-y-2">
          <p>
            <span class="font-medium">配送先:</span>
            {{ orderDetail.order.pickupLocation || '配送先情報なし' }}
          </p>
        </div>
      </div>

      <!-- 支払い情報 -->
      <div class="bg-white border border-gray-200 rounded-lg p-6">
        <h3 class="text-lg font-semibold mb-4">
          支払い情報
        </h3>
        <div class="space-y-3">
          <div class="flex justify-between">
            <span>商品合計:</span>
            <span>{{ formatPrice(orderDetail.order.payment.subtotal) }}</span>
          </div>
          <div
            v-if="orderDetail.order.payment.discount > 0"
            class="flex justify-between text-red-600"
          >
            <span>割引:</span>
            <span>-{{ formatPrice(orderDetail.order.payment.discount) }}</span>
          </div>
          <div
            v-if="orderDetail.order.payment.shippingFee > 0"
            class="flex justify-between"
          >
            <span>送料:</span>
            <span>{{ formatPrice(orderDetail.order.payment.shippingFee) }}</span>
          </div>
          <div class="border-t border-gray-200 pt-3">
            <div class="flex justify-between text-lg font-bold">
              <span>合計金額:</span>
              <span class="text-blue-600">{{ formatPrice(orderDetail.order.payment.total) }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- プロモーション情報 -->
      <div
        v-if="orderDetail.promotion"
        class="bg-white border border-gray-200 rounded-lg p-6"
      >
        <h3 class="text-lg font-semibold mb-3">
          適用プロモーション
        </h3>
        <div class="space-y-2">
          <p class="font-medium">
            {{ orderDetail.promotion.title }}
          </p>
          <p
            v-if="orderDetail.promotion.description"
            class="text-sm text-gray-600"
          >
            {{ orderDetail.promotion.description }}
          </p>
        </div>
      </div>
    </div>

    <!-- 注文が見つからない場合 -->
    <div
      v-else
      class="text-center py-12"
    >
      <p class="text-gray-500 text-lg mb-4">
        注文が見つかりませんでした
      </p>
      <button
        class="bg-blue-500 text-white px-6 py-3 rounded-lg hover:bg-blue-600 transition-colors"
        @click="goBackToOrders"
      >
        注文一覧へ戻る
      </button>
    </div>
  </div>
</template>
