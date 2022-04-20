<template>
  <div>
    <v-card outlined class="mx-auto" :max-width="isMobile ? 360 : 440">
      <v-card-text class="pa-md-12 pa-sm-4">
        <form @submit.prevent="handleSubmit">
          <v-text-field type="tel" :label="t('tel')" outlined dense required />
          <v-text-field
            type="email"
            :label="t('email')"
            outlined
            dense
            required
          />
          <v-text-field
            type="password"
            :label="t('password')"
            outlined
            dense
            required
          />
          <v-text-field
            type="password"
            :label="t('passwordConfirm')"
            outlined
            dense
            required
          />
          <the-submit-button :is-mobile="isMobile">
            {{ t('signUp') }}
          </the-submit-button>
        </form>
      </v-card-text>
    </v-card>
    <div class="text-center mt-10">
      <nuxt-link to="/">{{ t('alreadyHas') }}</nuxt-link>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, useRouter } from '@nuxtjs/composition-api'
import { I18n } from '~/types/locales'
import { useIsMobile, useI18n } from '~/lib/hooks'

export default defineComponent({
  layout: 'auth',
  setup() {
    const { isMobile } = useIsMobile()
    const translate = useI18n()

    const t = (str: keyof I18n['auth']['signUp']) => {
      return translate.t(`auth.signUp.${str}`)
    }
    const router = useRouter()

    const handleSubmit = () => {
      router.push('/verify')
    }

    return {
      isMobile,
      handleSubmit,
      t,
    }
  },
})
</script>
