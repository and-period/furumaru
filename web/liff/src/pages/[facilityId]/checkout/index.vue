<script setup lang="ts">
import { FmOrderSummary } from '@furumaru/shared';
import { useShoppingCartStore } from '~/stores/shopping';

const router = useRouter();
const shoppingCartStore = useShoppingCartStore();

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
