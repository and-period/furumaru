<script lang="ts" setup>
import useVuelidate from '@vuelidate/core'
import type { AlertType } from '~/lib/hooks'
import { getErrorMessage } from '~/lib/validations'
import type { CreateAdministratorRequest } from '~/types/api'
import { CreateAdministratorValidationRules } from '~/types/validations'

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
    type: Object as PropType<CreateAdministratorRequest>,
    default: (): CreateAdministratorRequest => ({
      lastname: '',
      lastnameKana: '',
      firstname: '',
      firstnameKana: '',
      email: '',
      phoneNumber: ''
    })
  }
})

const emit = defineEmits<{
  (e: 'update:form-data', formData: CreateAdministratorRequest): void
  (e: 'submit'): void
}>()

const formDataValue = computed({
  get: (): CreateAdministratorRequest => props.formData,
  set: (formData: CreateAdministratorRequest): void => emit('update:form-data', formData)
})

const validate = useVuelidate(CreateAdministratorValidationRules, formDataValue)

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

  <v-card>
    <v-card-title>管理者登録</v-card-title>

    <v-form @submit.prevent="onSubmit">
      <v-card-text>
        <v-row>
          <v-col>
            <v-text-field
              v-model="validate.lastname.$model"
              :error-messages="getErrorMessage(validate.lastname.$errors)"
              class="mr-4"
              label="管理者名:姓"
            />
          </v-col>
          <v-col>
            <v-text-field
              v-model="validate.firstname.$model"
              :error-messages="getErrorMessage(validate.firstname.$errors)"
              label="管理者名:名"
            />
          </v-col>
        </v-row>
        <v-row>
          <v-col>
            <v-text-field
              v-model="validate.lastnameKana.$model"
              :error-messages="getErrorMessage(validate.lastnameKana.$errors)"
              class="mr-4"
              label="管理者名:姓（ふりがな）"
            />
          </v-col>
          <v-col>
            <v-text-field
              v-model="validate.firstnameKana.$model"
              :error-messages="getErrorMessage(validate.firstnameKana.$errors)"
              label="管理者名:名（ふりがな）"
            />
          </v-col>
        </v-row>
        <v-text-field
          v-model="validate.email.$model"
          :error-messages="getErrorMessage(validate.email.$errors)"
          label="連絡先（メールアドレス）"
          type="email"
        />
        <v-text-field
          v-model="validate.phoneNumber.$model"
          :error-messages="getErrorMessage(validate.phoneNumber.$errors)"
          label="連絡先（電話番号）"
        />
      </v-card-text>

      <v-card-actions>
        <v-btn :loading="loading" block type="submit" variant="outlined" color="primary">
          登録
        </v-btn>
      </v-card-actions>
    </v-form>
  </v-card>
</template>
