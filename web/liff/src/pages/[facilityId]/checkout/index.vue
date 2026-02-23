<script setup lang="ts">
import { FmOrderSummary, FmCreditCardForm } from '@furumaru/shared';
import type { CreditCardData } from '@furumaru/shared';
import { useShoppingCartStore } from '~/stores/shopping';
import { useCheckoutStore } from '~/stores/checkout';
import { PaymentMethodType } from '~/types/api/v1';
import type { CalcCartResponse, Product } from '~/types/api/facility/models';

const router = useRouter();
const route = useRoute();
const shoppingCartStore = useShoppingCartStore();
const checkoutStore = useCheckoutStore();

const userStore = useUserStore();
const { lastCheckInAt, profile } = storeToRefs(userStore);

const checkinData = computed(() => {
  if (!profile.value || !profile.value.lastCheckInAt) return null;
  const datetime = new Date(profile.value.lastCheckInAt * 1000);
  return datetime.toLocaleDateString('ja-JP', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
  });
});

// チェックイン翌日の日付（表示用）
const nextDayAfterCheckinData = computed(() => {
  if (!profile.value || !profile.value.lastCheckInAt) return null;
  const datetime = new Date(profile.value.lastCheckInAt * 1000);
  // 翌日
  datetime.setDate(datetime.getDate() + 1);
  return datetime.toLocaleDateString('ja-JP', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
  });
});

const facilityId = computed<string>(() => String(route.params.facilityId || ''));

onMounted(async () => {
  await initializeCart();
});

const summary = computed(() => {
  const carts = shoppingCartStore.shoppingCart.carts || [];
  const items = carts.flatMap(c => c.items);
  const coordinator = carts.find(c => c.coordinator)?.coordinator;
  const subtotal = items.reduce((sum, item) => sum + (item.product?.price ?? 0) * item.quantity, 0);
  const boxCarts = carts.map((c, idx) => ({ id: String(c.number ?? idx + 1) }));
  const coordinatorId = carts[0]?.coordinatorId ?? '';

  return {
    items,
    coordinator,
    carts: boxCarts,
    subtotal,
    discount: 0,
    total: subtotal,
    requestId: undefined as string | undefined,
    coordinatorId,
  };
});

const calculatedCart = ref<CalcCartResponse | null>(null);

const calculatedSummary = computed(() => {
  if (!calculatedCart.value) return null;

  const productMap = new Map<string, Product>(
    calculatedCart.value.products.map(product => [product.id, product]),
  );

  const items = calculatedCart.value.items
    .map((item) => {
      const product = productMap.get(item.productId);
      return {
        ...item,
        product: product
          ? {
              ...product,
              thumbnail: product.media.find(media => media.isThumbnail),
            }
          : undefined,
      };
    });

  const carts = calculatedCart.value.carts.map((cart, idx) => ({ id: String(cart.number ?? idx + 1) }));

  return {
    items,
    coordinator: calculatedCart.value.coordinator,
    carts,
    subtotal: calculatedCart.value.subtotal,
    discount: calculatedCart.value.discount,
    total: calculatedCart.value.total,
    requestId: calculatedCart.value.requestId,
    coordinatorId: calculatedCart.value.coordinator.id,
  };
});

const orderSummary = computed(() => calculatedSummary.value ?? summary.value);

const promotionCodeFormValue = ref('');
const validPromotion = ref(false);
const invalidPromotion = ref(false);
const isApplyingPromotion = ref(false);
const recalculateErrorMessage = ref<string | null>(null);

const recalculateCart = async (promotionCode?: string) => {
  if (!summary.value.coordinatorId) {
    calculatedCart.value = null;
    return;
  }

  calculatedCart.value = await shoppingCartStore.calcCartByCoordinatorId(
    facilityId.value,
    summary.value.coordinatorId,
    undefined,
    promotionCode,
  );
};

const initializeCart = async () => {
  try {
    await shoppingCartStore.getCart(facilityId.value);
    await recalculateCart();
    recalculateErrorMessage.value = null;
  }
  catch (e) {
    console.error('Failed to initialize cart:', e);
    recalculateErrorMessage.value = '金額の再計算に失敗しました。時間をおいて再度お試しください。';
  }
};

const handleClickUsePromotionCodeButton = async () => {
  const code = promotionCodeFormValue.value.trim();
  if (!code) {
    invalidPromotion.value = false;
    validPromotion.value = false;
    return;
  }

  try {
    isApplyingPromotion.value = true;
    await recalculateCart(code);
    promotionCodeFormValue.value = code;
    invalidPromotion.value = false;
    validPromotion.value = true;
    recalculateErrorMessage.value = null;
  }
  catch {
    invalidPromotion.value = true;
    validPromotion.value = false;
    try {
      await recalculateCart();
      recalculateErrorMessage.value = null;
    }
    catch (e) {
      console.error('Failed to reset cart pricing after invalid coupon:', e);
      recalculateErrorMessage.value = '金額の再計算に失敗しました。時間をおいて再度お試しください。';
    }
  }
  finally {
    isApplyingPromotion.value = false;
  }
};

const handleClickCancelPromotionCodeButton = async () => {
  promotionCodeFormValue.value = '';
  invalidPromotion.value = false;
  validPromotion.value = false;
  try {
    await recalculateCart();
    recalculateErrorMessage.value = null;
  }
  catch (e) {
    console.error('Failed to recalculate cart after coupon cancellation:', e);
    calculatedCart.value = null;
    recalculateErrorMessage.value = '金額の再計算に失敗しました。時間をおいて再度お試しください。';
  }
};

// 以前の「決済プラン情報」表示は削除し、支払いフォームを表示します。

