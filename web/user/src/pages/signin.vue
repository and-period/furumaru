<template>
  <div class="text-center px-4">
    <the-marche-logo class="mb-6" />
    <v-row justify="center">
      <v-col cols="12" sm="12" md="6">
        <v-form>
          <v-text-field
            :label="t('email')"
            :placeholder="t('email')"
            outlined
            required
            type="email"
            dense
          />
          <v-text-field
            :label="t('password')"
            :placeholder="t('password')"
            outlined
            required
            dense
          />
          <div>
            <nuxt-link to="/">{{ t('forgetPasswordLink') }}</nuxt-link>
          </div>
          <the-submit-button class="mt-6" :is-mobile="isMobile">{{
            t('signIn')
          }}</the-submit-button>
        </v-form>
      </v-col>
    </v-row>
    <v-row justify="center">
      <v-col cols="12" sm="12" md="6">
        <the-stack :space="32" class="my-10 px-4">
          <v-btn
            outlined
            block
            class="base rounded-xl google-btn"
            :height="isMobile ? 40 : 64"
            :x-large="!isMobile"
            @click="handleClickGoogleSingInButton"
          >
            <img
              src="/g-normal.png"
              height="18"
              width="18"
              alt="google icon"
              class="mr-6"
            />
            Log In with Google
          </v-btn>
          <v-btn
            block
            outlined
            class="facebook white--text rounded-xl"
            :height="isMobile ? 40 : 64"
            :x-large="!isMobile"
            @click="handleClickFacebookSingInButton"
          >
            <v-icon left class="mr-6">fa-brands fa-facebook</v-icon>
            Log In with Facebook
          </v-btn>
          <v-btn
            block
            outlined
            class="line white--text rounded-xl"
            :height="isMobile ? 40 : 64"
            :x-large="!isMobile"
            @click="handleClickLineSingInButton"
          >
            <v-icon left class="mr-6">fa-brands fa-line</v-icon>
            Log In with LINE
          </v-btn>
        </the-stack>
      </v-col>
    </v-row>
    <p>
      {{ t('dontHaveAccount') }}
      <nuxt-link :to="localePath('/signup')">{{ t('signUpLink') }}</nuxt-link>
    </p>
  </div>
</template>

<script lang="ts">
import { defineComponent } from '@nuxtjs/composition-api'

import { useI18n, useIsMobile } from '~/lib/hooks'
import { I18n } from '~/types/locales'

export default defineComponent({
  layout: 'auth',
  setup() {
    const { isMobile } = useIsMobile()
    const { i18n } = useI18n()

    const t = (str: keyof I18n['auth']['signIn']) => {
      return i18n.t(`auth.signIn.${str}`)
    }

    const handleClickGoogleSingInButton = () => {
      console.log('NOT IMPLEMENTED')
    }
    const handleClickFacebookSingInButton = () => {
      console.log('NOT IMPLEMENTED')
    }
    const handleClickLineSingInButton = () => {
      console.log('NOT IMPLEMENTED')
    }

    return {
      isMobile,
      t,
      handleClickGoogleSingInButton,
      handleClickFacebookSingInButton,
      handleClickLineSingInButton,
    }
  },
})
</script>

<style lang="scss" scoped>
.google-btn {
  border-color: rgba(0, 0, 0, 0.3) !important;
  color: rgba(0, 0, 0, 0.54);
}
</style>
