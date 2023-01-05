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
          <template #[`item.payment.status`]="{ item }">
            <v-chip small :color="getStatusColor(item.payment.status)">
              {{ getStatus(item.payment.status) }}
            </v-chip>
          </template>
          <template #[`item.fulfillment.shippingMethod`]="{ item }">
            {{ getShippingMethod(item.fulfillment.shippingMethod) }}
          </template>
          <template #[`item.actions`]="{ item }">
            <v-btn outlined color="primary" small @click="handleEdit(item)">
              <v-icon small>mdi-pencil</v-icon>
              詳細
            </v-btn>
        </template>
        </v-data-table>
      </v-card-text>
    </v-card>
  </div>
</template>

<script lang="ts">
import { defineComponent, useFetch, useRouter } from '@nuxtjs/composition-api'
import { storeToRefs } from 'pinia'
import { DataTableHeader } from 'vuetify'

import { usePagination } from '~/lib/hooks'
import { useOrderStore } from '~/store/orders'
import { DeliveryType, OrderResponse, PaymentStatus } from '~/types/api'

export default defineComponent({
  setup() {
    const router = useRouter()
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
        value: 'payment.status',
      },
      {
        text: '購入日時',
        value: 'orderedAt',
      },
      {
        text: '配送方法',
        value: 'fulfillment.shippingMethod',
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

    const getStatusColor = (status: PaymentStatus): string => {
      switch (status) {
        case PaymentStatus.UNKNOWN:
          return 'unkown'
        case PaymentStatus.UNPAID:
          return 'secondary'
        case PaymentStatus.PENDING:
          return 'secondary'
        case PaymentStatus.AUTHORIZED:
          return 'info'
        case PaymentStatus.PAID:
          return 'primary'
        case PaymentStatus.REFUNDED:
          return 'primary'
        case PaymentStatus.EXPIRED:
          return 'error'
        default:
          return 'unkown'
      }
    }

    const getStatus = (status: PaymentStatus): string => {
      switch (status) {
        case PaymentStatus.UNKNOWN:
          return '不明'
        case PaymentStatus.UNPAID:
          return '未払い'
        case PaymentStatus.PENDING:
          return '保留中'
        case PaymentStatus.AUTHORIZED:
          return 'オーソリ済み'
        case PaymentStatus.PAID:
          return '支払い済み'
        case PaymentStatus.REFUNDED:
          return '返金済み'
        case PaymentStatus.EXPIRED:
          return '期限切れ'
        default:
          return '不明'
      }
    }

    const getShippingMethod = (shippingMethod: DeliveryType): string => {
      switch (shippingMethod) {
        case DeliveryType.UNKNOWN:
          return '不明'
        case DeliveryType.NORMAL:
          return '通常便'
        case DeliveryType.REFRIGERATED:
          return '冷蔵便'
        case DeliveryType.FROZEN:
          return '冷凍便'
      }
    }

    const handleEdit = (item: OrderResponse) => {
      router.push(`/orders/${item.id}`)
    }

    return {
      headers,
      orders,
      totalItems,
      fetchState,
      options,
      handleUpdateItemsPerPage,
      handleUpdatePage,
      handleEdit,
      getStatusColor,
      getStatus,
      getShippingMethod,
    }
  },
})
</script>
