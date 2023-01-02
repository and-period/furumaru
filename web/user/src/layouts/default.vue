<template>
  <v-app>
    <organisms-the-app-header
      :cart-item-count="0"
      :cart-empty-message="t('cartEmptyMessage')"
      :cart-not-empty-message="t('cartNotEmptyMessage')"
      :menu-list="headerMenuList"
      @click:cart="handleCartClick"
    >
      <nuxt-link to="/" class="mr-4 header-link">
        {{ t('becomeShopOwner') }}
      </nuxt-link>
    </organisms-the-app-header>
    <v-main class="bg-color">
      <v-container>
        <slot />
      </v-container>
    </v-main>
    <v-footer app>
      <span>&copy; {{ new Date().getFullYear() }}</span>
    </v-footer>
  </v-app>
</template>

<script lang="ts" setup>
import { I18n } from '~/types/locales'
import { HeaderMenuItem } from '~/types/props'

const router = useRouter()
const { $i18n } = useNuxtApp()

const t = (str: keyof I18n['layout']['header']) => {
  return $i18n.t(`layout.header.${str}`)
}

const handleCartClick = (): void => {
  console.log('NOT IMPLEMENTED')
}

const localeRef = computed<string>(() => {
  return $i18n.locale === $i18n.defaultLocale ? '' : $i18n.locale
})
const headerMenuList = computed<HeaderMenuItem[]>(() => [
  {
    name: t('signUp'),
    onClick: () => {
      router.push(`${localeRef.value}/signup`)
    },
  },
  {
    name: t('signIn'),
    onClick: () => {
      router.push(`${localeRef.value}/signin`)
    },
  },
  {
    name: t('changeLocaleText'),
    onClick: () => {
      const targetLocale = $i18n.localeCodes.find((code: string) => code !== $i18n.locale)
      targetLocale && $i18n.setLocale(targetLocale)
    },
  },
])
</script>

<style scoped>
.bg-color {
  background-color: #faf2e2;
}

.header-link {
  text-decoration: none;
  color: #1b1b22;
}
</style>
