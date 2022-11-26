<template>
  <div>
    <v-card-title>注文</v-card-title>

    <v-data-table
      show-select
      :headers="headers"
      :items="orders"
      :server-items-length="totalItems"
      :footer-props="options"
      no-data-text="表示する注文がありません"
      @update:items-per-page="handleUpdateItemsPerPage"
      @update:page="handleUpdatePage"
    />
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
        text: 'ID',
        value: 'id',
      },
      {
        text: '注文者',
        value: 'userName',
      },
      {
        text: '配送ステータス',
        value: '',
      },
      {
        text: '購入日時',
        value: 'orderedAt',
      },
      {
        text: '配送方法',
        value: '',
      },
      {
        text: '購入金額',
        value: 'payment.total',
      },
      {
        text: '伝票番号',
        value: 'payment.paymentId',
      },
    ]

    const handleUpdatePage = async (page: number) => {
      updateCurrentPage(page)
      await orderStore.fetchOrders(itemsPerPage.value, offset.value)
    }

    return {
      headers,
      orders,
      totalItems,
      fetchState,
      options,
      handleUpdateItemsPerPage,
      handleUpdatePage,
    }
  },
})
</script>
