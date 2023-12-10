<script lang="ts" setup>
import useVuelidate from '@vuelidate/core'
import { helpers } from '@vuelidate/validators'
import type { VerifyAuthEmailRequest } from '~/types/api'
import type { AlertType } from '~/lib/hooks'
import { required, getErrorMessage, minLength, maxLength } from '~/lib/validations'

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
  email: {
    type: String,
    default: ''
  },
  formData: {
    type: Object as PropType<VerifyAuthEmailRequest>,
    default: (): VerifyAuthEmailRequest => ({
      verifyCode: ''
    })
  }
})

const emit = defineEmits<{
  (e: 'click:resend-email'): void
  (e: 'update:from-data', formData: VerifyAuthEmailRequest): void
  (e: 'submit'): void
}>()

const rules = computed(() => ({
  verifyCode: {
    required,
    minLength: helpers.withMessage('検証コードは6文字で入力してください。', minLength(6)),
    maxLength: helpers.withMessage('検証コードは6文字で入力してください。', maxLength(6))
  }
}))
const formDataValue = computed({
  get: () => props.formData,
  set: (formData: VerifyAuthEmailRequest) => emit('update:from-data', formData)
})

const validate = useVuelidate(rules, formDataValue)

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
  <v-alert v-show="props.isAlert" :type="props.alertType" v-text="props.alertText" />

  <v-card elevation="0">
    <v-card-title>二要素認証</v-card-title>

    <v-card-text>
      <p class="text-center">
        認証コードが{{ props.email }}に送信されました
      </p>
      <div class="ma-auto" style="max-width: 300px">
        <!-- vuetifyが対応し次第、改修する -->
        <v-text-field v-model="validate.verifyCode.$model" :error-messages="getErrorMessage(validate.verifyCode.$errors)" />
      </div>
      <div class="text-center">
        <a class="orange--text text--darken-4" @click="onClickResendEmail">
          認証コードを再送する
        </a>
      </div>
    </v-card-text>
    <v-card-actions>
      <v-btn :loading="loading" block variant="outlined" color="primary" @click="onSubmit">
        認証
      </v-btn>
    </v-card-actions>
  </v-card>
</template>
