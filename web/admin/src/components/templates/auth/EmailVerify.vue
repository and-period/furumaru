<script lang="ts" setup>
import { mdiEmailCheck, mdiArrowLeft, mdiTimerSand, mdiRefresh } from '@mdi/js'
import useVuelidate from '@vuelidate/core'
import type { VerifyAuthEmailRequest } from '~/types/api/v1'
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
  <v-container class="pa-6">
    <atoms-app-alert
      :show="props.isAlert"
      :type="props.alertType"
      :text="props.alertText"
      class="mb-6"
    />

    <div class="mb-6">
      <v-btn
        variant="text"
        :icon="mdiArrowLeft"
        class="mb-4"
        @click="$router.back()"
      >
        戻る
      </v-btn>
      <h1 class="text-h4 font-weight-bold mb-2">
        <v-icon
          :icon="mdiEmailCheck"
          size="32"
          class="mr-3 text-success"
        />
        メールアドレス認証
      </h1>
      <p class="text-body-1 text-grey-darken-1">
        送信された認証コードを入力して、メールアドレスの変更を完了してください。
      </p>
    </div>

    <v-card
      elevation="2"
      class="verify-form-card"
    >
      <v-card-text class="pa-8 text-center">
        <div class="mb-6">
          <v-icon
            :icon="mdiEmailCheck"
            size="64"
            class="text-success mb-4"
          />
          <h3 class="text-h6 font-weight-medium mb-2">
            認証コードを送信しました
          </h3>
          <p class="text-body-2 text-grey-darken-1">
            <strong>{{ props.email }}</strong> に6桁の認証コードを送信しました
          </p>
        </div>

        <div class="mb-6">
          <h4 class="text-subtitle-1 font-weight-medium mb-4 text-grey-darken-1">
            認証コードを入力
          </h4>
          <div class="otp-container">
            <v-otp-input
              v-model="validate.verifyCode.$model"
              focus-all
              variant="outlined"
              :loading="loading"
              :error-messages="getErrorMessage(validate.verifyCode.$errors)"
              class="otp-input"
            />
          </div>
        </div>

        <v-alert
          type="info"
          variant="outlined"
          class="mb-6 text-left"
        >
          <div class="text-body-2">
            <v-icon
              :icon="mdiTimerSand"
              class="mr-2"
              size="16"
            />
            <strong>認証コードについて:</strong><br>
            • 認証コードの有効期限は15分です<br>
            • 6桁をすべて入力すると自動で認証が実行されます<br>
            • メールが届かない場合は迷惑メールフォルダもご確認ください
          </div>
        </v-alert>

        <div class="text-center">
          <v-btn
            variant="outlined"
            color="primary"
            :loading="loading"
            @click="onClickResendEmail"
          >
            <v-icon
              :icon="mdiRefresh"
              start
            />
            認証コードを再送
          </v-btn>
        </div>
      </v-card-text>
    </v-card>
  </v-container>
</template>

<!-- stylelint-disable selector-class-pattern -->
<style scoped>
.verify-form-card {
  border-radius: 12px;
  max-width: 600px;
  margin: 0 auto;
}

.otp-container {
  display: flex;
  justify-content: center;
  margin: 0 auto;
  max-width: 400px;
}

/* stylelint-disable-next-line selector-class-pattern */
.otp-input :deep(.v-otp-input__content) {
  gap: 12px;
}

/* stylelint-disable-next-line selector-class-pattern */
.otp-input :deep(.v-field__input) {
  font-family: Monaco, Menlo, Consolas, monospace;
  font-size: 24px;
  font-weight: 700;
  text-align: center;
}

/* stylelint-disable-next-line selector-class-pattern */
.otp-input :deep(.v-field--variant-outlined) {
  border-radius: 8px;
  border: 2px solid rgb(33 150 243 / 30%);
}

/* stylelint-disable-next-line selector-class-pattern */
.otp-input :deep(.v-field--focused) {
  border-color: rgb(33 150 243);
  box-shadow: 0 0 0 2px rgb(33 150 243 / 20%);
}

@media (width <= 600px) {
  .verify-form-card {
    margin: 0;
  }

  .form-section {
    padding: 16px;
  }

  .action-buttons {
    padding: 16px !important;
    flex-direction: column;
    gap: 12px;
  }

  .action-button-group {
    display: flex;
    flex-direction: column;
    gap: 12px;
    width: 100%;
  }

  .cancel-btn,
  .submit-btn {
    width: 100%;
    font-size: 16px;
  }

  .submit-btn {
    order: -1;
  }
}

@media (width >= 601px) {
  .action-button-group {
    display: flex;
    gap: 16px;
    align-items: center;
  }

  .submit-btn {
    padding-left: 32px;
    padding-right: 32px;
  }
}
</style>
