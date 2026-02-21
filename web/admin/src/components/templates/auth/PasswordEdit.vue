<script lang="ts" setup>
import { mdiEye, mdiEyeOff, mdiLock, mdiArrowLeft } from '@mdi/js'
import useVuelidate from '@vuelidate/core'
import type { ValidationArgs } from '@vuelidate/core'

import type { AlertType } from '~/lib/hooks'
import { getErrorMessage } from '~/lib/validations'
import type { UpdateAuthPasswordRequest } from '~/types/api/v1'
import { UpdateAuthPasswordValidationRules } from '~/types/validations'

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
    type: Object as PropType<UpdateAuthPasswordRequest>,
    default: () => ({
      oldPassword: '',
      newPassword: '',
      passwordConfirmation: '',
    }),
  },
})

const emit = defineEmits<{
  (e: 'update:form-data', formData: UpdateAuthPasswordRequest): void
  (e: 'submit'): void
}>()

const rules = computed<ValidationArgs>(() => UpdateAuthPasswordValidationRules(props.formData.newPassword))
const formDataValue = computed({
  get: (): UpdateAuthPasswordRequest => props.formData,
  set: (formData: UpdateAuthPasswordRequest): void => emit('update:form-data', formData),
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
          :icon="mdiLock"
          size="32"
          class="mr-3 text-warning"
        />
        パスワード変更
      </h1>
      <p class="text-body-1 text-grey-darken-1">
        セキュリティ保護のため、定期的なパスワード変更を推奨します
      </p>
    </div>

    <v-card
      elevation="2"
      class="password-form-card"
    >
      <v-form @submit.prevent="onSubmit">
        <v-card-text class="pa-8">
          <div class="form-section mb-6">
            <h3 class="text-subtitle-1 font-weight-medium mb-4 text-grey-darken-1">
              現在の認証情報
            </h3>
            <v-text-field
              v-model="validate.oldPassword.$model"
              label="現在のパスワード"
              :type="showOldPassword ? 'text' : 'password'"
              :append-icon="showOldPassword ? mdiEye : mdiEyeOff"
              :error-messages="getErrorMessage(validate.oldPassword.$errors)"
              variant="outlined"
              class="mb-4"
              @click:append="onChangeOldPasswordFieldType"
            />
          </div>
          <div class="form-section">
            <h3 class="text-subtitle-1 font-weight-medium mb-4 text-grey-darken-1">
              新しいパスワード
            </h3>
            <v-text-field
              v-model="validate.newPassword.$model"
              label="新しいパスワード"
              :type="showNewPassword ? 'text' : 'password'"
              :append-icon="showNewPassword ? mdiEye : mdiEyeOff"
              :error-messages="getErrorMessage(validate.newPassword.$errors)"
              variant="outlined"
              class="mb-4"
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
              variant="outlined"
              @click:append="onChangePasswordConfirmationFieldType"
            />

            <v-alert
              type="info"
              variant="outlined"
              class="mt-4"
            >
              <div class="text-body-2">
                <strong>パスワード要件:</strong><br>
                • 8文字以上<br>
                • 英字・数字を含む<br>
                • 記号の使用を推奨
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
              color="warning"
              variant="elevated"
              size="large"
              class="submit-btn"
            >
              パスワードを更新
            </v-btn>
          </div>
        </v-card-actions>
      </v-form>
    </v-card>
  </v-container>
</template>

<style scoped>
.password-form-card {
  border-radius: 12px;
  max-width: 800px;
  margin: 0 auto;
}

.form-section {
  border-radius: 8px;
  padding: 20px;
  background: rgb(248 250 252);
  border-left: 4px solid rgb(255 193 7);
}

.form-section h3 {
  border-bottom: 1px solid rgb(224 224 224);
  padding-bottom: 8px;
}

@media (width <= 600px) {
  .password-form-card {
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
