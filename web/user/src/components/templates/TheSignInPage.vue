<script lang="ts" setup>
import type { SignInRequest } from '~/types/api'
import type { LinkItem } from '~/types/props'

interface Props {
  pageName: string
  modelValue: SignInRequest
  buttonText: string
  hasError: boolean
  errorMessage: string
  usernameLabel: string
  usernamePlaceholder: string
  usernameErrorMessage: string
  passwordLabel: string
  passwordPlaceholder: string
  passwordErrorMessage: string
  dontHaveAccountText: string
  googleButtonText: string
  facebookButtonText: string
  lineButtonText: string
  forgetPasswordLink: LinkItem
  signUpLink: LinkItem
}

interface Emits {
  (e: 'update:modelValue', val: SignInRequest): void
  (e: 'submit'): void
  (e: 'click:googleSingInButton'): void
  (e: 'click:facebookSingInButton'): void
  (e: 'click:lineSingInButton'): void
}

const props = defineProps<Props>()
const emits = defineEmits<Emits>()

const formData = computed({
  get: () => props.modelValue,
  set: (val: SignInRequest) => emits('update:modelValue', val),
})

const handleSubmit = () => {
  emits('submit')
}

const handleClickGoogleSingInButton = () => {
  emits('click:googleSingInButton')
}
const handleClickFacebookSingInButton = () => {
  emits('click:facebookSingInButton')
}
const handleClickLineSingInButton = () => {
  emits('click:lineSingInButton')
}
</script>

<template>
  <div class="mx-auto block sm:min-w-[560px]">
    <the-marche-logo class="mb-10" />
    <the-card>
      <the-card-title>{{ pageName }}</the-card-title>
      <the-card-content>
        <the-alert v-show="hasError" class="mb-2">
          {{ errorMessage }}
        </the-alert>

        <the-stack>
          <the-sign-in-form
            v-model="formData"
            :button-text="buttonText"
            :username-label="usernameLabel"
            :username-placeholder="usernamePlaceholder"
            :username-error-message="usernameErrorMessage"
            :password-label="passwordLabel"
            :password-placeholder="passwordPlaceholder"
            :password-error-message="passwordErrorMessage"
            @submit="handleSubmit"
          />

          <p class="my-6 underline">
            <nuxt-link :to="forgetPasswordLink.href">
              {{ forgetPasswordLink.text }}
            </nuxt-link>
          </p>

          <the-google-auth-button
            :button-text="googleButtonText"
            @click="handleClickGoogleSingInButton"
          />
          <the-facebook-auth-button
            :button-text="facebookButtonText"
            @click="handleClickFacebookSingInButton"
          />
          <the-line-auth-button
            :button-text="lineButtonText"
            @click="handleClickLineSingInButton"
          />

          <div class="my-6">
            <p class="mb-2">{{ dontHaveAccountText }}<br /></p>
            <nuxt-link :to="signUpLink.href" class="underline">
              {{ signUpLink.text }}
            </nuxt-link>
          </div>
        </the-stack>
      </the-card-content>
    </the-card>
  </div>
</template>
