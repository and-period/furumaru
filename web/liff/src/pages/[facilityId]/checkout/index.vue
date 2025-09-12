<script setup lang="ts">
import { FmOrderSummary, FmCreditCardForm } from '@furumaru/shared';
import type { CreditCardData } from '@furumaru/shared';
import { useShoppingCartStore } from '~/stores/shopping';
import { useCheckoutStore } from '~/stores/checkout';
import { PaymentMethodType } from '~/types/api/v1';

const router = useRouter();
const route = useRoute();
const shoppingCartStore = useShoppingCartStore();
const checkoutStore = useCheckoutStore();

onMounted(() => {
  shoppingCartStore.getCart();
});

const summary = computed(() => {
  const carts = shoppingCartStore.shoppingCart.carts || [];
  const items = carts.flatMap(c => c.items);
  const coordinator = carts.find(c => c.coordinator)?.coordinator;
  const subtotal = items.reduce((sum, item) => sum + (item.product?.price ?? 0) * item.quantity, 0);
  const boxCarts = carts.map((c, idx) => ({ id: String(c.number ?? idx + 1) }));

  return { items, coordinator, carts: boxCarts, subtotal };
});

// 以前の「決済プラン情報」表示は削除し、支払いフォームを表示します。

// 支払いフォーム（クレジットカード）
const creditCard = ref<CreditCardData>({
  name: '',
  number: '',
  month: 0,
  year: 0,
  verificationValue: '',
});

const isSubmitting = ref(false);
const submitError = ref<string | null>(null);
const PAYMENT_METHOD_CARD = PaymentMethodType.PaymentMethodTypeCreditCard; // クレジットカード決済（仮のコード）

const handlePay = async () => {
  submitError.value = null;
  if (!summary.value.coordinator) return;

  // 簡易バリデーション
  if (!creditCard.value.number || !creditCard.value.name || !creditCard.value.month || !creditCard.value.year || !creditCard.value.verificationValue) {
    submitError.value = 'カード情報を入力してください';
    return;
  }

  try {
    isSubmitting.value = true;
    const facilityId = String(route.params.facilityId || '');
    const callbackUrl = `${window.location.origin}/${facilityId}/complete`;

    const res = await checkoutStore.startCheckout({
      callbackUrl,
      paymentMethod: PAYMENT_METHOD_CARD,
      creditCard: {
        name: creditCard.value.name,
        number: creditCard.value.number,
        month: creditCard.value.month,
        year: creditCard.value.year,
        verificationValue: creditCard.value.verificationValue,
      },
      total: summary.value.subtotal,
    });

    const url = res.url || checkoutStore.redirectUrl;
    if (url) {
      window.location.href = url;
    }
    else {
      submitError.value = '決済URLを取得できませんでした。';
    }
  }
  catch (e) {
    console.error(e);
    submitError.value = e instanceof Error ? e.message : '決済に失敗しました';
  }
  finally {
    isSubmitting.value = false;
  }
};
</script>

<template>
  <div class="mt-[50px]">
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

    <div class="px-4 py-6">
      <h2 class="mt-2 font-semibold font-inter text-center w-full">
        買い物カゴ
      </h2>

      <div
        v-if="shoppingCartStore.cartIsEmpty"
        class="mt-6 text-center text-gray-600"
      >
        カートに商品がありません。
      </div>

      <div
        v-else
        class="mt-6"
      >
        <template v-if="summary.coordinator">
          <fm-order-summary
            :items="summary.items"
            :coordinator="summary.coordinator"
            :carts="summary.carts"
            :subtotal="summary.subtotal"
            :discount="0"
            :total="summary.subtotal"
          />

          <!-- 支払いフォーム -->
          <div class="mt-6 bg-white rounded-md border p-4">
            <div class="text-sm font-semibold mb-2">
              お支払い情報
            </div>
            <fm-credit-card-form v-model="creditCard" />
            <div
              v-if="submitError"
              class="mt-3 text-sm text-red-600"
            >
              {{ submitError }}
            </div>
            <button
              class="mt-4 w-full bg-orange text-white py-3 px-4 rounded-lg font-semibold hover:bg-orange/[0.85] disabled:bg-gray-300 disabled:cursor-not-allowed"
              :disabled="isSubmitting || checkoutStore.isLoading"
              @click="handlePay"
            >
              {{ isSubmitting || checkoutStore.isLoading ? '処理中...' : '支払う' }}
            </button>
          </div>
        </template>
        <template v-else>
          <div class="text-center text-gray-600">
            コーディネーター情報が見つかりませんでした。
          </div>
        </template>
      </div>
    </div>
  </div>
</template>
