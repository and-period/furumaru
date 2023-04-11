<script lang="ts" setup>
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

const headers = [
  {
    text: '名前',
    value: 'name'
  },
  {
    text: '電話番号',
    value: 'phoneNumber'
  },
  {
    text: '購入数',
    value: 'totalOrder'
  },
  {
    text: '購入金額',
    value: 'totalAmount'
  },
  {
    text: 'アカウントの有無',
    value: 'registered'
  },
  {
    text: 'Action',
    value: 'action'
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
    <v-card>
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
            {{ `${item.lastname} ${item.firstname}` }}
          </template>
          <template #[`item.totalAmount`]="{ item }">
            {{ `${item.totalAmount}` }} 円
          </template>
          <template #[`item.registered`]="{ item }">
            <v-chip size="small" :color="getStatusColor(item.registered)">
              {{ registerStatus(item.registered) }}
            </v-chip>
          </template>
          <template #[`item.action`]>
            <v-btn variant="outlined" color="primary" size="small" @click="handleEdit()">
              <v-icon size="small">
                mdi-pencil
              </v-icon>
              詳細
            </v-btn>
            <v-btn variant="outlined" color="primary" size="small">
              <v-icon size="small">
                mdi-delete
              </v-icon>
              削除
            </v-btn>
          </template>
        </v-data-table>
      </v-card-text>
    </v-card>
  </div>
</template>
