<script lang="ts" setup>
import { VerifyAuthEmailRequest } from '~/types/api'
import { AlertType } from '~/lib/hooks'

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
  email: {
    type: String,
    default: ''
  },
  formData: {
    type: Object as PropType<VerifyAuthEmailRequest>,
    default: () => ({
      verifyCode: ''
    })
  }
})

const emit = defineEmits<{
  (e: 'click:resend-email'): void
  (e: 'submit'): void
}>()

const onClickResendEmail = (): void => {
  emit('click:resend-email')
}

const onSubmit = (): void => {
  emit('submit')
}
</script>

<template>
  <v-alert v-model="props.isAlert" :type="props.alertType" :text="props.alertText" />
  <v-card elevation="0">
    <v-card-title>二要素認証</v-card-title>
    <v-card-text>
      <p class="text-center">
        認証コードが{{ props.email }}に送信されました
      </p>
      <div class="ma-auto" style="max-width: 300px">
        <!-- vuetifyが対応し次第、改修する -->
        <!-- <v-otp-input v-model="props.formData.verifyCode" type="number" length="6" /> -->
        <v-text-field v-model="props.formData.verifyCode" />
      </div>
      <div class="text-center">
        <a
          class="orange--text text--darken-4"
          @click="onClickResendEmail"
        >認証コードを再送する</a>
      </div>
    </v-card-text>
    <v-card-actions>
      <v-btn block variant="outlined" color="primary" @click="onSubmit">
        認証
      </v-btn>
    </v-card-actions>
  </v-card>
</template>
