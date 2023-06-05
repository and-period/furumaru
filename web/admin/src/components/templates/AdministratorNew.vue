<script lang="ts" setup>
import { AlertType } from '~/lib/hooks'
import { CreateAdministratorRequest } from '~/types/api'

const props = defineProps({
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

const onSubmit = (): void => {
  emit('submit')
}
</script>

<template>
  <v-alert v-show="props.isAlert" :type="props.alertType" v-text="props.alertText" />

  <v-card>
    <v-card-title>管理者登録</v-card-title>

    <form @submit.prevent="onSubmit">
      <v-card-text>
        <div class="d-flex">
          <v-text-field
            v-model="formDataValue.lastname"
            class="mr-4"
            label="管理者名:姓"
            maxlength="16"
            required
          />
          <v-text-field
            v-model="formDataValue.firstname"
            label="管理者名:名"
            maxlength="16"
            required
          />
        </div>
        <div class="d-flex">
          <v-text-field
            v-model="formDataValue.lastnameKana"
            class="mr-4"
            label="管理者名:姓（ふりがな）"
            maxlength="32"
            required
          />
          <v-text-field
            v-model="formDataValue.firstnameKana"
            label="管理者名:名（ふりがな）"
            maxlength="32"
            required
          />
        </div>
        <v-text-field
          v-model="formDataValue.email"
          label="連絡先（メールアドレス）"
          type="email"
          required
        />
        <v-text-field
          v-model="formDataValue.phoneNumber"
          label="連絡先（電話番号）"
          required
        />
      </v-card-text>

      <v-card-actions>
        <v-btn block variant="outlined" color="primary" type="submit">
          登録
        </v-btn>
      </v-card-actions>
    </form>
  </v-card>
</template>
