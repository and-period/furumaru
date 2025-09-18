<script setup lang="ts">
import { storeToRefs } from 'pinia';
import { useOrderStore } from '~/stores/order';
import { useAuthStore } from '~/stores/auth';

const route = useRoute();
const router = useRouter();

const facilityId = computed<string>(() => String(route.params.facilityId || ''));

// ストア
const orderStore = useOrderStore();
const authStore = useAuthStore();
const { orders, isLoading, error } = storeToRefs(orderStore);

// 認証チェック
const checkAuth = () => {
  if (!authStore.isAuthenticated) {
    router.push(`/${facilityId.value}`);
    return false;
  }
  return true;
};

// 注文一覧取得
const fetchOrders = async () => {
  if (!checkAuth()) return;

  try {
    await orderStore.getOrders(facilityId.value, 50, 0);
  }
  catch (e) {
    console.error('Failed to fetch orders:', e);
  }
};

// 注文詳細ページへ遷移
const goToOrderDetail = (orderId: string) => {
  router.push(`/${facilityId.value}/orders/${orderId}`);
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
    2: 'bg-orange text-white',
    3: 'bg-purple-100 text-purple-800',
    4: 'bg-green-100 text-green-800',
    5: 'bg-red-100 text-red-800',
  };
  return classMap[status] || 'bg-gray-100 text-gray-800';
};

// 商品情報を取得する関数
const getProductInfo = (productId: string) => {
  return orders.value?.products?.find(product => product.id === productId);
};

// マウント時に注文一覧を取得
onMounted(() => {
  fetchOrders();
});
</script>

<template>
  <div class="container mx-auto px-4 pb-6">
    <h2 class="mt-6 font-semibold font-inter text-center w-full">
      注文一覧
    </h2>

    <!-- Loading state -->
    <div
      v-if="isLoading"
      class="text-center py-8"
    >
      <p class="text-gray-600">
        注文情報を読み込み中...
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
        class="bg-main text-white px-4 py-2 rounded hover:bg-orange transition-colors"
        @click="fetchOrders"
      >
        再試行
      </button>
    </div>

    <!-- 注文一覧 -->
    <div
      v-else-if="orders?.orders && orders.orders.length > 0"
      class="space-y-4 mt-6"
    >
      <div
        v-for="order in orders.orders"
        :key="order.id"
        class="border border-gray-200 rounded-lg p-4 cursor-pointer hover:bg-gray-50 transition-colors"
        @click="goToOrderDetail(order.id)"
      >
        <div class="flex justify-between items-start mb-2">
          <div>
            <div class="font-semibold text-md">
              注文 #{{ order.id.slice(-8) }}
            </div>
            <p class="text-sm text-gray-600">
              {{ formatDate(order.payment.orderedAt) }}
            </p>
          </div>
          <span
            class="px-2 py-1 rounded-full text-xs font-medium"
            :class="getOrderStatusClass(order.payment.status)"
          >
            {{ getOrderStatusText(order.payment.status) }}
          </span>
        </div>

        <div class="space-y-2">
          <!-- 商品情報の表示 -->
          <div
            v-if="order.items && order.items.length > 0"
            class="text-sm"
          >
            <p class="text-gray-600">
              商品数: {{ order.items.length }}点
            </p>
            <div class="text-gray-800">
              <div
                v-for="item in order.items.slice(0, 2)"
                :key="item.productId"
                class="truncate"
              >
                {{ getProductInfo(item.productId)?.name || `商品ID: ${item.productId}` }} × {{ item.quantity }}
              </div>
              <div
                v-if="order.items.length > 2"
                class="text-gray-500"
              >
                他{{ order.items.length - 2 }}点
              </div>
            </div>
          </div>

          <!-- 合計金額 -->
          <div class="flex justify-between items-center pt-2 border-t border-gray-100">
            <span class="font-semibold text-gray-800">合計金額:</span>
            <span class="font-bold text-lg text-main">
              {{ formatPrice(order.payment.total) }}
            </span>
          </div>
        </div>
      </div>
    </div>

    <!-- 注文がない場合 -->
    <div
      v-else
      class="text-center py-12"
    >
      <p class="text-gray-500 text-lg mb-4">
        注文履歴がありません
      </p>
      <nuxt-link
        :to="`/${facilityId}`"
        class="inline-block bg-blue-500 text-white px-6 py-3 rounded-lg hover:bg-blue-600 transition-colors"
      >
        商品一覧へ戻る
      </nuxt-link>
    </div>

    <!-- リフレッシュボタン -->
    <div
      v-if="orders?.orders && orders.orders.length > 0"
      class="text-center mt-6"
    >
      <button
        class="border border-orange text-orange px-4 py-2 rounded transition-colors disabled:opacity-50 hover:bg-orange/10"
        :disabled="isLoading"
        @click="fetchOrders"
      >
        {{ isLoading ? '読み込み中...' : '更新' }}
      </button>
    </div>
  </div>
</template>
