<script lang="ts" setup>
import useVuelidate from '@vuelidate/core'
import { UpdateAuthEmailRequest } from '~/types/api'
import { AlertType } from '~/lib/hooks'
import { required, email, getErrorMessage } from '~/lib/validations'

interface Props {
  isAlert: boolean,
  alertType: AlertType,
  alertText: string,
  formData: UpdateAuthEmailRequest
}

const props = defineProps<Props>()

const emits = defineEmits<{
  (e: 'submit'): void,
  (e: 'update:isAlert', val: boolean): void
  (e: 'update:formData', val: UpdateAuthEmailRequest): void
}>()

const isAlertValue = computed({
  get: () => props.isAlert,
  set: (val: boolean) => emits('update:isAlert', val)
})

const formDataValue = computed({
  get: () => props.formData,
  set: (val: UpdateAuthEmailRequest) => emits('update:formData', val)
})

const rules = computed(() => {
  return {
    email: { required, email }
  }
})

const v$ = useVuelidate(rules, formDataValue)

const onSubmit = async () => {
  const result = await v$.value.$validate()
  if (!result) {
    return
  }
  emits('submit')
}
</script>

<template>
  <v-alert v-model="isAlertValue" :type="props.alertType" :text="props.alertText" class="mb-2" />
  <v-card elevation="0">
    <v-card-title>メールアドレス変更</v-card-title>
    <v-form @submit.prevent="onSubmit">
      <v-card-text>
        <v-text-field v-model="v$.email.$model" :error-messages="getErrorMessage(v$.email.$errors)" label="新規メールアドレス" />
      </v-card-text>
      <v-card-actions>
        <v-btn type="submit" block color="primary" variant="outlined">
          変更
        </v-btn>
      </v-card-actions>
    </v-form>
  </v-card>
</template>
