<script setup lang="ts">
import { useCheckoutStore } from '~/stores/checkout';

const route = useRoute();
const router = useRouter();
const checkoutStore = useCheckoutStore();

const isLoading = ref(false);
const error = ref<string | null>(null);
const orderId = ref<string | null>(null);

onMounted(async () => {
  error.value = null;
  // 支払い完了後のリダイレクト時に付与される想定のクエリを確認
  const q = route.query as Record<string, string | string[]>;
  const qVal = (key: string) => {
    const v = q[key];
    return Array.isArray(v) ? v[0] : v;
  };

  const orderIdQuery = qVal('orderId') || qVal('order_id');
  if (orderIdQuery) {
    orderId.value = String(orderIdQuery);
    return;
  }

  const transactionId = qVal('transactionId') || qVal('transaction_id') || qVal('tid');
  if (!transactionId) return;

  try {
    isLoading.value = true;
    const res = await checkoutStore.fetchCheckoutState(String(transactionId));
    orderId.value = res.orderId ?? null;
  }
  catch (e) {
    console.error(e);
    error.value = e instanceof Error ? e.message : '注文情報の取得に失敗しました';
  }
  finally {
    isLoading.value = false;
  }
});

const goToTop = () => {
  const facilityId = String(route.params.facilityId || '');
  const path = facilityId ? `/${facilityId}` : '/';
  router.push(path);
};
</script>

<template>
  <div>
    <div class="flex flex-col items-center justify-center bg-white">
      <img
        src="/complete.svg"
        alt="完了"
        class="w-40 h-40 mt-8"
      >
      <p class=" mt-8 text-base text-main">
        ご注文ありがとうございます！
      </p>
      <p class=" mt-2 text-base text-main">
        ご購入手続きが完了しました。
      </p>
      <div class=" mt-8 text-sm text-main">
        <template v-if="isLoading">
          注文番号を取得しています...
        </template>
        <template v-else-if="error">
          申し訳ありません。{{ error }}
        </template>
        <template v-else-if="orderId">
          お客様の注文番号は「{{ orderId }}」です
        </template>
        <template v-else>
          注文番号の取得ができませんでした。
        </template>
      </div>
    </div>
    <div class="px-4 text-center">
      <p class="mt-8 text-sm text-main">
        ご注文内容の詳細は、注文履歴からご確認いただけます。
      </p>
      <p class="mt-2 text-sm text-main">
        ご不明点がございましたら、LINEトーク画面からお問い合わせください。
      </p>
      <button
        class="mt-8 w-full bg-orange text-white py-3 px-4 rounded-lg font-semibold hover:bg-orange/[0.85]"
        @click="goToTop"
      >
        トップに戻る
      </button>
    </div>
  </div>
</template>
