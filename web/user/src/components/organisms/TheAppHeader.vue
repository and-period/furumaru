<template>
  <div color="base">
    aaa
  </div>
</template>

<script lang="ts" setup>
import { PropType } from 'vue'

import { HeaderMenuItem } from '~/types/props'

const props = defineProps({
  profileImgUrl: {
    type: String,
    required: false,
    default: null
  },
  cartItemCount: {
    type: Number,
    required: true,
    default: 0
  },
  cartEmptyMessage: {
    type: String,
    required: true
  },
  cartNotEmptyMessage: {
    type: String,
    required: true
  },
  menuList: {
    type: Array as PropType<HeaderMenuItem[]>,
    default: () => {
      return []
    }
  }
})

const emits = defineEmits<{(name: 'click:cart'): void}>()

const cartItemCount = computed<number>(() => {
  return props.cartItemCount
})
const _cartIsEmpty = computed<boolean>(() => {
  return cartItemCount.value === 0
})
const _cartContent = computed<number | string>(() => {
  return cartItemCount.value > 99 ? '99+' : cartItemCount.value
})

const _handleCartClick = (): void => {
  emits('click:cart')
}
</script>
