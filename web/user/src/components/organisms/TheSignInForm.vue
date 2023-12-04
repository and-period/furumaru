<script lang="ts" setup>
import type { SignInRequest } from '~/types/api'

interface Props {
  modelValue: SignInRequest
  buttonText: string
  usernameLabel: string
  usernamePlaceholder: string
  usernameErrorMessage: string
  passwordLabel: string
  passwordPlaceholder: string
  passwordErrorMessage: string
}

const props = defineProps<Props>()
const emits = defineEmits<{
  (e: 'update:modelValue', val: SignInRequest): void
  (e: 'submit'): void
}>()

const formData = computed({
  get: () => props.modelValue,
  set: (val: SignInRequest) => emits('update:modelValue', val),
})

const handleSubmit = () => {
  emits('submit')
}
</script>

<template>
  <form @submit.prevent="handleSubmit">
    <the-stack>
      <the-text-input
        v-model="formData.username"
        :label="usernameLabel"
        :placeholder="usernamePlaceholder"
        :with-label="false"
        :error-message="usernameErrorMessage"
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
      <the-submit-button class="mt-4">
        {{ buttonText }}
      </the-submit-button>
    </the-stack>
  </form>
</template>
