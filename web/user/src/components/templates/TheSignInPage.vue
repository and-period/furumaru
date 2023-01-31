<script lang="ts" setup>
import { SignInRequest } from '~/types/api'
import { LinkItem } from '~/types/props'

interface Props {
  pageName: string
  modelValue: SignInRequest
  buttonText: string
  hasError: boolean
  usernameLabel: string
  usernamePlaceholder: string
  usernameErrorMessage: string
  passwordLabel: string
  passwordPlaceholder: string
  passwordErrorMessage: string
  dontHaveAccountText: string
  forgetPasswordLink: LinkItem
  signUpLink: LinkItem
}

interface Emits {
  (e: 'update:modelValue', val: SignInRequest): void;
  (e: 'submit'): void;
  (e: 'click:googleSingInButton'): void;
  (e: 'click:facebookSingInButton'): void;
  (e: 'click:lineSingInButton'): void;
}

const props = defineProps<Props>()
const emits = defineEmits<Emits>()

const formData = computed({
  get: () => props.modelValue,
  set: (val: SignInRequest) => emits('update:modelValue', val)
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
  <div class="block m-auto sm:min-w-[560px]">
    <the-marche-logo class="mb-10" />
    <the-card>
      <the-card-title>{{ pageName }}</the-card-title>
      <the-card-content class="sm:px-16 sm:px-6 text-center">
        <the-alert v-show="hasError" class="mb-2">
          メールアドレスかパスワードが間違っています。
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

          <p class="underline my-3">
            <nuxt-link :to="forgetPasswordLink.href">
              {{ forgetPasswordLink.text }}
            </nuxt-link>
          </p>

          <the-google-auth-button @click="handleClickGoogleSingInButton" />
          <the-facebook-auth-button @click="handleClickFacebookSingInButton" />
          <the-line-auth-button @click="handleClickLineSingInButton" />

          <p class="my-2">
            {{ dontHaveAccountText }}<br>
            <nuxt-link :to="signUpLink.href" class="underline">
              {{ signUpLink.text }}
            </nuxt-link>
          </p>
        </the-stack>
      </the-card-content>
    </the-card>
  </div>
</template>
