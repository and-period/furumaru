<script lang="ts" setup>
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
  <v-alert
    v-show="props.isAlert"
    :type="props.alertType"
    v-text="props.alertText"
  />

  <v-card elevation="0">
    <v-card-title>メールアドレス変更</v-card-title>

    <v-form @submit.prevent="onSubmit">
      <v-card-text>
        <v-text-field
          v-model="validate.email.$model"
          :error-messages="getErrorMessage(validate.email.$errors)"
          label="新規メールアドレス"
          type="email"
        />
      </v-card-text>
      <v-card-actions>
        <v-btn
          :loading="loading"
          block
          type="submit"
          color="primary"
          variant="outlined"
        >
          変更
        </v-btn>
      </v-card-actions>
    </v-form>
  </v-card>
</template>
