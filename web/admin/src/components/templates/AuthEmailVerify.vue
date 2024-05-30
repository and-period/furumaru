<script lang="ts" setup>
import useVuelidate from '@vuelidate/core'
import type { VerifyAuthEmailRequest } from '~/types/api'
import type { AlertType } from '~/lib/hooks'
import { getErrorMessage } from '~/lib/validations'
import { VerifyAuthEmailValidationRules } from '~/types/validations'

const props = defineProps({
  loading: {
    type: Boolean,
    default: false,
  },
  isAlert: {
    type: Boolean,
    default: false,
  },
  alertType: {
    type: String as PropType<AlertType>,
    default: undefined,
  },
  alertText: {
    type: String,
    default: '',
  },
  email: {
    type: String,
    default: '',
  },
  formData: {
    type: Object as PropType<VerifyAuthEmailRequest>,
    default: (): VerifyAuthEmailRequest => ({
      verifyCode: '',
    }),
  },
})

const emit = defineEmits<{
  (e: 'click:resend-email'): void
  (e: 'update:from-data', formData: VerifyAuthEmailRequest): void
  (e: 'submit'): void
}>()

const formDataValue = computed({
  get: () => props.formData,
  set: (formData: VerifyAuthEmailRequest) => emit('update:from-data', formData),
})

const validate = useVuelidate(VerifyAuthEmailValidationRules, formDataValue)

const onClickResendEmail = (): void => {
  emit('click:resend-email')
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
  <v-alert
    v-show="props.isAlert"
    :type="props.alertType"
    v-text="props.alertText"
  />

  <v-card elevation="0">
    <v-card-text>
      <p class="text-center">
        認証コードが{{ props.email }}に送信されました
      </p>
      <div
        class="ma-auto"
        style="max-width: 300px"
      >
        <v-otp-input
          v-model="validate.verifyCode.$model"
          focus-all
          variant="solo-filled"
          :loading="loading"
          :error-messages="getErrorMessage(validate.verifyCode.$errors)"
        />
      </div>
      <div class="text-center mt-4">
        <v-btn
          class="orange--text text--darken-4"
          @click="onClickResendEmail"
        >
          認証コードを再送する
        </v-btn>
      </div>
    </v-card-text>
  </v-card>
</template>
