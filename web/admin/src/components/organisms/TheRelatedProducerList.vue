<script lang="ts" setup>
import { DataTableHeader } from 'vuetify'

import { useCoordinatorStore } from '~/store'
import { ProducersResponseProducersInner } from '~/types/api'

const props = defineProps({
  tableFooterProps: {
    type: Object,
    default: () => {}
  }
})

const emit = defineEmits<{
  (e: 'update:items-per-page', page: number): void
  (e: 'update:page', page: number): void
}>()

const router = useRouter()
const coordinatorStore = useCoordinatorStore()

const relateProducers = computed(() => {
  console.log(coordinatorStore.producers)
  return coordinatorStore.producers
})

const totalItems = computed(() => {
  return coordinatorStore.totalItems
})

const producerHeaders: DataTableHeader[] = [
  {
    text: 'サムネイル',
    value: 'thumbnailUrl'
  },
  {
    text: '生産者名',
    value: 'name'
  },
  {
    text: '店舗名',
    value: 'storeName'
  },
  {
    text: 'Email',
    value: 'email'
  },
  {
    text: '電話番号',
    value: 'phoneNumber'
  },
  {
    text: 'Actions',
    value: 'actions',
    sortable: false
  }
]

const handleUpdateItemsPerPage = (page: number) => {
  emit('update:items-per-page', page)
}

const handleUpdatePage = (page: number) => {
  emit('update:page', page)
}

const convertPhone = (phoneNumber: string): string => {
  return phoneNumber.replace('+81', '0')
}

const handleEdit = (item: ProducersResponseProducersInner) => {
  router.push(`/producers/edit/${item.id}`)
}
</script>

<template>
  <div>
    <v-data-table
      :headers="producerHeaders"
      :items="relateProducers"
      :server-items-length="totalItems"
      :footer-props="props.tableFooterProps"
      no-data-text="登録されている生産者がいません。"
      @update:items-per-page="handleUpdateItemsPerPage"
      @update:page="handleUpdatePage"
    >
      <template #[`item.thumbnailUrl`]="{ item }">
        <v-avatar>
          <img
            v-if="item.thumbnailUrl !== ''"
            :src="item.thumbnailUrl"
            :alt="`${item.storeName}-profile`"
          >
          <v-icon v-else>
            mdi-account
          </v-icon>
        </v-avatar>
      </template>
      <template #[`item.name`]="{ item }">
        {{ `${item.lastname} ${item.firstname}` }}
      </template>
      <template #[`item.phoneNumber`]="{ item }">
        {{ convertPhone(item.phoneNumber) }}
      </template>
      <template #[`item.actions`]="{ item }">
        <v-btn variant="outlined" color="primary" size="small" @click="handleEdit(item)">
          <v-icon size="small">
            mdi-pencil
          </v-icon>
          編集
        </v-btn>
      </template>
    </v-data-table>
  </div>
</template>
