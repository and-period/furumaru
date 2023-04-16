<script lang="ts" setup>
import { mdiEye, mdiEyeOff } from '@mdi/js'

import { AlertType } from '~/lib/hooks'
import { SignInRequest } from '~/types/api'

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
    type: Object as PropType<SignInRequest>,
    default: () => ({
      username: '',
      password: ''
    })
  }
})

const emit = defineEmits<{
  (e: 'submit'): void
}>()

const showPassword = ref<boolean>(false)

const onChangePasswordFieldType = (): void => {
  showPassword.value = !showPassword.value
}

const onSubmit = (): void => {
  emit('submit')
}
</script>

<template>
  <v-alert v-model="props.isAlert" :type="props.alertType" :text="props.alertText" />
  <div class="py-8">
    <atoms-app-logo-with-title :width="450" class="mx-auto" />
  </div>
  <v-card>
    <v-form @submit.prevent="onSubmit">
      <v-card-title>ログイン</v-card-title>
      <v-card-text>
        <v-text-field
          v-model="props.formData.username"
          required
          label="ユーザーID（メールアドレス)"
          type="email"
        />
        <v-text-field
          v-model="props.formData.password"
          required
          label="パスワード"
          :type="showPassword ? 'text' : 'password'"
          :append-icon="showPassword ? mdiEye : mdiEyeOff"
          @click:append="onChangePasswordFieldType"
        />
      </v-card-text>
      <v-card-actions>
        <v-btn type="submit" block color="primary" variant="outlined">
          ログイン
        </v-btn>
      </v-card-actions>
    </v-form>
  </v-card>
</template>
