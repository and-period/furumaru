<script lang="ts" setup>
import type { AuthUserResponse } from '~/types/api'
import type { FooterMenuItem, HeaderMenuItem, LinkItem } from '~/types/props'
import type { ShoppingCart } from '~/types/store'

interface Props {
  homePath: string
  menuItems: HeaderMenuItem[]
  isScrolled: boolean
  isAuthenticated: boolean
  user: AuthUserResponse | undefined
  authenticatedAccountMenuItem: LinkItem[]
  noAuthenticatedAccountMenuItem: LinkItem[]
  notificationTitle: string
  noNotificationItemText: string
  notificationItems: any[]
  cartIsEmpty: boolean
  cartMenuMessage: string
  cartTotalPriceText: string
  cartTotalPriceTaxIncludedText: string
  totalPrice: number
  cartItems: ShoppingCart[]
  signInLinkText: string
  myPageLinkText: string
  allItemLinkText: string
  allMarcheLinkText: string
  volunteerLinkText: string
  aboutLinkText: string
  viewMycartText: string
  numberOfCartsText: string
  shippingFeeAnnotation: string
  shippingFeeAnnotationLinkText: string
  shippingFeeAnnotationCheckText: string
  spMenuItems: LinkItem[]
  footerMenuItems: FooterMenuItem[]
}

const props = defineProps<Props>()

interface Emits {
  (e: 'click:buyButton'): void
  (e: 'click:removeItemFromCart', cartNumber: number, id: string): void
  (e: 'click:myPageButton'): void
  (e: 'click:logoutButton'): void
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

const handleClickRemoveItemFromCartButton = (
  cartNumber: number,
  id: string,
) => {
  emits('click:removeItemFromCart', cartNumber, id)
}

const SP_MENU_ITEMS = computed(() => [
  {
    icon: 'account',
    text: props.isAuthenticated ? props.myPageLinkText : props.signInLinkText,
    to: props.isAuthenticated ? '/account' : '/signin',
  },
  // {
  //   icon: 'ring',
  //   text: 'お知らせ',
  //   to: '/',
  // },
  {
    icon: 'cart',
    text: props.viewMycartText,
    to: '/purchase',
  },
  // {
  //   icon: 'search',
  //   text: '商品をさがす',
  //   to: '/',
  // },
  {
    icon: 'fruits',
    text: props.allItemLinkText,
    to: '/items',
  },
  {
    icon: 'flag',
    text: props.allMarcheLinkText,
    to: '/marches',
  },
  {
    icon: 'furumaru',
    text: props.volunteerLinkText,
    to: '/volunteer',
  },
  {
    icon: 'furumaru',
    text: props.aboutLinkText,
    to: '/about',
  },
])

const handleClickMyPageButton = () => {
  emits('click:myPageButton')
}

const handleClickLogoutButton = () => {
  emits('click:logoutButton')
}
</script>

<template>
  <div
    :class="{
      'relative flex h-[64px] items-center justify-between bg-base px-4 py-2 duration-300 ease-in-out md:px-10 md:py-4': true,
      'md:h-[80px]': isScrolled,
      'md:h-[116px]': !isScrolled,
    }"
  >
    <div class="md:hidden">
      <the-icon-button
        class="h-10 w-10 md:hidden"
        @click="handleClickMenuIconButton"
      >
        <the-outline-close-icon v-if="spMenuOpen" />
        <the-menu-icon v-else />
      </the-icon-button>
    </div>

    <nuxt-link
      :to="homePath"
      class="flex h-full"
      @click="closeSpMenu"
    >
      <the-marche-logo class="m-0 max-h-full max-w-full items-center" />
    </nuxt-link>

    <div class="flex items-center text-main">
      <nav class="mr-16 hidden xl:block">
        <ul class="flex list-none gap-x-10">
          <li
            v-for="(item, i) in menuItems"
            :key="i"
          >
            <nuxt-link
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
          class="hidden md:block"
          :is-authenticated="isAuthenticated"
          :user="user"
          :authenticated-menu-items="authenticatedAccountMenuItem"
          :no-authenticated-menu-items="noAuthenticatedAccountMenuItem"
          @click:logout-button="handleClickLogoutButton"
          @click:my-page-button="handleClickMyPageButton"
        />

        <the-pc-notification-menu
          class="hidden md:block"
          :title="notificationTitle"
          :no-item-text="noNotificationItemText"
          :notification-items="notificationItems"
        />

        <the-pc-cart-menu
          :is-authenticated="isAuthenticated"
          :cart-is-empty="cartIsEmpty"
          :cart-menu-message="cartMenuMessage"
          :cart-total-price-text="cartTotalPriceText"
          :cart-total-price-tax-included-text="cartTotalPriceTaxIncludedText"
          :cart-items="cartItems"
          :view-mycart-text="viewMycartText"
          :number-of-carts-text="numberOfCartsText"
          :shipping-fee-annotation="shippingFeeAnnotation"
          :shipping-fee-annotation-link-text="shippingFeeAnnotationLinkText"
          :shipping-fee-annotation-check-text="shippingFeeAnnotationCheckText"
          :total-price="totalPrice"
          @click:buy-button="handleClickBuyButton"
          @click:remove-item-from-cart="handleClickRemoveItemFromCartButton"
        />

        <the-icon-button
          class="hidden h-10 w-10 md:block xl:hidden"
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
      'absolute z-[50] block h-screen w-full bg-base p-4 transition duration-300 md:px-16': true,
      'opacity-100': spMenuOpen,
      'invisible opacity-0': !spMenuOpen,
    }"
  >
    <div class="flex flex-col gap-4 text-center text-main">
      <div class="flex flex-col text-left">
        <nuxt-link
          v-for="(item, i) in SP_MENU_ITEMS"
          :key="i"
          :to="item.to"
          class="grid grid-cols-12 items-center border-b border-dashed border-main py-4"
          @click="handleClickMenuItem"
        >
          <div class="col-span-1 flex justify-center">
            <the-icon-wrapper
              :icon="item.icon"
              class="h-6 w-6"
            />
          </div>
          <div class="col-span-10 pl-4">
            {{ item.text }}
          </div>
          <div class="col-span-1">
            <the-right-arrow-icon class="h-4 w-4" />
          </div>
        </nuxt-link>
      </div>

      <div class="my-12 flex justify-center gap-x-3">
        <the-icon-button>
          <the-instagram-icon
            id="header-instagram-icon"
            fill="#604C3F"
          />
        </the-icon-button>
        <the-icon-button>
          <the-facebook-icon
            id="header-facebook-icon"
            fill="#604C3F"
          />
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
