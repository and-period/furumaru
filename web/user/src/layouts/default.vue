<template>
  <v-app>
    <the-app-header
      :cart-item-count="0"
      :cart-empty-message="t('cartEmptyMessage')"
      :cart-not-empty-message="t('cartNotEmptyMessage')"
      @click:cart="handleCartClick"
    >
      <nuxt-link to="/" class="mr-4 header-link">
        {{ t('becomeShopOwner') }}
      </nuxt-link>
    </the-app-header>
    <v-main class="bg-color">
      <v-container>
        <Nuxt />
      </v-container>
    </v-main>
    <v-footer app>
      <span>&copy; {{ new Date().getFullYear() }}</span>
    </v-footer>
  </v-app>
</template>

<script lang="ts">
import { defineComponent } from '@nuxtjs/composition-api'

import { useI18n } from '~/lib/hooks'
import { I18n } from '~/types/locales'

export default defineComponent({
  setup() {
    const { i18n } = useI18n()

    const t = (str: keyof I18n['default']['header']) => {
      return i18n.t(`default.header.${str}`)
    }

    console.log('aaaa', t('becomeShopOwner'), i18n.t('auth.signUp.tel'))

    const handleCartClick = (): void => {
      console.log('NOT IMPLEMENTED')
    }

    return {
      handleCartClick,
      t,
    }
  },
})
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
