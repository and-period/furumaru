<script lang="ts" setup>
import dayjs from 'dayjs'
import { storeToRefs } from 'pinia'
import { DataTableHeader } from 'vuetify'

import { usePagination } from '~/lib/hooks'
import { useOrderStore } from '~/store/orders'
import { DeliveryType, OrderResponse, PaymentStatus } from '~/types/api'

const router = useRouter()
const orderStore = useOrderStore()
const { orders, totalItems } = storeToRefs(orderStore)

const importCompany = ref<string>('')
const exportCompany = ref<string>('')
const importDialog = ref<boolean>(false)
const exportDialog = ref<boolean>(false)

const {
  updateCurrentPage,
  itemsPerPage,
  handleUpdateItemsPerPage,
  options,
  offset,
} = usePagination()

const _ = useAsyncData(() => {
  return orderStore.fetchOrders(itemsPerPage.value, offset.value)
})

const headers: DataTableHeader[] = [
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
  },
  {
    text: '注文ID',
    value: 'id',
  },
]

const handleUpdatePage = async (page: number) => {
  updateCurrentPage(page)
  await orderStore.fetchOrders(itemsPerPage.value, offset.value)
}

const getStatusColor = (status: PaymentStatus): string => {
  switch (status) {
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
    case DeliveryType.NORMAL:
      return '通常便'
    case DeliveryType.REFRIGERATED:
      return '冷蔵便'
    case DeliveryType.FROZEN:
      return '冷凍便'
    default:
      return '不明'
  }
}

const getDay = (unixTime: number): string => {
  return dayjs.unix(unixTime).format('YYYY/MM/DD HH:mm')
}

const toggleImportDialog = (): void => {
  importDialog.value = !importDialog.value
}

const toggleExportDialog = (): void => {
  exportDialog.value = !exportDialog.value
}

const deliveryCompanyList = [
  { deliveryCompany: '佐川急便', value: '佐川急便' },
  { deliveryCompany: 'ヤマト運輸', value: 'ヤマト運輸' },
]

const handleEdit = (item: OrderResponse) => {
  router.push(`/orders/${item.id}`)
}

const handleImport = () => {
  // TODO: APIの実装が完了後に対応
}

const handleExport = () => {
  // TODO: APIの実装が完了後に対応
}
</script>

<template>
  <div>
    <v-card-title>
      注文
      <v-spacer />
      <v-btn outlined color="primary" @click="toggleImportDialog">
        <v-icon left>mdi-import</v-icon>
        Import
      </v-btn>
      <v-btn
        outlined
        class="ml-4"
        color="secondary"
        @click="toggleExportDialog"
      >
        <v-icon left>mdi-export</v-icon>
        Export
      </v-btn>
    </v-card-title>

    <v-dialog v-model="importDialog" width="500">
      <v-card>
        <v-card-title class="text-h6 primaryLight">
          ファイルの取り込み
        </v-card-title>

        <v-select
          v-model="importCompany"
          label="配送会社"
          class="mr-2 ml-2"
          :items="deliveryCompanyList"
          item-text="deliveryCompany"
          item-value="value"
        />
        <v-file-input class="mr-2" label="CSVを選択" />
        <v-divider />
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn color="error" text @click="toggleImportDialog">
            キャンセル
          </v-btn>
          <v-btn color="primary" outlined @click="handleImport"> 登録 </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
    <v-dialog v-model="exportDialog" width="500">
      <v-card>
        <v-card-title class="text-h6 primaryLight">
          ファイルの出力
        </v-card-title>
        <v-divider />

        <v-select
          v-model="exportCompany"
          label="配送会社"
          class="mr-2 ml-2"
          :items="deliveryCompanyList"
          item-text="deliveryCompany"
          item-value="value"
        />
        <v-file-input class="mr-2" label="CSVを選択" />
        <v-divider />
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn color="error" text @click="toggleExportDialog">
            キャンセル
          </v-btn>
          <v-btn color="primary" outlined @click="handleExport"> 登録 </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
    <v-card class="mt-4" flat>
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
          <template #[`item.orderedAt`]="{ item }">
            {{ getDay(item.orderedAt) }}
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
