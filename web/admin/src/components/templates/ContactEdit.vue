<script lang="ts" setup>
import { AlertType } from '~/lib/hooks'
import { ContactPriority, ContactResponse, ContactStatus, UpdateContactRequest } from '~/types/api'

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
  contact: {
    type: Object as PropType<ContactResponse>,
    default: () => ({
      id: '',
      title: '',
      content: '',
      username: '',
      email: '',
      phoneNumber: '',
      status: ContactStatus.UNKNOWN,
      priority: ContactPriority.UNKNOWN,
      note: '',
      createdAt: 0,
      updatedAt: 0
    })
  },
  formData: {
    type: Object as PropType<UpdateContactRequest>,
    default: () => ({
      status: ContactStatus.UNKNOWN,
      proprity: ContactPriority.UNKNOWN,
      note: ''
    })
  }
})

const emit = defineEmits<{
  (e: 'submit'): void
}>()

const priorities = [
  { title: '高', value: ContactPriority.HIGH },
  { title: '中', value: ContactPriority.MIDDLE },
  { title: '低', value: ContactPriority.LOW },
  { title: '未設定', value: ContactPriority.UNKNOWN }
]
const statuses = [
  { title: '未着手', value: ContactStatus.TODO },
  { title: '進行中', value: ContactStatus.INPROGRESS },
  { title: '完了', value: ContactStatus.DONE },
  { title: '対応不要', value: ContactStatus.DISCARD },
  { title: '不明', value: ContactStatus.UNKNOWN }
]

const convertPhoneNumber = computed<string>(() => {
  const phoneNumber = props.contact.phoneNumber.replace('+81', '0')
  return phoneNumber
})

const onSubmit = (): void => {
  emit('submit')
}
</script>

<template>
  <v-alert v-show="props.isAlert" :type="props.alertType" v-text="props.alertText" />
  <v-card elevation="0">
    <v-card-title>お問合せ管理</v-card-title>
    <v-card-text>
      <v-text-field
        v-model="props.contact.username"
        name="name"
        label="名前"
        readonly
      />

      <v-text-field
        v-model="props.contact.title"
        name="title"
        label="件名"
        readonly
      />

      <v-textarea
        v-model="props.contact.content"
        name="contact"
        label="お問合せ内容"
        readonly
      />

      <v-select
        v-model="props.formData.priority"
        :items="priorities"
        item-title="title"
        item-value="value"
        label="優先度"
      />

      <v-select
        v-model="props.formData.status"
        :items="statuses"
        item-title="title"
        item-value="value"
        label="ステータス"
      />

      <v-text-field
        v-model="props.contact.email"
        name="mailAddress"
        label="メールアドレス"
        readonly
      />

      <v-text-field
        v-model="convertPhoneNumber"
        name="phoneNumber"
        label="電話番号"
        readonly
      />

      <v-textarea
        v-model="props.formData.note"
        name="note"
        label="メモ"
      />
    </v-card-text>
    <v-card-actions>
      <v-btn block variant="outlined" color="primary" @click="onSubmit">
        更新
      </v-btn>
    </v-card-actions>
  </v-card>
</template>
