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
          name="subject"
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
          :items="priority"
          label="優先度"
          :value="getPriority(formData.priority)"
        ></v-select>

        <v-select
          :items="status"
          label="ステータス"
          :value="getStatus(formData.status)"
        ></v-select>

        <v-text-field
          name="mailAddress"
          label="メールアドレス"
          :value="formData.email"
          readonly
        ></v-text-field>

        <v-text-field
          name="telephoneNumber"
          label="電話番号"
          :value="formData.phoneNumber"
          readonly
        ></v-text-field>

        <v-textarea name="memo" label="メモ"></v-textarea>
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
  useFetch,
  useRoute,
  useRouter,
} from '@nuxtjs/composition-api'

import { useContactStore } from '~/store/contact'
import { ContactResponse } from '~/types/api'

export default defineComponent({
  setup() {
    const route = useRoute()
    const router = useRouter()
    const id = route.value.params.id

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
    })

    const getPriorityColor = (priority: string): string => {
      switch (priority) {
        case 'High':
          return 'red'
        case 'Middle':
          return 'orange'
        case 'Low':
          return 'blue'
        default:
          return ''
      }
    }

    const getPriority = (priority: number): string => {
      switch (priority) {
        case 1:
          return 'High'
        case 2:
          return 'Middle'
        case 3:
          return 'Low'
        default:
          return 'Unknown'
      }
    }

    const getStatusColor = (status: number): string => {
      switch (status) {
        case 1:
          return 'red'
        case 2:
          return 'orange'
        case 3:
          return 'blue'
        default:
          return ''
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
          return '不明'
      }
    }

    const RegisterBtn = async (): Promise<void> => {
      try {
        await contactStore.contactUpdate(formData, id)
        router.push('/')
      } catch (error) {
        console.log(error)
      }
    }

    return {
      priority: ['High', 'Middle', 'Low', 'Unknown'],
      status: ['未着手', '進行中', '完了', '不明'],
      id,
      contacts,
      fetchState,
      formData,
      getStatus,
      getPriority,
      getPriorityColor,
      getStatusColor,
      RegisterBtn,
    }
  },
})
</script>
