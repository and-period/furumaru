<template>
  <div>
    <v-card-title>お問い合わせ管理</v-card-title>
    <v-card>
      <v-card-text>
        <v-data-table
          :headers="headers"
          :items="contacts"
          :items-per-page="itemsPerPage"
          :server-items-length="total"
          :footer-props="options"
          @update:page="handleUpdatePage"
          @update:items-per-page="handleUpdateItemsPerPage"
        >
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
  watch,
} from '@nuxtjs/composition-api'
import { DataTableHeader } from 'vuetify'

import { usePagination } from '~/lib/hooks'
import { useContactStore } from '~/store/contact'
import {
  ContactPriority,
  ContactsResponseContactsInner,
  ContactStatus,
} from '~/types/api'

export default defineComponent({
  setup() {
    const router = useRouter()
    const contactStore = useContactStore()
    const {
      itemsPerPage,
      offset,
      options,
      updateCurrentPage,
      handleUpdateItemsPerPage,
    } = usePagination()

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

    useFetch(async () => {
      await fetchContacts()
    })

    const contacts = computed(() => {
      return contactStore.contacts
    })
    const total = computed(() => {
      return contactStore.total
    })

    watch(itemsPerPage, () => {
      fetchContacts()
    })

    const handleUpdatePage = async (page: number) => {
      updateCurrentPage(page)
      await fetchContacts()
    }

    const fetchContacts = async () => {
      try {
        await contactStore.fetchContacts(itemsPerPage.value, offset.value)
      } catch (err) {
        console.log(err)
      }
    }

    const getPriorityColor = (priority: ContactPriority): string => {
      switch (priority) {
        case ContactPriority.LOW:
          return 'primary'
        case ContactPriority.MIDDLE:
          return 'secondary'
        case ContactPriority.HIGH:
          return 'error'
        default:
          return 'unknown'
      }
    }

    const getPriority = (priority: ContactPriority): string => {
      switch (priority) {
        case ContactPriority.LOW:
          return '低'
        case ContactPriority.MIDDLE:
          return '中'
        case ContactPriority.HIGH:
          return '高'
        default:
          return '未設定'
      }
    }

    const getStatusColor = (status: ContactStatus): string => {
      switch (status) {
        case ContactStatus.TODO:
          return 'error'
        case ContactStatus.INPROGRESS:
          return 'secondary'
        case ContactStatus.DONE:
          return 'primary'
        case ContactStatus.DISCARD:
          return 'info'
        default:
          return 'unknown'
      }
    }

    const getStatus = (status: ContactStatus): string => {
      switch (status) {
        case ContactStatus.TODO:
          return '未着手'
        case ContactStatus.INPROGRESS:
          return '進行中'
        case ContactStatus.DONE:
          return '対応完了'
        case ContactStatus.DISCARD:
          return '対応不要'
        default:
          return '不明'
      }
    }

    const handleEdit = (item: ContactsResponseContactsInner) => {
      router.push(`/contacts/edit/${item.id}`)
    }

    return {
      headers,
      contacts,
      itemsPerPage,
      total,
      options,
      getPriority,
      getPriorityColor,
      getStatus,
      getStatusColor,
      handleEdit,
      handleUpdatePage,
      handleUpdateItemsPerPage,
    }
  },
})
</script>
