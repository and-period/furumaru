<template>
  <div>
    <v-card-title>お問い合わせ管理</v-card-title>
    <v-card>
      <v-card-text>
        <v-data-table :headers="headers" :items="contacts" :items-per-page="5">
          <template #[`item.priority`]="{ item }">
            <v-chip :color="getPriorityColor(item.priority)" small dark>
              {{ getPriority(item.priority) }}
            </v-chip>
          </template>
          <template #[`item.status`]="{ item }">
            <v-chip :color="getStatusColor(item.status)" small dark>
              {{ getStatus(item.status) }}
            </v-chip>
          </template>
          <template #[`item.actions`]="{ item }">
            <v-btn outlined color="primary" small @click="handleEdit(item)">
              <v-icon small>mdi-pencil</v-icon>
              編集
            </v-btn>
          </template>
        </v-data-table>
      </v-card-text>
    </v-card>
  </div>
</template>

<script lang="ts">
import {
  computed,
  defineComponent,
  useFetch,
  useRouter,
} from '@nuxtjs/composition-api'
import { DataTableHeader } from 'vuetify'

import { useContactStore } from '~/store/contact'
import { ContactsResponseContactsInner } from '~/types/api'

export default defineComponent({
  setup() {
    const router = useRouter()
    const contactStore = useContactStore()
    const contacts = computed(() => {
      return contactStore.contacts
    })
    const headers: DataTableHeader[] = [
      {
        text: '件名',
        value: 'title',
      },
      {
        text: 'メールアドレス',
        value: 'email',
      },
      {
        text: '優先度',
        value: 'priority',
      },
      {
        text: 'ステータス',
        value: 'status',
      },
      {
        text: 'Actions',
        value: 'actions',
        sortable: false,
      },
    ]

    const getPriorityColor = (priority: any): string => {
      switch (priority) {
        case 1:
          return 'error'
        case 2:
          return 'secondary'
        case 3:
          return 'primary'
        default:
          return 'unknown'
      }
    }

    const getPriority = (priority: any): string => {
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

    const getStatusColor = (status: any): string => {
      switch (status) {
        case 1:
          return 'error'
        case 2:
          return 'secondary'
        case 3:
          return 'primary'
        default:
          return 'unknown'
      }
    }

    const getStatus = (status: any): string => {
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

    const handleEdit = (item: ContactsResponseContactsInner) => {
      router.push(`/contacts/edit/${item.id}`)
    }

    const { fetchState } = useFetch(async () => {
      try {
        await contactStore.fetchContacts()
      } catch (err) {
        console.log(err)
      }
    })

    return {
      headers,
      contacts,
      fetchState,
      getPriority,
      getPriorityColor,
      getStatus,
      getStatusColor,
      handleEdit,
    }
  },
})
</script>
