<script lang="ts" setup>
import { I18n } from '~/types/locales'
import { HeaderMenuItem } from '~/types/props'

const router = useRouter()
const route = useRoute()
const i18n = useI18n()
const localePath = useLocalePath()

const t = (str: keyof I18n['layout']['header']) => {
  return i18n.t(`layout.header.${str}`)
}

const handleCartClick = (): void => {
  console.log('NOT IMPLEMENTED')
}

const _localeRef = computed(() => {
  return i18n.locale === i18n.fallbackLocale ? '' : i18n.locale
})

const navbarMenuList = computed<HeaderMenuItem[]>(() => [
  {
    text: t('topLinkText'),
    onClick: () => router.push(localePath('/')),
    active: route.path === localePath('/')
  },
  {
    text: t('searchItemLinkText'),
    onClick: () => router.push(localePath('/search')),
    active: route.path === localePath('/search')
  },
  {
    text: t('allItemLinkText'),
    onClick: () => router.push(localePath('/items')),
    active: route.path === localePath('/search')
  },
  {
    text: t('aboutLinkText'),
    onClick: () => router.push(localePath('/about')),
    active: route.path === localePath('/about')
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
    <the-app-footer />
  </div>
</template>
