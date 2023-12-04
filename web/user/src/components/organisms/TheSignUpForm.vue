<script lang="ts" setup>
import type { CreateAuthRequest } from '~/types/api'

interface Props {
  modelValue: CreateAuthRequest
  buttonText: string
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
  <form @submit.prevent="handleSubmit">
    <the-stack>
      <the-text-input
        v-model="formData.phoneNumber"
        :label="telLabel"
        :placeholder="telPlaceholder"
        :with-label="false"
        :error-message="telErrorMessage"
        type="tel"
        required
      />
      <the-text-input
        v-model="formData.email"
        :label="emailLabel"
        :placeholder="emailPlaceholder"
        :with-label="false"
        :error-message="emailErrorMessage"
        type="email"
        required
      />
      <the-text-input
        v-model="formData.password"
        :label="passwordLabel"
        :placeholder="passwordPlaceholder"
        :with-label="false"
        :error-message="passwordErrorMessage"
        type="password"
        required
      />
      <the-text-input
        v-model="formData.passwordConfirmation"
        :label="passwordConfirmLabel"
        :placeholder="passwordConfirmPlaceholder"
        :with-label="false"
        :error-message="passwordConfirmErrorMessage"
        type="password"
        required
      />
      <the-submit-button class="mt-4">
        {{ buttonText }}
      </the-submit-button>
    </the-stack>
  </form>
</template>
