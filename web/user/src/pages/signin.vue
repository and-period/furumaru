<template>
  <div class="block m-auto sm:min-w-[560px]">
    <the-marche-logo class="mb-10" />
    <the-card>
      <the-card-title>ログイン</the-card-title>
      <the-card-content class="sm:px-16 sm:px-6 text-center">
        <the-alert class="mb-2">メールアドレスかパスワードが間違っています。</the-alert>
        <the-stack>
          <the-sign-in-form
            v-model="formData"
            :button-text="t('signIn')"
            :username-label="t('email')"
            :username-placeholder="t('email')"
            username-error-message=""
            :password-label="t('password')"
            :password-placeholder="t('password')"
            password-error-message=""
          />

          <p class="underline my-3">
            <nuxt-link to="/">{{ t('forgetPasswordLink') }}</nuxt-link>
          </p>

          <the-google-auth-button @click="handleClickGoogleSingInButton" />
          <the-facebook-auth-button @click="handleClickFacebookSingInButton" />
          <the-line-auth-button @click="handleClickLineSingInButton" />

          <p class="my-2">
            {{ t('dontHaveAccount') }}<br />
            <nuxt-link :to="localePath('/signup')" class="underline">{{ t('signUpLink') }}</nuxt-link>
          </p>
        </the-stack>
      </the-card-content>
    </the-card>
  </div>
</template>

<script lang="ts" setup>
import { SignInRequest } from '~/types/api'
import { I18n } from '~/types/locales'

definePageMeta({
  layout: 'auth',
})

const { $i18n } = useNuxtApp()

const t = (str: keyof I18n['auth']['signIn']): string => {
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

const formData = reactive<SignInRequest>({
  username: '',
  password: '',
})
</script>
