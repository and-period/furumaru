<script lang="ts" setup>
import { mdiEye, mdiEyeOff } from '@mdi/js'
import useVuelidate, { ValidationArgs } from '@vuelidate/core'

import { AlertType } from '~/lib/hooks'
import {
  required,
  minLength,
  maxLength,
  sameAs,
  getErrorMessage
} from '~/lib/validations'
import { UpdateAuthPasswordRequest } from '~/types/api'

const props = defineProps({
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
  (e: 'submit'): void
}>()

const rules = computed<ValidationArgs>(() => ({
  oldPassword: {
    required
  },
  newPassword: {
    required,
    minLength: minLength(8),
    maxLength: maxLength(32)
  },
  passwordConfirmation: {
    required,
    sameAs: sameAs(props.formData.newPassword)
  }
}))

const validate = useVuelidate(rules, props.formData)

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
  const result = await validate.value.$validate()
  if (!result) {
    return
  }

  emit('submit')
}
</script>

<template>
  <v-alert v-show="props.isAlert" :type="props.alertType" v-text="props.alertText" />
  <v-form @submit.prevent="onSubmit">
    <v-card>
      <v-card-title>パスワード変更</v-card-title>
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
        <v-btn type="submit" block color="primary" variant="outlined">
          変更
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-form>
</template>
