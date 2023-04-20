<script lang="ts" setup>
import { UpdateAuthEmailRequest } from '~/types/api'
import { AlertType } from '~/lib/hooks'

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
    type: Object as PropType<UpdateAuthEmailRequest>,
    default: () => ({
      email: ''
    })
  }
})

const emit = defineEmits<{
  (e: 'submit'): void
}>()

const onSubmit = (): void => {
  emit('submit')
}
</script>

<template>
  <v-alert v-model="props.isAlert" :type="props.alertType" :text="props.alertText" />
  <v-card elevation="0">
    <v-card-title>メールアドレス変更</v-card-title>
    <v-form @submit.prevent="onSubmit">
      <v-card-text>
        <v-text-field v-model="props.formData.email" label="新規メールアドレス" />
      </v-card-text>
      <v-card-actions>
        <v-btn type="submit" block color="primary" variant="outlined">
          変更
        </v-btn>
      </v-card-actions>
    </v-form>
  </v-card>
</template>
