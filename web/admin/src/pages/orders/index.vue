<template>
  <div>
    <v-card-title>注文</v-card-title>

    <v-card>
      <v-card-text>
        <v-data-table
          show-select
          :headers="headers"
          :items="orders"
          :server-items-length="totalItems"
          :footer-props="options"
          no-data-text="表示する注文がありません"
          @update:items-per-page="handleUpdateItemsPerPage"
          @update:page="handleUpdatePage"
        >
          <template #[`item.status`]="{ item }">
            <v-chip samll :color="getSatusColor(item.status)">
              {{ getStatus(item.status) }}
            </v-chip>
          </template>
        </v-data-table>
      </v-card-text>
    </v-card>
  </div>
</template>

<script lang="ts">
import { defineComponent, useFetch } from '@nuxtjs/composition-api'
import { storeToRefs } from 'pinia'
import { DataTableHeader } from 'vuetify'

import { usePagination } from '~/lib/hooks'
import { useOrderStore } from '~/store/orders'

export default defineComponent({
  setup() {
    const orderStore = useOrderStore()
    const { orders, totalItems } = storeToRefs(orderStore)

    const {
      updateCurrentPage,
      itemsPerPage,
      handleUpdateItemsPerPage,
      options,
      offset,
    } = usePagination()

    const { fetchState } = useFetch(() => {
      orderStore.fetchOrders(itemsPerPage.value, offset.value)
    })

    const headers: DataTableHeader[] = [
      {
        text: '注文ID',
        value: 'id',
      },
      {
        text: '注文者',
        value: 'userName',
      },
      {
        text: '配送ステータス',
        value: 'status',
      },
      {
        text: '購入日時',
        value: 'orderedAt',
      },
      {
        text: '配送方法',
        value: 'payment.paymentType',
      },
      {
        text: '購入金額',
        value: 'payment.total',
      },
      {
        text: '伝票番号',
        value: 'payment.paymentId',
      },
      {
        text: 'Actions',
        value: 'actions',
        sortable: false,
      }
    ]

    const handleUpdatePage = async (page: number) => {
      updateCurrentPage(page)
      await orderStore.fetchOrders(itemsPerPage.value, offset.value)
    }

    const getStatusColor = (status: number): string => {
      switch (status) {
        case 1:
          return 'accent'
        case 2:
          return ''
        case 3:
          return ''
        case 4:
          return ''
        case 5:
          return ''
        case 6:
          return ''
        default:
          return 'accentDarken'
      }
    }

    return {
      headers,
      orders,
      totalItems,
      fetchState,
      options,
      handleUpdateItemsPerPage,
      handleUpdatePage,
      getStatusColor,
    }
  },
})
</script>
