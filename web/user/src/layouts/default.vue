<script lang="ts" setup>
import { storeToRefs } from 'pinia'
import { MOCK_CART_ITEMS } from '~/constants/mock'
import { useAuthStore } from '~/store/auth'
import { useNotificationStore } from '~/store/notification'
import { useShoppingCartStore } from '~/store/shopping'
import { I18n } from '~/types/locales'
import { FooterMenuItem, HeaderMenuItem, LinkItem } from '~/types/props'

const router = useRouter()
const route = useRoute()
const i18n = useI18n()
const localePath = useLocalePath()

const notificationStore = useNotificationStore()
const { notifications } = storeToRefs(notificationStore)

const authStore = useAuthStore()
const { isAuthenticated } = storeToRefs(authStore)

const shoppingStore = useShoppingCartStore()
const { getCart } = shoppingStore
const { cartIsEmpty } = storeToRefs(shoppingStore)
const cartItems = MOCK_CART_ITEMS

getCart()

const ht = (str: keyof I18n['layout']['header']) => {
  return i18n.t(`layout.header.${str}`)
}

const ft = (str: keyof I18n['layout']['footer']) => {
  return i18n.t(`layout.footer.${str}`)
}

const cartMenuMessage = computed<string>(() => {
  return i18n.t('layout.header.cartMenuMessage', { count: cartItems.length })
})

const navbarMenuList = computed<HeaderMenuItem[]>(() => [
  {
    text: ht('topLinkText'),
    to: localePath('/'),
    active: route.path === localePath('/'),
  },
  {
    text: ht('searchItemLinkText'),
    to: localePath('/search'),
    active: route.path === localePath('/search'),
  },
  {
    text: ht('allItemLinkText'),
    to: localePath('/items'),
    active: route.path === localePath('/items'),
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
    href: localePath('/mypage'),
  },
  {
    text: ht('viewMyCartText'),
    href: localePath('/cart'),
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
    to: '',
  },
  {
    text: ft('lawLinkText'),
    to: '/legal-notice',
  },
  {
    text: ft('inquiryLinkText'),
    to: '',
  },
])

const isScrolled = ref<boolean>(false)

const onScroll = () => {
  if (!isScrolled.value && window.scrollY > 50) {
    isScrolled.value = true
  } else if (isScrolled.value && window.scrollY < 30) {
    isScrolled.value = false
  }
}

onMounted(() => {
  window.addEventListener('scroll', onScroll)
})

onUnmounted(() => {
  window.removeEventListener('scroll', onScroll)
})

const handleClickBuyButton = () => {
  router.push('/purchase')
}
</script>

<template>
  <div class="sticky top-0 z-[60]">
    <the-app-header
      :home-path="localePath('/')"
      :is-authenticated="isAuthenticated"
      :is-scrolled="isScrolled"
      :authenticated-account-menu-item="authenticatedMenuItems"
      :no-authenticated-account-menu-item="noAuthenticatedMenuItems"
      :menu-items="navbarMenuList"
      :notification-title="ht('notificationTitle')"
      :no-notification-item-text="ht('noNotificationItemText')"
      :notification-items="notifications"
      :cart-is-empty="cartIsEmpty"
      :cart-items="cartItems"
      :cart-menu-message="cartMenuMessage"
      :sp-menu-items="spModeMenuItems"
      :footer-menu-items="footerMenuList"
      @click:buy-button="handleClickBuyButton"
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
