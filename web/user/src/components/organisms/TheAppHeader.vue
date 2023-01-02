<template>
  <v-app-bar flat color="base">
    <v-toolbar-title class="text-accent text-h6">Online Marche</v-toolbar-title>
    <v-spacer />
    <slot />
    <v-tooltip location="bottom">
      <template #activator="{ on, attrs }">
        <v-btn icon class="mr-4" v-bind="attrs" v-on="on" @click="handleCartClick">
          <v-badge :content="cartContent" :value="!cartIsEmpty">
            <v-icon>mdi-cart-outline</v-icon>
          </v-badge>
        </v-btn>
      </template>
      <span>{{ cartIsEmpty ? props.cartEmptyMessage : props.cartNotEmptyMessage }}</span>
    </v-tooltip>
    <v-menu>
      <template #activator="{ on }">
        <molecules-the-header-menu-button :img-src="props.profileImgUrl" v-on="on" />
      </template>
      <v-list>
        <v-list-item v-for="(menuItem, i) in props.menuList" :key="i" @click="menuItem.onClick">
          <v-list-item-title>{{ menuItem.name }}</v-list-item-title>
        </v-list-item>
      </v-list>
    </v-menu>
  </v-app-bar>
</template>

<script lang="ts" setup>
import { PropType } from 'vue'

import { HeaderMenuItem } from '~/types/props'

const props = defineProps({
  profileImgUrl: {
    type: String,
    required: false,
    default: null,
  },
  cartItemCount: {
    type: Number,
    required: true,
    default: 0,
  },
  cartEmptyMessage: {
    type: String,
    required: true,
  },
  cartNotEmptyMessage: {
    type: String,
    required: true,
  },
  menuList: {
    type: Array as PropType<HeaderMenuItem[]>,
    default: () => {
      return []
    },
  },
})

const emits = defineEmits<{
  (name: 'click:cart'): void
}>()

const cartItemCount = computed<number>(() => {
  return props.cartItemCount
})
const cartIsEmpty = computed<boolean>(() => {
  return cartItemCount.value === 0
})
const cartContent = computed<number | string>(() => {
  return cartItemCount.value > 99 ? '99+' : cartItemCount.value
})

const handleCartClick = (): void => {
  emits('click:cart')
}
</script>

<style lang="scss" scoped>
.title {
  font-family: 'STIX Two Text' !important;
  font-style: normal;
  font-weight: 400;
}
</style>
