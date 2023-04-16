<script lang="ts" setup>
import { VDataTable } from 'vuetify/labs/components'
import { mdiPencil, mdiDelete } from '@mdi/js'

import { usePagination } from '~/lib/hooks'
import { useUserStore } from '~/store/customer'

const router = useRouter()
const userStore = useUserStore()
const {
  itemsPerPage,
  offset,
  options,
  updateCurrentPage,
  handleUpdateItemsPerPage
} = usePagination()
const id = 'ThisIsID'

const headers: VDataTable['headers'] = [
  {
    title: '名前',
    key: 'name'
  },
  {
    title: '電話番号',
    key: 'phoneNumber'
  },
  {
    title: '購入数',
    key: 'totalOrder'
  },
  {
    title: '購入金額',
    key: 'totalAmount'
  },
  {
    title: 'アカウントの有無',
    key: 'registered'
  },
  {
    title: 'Action',
    key: 'action'
  }
]

const fetchState = useAsyncData(async () => {
  await fetchUsers()
})

const users = computed(() => {
  return userStore.users
})
const total = computed(() => {
  return userStore.totalItems
})

watch(itemsPerPage, () => {
  fetchUsers()
})

const handleUpdatePage = async (page: number) => {
  updateCurrentPage(page)
  await fetchUsers()
}

const fetchUsers = async () => {
  try {
    await userStore.fetchUsers(itemsPerPage.value, offset.value)
  } catch (err) {
    console.log(err)
  }
}

const getStatusColor = (account: boolean): string => {
  return account ? 'primary' : 'red'
}

const registerStatus = (registered: boolean): string => {
  return registered ? '有' : '無'
}

const handleEdit = () => {
  router.push(`/customers/edit/${id}`)
}

try {
  await fetchState.execute()
} catch (err) {
  console.log('failed to setup', err)
}
</script>

<template>
  <div>
    <v-card-title>顧客管理</v-card-title>
    <v-card flat>
      <v-card-text>
        <v-data-table
          :headers="headers"
          :items="users"
          :items-per-page="itemsPerPage"
          :footer-props="options"
          no-data-text="登録されている顧客情報がありません"
          @update:page="handleUpdatePage"
          @update:items-per-page="handleUpdateItemsPerPage"
        >
          <template #[`item.name`]="{ item }">
            {{ `${item.raw.lastname} ${item.raw.firstname}` }}
          </template>
          <template #[`item.totalAmount`]="{ item }">
            {{ `${item.raw.totalAmount}` }} 円
          </template>
          <template #[`item.registered`]="{ item }">
            <v-chip size="small" :color="getStatusColor(item.raw.registered)">
              {{ registerStatus(item.raw.registered) }}
            </v-chip>
          </template>
          <template #[`item.action`]>
            <v-btn class="mr-2" variant="outlined" color="primary" size="small" @click="handleEdit()">
              <v-icon size="small" :icon="mdiPencil" />
              詳細
            </v-btn>
            <v-btn variant="outlined" color="primary" size="small">
              <v-icon size="small" :icon="mdiDelete" />
              削除
            </v-btn>
          </template>
        </v-data-table>
      </v-card-text>
    </v-card>
  </div>
</template>
