<script lang="ts" setup>
import { HeaderMenuItem, LinkItem } from '~/types/props'

interface Props {
  menuItems: HeaderMenuItem[]
  isAuthenticated: boolean
  authenticatedAccountMenuItem: LinkItem[]
  noAuthenticatedAccountMenuItem: LinkItem[]
  notificationTitle: string
  noNotificationItemText: string
  notificationItems: any[]
  cartIsEmpty: boolean
  cartMenuMessage: string
  cartItems: any[]
}

defineProps<Props>()
</script>

<template>
  <div class="flex px-10 py-4 bg-base items-center">
    <the-marche-logo class="m-0 " />
    <div class="flex items-center justify-end w-full text-main">
      <nav class="mr-16">
        <ul class="list-none flex gap-x-10">
          <li v-for="item, i in menuItems" :key="i">
            <a href="#" :class="{'border-b border-main pb-1': item.active}" @click="item.onClick">
              {{ item.text }}
            </a>
          </li>
        </ul>
      </nav>

      <div class="flex gap-x-8 items-center">
        <the-pc-account-menu
          :is-authenticated="isAuthenticated"
          :authenticated-menu-items="authenticatedAccountMenuItem"
          :no-authenticated-menu-items="noAuthenticatedAccountMenuItem"
        />

        <the-pc-notification-menu
          :title="notificationTitle"
          :no-item-text="noNotificationItemText"
          :notification-items="notificationItems"
        />

        <the-pc-cart-menu
          :is-authenticated="isAuthenticated"
          :cart-is-empty="cartIsEmpty"
          :cart-menu-message="cartMenuMessage"
          :cart-items="cartItems"
        />
      </div>
    </div>
  </div>
</template>
