<script lang="ts" setup>
import { storeToRefs } from 'pinia'
import { useAuthStore } from '~/store/auth'
import { useNotificationStore } from '~/store/notification'
import { useShoppingStore } from '~/store/shopping'
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

const shoppingStore = useShoppingStore()
const { cartIsEmpty, cartItems } = storeToRefs(shoppingStore)

const ht = (str: keyof I18n['layout']['header']) => {
  return i18n.t(`layout.header.${str}`)
}

const ft = (str: keyof I18n['layout']['footer']) => {
  return i18n.t(`layout.footer.${str}`)
}

const cartMenuMessage = computed<string>(() => {
  return i18n.t('layout.header.cartMenuMessage', { count: cartItems.value.length })
})

const navbarMenuList = computed<HeaderMenuItem[]>(() => [
  {
    text: ht('topLinkText'),
    onClick: () => router.push(localePath('/')),
    active: route.path === localePath('/')
  },
  {
    text: ht('searchItemLinkText'),
    onClick: () => router.push(localePath('/search')),
    active: route.path === localePath('/search')
  },
  {
    text: ht('allItemLinkText'),
    onClick: () => router.push(localePath('/items')),
    active: route.path === localePath('/search')
  },
  {
    text: ht('aboutLinkText'),
    onClick: () => router.push(localePath('/about')),
    active: route.path === localePath('/about')
  }
])

const authenticatedMenuItems = computed<LinkItem[]>(() => [])

const noAuthenticatedMenuItems = computed<LinkItem[]>(() => [
  {
    text: ht('signIn'),
    href: localePath('/signin')
  },
  {
    text: ht('signUp'),
    href: localePath('/signup')
  }
])

const footerMenuList = computed<FooterMenuItem[]>(() => [
  {
    text: ft('qaLinkText'),
    onClick: () => {}
  },
  {
    text: ft('privacyPolicyLinkText'),
    onClick: () => {}
  },
  {
    text: ft('lawLinkText'),
    onClick: () => {}
  },
  {
    text: ft('inquiryLinkText'),
    onClick: () => {}
  }
])
</script>

<template>
  <div class="flex flex-col min-h-screen bg-base">
    <the-app-header
      :is-authenticated="isAuthenticated"
      :authenticated-account-menu-item="authenticatedMenuItems"
      :no-authenticated-account-menu-item="noAuthenticatedMenuItems"
      :menu-items="navbarMenuList"
      :notification-title="ht('notificationTitle')"
      :no-notification-item-text="ht('noNotificationItemText')"
      :notification-items="notifications"
      :cart-is-empty="cartIsEmpty"
      :cart-items="cartItems"
      :cart-menu-message="cartMenuMessage"
    />
    <main class="flex-grow overflow-hidden">
      <div class="container pb-10 mx-auto">
        <slot />
      </div>
    </main>
    <the-app-footer :menu-items="footerMenuList" />
  </div>
</template>
