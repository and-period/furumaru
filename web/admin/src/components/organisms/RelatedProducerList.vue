<script lang="ts" setup>
import { mdiAccount, mdiPencil } from '@mdi/js'
import { VDataTable } from 'vuetify/lib/labs/components.mjs'

import { useCoordinatorStore } from '~/store'
import { Producer } from '~/types/api'

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
  return coordinatorStore.producers
})

const totalItems = computed(() => {
  return coordinatorStore.totalItems
})

const producerHeaders: VDataTable['headers'] = [
  {
    title: 'サムネイル',
    key: 'thumbnailUrl'
  },
  {
    title: '生産者名',
    key: 'name'
  },
  {
    title: '店舗名',
    key: 'storeName'
  },
  {
    title: 'Email',
    key: 'email'
  },
  {
    title: '電話番号',
    key: 'phoneNumber'
  },
  {
    title: '',
    key: 'actions',
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

const handleEdit = (item: Producer) => {
  router.push(`/producers/${item.id}`)
}
</script>

<template>
  <div>
    <v-data-table-server
      :headers="producerHeaders"
      :items="relateProducers"
      :items-length="totalItems"
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
          <v-icon v-else :icon="mdiAccount" />
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
          <v-icon size="small" :icon="mdiPencil" />
          編集
        </v-btn>
      </template>
    </v-data-table-server>
  </div>
</template>
