<script lang="ts" setup>
import useVuelidate from '@vuelidate/core'
import { convertI18nToJapanesePhoneNumber } from '~/lib/formatter'
import type { AlertType } from '~/lib/hooks'
import { getErrorMessage, maxLength } from '~/lib/validations'
import { ContactStatus } from '~/types/api/v1'
import type { ContactResponse, UpdateContactRequest } from '~/types/api/v1'

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
  contact: {
    type: Object as PropType<ContactResponse>,
    default: () => ({
      id: '',
      title: '',
      content: '',
      username: '',
      email: '',
      phoneNumber: '',
      status: ContactStatus.ContactStatusUnknown,
      note: '',
      createdAt: 0,
      updatedAt: 0,
    }),
  },
  formData: {
    type: Object as PropType<UpdateContactRequest>,
    default: () => ({
      status: ContactStatus.ContactStatusUnknown,
      note: '',
    }),
  },
})

const emit = defineEmits<{
  (e: 'update:contact', contact: ContactResponse): void
  (e: 'update:form-data', formData: UpdateContactRequest): void
  (e: 'submit'): void
}>()

const statuses = [
  { title: '未着手', value: ContactStatus.ContactStatusWaiting },
  { title: '進行中', value: ContactStatus.ContactStatusInprogress },
  { title: '完了', value: ContactStatus.ContactStatusDone },
  { title: '対応不要', value: ContactStatus.ContactStatusDiscard },
  { title: '不明', value: ContactStatus.ContactStatusUnknown },
]

const rules = computed(() => ({
  status: {},
  note: { maxLength: maxLength(2000) },
}))
const formDataValue = computed({
  get: (): UpdateContactRequest => props.formData,
  set: (v: UpdateContactRequest): void => emit('update:form-data', v),
})
const contactValue = computed((): ContactResponse => {
  return props.contact
})
const phoneNumber = computed((): string => {
  return convertI18nToJapanesePhoneNumber(props.contact.phoneNumber)
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
  <atoms-app-alert
    :show="props.isAlert"
    :type="props.alertType"
    :text="props.alertText"
  />

  <v-card elevation="0">
    <v-card-title>お問合せ管理</v-card-title>

    <v-form @submit.prevent="onSubmit">
      <v-card-text>
        <v-text-field
          v-model="contactValue.username"
          name="name"
          label="名前"
          readonly
        />
        <v-text-field
          v-model="contactValue.title"
          name="title"
          label="件名"
          readonly
        />
        <v-textarea
          v-model="contactValue.content"
          name="contact"
          label="お問合せ内容"
          readonly
        />
        <v-select
          v-model="validate.status.$model"
          :error-messages="getErrorMessage(validate.status.$errors)"
          :items="statuses"
          item-title="title"
          item-value="value"
          label="ステータス"
        />
        <v-text-field
          v-model="contactValue.email"
          name="mailAddress"
          label="メールアドレス"
          readonly
        />
        <v-text-field
          v-model="phoneNumber"
          name="phoneNumber"
          label="電話番号"
          readonly
        />
        <v-textarea
          v-model="validate.note.$model"
          :error-messages="getErrorMessage(validate.note.$errors)"
          name="note"
          label="メモ"
        />
      </v-card-text>

      <v-card-actions>
        <v-btn
          :loading="loading"
          block
          variant="outlined"
          color="primary"
        >
          更新
        </v-btn>
      </v-card-actions>
    </v-form>
  </v-card>
</template>
