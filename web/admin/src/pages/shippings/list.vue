<script lang="ts" setup>
import { useAuthStore, useShippingStore } from '~/store'
import { dateTimeFormatter } from '~/lib/formatter/day'

const authStore = useAuthStore()
const { adminId } = storeToRefs(authStore)

const shippingStore = useShippingStore()
const { fetchShippings } = shippingStore
const router = useRouter()

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

const itemsPerPage = ref<number>(20)

const { data, status, error } = useAsyncData(async () => {
  return fetchShippings(adminId.value)
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
</script>

<template>
  <v-card
    :loading="status === 'pending'"
  >
    <v-card-title>配送先一覧</v-card-title>
    <v-card-text>
      <v-data-table-server
        hover
        :headers="headers"
        :items="shippings"
        :items-per-page-options="[5, 10, 20, 100]"
        :total-items="totalItems"
        :items-length="totalItems"
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
</template>
