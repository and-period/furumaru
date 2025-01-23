<script lang="ts" setup>
import { storeToRefs } from 'pinia'
import { useEventBus } from '@vueuse/core'
import { useAuthStore } from '~/store/auth'
import { useNotificationStore } from '~/store/notification'
import { useShoppingCartStore } from '~/store/shopping'
import type { I18n } from '~/types/locales'
import type { FooterMenuItem, HeaderMenuItem, LinkItem } from '~/types/props'

const router = useRouter()
const route = useRoute()
const i18n = useI18n()
const localePath = useLocalePath()

const notificationStore = useNotificationStore()
const { notifications } = storeToRefs(notificationStore)

const authStore = useAuthStore()
const { logout } = authStore
const { isAuthenticated, user } = storeToRefs(authStore)

const shoppingStore = useShoppingCartStore()
const { getCart, removeProductFromCart } = shoppingStore
const { cartIsEmpty, shoppingCart, totalPrice } = storeToRefs(shoppingStore)

getCart()

const appHeaderRef = ref<{ openCartMenu: () => void }>({
  openCartMenu: () => {},
})

const { on } = useEventBus('add-to-cart')
on(() => {
  appHeaderRef.value.openCartMenu()
})

const ht = (str: keyof I18n['layout']['header']) => {
  return i18n.t(`layout.header.${str}`)
}

const ft = (str: keyof I18n['layout']['footer']) => {
  return i18n.t(`layout.footer.${str}`)
}

const cartMenuMessage = computed<string>(() => {
  return i18n.t('layout.header.cartMenuMessage', {
    count: shoppingCart.value.carts?.length ?? 0,
  })
})

const navbarMenuList = computed<HeaderMenuItem[]>(() => [
  {
    text: ht('topLinkText'),
    to: localePath('/'),
    active: route.path === localePath('/'),
  },
  // {
  //   text: ht('searchItemLinkText'),
  //   to: localePath('/search'),
  //   active: route.path === localePath('/search'),
  // },
  {
    text: ht('allItemLinkText'),
    to: localePath('/items'),
    active: route.path === localePath('/items'),
  },
  {
    text: ht('allExperienceLinkText'),
    to: localePath('/experiences'),
    active: route.path === localePath('/experiences'),
  },
  {
    text: ht('allMarcheLinkText'),
    to: localePath('/marches'),
    active: route.path === localePath('/marches'),
  },
  {
    text: ht('blogLinkText'),
    to: localePath('/volunteer'),
    active: route.path === localePath('/volunteer'),
  },
  {
    text: ht('aboutLinkText'),
    to: localePath('/about'),
    active: route.path === localePath('/about'),
  },
])

const spModeMenuItems = computed<LinkItem[]>(() => [
  {
    text: ht('myPageLinkText'),
    href: localePath('/account'),
  },
  {
    text: ht('viewMyCartText'),
    href: localePath('/purchase'),
  },
])

const authenticatedMenuItems = computed<LinkItem[]>(() => [])

const noAuthenticatedMenuItems = computed<LinkItem[]>(() => [
  {
    text: ht('signIn'),
    href: localePath('/signin'),
  },
  {
    text: ht('signUp'),
    href: localePath('/signup'),
  },
])

const footerMenuList = computed<FooterMenuItem[]>(() => [
  {
    text: ft('qaLinkText'),
    to: '',
  },
  {
    text: ft('privacyPolicyLinkText'),
    to: '/privacy',
  },
  {
    text: ft('lawLinkText'),
    to: '/legal-notice',
  },
  {
    text: ft('inquiryLinkText'),
    to: '/contact',
  },
])

const isScrolled = ref<boolean>(false)

const onScroll = () => {
  if (!isScrolled.value && window.scrollY > 50) {
    isScrolled.value = true
  }
  else if (isScrolled.value && window.scrollY < 5) {
    isScrolled.value = false
  }
}

const handleClickBuyButton = () => {
  router.push('/purchase')
}

const handleClickRemoveItemFromCartButton = (
  cartNumber: number,
  id: string,
) => {
  removeProductFromCart(cartNumber, id)
}

const handleClickMyPageButton = () => {
  router.push('/account')
}

onMounted(() => {
  window.addEventListener('scroll', onScroll)
})

onUnmounted(() => {
  window.removeEventListener('scroll', onScroll)
})
</script>

<template>
  <div class="sticky top-0 z-[60]">
    <the-app-header
      ref="appHeaderRef"
      :home-path="localePath('/')"
      :is-authenticated="isAuthenticated"
      :user="user"
      :is-scrolled="isScrolled"
      :authenticated-account-menu-item="authenticatedMenuItems"
      :no-authenticated-account-menu-item="noAuthenticatedMenuItems"
      :menu-items="navbarMenuList"
      :notification-title="ht('notificationTitle')"
      :no-notification-item-text="ht('noNotificationItemText')"
      :notification-items="notifications"
      :total-price="totalPrice"
      :cart-is-empty="cartIsEmpty"
      :cart-items="shoppingCart.carts"
      :cart-menu-message="cartMenuMessage"
      :cart-total-price-text="ht('cartTotalPriceText')"
      :cart-total-price-tax-included-text="ht('cartTotalPriceTaxIncludedText')"
      :sign-in-link-text="ht('signIn')"
      :my-page-link-text="ht('myPageLinkText')"
      :view-mycart-text="ht('viewMyCartText')"
      :all-item-link-text="ht('allItemLinkText')"
      :all-marche-link-text="ht('allMarcheLinkText')"
      :blog-link-text="ht('blogLinkText')"
      :about-link-text="ht('aboutLinkText')"
      :number-of-carts-text="ht('numberOfCartsText')"
      :shipping-fee-annotation="ht('shippingFeeAnnotation')"
      :shipping-fee-annotation-link-text="ht('shippingFeeAnnotationLinkText')"
      :shipping-fee-annotation-check-text="ht('shippingFeeAnnotationCheckText')"
      :sp-menu-items="spModeMenuItems"
      :footer-menu-items="footerMenuList"
      @click:logout-button="logout"
      @click:my-page-button="handleClickMyPageButton"
      @click:buy-button="handleClickBuyButton"
      @click:remove-item-from-cart="handleClickRemoveItemFromCartButton"
    />
    <div
      id="snackbar"
      class="absolute top-[28px] z-[70] flex min-h-min w-full flex-col gap-2 bg-transparent md:top-[120px]"
    />
  </div>
  <div class="flex min-h-screen flex-col bg-base">
    <main class="grow overflow-hidden">
      <div class="mx-auto pb-16">
        <slot />
      </div>
    </main>
    <the-app-footer :menu-items="footerMenuList" />
  </div>
</template>
