<script lang="ts" setup>
import useVuelidate from '@vuelidate/core'
import { helpers } from '@vuelidate/validators'
import { VerifyAuthEmailRequest } from '~/types/api'
import { AlertType } from '~/lib/hooks'
import { required, getErrorMessage, minLength, maxLength } from '~/lib/validations'

interface Props {
  isAlert: boolean,
  alertType: AlertType
  alertText: string
  email: string
  formData: VerifyAuthEmailRequest
}

const props = defineProps<Props>()

interface Emits {
  (e: 'update:fromData', val: VerifyAuthEmailRequest): void
  (e: 'update:isAlert', val: boolean): void,
  (e: 'click:resend-email'): void
  (e: 'submit'): void
}

const emits = defineEmits<Emits>()

const isAlertValue = computed({
  get: () => props.isAlert,
  set: (val: boolean) => emits('update:isAlert', val)
})

const formDataValue = computed({
  get: () => props.formData,
  set: (val: VerifyAuthEmailRequest) => emits('update:fromData', val)
})

const rules = computed(() => {
  return {
    verifyCode: {
      required,
      minLength: helpers.withMessage('検証コードは6文字で入力してください。', minLength(6)),
      maxLength: helpers.withMessage('検証コードは6文字で入力してください。', maxLength(6))
    }
  }
})

const v$ = useVuelidate(rules, formDataValue)

const onClickResendEmail = (): void => {
  emits('click:resend-email')
}

const onSubmit = async () => {
  const result = await v$.value.$validate()
  if (!result) {
    return
  }
  emits('submit')
}
</script>

<template>
  <v-alert v-model="isAlertValue" class="mb-2" :type="alertType" :text="alertText" />

  <v-card elevation="0">
    <v-card-title>二要素認証</v-card-title>

    <v-card-text>
      <p class="text-center">
        認証コードが{{ props.email }}に送信されました
      </p>
      <div class="ma-auto" style="max-width: 300px">
        <!-- vuetifyが対応し次第、改修する -->
        <v-text-field v-model="v$.verifyCode.$model" :error-messages="getErrorMessage(v$.verifyCode.$errors)" />
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
