<script lang="ts" setup>
import { mdiPencil } from '@mdi/js'
import { VDataTable } from 'vuetify/lib/labs/components'

import { usePagination } from '~/lib/hooks'
import { useContactStore } from '~/store'
import {
  ContactPriority,
  ContactsResponseContactsInner,
  ContactStatus
} from '~/types/api'

const router = useRouter()
const contactStore = useContactStore()
const {
  itemsPerPage,
  offset,
  options,
  updateCurrentPage,
  handleUpdateItemsPerPage
} = usePagination()

const headers: VDataTable['headers'] = [
  {
    title: '件名',
    key: 'title'
  },
  {
    title: 'メールアドレス',
    key: 'email'
  },
  {
    title: '優先度',
    key: 'priority'
  },
  {
    title: 'ステータス',
    key: 'status'
  },
  {
    title: 'Actions',
    key: 'actions',
    sortable: false
  }
]

const fetchState = useAsyncData(async () => {
  await fetchContacts()
})

const sortBy = ref<VDataTable['sortBy']>([])

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
    const orders: string[] = sortBy.value?.map((item) => {
      switch (item.order) {
        case 'asc':
          return item.key
        case 'desc':
          return `-${item.key}`
        default:
          return item.order ? item.key : `-${item.key}`
      }
    }) || []

    await contactStore.fetchContacts(itemsPerPage.value, offset.value, orders)
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

try {
  await fetchState.execute()
} catch (err) {
  console.log('failed to setup', err)
}
</script>

<template>
  <div>
    <v-card-title>お問い合わせ管理</v-card-title>
    <v-card>
      <v-card-text>
        <v-data-table-server
          v-model:sort-by="sortBy"
          :headers="headers"
          :items="contacts"
          :items-per-page="itemsPerPage"
          :items-length="total"
          :footer-props="options"
          :multi-sort="true"
          @update:page="handleUpdatePage"
          @update:items-per-page="handleUpdateItemsPerPage"
          @update:sort-by="fetchState.refresh"
          @update:sort-desc="fetchState.refresh"
        >
          <template #[`item.priority`]="{ item }">
            <v-chip :color="getPriorityColor(item.raw.priority)" size="small">
              {{ getPriority(item.raw.priority) }}
            </v-chip>
          </template>
          <template #[`item.status`]="{ item }">
            <v-chip :color="getStatusColor(item.raw.status)" size="small">
              {{ getStatus(item.raw.status) }}
            </v-chip>
          </template>
          <template #[`item.actions`]="{ item }">
            <v-btn variant="outlined" color="primary" size="small" @click="handleEdit(item.raw)">
              <v-icon size="small" :icon="mdiPencil" />
              編集
            </v-btn>
          </template>
        </v-data-table-server>
      </v-card-text>
    </v-card>
  </div>
</template>
