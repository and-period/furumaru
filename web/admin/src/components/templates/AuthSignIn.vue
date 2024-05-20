<script lang="ts" setup>
import { mdiEye, mdiEyeOff } from '@mdi/js'

import type { AlertType } from '~/lib/hooks'
import type { SignInRequest } from '~/types/api'

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
    type: Object as PropType<SignInRequest>,
    default: () => ({
      username: '',
      password: ''
    })
  }
})

const emit = defineEmits<{
  (e: 'update:form-data', formData: SignInRequest): void
  (e: 'submit'): void
}>()

const showPassword = ref<boolean>(false)

const formDataValue = computed({
  get: (): SignInRequest => props.formData,
  set: (v: SignInRequest): void => emit('update:form-data', v)
})

const onChangePasswordFieldType = (): void => {
  showPassword.value = !showPassword.value
}

const onSubmit = (): void => {
  emit('submit')
}
</script>

<template>
  <v-alert v-show="props.isAlert" :type="props.alertType" :text="props.alertText" />

  <atoms-app-logo-with-title :width="450" class="mx-auto py-8" />

  <v-card>
    <v-card-title>ログイン</v-card-title>

    <v-form @submit.prevent="onSubmit">
      <v-card-text>
        <v-text-field
          v-model="formDataValue.username"
          required
          label="ユーザーID（メールアドレス)"
          type="email"
        />
        <v-text-field
          v-model="formDataValue.password"
          required
          label="パスワード"
          :type="showPassword ? 'text' : 'password'"
          :append-icon="showPassword ? mdiEye : mdiEyeOff"
          @click:append="onChangePasswordFieldType"
        />
        <nuxt-link to="/recover">
          パスワードを忘れた場合
        </nuxt-link>
      </v-card-text>
      <v-card-actions>
        <v-btn :loading="loading" type="submit" block color="primary" variant="outlined">
          ログイン
        </v-btn>
      </v-card-actions>
    </v-form>
  </v-card>
  <v-container class="text-right">
    <v-row>
      <v-col cols="12" sm="1" md="12">
        <nuxt-link class="text-body-2" to="/privacy">
          プライバシーポリシー
        </nuxt-link>
      </v-col>
      <v-col cols="12" sm="1" md="12">
        <nuxt-link class="text-body-2" to="/legal-notice">
          特商取引法に基づく表記
        </nuxt-link>
      </v-col>
    </v-row>
  </v-container>
</template>
