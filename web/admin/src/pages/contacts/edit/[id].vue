<script lang="ts" setup>
import { useContactStore } from '~/store'
import {
  ContactPriority,
  ContactResponse,
  ContactStatus,
  UpdateContactRequest
} from '~/types/api'

const route = useRoute()
const router = useRouter()
const id = route.params.id as string

const priorities = ['High', 'Middle', 'Low']
const statuses = ['未着手', '進行中', '完了']

const formNote = ref<string>('')
const formPriority = ref<string>('')
const formStatus = ref<string>('')

const { getContact } = useContactStore()

const contactStore = useContactStore()
const contacts = computed(() => {
  return contactStore.contacts
})

const formData = reactive<ContactResponse>({
  id,
  title: '',
  content: '',
  username: '',
  email: '',
  phoneNumber: '',
  status: 0,
  priority: 0,
  note: '',
  createdAt: 0,
  updatedAt: 0
})

useAsyncData(async () => {
  const contact = await getContact(id)
  formData.title = contact.title
  formData.content = contact.content
  formData.username = contact.username
  formData.email = contact.email
  formData.phoneNumber = convertPhoneNumber(contact.phoneNumber)
  formData.status = contact.status
  formData.priority = contact.priority
  formData.note = contact.note
  formData.createdAt = contact.createdAt
  formData.updatedAt = contact.updatedAt
  formPriority.value = getPriority(formData.priority)
  formStatus.value = getStatus(formData.status)
  formNote.value = formData.note
})

const convertPhoneNumber = (phoneNumber: string): string => {
  return phoneNumber.replace('+81', '0')
}

const getPriority = (priority: ContactPriority): string => {
  switch (priority) {
    case ContactPriority.LOW:
      return 'High'
    case ContactPriority.MIDDLE:
      return 'Middle'
    case ContactPriority.HIGH:
      return 'Low'
    default:
      return 'Middle'
  }
}

const getStatus = (status: ContactStatus): string => {
  switch (status) {
    case ContactStatus.TODO:
      return '未着手'
    case ContactStatus.INPROGRESS:
      return '進行中'
    case ContactStatus.DONE:
      return '完了'
    default:
      return '未着手'
  }
}

const handleEdit = async (): Promise<void> => {
  try {
    const payload = reactive<UpdateContactRequest>({
      status: getStatusID(formStatus.value),
      priority: getPriorityID(formPriority.value),
      note: formNote.value
    })

    await contactStore.contactUpdate(payload, id)
    router.push('/contacts')
  } catch (error) {
    console.log(error)
  }
}

const getPriorityID = (priority: string): ContactPriority => {
  switch (priority) {
    case 'High':
      return ContactPriority.HIGH
    case 'Middle':
      return ContactPriority.MIDDLE
    case 'Low':
      return ContactPriority.LOW
    default:
      return ContactPriority.MIDDLE
  }
}

const getStatusID = (status: string): ContactStatus => {
  switch (status) {
    case '未着手':
      return ContactStatus.TODO
    case '進行中':
      return ContactStatus.INPROGRESS
    case '完了':
      return ContactStatus.DONE
    default:
      return ContactStatus.TODO
  }
}
</script>

<template>
  <div>
    <v-card-title>お問合せ管理</v-card-title>
    <v-card elevation="0">
      <v-card-text>
        <v-text-field
          v-model="formData.username"
          name="name"
          label="名前"
          readonly
        />

        <v-text-field
          v-model="formData.title"
          name="title"
          label="件名"
          readonly
        />

        <v-textarea
          v-model="formData.content"
          name="contact"
          label="お問合せ内容"
          readonly
        />

        <v-select
          v-model="formPriority"
          :items="priorities"
          label="優先度"
          :item-value="getPriority(formData.priority)"
        />

        <v-select
          v-model="formStatus"
          :items="statuses"
          label="ステータス"
          :item-value="getStatus(formData.status)"
        />

        <v-text-field
          v-model="formData.email"
          name="mailAddress"
          label="メールアドレス"
          readonly
        />

        <v-text-field
          v-model="formData.phoneNumber"
          name="phoneNumber"
          label="電話番号"
          readonly
        />

        <v-textarea
          v-model="formNote"
          name="note"
          label="メモ"
        />
      </v-card-text>
      <v-card-actions>
        <v-btn block outlined color="primary" @click="handleEdit">
          更新
        </v-btn>
      </v-card-actions>
    </v-card>
  </div>
</template>
