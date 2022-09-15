<template>
  <div>
    <v-card-title>お問合せ管理</v-card-title>
    <v-card elevation="0">
      <v-card-text>
        <v-text-field
          name="name"
          label="名前"
          :value="formData.username"
          readonly
        ></v-text-field>

        <v-text-field
          name="title"
          label="件名"
          :value="formData.title"
          readonly
        ></v-text-field>

        <v-textarea
          name="contact"
          label="お問合せ内容"
          :value="formData.content"
          readonly
        ></v-textarea>

        <v-select
          v-model="formPriority"
          :items="priority"
          label="優先度"
          :value="getPriority(formData.priority)"
        ></v-select>

        <v-select
          v-model="formStatus"
          :items="status"
          label="ステータス"
          :item-value="getStatus(formData.status)"
        ></v-select>

        <v-text-field
          name="mailAddress"
          label="メールアドレス"
          :value="formData.email"
          readonly
        ></v-text-field>

        <v-text-field
          name="phoneNumber"
          label="電話番号"
          :value="formData.phoneNumber"
          readonly
        ></v-text-field>

        <v-textarea
          v-model="formNote"
          name="note"
          label="メモ"
          :value="formData.note"
        ></v-textarea>
      </v-card-text>
      <v-card-actions>
        <v-btn block outlined color="primary" @click="RegisterBtn">登録</v-btn>
      </v-card-actions>
    </v-card>
  </div>
</template>

<script lang="ts">
import {
  computed,
  defineComponent,
  reactive,
  ref,
  useFetch,
  useRoute,
  useRouter,
} from '@nuxtjs/composition-api'

import { useContactStore } from '~/store/contact'
import { ContactResponse, UpdateContactRequest } from '~/types/api'

export default defineComponent({
  setup() {
    const route = useRoute()
    const router = useRouter()
    const id = route.value.params.id

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
      updatedAt: 0,
    })

    const { fetchState } = useFetch(async () => {
      const contact = await getContact(id)
      formData.title = contact.title
      formData.content = contact.content
      formData.username = contact.username
      formData.email = contact.email
      formData.phoneNumber = contact.phoneNumber
      formData.status = contact.status
      formData.priority = contact.priority
      formData.note = contact.note
      formData.createdAt = contact.createdAt
      formData.updatedAt = contact.updatedAt
      formPriority.value = getPriority(formData.priority)
      formStatus.value = getStatus(formData.status)
      formNote.value = formData.note
    })

    const getPriority = (priority: number): string => {
      switch (priority) {
        case 1:
          return 'High'
        case 2:
          return 'Middle'
        case 3:
          return 'Low'
        default:
          return 'Middle'
      }
    }

    const getStatus = (status: number): string => {
      switch (status) {
        case 1:
          return '未着手'
        case 2:
          return '進行中'
        case 3:
          return '完了'
        default:
          return '未着手'
      }
    }

    const RegisterBtn = async (): Promise<void> => {
      try {
        const payload = reactive<UpdateContactRequest>({
          status: getStatusID(formStatus.value),
          priority: getPriorityID(formPriority.value),
          note: formNote.value,
        })

        await contactStore.contactUpdate(payload, id)
        router.push('/')
      } catch (error) {
        console.log(error)
      }
    }

    const getPriorityID = (priority: string): number => {
      switch (priority) {
        case 'High':
          return 1
        case 'Middle':
          return 2
        case 'Low':
          return 3
        default:
          return 2
      }
    }

    const getStatusID = (status: string): number => {
      switch (status) {
        case '未着手':
          return 1
        case '進行中':
          return 2
        case '完了':
          return 3
        default:
          return 1
      }
    }

    return {
      priority: ['High', 'Middle', 'Low'],
      status: ['未着手', '進行中', '完了'],
      id,
      contacts,
      fetchState,
      formData,
      getStatus,
      getPriority,
      RegisterBtn,
      formNote,
      formPriority,
      formStatus,
      getPriorityID,
      getStatusID,
    }
  },
})
</script>
