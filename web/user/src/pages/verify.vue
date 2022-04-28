<template>
  <div>
    <the-marche-logo class="text-center mb-6" />
    <v-card
      elevation="0"
      class="mx-auto bg-transparent"
      :max-width="isMobile ? 360 : 440"
    >
      <v-card-text class="pa-md-12 pa-sm-4">
        <form @submit.prevent="handleSubmit">
          <the-stack :space="24" class="text-center">
            <p class="text-h6">{{ t('message') }}</p>
            <the-pincode-input value="" :length="6" class="mt-6" />
            <the-submit-button :is-mobile="isMobile">
              {{ t('btnText') }}
            </the-submit-button>
          </the-stack>
        </form>
      </v-card-text>
    </v-card>
  </div>
</template>

<script lang="ts">
import { defineComponent } from '@nuxtjs/composition-api'
import { useIsMobile, useI18n } from '~/lib/hooks'
import { I18n } from '~/types/locales'

export default defineComponent({
  layout: 'auth',
  setup() {
    const { isMobile } = useIsMobile()
    const { i18n } = useI18n()

    const t = (str: keyof I18n['auth']['verify']) => {
      return i18n.t(`auth.verify.${str}`)
    }

    const handleSubmit = () => {
      console.log('submit not implemented')
    }

    return {
      isMobile,
      t,
      handleSubmit,
    }
  },
})
</script>
