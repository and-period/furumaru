<script lang="ts" setup>
import type { AlertType } from '~/lib/hooks'
import type { ResetAuthPasswordRequest } from '~/types/api/v1'

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
    type: Object as PropType<ResetAuthPasswordRequest>,
    default: (): ResetAuthPasswordRequest => ({
      email: '',
      verifyCode: '',
      password: '',
      passwordConfirmation: '',
    }),
  },
})

const emit = defineEmits<{
  (e: 'update:form-data', formData: ResetAuthPasswordRequest): void
  (e: 'click:cancel'): void
  (e: 'submit'): void
}>()

const formDataValue = computed({
  get: (): ResetAuthPasswordRequest => props.formData,
  set: (v: ResetAuthPasswordRequest): void => emit('update:form-data', v),
})

const onClickCancel = (): void => {
  emit('click:cancel')
}

const onSubmit = (): void => {
  emit('submit')
}
</script>

<template>
  <atoms-app-alert
    :show="props.isAlert"
    :type="props.alertType"
    :text="props.alertText"
  />

  <atoms-app-logo-with-title
    :width="450"
    class="mx-auto py-8"
  />

  <v-card>
    <v-card-title>パスワードリセット</v-card-title>

    <v-form @submit.prevent="onSubmit">
      <v-card-text>
        <v-text-field
          v-model="formDataValue.email"
          required
          label="メールアドレス"
          type="email"
        />
      </v-card-text>
      <v-card-actions>
        <v-spacer />
        <v-btn
          color="error"
          variant="outlined"
          @click="onClickCancel"
        >
          サインイン画面にもどる
        </v-btn>
        <v-btn
          :loading="loading"
          type="submit"
          color="primary"
          variant="outlined"
        >
          次へ
        </v-btn>
      </v-card-actions>
    </v-form>
  </v-card>
</template>
