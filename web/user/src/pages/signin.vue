<template>
  <div class="text-center px-4">
    <atoms-the-marche-logo class="mb-6" />
    <v-row justify="center">
      <v-col cols="12" sm="12" md="6">
        <v-form>
          <v-text-field :label="t('email')" :placeholder="t('email')" variant="outlined" required type="email" dense />
          <v-text-field :label="t('password')" :placeholder="t('password')" variant="outlined" required dense />
          <div>
            <nuxt-link to="/">{{ t('forgetPasswordLink') }}</nuxt-link>
          </div>
          <molecules-the-submit-button class="mt-6" :is-mobile="isMobile">{{
            t('signIn')
          }}</molecules-the-submit-button>
        </v-form>
      </v-col>
    </v-row>
    <v-row justify="center">
      <v-col cols="12" sm="12" md="6">
        <atoms-the-stack :space="32" class="my-10 px-4">
          <v-btn
            variant="outlined"
            block
            class="base rounded-xl google-btn"
            :height="isMobile ? 40 : 64"
            :size="!isMobile && 'x-large'"
            @click="handleClickGoogleSingInButton"
          >
            <img src="~/assets/img/g-normal.png" height="18" width="18" alt="google icon" class="mr-6" />
            Log In with Google
          </v-btn>
          <v-btn
            block
            variant="outlined"
            class="facebook text-white rounded-xl"
            :height="isMobile ? 40 : 64"
            :size="!isMobile && 'x-large'"
            @click="handleClickFacebookSingInButton"
          >
            <v-icon start class="mr-6">fa-brands fa-facebook</v-icon>
            Log In with Facebook
          </v-btn>
          <v-btn
            block
            variant="outlined"
            class="line text-white rounded-xl"
            :height="isMobile ? 40 : 64"
            :size="!isMobile && 'x-large'"
            @click="handleClickLineSingInButton"
          >
            <v-icon start class="mr-6">fa-brands fa-line</v-icon>
            Log In with LINE
          </v-btn>
        </atoms-the-stack>
      </v-col>
    </v-row>
    <p>
      {{ t('dontHaveAccount') }}
      <nuxt-link :to="localePath('/signup')">{{ t('signUpLink') }}</nuxt-link>
    </p>
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

const t = (str: keyof I18n['auth']['signIn']) => {
  return $i18n.t(`auth.signIn.${str}`)
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
</script>

<style lang="scss" scoped>
.google-btn {
  border-color: rgba(0, 0, 0, 0.3) !important;
  color: rgba(0, 0, 0, 0.54);
}
</style>
