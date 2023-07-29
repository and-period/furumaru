<script lang="ts" setup>
interface Props {
  isAuthenticated: boolean
  cartIsEmpty: boolean
  cartMenuMessage: string
  cartItems: any[]
}

defineProps<Props>()

const totalPrice = computed(() => {
  return new Intl.NumberFormat('ja-JP', {
    style: 'currency',
    currency: 'JPY',
  }).format(18000)
})
</script>

<template>
  <the-dropdown-with-icon>
    <template #icon>
      <the-cart-icon id="header-cart-icon" fill="#604C3F" />
    </template>
    <template #content>
      <div
        class="flex max-h-[calc(100vh_-_150px)] flex-col gap-y-4 overflow-auto p-4 leading-8"
      >
        <p v-html="cartMenuMessage" />
        <hr class="border-main" />
        <div>
          合計金額:
          <p
            class="font-bold after:ml-2 after:text-[16px] after:content-['(税込)']"
          >
            {{ totalPrice }}
          </p>
        </div>
        <button class="w-full bg-main py-1 text-white">ログインして購入</button>

        <div class="border border-orange p-3 text-sm text-orange">
          現在のカゴの数: {{ cartItems.length }}
          <p>
            買い物カゴごとに送料がかかります。 詳しくは<a
              href="#"
              class="underline"
              >こちら</a
            >からご確認ください。
          </p>
        </div>

        <the-cart-item
          v-for="(item, i) in cartItems"
          :key="i"
          :cart-number="i + 1"
          :marche-name="item.marche"
          :box-type="item.boxType"
          :box-size="item.boxSize"
          :items="item.items"
        />
      </div>
    </template>
  </the-dropdown-with-icon>
</template>
