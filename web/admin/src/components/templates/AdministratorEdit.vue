<script lang="ts" setup>
import useVuelidate from '@vuelidate/core'
import { AlertType } from '~/lib/hooks'
import { getErrorMessage, maxLength, required, tel } from '~/lib/validations'
import { AdministratorResponse, UpdateAdministratorRequest } from '~/types/api'

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
  administrator: {
    type: Object as PropType<AdministratorResponse>,
    default: (): AdministratorResponse => ({
      id: '',
      lastname: '',
      lastnameKana: '',
      firstname: '',
      firstnameKana: '',
      email: '',
      phoneNumber: '',
      createdAt: 0,
      updatedAt: 0
    })
  },
  formData: {
    type: Object as PropType<UpdateAdministratorRequest>,
    default: (): UpdateAdministratorRequest => ({
      lastname: '',
      lastnameKana: '',
      firstname: '',
      firstnameKana: '',
      phoneNumber: ''
    })
  }
})

const emit = defineEmits<{
  (e: 'update:administrator', administrator: AdministratorResponse): void
  (e: 'update:form-data', formData: UpdateAdministratorRequest): void
  (e: 'submit'): void
}>()

const rules = computed(() => ({
  lastname: { required, maxLength: maxLength(16) },
  firstname: { required, maxLength: maxLength(16) },
  lastnameKana: { required, maxLength: maxLength(32) },
  firstnameKana: { required, maxLength: maxLength(32) },
  phoneNumber: { required, tel }
}))
const administratorValue = computed({
  get: (): AdministratorResponse => props.administrator,
  set: (administrator: AdministratorResponse): void => emit('update:administrator', administrator)
})
const formDataValue = computed({
  get: (): UpdateAdministratorRequest => props.formData,
  set: (formData: UpdateAdministratorRequest): void => emit('update:form-data', formData)
})

const validate = useVuelidate(rules, formDataValue)

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

  <v-card :loading="loading">
    <v-card-title>管理者編集</v-card-title>

    <form @submit.prevent="onSubmit">
      <v-card-text>
        <div class="d-flex">
          <v-text-field
            v-model="formDataValue.lastname"
            :error-messages="getErrorMessage(validate.lastname.$errors)"
            class="mr-4"
            label="管理者名:姓"
            maxlength="16"
            required
          />
          <v-text-field
            v-model="formDataValue.firstname"
            :error-messages="getErrorMessage(validate.firstname.$errors)"
            label="管理者名:名"
            maxlength="16"
            required
          />
        </div>
        <div class="d-flex">
          <v-text-field
            v-model="formDataValue.lastnameKana"
            :error-messages="getErrorMessage(validate.lastnameKana.$errors)"
            class="mr-4"
            label="管理者名:姓（ふりがな）"
            maxlength="32"
            required
          />
          <v-text-field
            v-model="formDataValue.firstnameKana"
            :error-messages="getErrorMessage(validate.firstnameKana.$errors)"
            label="管理者名:名（ふりがな）"
            maxlength="32"
            required
          />
        </div>
        <v-text-field
          v-model="administratorValue.email"
          label="連絡先（メールアドレス）"
          type="email"
          readonly
        />
        <v-text-field
          v-model="formDataValue.phoneNumber"
          :error-messages="getErrorMessage(validate.phoneNumber.$errors)"
          label="連絡先（電話番号）"
          required
        />
      </v-card-text>

      <v-card-actions>
        <v-btn block :loading="loading" type="submit" variant="outlined" color="primary">
          更新
        </v-btn>
      </v-card-actions>
    </form>
  </v-card>
</template>
