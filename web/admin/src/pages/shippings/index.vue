<script lang="ts" setup>
import { mdiPlus } from '@mdi/js'
import { useAuthStore, useShippingStore } from '~/store'
import { dateTimeFormatter } from '~/lib/formatter/day'
import { useAlert, usePagination } from '~/lib/hooks'

const authStore = useAuthStore()
const { adminId } = storeToRefs(authStore)

const shippingStore = useShippingStore()
const { fetchShippings } = shippingStore
const router = useRouter()

const { alertType, isShow, alertText, show, hide } = useAlert('error')

const headers: VDataTable['headers'] = [
  {
    title: 'id',
    key: 'id',
    sortable: false,
  },
  {
    title: 'デフォルト配送先',
    key: 'isDefault',
  },
  {
    title: '作成日',
    key: 'createdAt',
  },
  {
    title: '更新日',
    key: 'updatedAt',
  },
]

const pagination = usePagination()

const { data, status, error, refresh } = useAsyncData(async () => {
  hide()
  return fetchShippings(adminId.value, pagination.itemsPerPage.value, pagination.offset.value)
})

watch(error, (newError) => {
  if (newError) {
    if (newError instanceof Error) {
      show(newError.message)
    }
    console.log(newError)
  }
})

const shippings = computed(() => {
  if (data.value) {
    return data.value.shippings
  }
  else {
    return []
  }
})

const totalItems = computed(() => {
  if (data.value) {
    return data.value.total
  }
  else {
    return 0
  }
})

const handleUpdatePage = async (page: number): Promise<void> => {
  pagination.updateCurrentPage(page)
  await refresh()
}
</script>

<template>
  <div>
    <v-alert
      v-show="isShow"
      :type="alertType"
      v-text="alertText"
    />

    <v-card
      :loading="status === 'pending'"
    >
      <v-card-title class="d-flex">
        配送先一覧
        <v-spacer />
        <v-btn
          variant="outlined"
          to="/shippings/new"
        >
          <v-icon :icon="mdiPlus" />
          新規作成
        </v-btn>
      </v-card-title>
      <v-card-text>
        <v-data-table-server
          hover
          :headers="headers"
          :items="shippings"
          :items-per-page-options="[5, 10, 20, 100]"
          :total-items="totalItems"
          :items-length="totalItems"
          :table-items-per-page="pagination.itemsPerPage.value"
          @click:update-page="handleUpdatePage"
          @click:update-items-per-page="pagination.handleUpdateItemsPerPage"
          @click:row="(_: any, { item } : any) => {
            router.push(`/shippings/${item.id}`)
          }"
        >
          <template #[`item.createdAt`]="{ item }">
            {{ dateTimeFormatter(item.createdAt) }}
          </template>

          <template #[`item.updatedAt`]="{ item }">
            {{ dateTimeFormatter(item.updatedAt) }}
          </template>
        </v-data-table-server>
      </v-card-text>
    </v-card>
  </div>
</template>
