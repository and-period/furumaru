<script lang="ts" setup>
interface Props {
  isAuthenticated: boolean
  cartIsEmpty: boolean
  cartMenuMessage: string
  cartItems: any[]
}

defineProps<Props>()

const totalPrice = computed(() => {
  return new Intl.NumberFormat('ja-JP', { style: 'currency', currency: 'JPY' }).format(18000)
})
</script>

<template>
  <the-dropdown-with-icon>
    <template #icon>
      <the-cart-icon id="header-cart-icon" fill="#604C3F" />
    </template>
    <template #content>
      <div class="p-4 leading-8 flex flex-col gap-y-4 max-h-[calc(100vh_-_150px)] overflow-auto">
        <p v-html="cartMenuMessage" />
        <hr class=" border-main">
        <div>
          合計金額:
          <p class="font-bold after:ml-2 after:content-['(税込)'] after:text-[16px]">
            {{ totalPrice }}
          </p>
        </div>
        <button class="py-1 bg-main text-white w-full">
          ログインして購入
        </button>

        <div class="p-3 border border-orange text-orange text-sm">
          現在のカゴの数: {{ cartItems.length }}
          <p>
            買い物カゴごとに送料がかかります。
            詳しくは<a href="#" class="underline">こちら</a>からご確認ください。
          </p>
        </div>

        <the-cart-item
          v-for="item, i in cartItems"
          :key="i"
          :cart-number="i+1"
          :marche-name="item.marche"
          :box-type="item.boxType"
          :box-size="item.boxSize"
          :items="item.items"
        />
      </div>
    </template>
  </the-dropdown-with-icon>
</template>
