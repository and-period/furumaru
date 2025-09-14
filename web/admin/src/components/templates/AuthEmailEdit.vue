<script lang="ts" setup>
import { mdiEmail, mdiArrowLeft } from '@mdi/js'
import useVuelidate from '@vuelidate/core'
import type { UpdateAuthEmailRequest } from '~/types/api/v1'
import type { AlertType } from '~/lib/hooks'
import { getErrorMessage } from '~/lib/validations'
import { UpdateAuthEmailValidationRules } from '~/types/validations'

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
  formData: {
    type: Object as PropType<UpdateAuthEmailRequest>,
    default: (): UpdateAuthEmailRequest => ({
      email: '',
    }),
  },
})

const emits = defineEmits<{
  (e: 'update:form-data', formData: UpdateAuthEmailRequest): void
  (e: 'submit'): void
}>()

const formDataValue = computed({
  get: () => props.formData,
  set: (formData: UpdateAuthEmailRequest) => emits('update:form-data', formData),
})

const validate = useVuelidate(UpdateAuthEmailValidationRules, formDataValue)

const onSubmit = async (): Promise<void> => {
  const valid = await validate.value.$validate()
  if (!valid) {
    return
  }

  emits('submit')
}
</script>

<template>
  <v-container class="pa-6">
    <v-alert
      v-show="props.isAlert"
      :type="props.alertType"
      class="mb-6"
      v-text="props.alertText"
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
          :icon="mdiEmail"
          size="32"
          class="mr-3 text-primary"
        />
        メールアドレス変更
      </h1>
      <p class="text-body-1 text-grey-darken-1">
        ログイン用メールアドレスを変更します。変更後は認証が必要です。
      </p>
    </div>

    <v-card
      elevation="2"
      class="email-form-card"
    >
      <v-form @submit.prevent="onSubmit">
        <v-card-text class="pa-8">
          <div class="form-section">
            <h3 class="text-subtitle-1 font-weight-medium mb-4 text-grey-darken-1">
              新しいメールアドレス
            </h3>
            <v-text-field
              v-model="validate.email.$model"
              :error-messages="getErrorMessage(validate.email.$errors)"
              label="新しいメールアドレス"
              type="email"
              variant="outlined"
              prepend-inner-icon="mdi-email"
              placeholder="example@furumaru.com"
            />

            <v-alert
              type="info"
              variant="outlined"
              class="mt-4"
            >
              <div class="text-body-2">
                <strong>注意事項:</strong><br>
                • 変更後は新しいメールアドレスでのログインになります<br>
                • 認証コードが送信されます<br>
                • 認証完了までメールアドレスは変更されません
              </div>
            </v-alert>
          </div>
        </v-card-text>
        <v-card-actions class="pa-8 pt-0 action-buttons">
          <v-spacer class="d-none d-sm-flex" />
          <div class="action-button-group">
            <v-btn
              variant="text"
              class="cancel-btn"
              @click="$router.back()"
            >
              キャンセル
            </v-btn>
            <v-btn
              :loading="loading"
              type="submit"
              color="primary"
              variant="elevated"
              size="large"
              class="submit-btn"
            >
              認証コードを送信
            </v-btn>
          </div>
        </v-card-actions>
      </v-form>
    </v-card>
  </v-container>
</template>

<style scoped>
.email-form-card {
  border-radius: 12px;
  max-width: 800px;
  margin: 0 auto;
}

.form-section {
  border-radius: 8px;
  padding: 20px;
  background: rgb(248 250 252);
  border-left: 4px solid rgb(33 150 243);
}

.form-section h3 {
  border-bottom: 1px solid rgb(224 224 224);
  padding-bottom: 8px;
}

@media (width <= 600px) {
  .email-form-card {
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
