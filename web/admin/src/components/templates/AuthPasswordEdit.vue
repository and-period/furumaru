<script lang="ts" setup>
import { mdiEye, mdiEyeOff } from '@mdi/js'
import useVuelidate, { type ValidationArgs } from '@vuelidate/core'

import type { AlertType } from '~/lib/hooks'
import { getErrorMessage } from '~/lib/validations'
import type { UpdateAuthPasswordRequest } from '~/types/api'
import { UpdateAuthPasswordValidationRules } from '~/types/validations'

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
    type: Object as PropType<UpdateAuthPasswordRequest>,
    default: () => ({
      oldPassword: '',
      newPassword: '',
      passwordConfirmation: ''
    })
  }
})

const emit = defineEmits<{
  (e: 'update:form-data', formData: UpdateAuthPasswordRequest): void
  (e: 'submit'): void
}>()

const rules = computed<ValidationArgs>(() => UpdateAuthPasswordValidationRules(props.formData.newPassword))
const formDataValue = computed({
  get: (): UpdateAuthPasswordRequest => props.formData,
  set: (formData: UpdateAuthPasswordRequest): void => emit('update:form-data', formData)
})

const validate = useVuelidate(rules, formDataValue)

const showOldPassword = ref<boolean>(false)
const showNewPassword = ref<boolean>(false)
const showPasswordConfirmation = ref<boolean>(false)

const onChangeOldPasswordFieldType = (): void => {
  showOldPassword.value = !showOldPassword.value
}

const onChangeNewPasswordFieldType = (): void => {
  showNewPassword.value = !showNewPassword.value
}

const onChangePasswordConfirmationFieldType = (): void => {
  showPasswordConfirmation.value = !showPasswordConfirmation.value
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
  <v-alert v-show="props.isAlert" :type="props.alertType" v-text="props.alertText" />

  <v-card>
    <v-card-title>パスワード変更</v-card-title>

    <v-form @submit.prevent="onSubmit">
      <v-card-text>
        <v-text-field
          v-model="validate.oldPassword.$model"
          label="現在のパスワード"
          :type="showOldPassword ? 'text' : 'password'"
          :append-icon="showOldPassword ? mdiEye : mdiEyeOff"
          :error-messages="getErrorMessage(validate.oldPassword.$errors)"
          @click:append="onChangeOldPasswordFieldType"
        />
        <v-text-field
          v-model="validate.newPassword.$model"
          label="新しいパスワード"
          :type="showNewPassword ? 'text' : 'password'"
          :append-icon="showNewPassword ? mdiEye : mdiEyeOff"
          :error-messages="getErrorMessage(validate.newPassword.$errors)"
          @click:append="onChangeNewPasswordFieldType"
        />
        <v-text-field
          v-model="validate.passwordConfirmation.$model"
          label="新しいパスワード(確認用)"
          :type="showPasswordConfirmation ? 'text' : 'password'"
          :append-icon="showPasswordConfirmation ? mdiEye : mdiEyeOff"
          :error-messages="
            getErrorMessage(validate.passwordConfirmation.$errors) === ''
              ? ''
              : 'パスワードが一致しません。'
          "
          @click:append="onChangePasswordConfirmationFieldType"
        />
      </v-card-text>

      <v-card-actions>
        <v-btn :loading="loading" block type="submit" color="primary" variant="outlined">
          更新
        </v-btn>
      </v-card-actions>
    </v-form>
  </v-card>
</template>
