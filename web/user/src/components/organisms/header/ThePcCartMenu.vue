<script lang="ts" setup>
interface Props {
  isAuthenticated: boolean
  cartIsEmpty: boolean
  cartMenuMessage: string
  cartItems: any[]
}

defineProps<Props>()
</script>

<template>
  <the-dropdown-with-icon>
    <template #icon>
      <the-cart-icon id="header-cart-icon" fill="#604C3F" />
    </template>
    <template #content>
      <div class="p-4 leading-8">
        <p v-html="cartMenuMessage" />
        <div v-if="!cartIsEmpty" class="mt-4">
          <hr class=" border-main">
        </div>
      </div>
    </template>
  </the-dropdown-with-icon>
</template>
