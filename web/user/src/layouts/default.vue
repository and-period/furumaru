<script lang="ts" setup>
import { I18n } from '~/types/locales'
import { FooterMenuItem, HeaderMenuItem } from '~/types/props'

const router = useRouter()
const route = useRoute()
const i18n = useI18n()
const localePath = useLocalePath()

const ht = (str: keyof I18n['layout']['header']) => {
  return i18n.t(`layout.header.${str}`)
}

const ft = (str: keyof I18n['layout']['footer']) => {
  return i18n.t(`layout.footer.${str}`)
}

const handleCartClick = (): void => {
  console.log('NOT IMPLEMENTED')
}

const _localeRef = computed(() => {
  return i18n.locale === i18n.fallbackLocale ? '' : i18n.locale
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
  <div class="flex flex-col min-h-screen">
    <the-app-header
      :menu-items="navbarMenuList"
      @click:cart="handleCartClick"
    />
    <main class="bg-base flex-grow">
      <slot />
    </main>
    <the-app-footer :menu-items="footerMenuList" />
  </div>
</template>
