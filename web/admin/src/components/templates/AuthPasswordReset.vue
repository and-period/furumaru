<script lang="ts" setup>
import { mdiEye, mdiEyeOff } from '@mdi/js'
import useVuelidate, { type ValidationArgs } from '@vuelidate/core'

import type { AlertType } from '~/lib/hooks'
import { getErrorMessage } from '~/lib/validations'
import type { ResetAuthPasswordRequest } from '~/types/api'
import { ResetAuthPasswordValidationRules } from '~/types/validations'

const props = defineProps({
  loading: {
    type: Boolean,
    default: false
  },
  isAlert: {
    type: Boolean,
    default: false
  },
  alertType: {
    type: String as PropType<AlertType>,
    default: undefined
  },
  alertText: {
    type: String,
    default: ''
  },
  formData: {
    type: Object as PropType<ResetAuthPasswordRequest>,
    default: (): ResetAuthPasswordRequest => ({
      email: '',
      verifyCode: '',
      password: '',
      passwordConfirmation: ''
    })
  }
})

const emit = defineEmits<{
  (e: 'update:form-data', formData: ResetAuthPasswordRequest): void
  (e: 'click:cancel'): void
  (e: 'submit'): void
}>()

const showPassword = ref<boolean>(false)
const showPasswordConfirmation = ref<boolean>(false)

const rules = computed<ValidationArgs>(() => ResetAuthPasswordValidationRules(props.formData.password))
const formDataValue = computed({
  get: (): ResetAuthPasswordRequest => props.formData,
  set: (v: ResetAuthPasswordRequest): void => emit('update:form-data', v)
})

const validate = useVuelidate(rules, formDataValue)

const onChangePasswordFieldType = (): void => {
  showPassword.value = !showPassword.value
}

const onChangePasswordConfirmationFieldType = (): void => {
  showPasswordConfirmation.value = !showPasswordConfirmation.value
}

const onClickCancel = (): void => {
  emit('click:cancel')
}

const onSubmit = async (): Promise<void> => {
  const valid = await validate.value.$validate()
  if (!valid) {
    return
  }

  emit('submit')
}
</script>

<template>
  <v-alert v-show="props.isAlert" :type="props.alertType" :text="props.alertText" />

  <atoms-app-logo-with-title :width="450" class="mx-auto py-8" />

  <v-card>
    <v-card-title>パスワードリセット</v-card-title>

    <v-card-subtitle>
      {{ formData.email }}へコードを送信しました。メールを確認の上、コードと新しいパスワードを入力してください。
    </v-card-subtitle>

    <v-form @submit.prevent="onSubmit">
      <v-card-text>
        <v-otp-input
          v-model="validate.verifyCode.$model"
          focus-all
          label="認証コード"
          variant="solo-filled"
          :length="6"
          :error-messages="getErrorMessage(validate.verifyCode.$errors)"
        />
        <v-text-field
          v-model="validate.password.$model"
          label="新しいパスワード"
          :type="showPassword ? 'text' : 'password'"
          :append-icon="showPassword ? mdiEye : mdiEyeOff"
          :error-messages="getErrorMessage(validate.password.$errors)"
          @click:append="onChangePasswordFieldType"
        />
        <v-text-field
          v-model="validate.passwordConfirmation.$model"
          label="新しいパスワード（確認用）"
          :type="showPasswordConfirmation ? 'text' : 'password'"
          :append-icon="showPasswordConfirmation ? mdiEye : mdiEyeOff"
          :error-messages="getErrorMessage(validate.passwordConfirmation.$errors)"
          @click:append="onChangePasswordConfirmationFieldType"
        />
      </v-card-text>
      <v-card-actions>
        <v-spacer />
        <v-btn color="error" variant="outlined" @click="onClickCancel">
          サインイン画面にもどる
        </v-btn>
        <v-btn :loading="loading" type="submit" color="primary" variant="outlined">
          パスワードを更新する
        </v-btn>
      </v-card-actions>
    </v-form>
  </v-card>
</template>
