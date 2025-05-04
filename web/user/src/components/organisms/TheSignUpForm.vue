<script lang="ts" setup>
import type { CreateAuthUserRequest } from '~/types/api'
import type { I18n } from '~/types/locales'

interface Props {
  modelValue: CreateAuthUserRequest
  buttonText: string
  lastnameKanaErrorMessage: string
  firstnameKanaErrorMessage: string
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
  (e: 'update:modelValue', value: CreateAuthUserRequest): void
}

const i18n = useI18n()
const props = defineProps<Props>()

const emits = defineEmits<Emits>()

const t = (str: keyof I18n['auth']['signUp']) => {
  return i18n.t(`auth.signUp.${str}`)
}

const formData = computed({
  get: () => props.modelValue,
  set: (val: CreateAuthUserRequest) => emits('update:modelValue', val),
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
        :placeholder="t('usernamePlaceholder')"
        type="text"
        required
      />

      <the-text-input
        v-model="formData.accountId"
        label="ユーザーID"
        :with-label="false"
        :placeholder="t('userIdPlaceholder')"
        type="text"
        required
      />

      <div class="grid grid-cols-2 gap-2 gap-x-6">
        <the-text-input
          v-model="formData.lastname"
          label="姓"
          :with-label="false"
          :placeholder="t('lastNamePlaceholder')"
          type="text"
          required
        />
        <the-text-input
          v-model="formData.firstname"
          label="名"
          :with-label="false"
          :placeholder="t('firstNamePlaceholder')"
          type="text"
          required
        />
        <the-text-input
          v-model="formData.lastnameKana"
          label="姓（かな）"
          :with-label="false"
          :placeholder="t('lastNameKanaPlaceholder')"
          :error-message="lastnameKanaErrorMessage"
          type="text"
          required
        />

        <the-text-input
          v-model="formData.firstnameKana"
          label="名（かな）"
          :with-label="false"
          :placeholder="t('firstNameKanaPlaceholder')"
          :error-message="firstnameKanaErrorMessage"
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