// 支払いフォーム（クレジットカード）
const creditCard = ref<CreditCardData>({
  name: '',
  number: '',
  month: 0,
  year: 0,
  verificationValue: '',
});

// 受け取り日時
const pickupAtFormData = ref<string>('checkin'); // 受け取り日時

const pickupAt = computed<number>(() => {
  const ts = profile.value?.lastCheckInAt;
  if (!ts) return 0;

  // 基準: チェックイン日のローカル日付
  const base = new Date(ts * 1000);
  const year = base.getFullYear();
  const month = base.getMonth(); // 0-based
  const day = base.getDate();

  if (pickupAtFormData.value === 'checkin') {
    return ts;
  }
  if (pickupAtFormData.value === 'afternoon') {
    // 当日 17:00
    const d = new Date(year, month, day, 17, 0, 0, 0);
    return Math.floor(d.getTime() / 1000);
  }
  if (pickupAtFormData.value === 'morning') {
    // 翌日 08:00
    const d = new Date(year, month, day + 1, 8, 0, 0, 0);
    return Math.floor(d.getTime() / 1000);
  }
  return 0;
});

const isSubmitting = ref(false);
const submitError = ref<string | null>(null);
const PAYMENT_METHOD_CARD = PaymentMethodType.PaymentMethodTypeCreditCard; // クレジットカード決済（仮のコード）

const priceFormatter = (price: number): string => {
  return new Intl.NumberFormat('ja-JP', {
    style: 'currency',
    currency: 'JPY',
  }).format(price);
};

const handlePay = async () => {
  submitError.value = null;
  if (!orderSummary.value.coordinator) return;

  // 簡易バリデーション
  if (!creditCard.value.number || !creditCard.value.name || !creditCard.value.month || !creditCard.value.year || !creditCard.value.verificationValue) {
    submitError.value = 'カード情報を入力してください';
    return;
  }

  try {
    isSubmitting.value = true;
    const facilityId = String(route.params.facilityId || '');
    const callbackUrl = `${window.location.origin}/${facilityId}/complete`;

    const res = await checkoutStore.startCheckout(facilityId, {
      callbackUrl,
      paymentMethod: PAYMENT_METHOD_CARD,
      creditCard: {
        name: creditCard.value.name,
        number: creditCard.value.number,
        month: creditCard.value.month,
        year: creditCard.value.year,
        verificationValue: creditCard.value.verificationValue,
      },
      pickupAt: pickupAt.value,
      promotionCode: validPromotion.value ? promotionCodeFormValue.value : undefined,
      requestId: orderSummary.value.requestId,
      total: orderSummary.value.total,
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
            :items="orderSummary.items"
            :coordinator="orderSummary.coordinator"
            :carts="orderSummary.carts"
            :subtotal="orderSummary.subtotal"
            :discount="orderSummary.discount"
            :total="orderSummary.total"
          />

          <template v-if="validPromotion">
            <div class="mt-4 flex justify-between rounded-lg border border-orange p-2 text-sm text-orange">
              <div class="flex items-center gap-1">
                クーポンコードを適用しました
              </div>
              <button @click="handleClickCancelPromotionCodeButton">
                解除
              </button>
            </div>
            <div class="mt-2 rounded-md bg-orange/5 p-3 text-sm text-gray-700">
              <div class="flex items-center justify-between">
                <span>商品合計</span>
                <span>{{ priceFormatter(orderSummary.subtotal) }}</span>
              </div>
              <div class="mt-1 flex items-center justify-between text-orange">
                <span>クーポン割引</span>
                <span>-{{ priceFormatter(orderSummary.discount) }}</span>
              </div>
              <div class="mt-1 flex items-center justify-between font-semibold">
                <span>適用後合計</span>
                <span>{{ priceFormatter(orderSummary.total) }}</span>
              </div>
            </div>
          </template>

          <template v-else>
            <div class="mt-4 flex gap-2">
              <div class="grow">
                <input
                  v-model="promotionCodeFormValue"
                  type="text"
                  class="w-full border border-gray-300 bg-gray-50 p-2.5 text-sm"
                  placeholder="クーポンコード"
                >
              </div>
              <button
                class="whitespace-nowrap bg-orange px-4 py-2 text-sm text-white disabled:bg-gray-300"
                :disabled="isApplyingPromotion"
                @click="handleClickUsePromotionCodeButton"
              >
                {{ isApplyingPromotion ? '適用中...' : '適用' }}
              </button>
            </div>
            <div
              v-if="invalidPromotion"
              class="mt-2 px-1 text-xs text-red-600"
            >
              クーポンコードが無効です
            </div>
          </template>

          <div
            v-if="recalculateErrorMessage"
            class="mt-2 px-1 text-xs text-red-600"
          >
            {{ recalculateErrorMessage }}
          </div>

          <!-- 受け取り日時選択 -->
          <div
            v-if="false"
            class="flex flex-col px-2 gap-1"
          >
            <label class=" font-semibold text-sm block mb-2">配達時間</label>
            <div>
              <input
                id="checkin"
                v-model="pickupAtFormData"
                type="radio"
                value="checkin"
                class="mx-2 accent-orange"
              >
              <label for="checkin">
                チェックイン時間（{{ lastCheckInAt }}）
              </label>
            </div>
            <div>
              <input
                id="afternoon"
                v-model="pickupAtFormData"
                type="radio"
                value="afternoon"
                class="mx-2 accent-orange"
              >
              <label for="afternoon">
                当日午後（{{ checkinData }} 17:00）
              </label>
            </div>
            <div>
              <input
                id="morning"
                v-model="pickupAtFormData"
                type="radio"
                value="morning"
                class="mx-2 accent-orange"
              >
              <label for="morning">
                翌日午前（{{ nextDayAfterCheckinData }} 08:00）
              </label>
            </div>
          </div>

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
