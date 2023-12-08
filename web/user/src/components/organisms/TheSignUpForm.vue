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
        v-model="formData.username"
        label="ユーザー名"
        :with-label="false"
        placeholder="ユーザー名（ふるマル）"
        type="text"
        required
      />

      <the-text-input
        v-model="formData.accountId"
        label="ユーザーID"
        :with-label="false"
        placeholder="ユーザーID（@furumaruchan）"
        type="text"
        required
      />

      <div class="grid grid-cols-2 gap-2 gap-x-6">
        <the-text-input
          v-model="formData.lastname"
          label="姓"
          :with-label="false"
          placeholder="姓"
          type="text"
          required
        />
        <the-text-input
          v-model="formData.firstname"
          label="名"
          :with-label="false"
          placeholder="名"
          type="text"
          required
        />
        <the-text-input
          v-model="formData.lastnameKana"
          label="姓（かな）"
          :with-label="false"
          placeholder="姓（かな）"
          type="text"
          required
        />

        <the-text-input
          v-model="formData.firstnameKana"
          label="名（かな）"
          :with-label="false"
          placeholder="名（かな）"
          type="text"
          required
        />
      </div>

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
