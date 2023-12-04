<script lang="ts" setup>
import type { CreateAuthRequest } from '~/types/api'
import type { LinkItem } from '~/types/props'

interface Props {
  pageName: string
  errorMessage: string
  buttonText: string
  modelValue: CreateAuthRequest
  telLabel: string
  telPlaceholder: string
  telErrorMessage: string
  emailLabel: string
  emailPlaceholder: string
  emailErrorMessage: string
  passwordLabel: string
  passwordPlaceholder: string
  passwordErrorMessage: string
  passwordConfirmLabel: string
  passwordConfirmPlaceholder: string
  passwordConfirmErrorMessage: string
  alreadyHasLink: LinkItem
}

interface Emits {
  (e: 'submit'): void
  (e: 'update:modelValue', value: CreateAuthRequest): void
}

const props = defineProps<Props>()
const emits = defineEmits<Emits>()

const formData = computed({
  get: () => props.modelValue,
  set: (val: CreateAuthRequest) => emits('update:modelValue', val),
})

const handleSubmit = () => {
  emits('submit')
}
</script>

<template>
  <div class="mx-auto block sm:min-w-[560px]">
    <the-marche-logo class="mb-10" />
    <the-card>
      <the-card-title>
        {{ pageName }}
      </the-card-title>
      <the-card-content>
        <the-alert v-show="errorMessage" class="mb-2">
          {{ errorMessage }}
        </the-alert>

        <the-stack>
          <the-sign-up-form
            v-model="formData"
            :button-text="buttonText"
            :tel-label="telLabel"
            :tel-placeholder="telPlaceholder"
            :tel-error-message="telErrorMessage"
            :email-label="emailLabel"
            :email-placeholder="emailPlaceholder"
            :email-error-message="emailErrorMessage"
            :password-label="passwordLabel"
            :password-placeholder="passwordPlaceholder"
            :password-error-message="passwordErrorMessage"
            :password-confirm-label="passwordConfirmLabel"
            :password-confirm-placeholder="passwordConfirmPlaceholder"
            :password-confirm-error-message="passwordConfirmErrorMessage"
            @submit="handleSubmit"
          />

          <p class="my-6 underline">
            <nuxt-link :to="alreadyHasLink.href">
              {{ alreadyHasLink.text }}
            </nuxt-link>
          </p>
        </the-stack>
      </the-card-content>
    </the-card>
  </div>
</template>
