<template>
  <div>
    <atoms-the-marche-logo class="text-center mb-6" />
    <v-card variant="outlined" class="mx-auto" :max-width="isMobile ? 360 : 440">
      <v-card-text class="pa-md-12 pa-sm-4">
        <form @submit.prevent="handleSubmit">
          <v-text-field type="tel" :label="t('tel')" variant="outlined" dense required />
          <v-text-field type="email" :label="t('email')" variant="outlined" dense required />
          <v-text-field type="password" :label="t('password')" variant="outlined" dense required />
          <v-text-field type="password" :label="t('passwordConfirm')" variant="outlined" dense required />
          <molecules-the-submit-button :is-mobile="isMobile">
            {{ t('signUp') }}
          </molecules-the-submit-button>
        </form>
      </v-card-text>
    </v-card>
    <div class="text-center mt-10">
      <nuxt-link :to="localePath('/signin')">{{ t('alreadyHas') }}</nuxt-link>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { useIsMobile } from '~/lib/hooks'
import { I18n } from '~/types/locales'

definePageMeta({
  layout: 'auth',
})

const { isMobile } = useIsMobile()
const { $i18n } = useNuxtApp()
const router = useRouter()

const t = (str: keyof I18n['auth']['signUp']) => {
  return $i18n.t(`auth.signUp.${str}`)
}

const handleSubmit = () => {
  router.push('/verify')
}
</script>
