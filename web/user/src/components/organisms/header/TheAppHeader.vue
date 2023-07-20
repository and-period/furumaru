<script lang="ts" setup>
import { FooterMenuItem, HeaderMenuItem, LinkItem } from '~/types/props'

interface Props {
  homePath: string
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
  spMenuItems: LinkItem[]
  footerMenuItems: FooterMenuItem[]
}

defineProps<Props>()

const spMenuOpen = ref<boolean>(false)

const handleClickMenuIconButton = () => {
  spMenuOpen.value = !spMenuOpen.value
}

const closeSpMenu = () => {
  spMenuOpen.value = false
}

const handleClickMenuItem = (item: HeaderMenuItem | FooterMenuItem) => {
  item.onClick()
  closeSpMenu()
}
</script>

<template>
  <div class="flex md:px-10 px-4 md:py-4 py-2 bg-base items-center relative">
    <nuxt-link :to="homePath" @click="closeSpMenu">
      <the-marche-logo class="m-0 md:w-64 w-32" />
    </nuxt-link>
    <div class="flex items-center justify-end w-full text-main">
      <nav class="mr-16 xl:block hidden">
        <ul class="list-none flex gap-x-10">
          <li v-for="item, i in menuItems" :key="i">
            <a href="#" :class="{'border-b border-main pb-1': item.active}" @click="item.onClick">
              {{ item.text }}
            </a>
          </li>
        </ul>
      </nav>

      <div class="flex items-center md:gap-x-8 gap-x-2">
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

        <the-icon-button class="xl:hidden h-10 w-10" @click="handleClickMenuIconButton">
          <the-outline-close-icon v-if="spMenuOpen" />
          <the-menu-icon v-else />
        </the-icon-button>
      </div>
    </div>
  </div>

  <div
    :class="{
      'absolute z-[50] bg-base block w-full h-screen transition duration-300 px-16': true,
      'opacity-100': spMenuOpen,
      'opacity-0 invisible': !spMenuOpen
    }"
  >
    <div class="flex flex-col gap-4 text-center text-main">
      <a v-for="item, i in menuItems" :key="i" href="#" @click="handleClickMenuItem(item)">
        {{ item.text }}
      </a>
      <hr class=" border-dashed border-main my-6">
      <nuxt-link v-for="item ,i in spMenuItems" :key="i" :to="item.href">
        {{ item.text }}
      </nuxt-link>

      <div class="flex justify-center gap-x-3 my-12">
        <the-icon-button>
          <the-instagram-icon id="header-instagram-icon" fill="#604C3F" />
        </the-icon-button>
        <the-icon-button>
          <the-facebook-icon id="header-facebook-icon" fill="#604C3F" />
        </the-icon-button>
      </div>

      <a v-for="item, i in footerMenuItems" :key="i" href="#" @click="handleClickMenuItem(item)">
        {{ item.text }}
      </a>
    </div>
  </div>
</template>
