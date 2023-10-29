<script lang="ts" setup>
import { FooterMenuItem, HeaderMenuItem, LinkItem } from '~/types/props'

interface Props {
  homePath: string
  menuItems: HeaderMenuItem[]
  isScrolled: boolean
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

interface Emits {
  (e: 'click:buyButton'): void
}

const emits = defineEmits<Emits>()

const spMenuOpen = ref<boolean>(false)

const handleClickMenuIconButton = () => {
  spMenuOpen.value = !spMenuOpen.value
}

const closeSpMenu = () => {
  spMenuOpen.value = false
}

const handleClickMenuItem = (item: HeaderMenuItem | FooterMenuItem) => {
  closeSpMenu()
}

const handleClickBuyButton = () => {
  emits('click:buyButton')
  closeSpMenu()
}
</script>

<template>
  <div
    :class="{
      'relative flex items-center bg-base px-4 py-2 duration-300 ease-in-out md:px-10 md:py-4 ': true,
      'md:h-[80px]': isScrolled,
      'md:h-[116px]': !isScrolled,
    }"
  >
    <nuxt-link :to="homePath" class="flex h-full" @click="closeSpMenu">
      <the-marche-logo class="m-0 max-h-full max-w-full items-center" />
    </nuxt-link>

    <div class="flex w-full items-center justify-end text-main">
      <nav class="mr-16 hidden xl:block">
        <ul class="flex list-none gap-x-10">
          <li v-for="(item, i) in menuItems" :key="i">
            <nuxt-link
              href="#"
              :class="{ 'border-b border-main pb-1': item.active }"
              :to="item.to"
            >
              {{ item.text }}
            </nuxt-link>
          </li>
        </ul>
      </nav>

      <div class="flex items-center gap-x-2 md:gap-x-8">
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
          @click:buy-button="handleClickBuyButton"
        />

        <the-icon-button
          class="h-10 w-10 xl:hidden"
          @click="handleClickMenuIconButton"
        >
          <the-outline-close-icon v-if="spMenuOpen" />
          <the-menu-icon v-else />
        </the-icon-button>
      </div>
    </div>
  </div>

  <div
    :class="{
      'absolute z-[50] block h-screen w-full bg-base px-16 transition duration-300': true,
      'opacity-100': spMenuOpen,
      'invisible opacity-0': !spMenuOpen,
    }"
  >
    <div class="flex flex-col gap-4 text-center text-main">
      <a
        v-for="(item, i) in menuItems"
        :key="i"
        href="#"
        @click="handleClickMenuItem(item)"
      >
        {{ item.text }}
      </a>
      <hr class="my-6 border-dashed border-main" />
      <nuxt-link v-for="(item, i) in spMenuItems" :key="i" :to="item.href">
        {{ item.text }}
      </nuxt-link>

      <div class="my-12 flex justify-center gap-x-3">
        <the-icon-button>
          <the-instagram-icon id="header-instagram-icon" fill="#604C3F" />
        </the-icon-button>
        <the-icon-button>
          <the-facebook-icon id="header-facebook-icon" fill="#604C3F" />
        </the-icon-button>
      </div>

      <nuxt-link
        v-for="(item, i) in footerMenuItems"
        :key="i"
        :to="item.to"
        @click="handleClickMenuItem(item)"
      >
        {{ item.text }}
      </nuxt-link>
    </div>
  </div>
</template>
