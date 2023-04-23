<script lang="ts" setup>
import { mdiImport, mdiExport, mdiPencil } from '@mdi/js'
import * as dayjs from 'dayjs'
import { storeToRefs } from 'pinia'
import { VDataTable } from 'vuetify/lib/labs/components'

import { usePagination } from '~/lib/hooks'
import { useOrderStore } from '~/store'
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
  offset
} = usePagination()

const fetchState = useAsyncData(() => {
  return orderStore.fetchOrders(itemsPerPage.value, offset.value)
})

const headers: VDataTable['headers'] = [
  {
    title: '注文者',
    key: 'userName'
  },
  {
    title: '配送ステータス',
    key: 'payment.status'
  },
  {
    title: '購入日時',
    key: 'orderedAt'
  },
  {
    title: '配送方法',
    key: 'fulfillment.shippingMethod'
  },
  {
    title: '購入金額',
    key: 'payment.total'
  },
  {
    title: '伝票番号',
    key: 'payment.paymentId'
  },
  {
    title: '注文ID',
    key: 'id'
  }
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
  { deliveryCompany: 'ヤマト運輸', value: 'ヤマト運輸' }
]

const handleRowClick = (item: OrderResponse): void => {
  router.push(`/orders/${item.id}`)
}

const handleImport = () => {
  // TODO: APIの実装が完了後に対応
}

const handleExport = () => {
  // TODO: APIの実装が完了後に対応
}

try {
  await fetchState.execute()
} catch (err) {
  console.log('failed to setup', err)
}
</script>

<template>
  <div>
    <v-card-title>
      注文
      <v-spacer />
      <v-btn variant="outlined" color="primary" @click="toggleImportDialog">
        <v-icon start :icon="mdiImport" />
        Import
      </v-btn>
      <v-btn
        outlined
        class="ml-4"
        color="secondary"
        @click="toggleExportDialog"
      >
        <v-icon start :icon="mdiExport" />
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
          item-title="deliveryCompany"
          item-value="value"
        />
        <v-file-input class="mr-2" label="CSVを選択" />
        <v-divider />
        <v-card-actions>
          <v-spacer />
          <v-btn color="error" variant="text" @click="toggleImportDialog">
            キャンセル
          </v-btn>
          <v-btn color="primary" variant="outlined" @click="handleImport">
            登録
          </v-btn>
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
          item-title="deliveryCompany"
          item-value="value"
        />
        <v-file-input class="mr-2" label="CSVを選択" />
        <v-divider />
        <v-card-actions>
          <v-spacer />
          <v-btn color="error" variant="text" @click="toggleExportDialog">
            キャンセル
          </v-btn>
          <v-btn color="primary" variant="outlined" @click="handleExport">
            登録
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
    <v-card class="mt-4" flat>
      <v-card-text>
        <v-data-table-server
          :headers="headers"
          :items="orders"
          :items-length="totalItems"
          :footer-props="options"
          no-data-text="表示する注文がありません"
          hover
          @update:items-per-page="handleUpdateItemsPerPage"
          @update:page="handleUpdatePage"
          @click:row="(_: any, { item }: any) => handleRowClick(item.raw)"
        >
          <template #[`item.payment.status`]="{ item }">
            <v-chip size="small" :color="getStatusColor(item.raw.payment.status)">
              {{ getStatus(item.raw.payment.status) }}
            </v-chip>
          </template>
          <template #[`item.fulfillment.shippingMethod`]="{ item }">
            {{ getShippingMethod(item.raw.fulfillment.shippingMethod) }}
          </template>
          <template #[`item.orderedAt`]="{ item }">
            {{ getDay(item.raw.orderedAt) }}
          </template>
        </v-data-table-server>
      </v-card-text>
    </v-card>
  </div>
</template>
