<template>
  <div>
    <v-data-table
      :headers="producerHeaders"
      :items="relateProducers"
      :server-items-length="totalItems"
      :footer-props="tableFooterProps"
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
          />
          <v-icon v-else>mdi-account</v-icon>
        </v-avatar>
      </template>
      <template #[`item.name`]="{ item }">
        {{ `${item.lastname} ${item.firstname}` }}
      </template>
      <template #[`item.phoneNumber`]="{ item }">
        {{ convertPhone(item.phoneNumber) }}
      </template>
      <template #[`item.actions`]="{ item }">
        <v-btn outlined color="primary" small @click="handleEdit(item)">
          <v-icon small>mdi-pencil</v-icon>
          編集
        </v-btn>
      </template>
    </v-data-table>
  </div>
</template>

<script lang="ts">
import { useRouter } from '@nuxtjs/composition-api'
import { computed, defineComponent } from '@vue/composition-api'
import { DataTableHeader } from 'vuetify'

import { useCoordinatorStore } from '~/store/coordinator'
import { ProducersResponseProducersInner } from '~/types/api'

export default defineComponent({
  props: {
    tableFooterProps: {
      type: Object,
      default: () => {},
    },
  },
  setup(_, { emit }) {
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
        value: 'thumbnailUrl',
      },
      {
        text: '生産者名',
        value: 'name',
      },
      {
        text: '店舗名',
        value: 'storeName',
      },
      {
        text: 'Email',
        value: 'email',
      },
      {
        text: '電話番号',
        value: 'phoneNumber',
      },
      {
        text: 'Actions',
        value: 'actions',
        sortable: false,
      },
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

    return {
      relateProducers,
      totalItems,
      producerHeaders,
      handleUpdateItemsPerPage,
      handleUpdatePage,
      convertPhone,
      handleEdit,
    }
  },
})
</script>
