<script setup lang="ts">
type IconName =
  | 'account'
  | 'ring'
  | 'cart'
  | 'search'
  | 'fruits'
  | 'flag'
  | 'furumaru'

interface Props {
  icon: IconName
}

defineProps<Props>()
</script>

<template>
  <template v-if="icon === 'account'">
    <the-account-icon />
  </template>

  <template v-if="icon === 'ring'">
    <the-ring-icon />
  </template>

  <template v-if="icon === 'cart'">
    <the-cart-icon
      id="cart-icon-xxx"
      fill="#604C3F"
    />
  </template>

  <template v-if="icon === 'search'">
    <the-search-icon />
  </template>

  <template v-if="icon === 'fruits'">
    <the-fruits-icon />
  </template>

  <template v-if="icon === 'flag'">
    <the-flag-icon />
  </template>

  <template v-if="icon === 'furumaru'">
    <the-furumaru-icon />
  </template>
</template>
