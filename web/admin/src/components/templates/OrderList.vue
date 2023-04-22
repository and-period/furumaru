<script lang="ts" setup>
import { mdiImport, mdiExport, mdiPencil } from '@mdi/js'
import { VDataTable } from 'vuetify/lib/labs/components'
import { unix } from 'dayjs'

import { AlertType } from '~/lib/hooks'
import { DeliveryType, OrdersResponse, PaymentStatus } from '~/types/api'
import { Order } from '~/types/props'

const props = defineProps({
  loading: {
    type: Boolean,
    default: false
  },
  dialog: {
    type: Object,
    default: () => ({
      import: false,
      export: false
    })
  },
  isAlert: {
    type: Boolean,
    default: false
  },
  alertType: {
    type: String as PropType<AlertType>,
    default: undefined
  },
  alertText: {
    type: String,
    default: ''
  },
  orders: {
    type: Array<OrdersResponse['orders']>,
    default: () => []
  },
  tableItemsPerPage: {
    type: Number,
    default: 20
  },
  tableItemsTotal: {
    type: Number,
    default: 0,
  },
  importFormData: {
    type: Object, // TODO: API設計が決まり次第型定義の厳格化
    default: () => ({
      company: false
    })
  },
  exportFormData: {
    type: Object, // TODO: API設計が決まり次第型定義の厳格化
    default: () => ({
      company: false
    })
  }
})

const emit = defineEmits<{
  (e: 'click:edit', orderId: string): void
  (e: 'click:update-page', page: number): void
  (e: 'click:update-items-per-page', page: number): void
  (e: 'submit:import'): void
  (e: 'submit:export'): void
}>()

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
    title: 'Actions',
    key: 'actions',
    sortable: false
  },
  {
    title: '注文ID',
    key: 'id'
  }
]
const fulfillmentCompanies: Order[] = [
  { name: '佐川急便', value: '佐川急便' },
  { name: 'ヤマト運輸', value: 'ヤマト運輸' }
]

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

const getOrderdAt = (orderdAt: number): string => {
  return unix(orderdAt).format('YYYY/MM/DD HH:mm')
}

const toggleImportDialog = (): void => {
  props.dialog.import = !props.dialog.import
}

const toggleExportDialog = (): void => {
  props.dialog.export = !props.dialog.export
}

const onClickUpdatePage = (page: number): void => {
  emit('click:update-page', page)
}

const onClickUpdateItemsPerPage = (page: number): void => {
  emit('click:update-items-per-page', page)
}

const onEdit = (orderId: string): void => {
  emit('click:edit', orderId)
}

const onSubmitImport = (): void => {
  emit('submit:import')
}

const onSubmitExport = (): void => {
  emit('submit:export')
}
</script>

<template>
  <v-alert v-show="props.isAlert" :type="props.alertType" v-text="props.alertText" />
  <v-dialog v-model="props.dialog.import" width="500">
    <v-card>
      <v-card-title class="text-h6 primaryLight">
        ファイルの取り込み
      </v-card-title>

      <v-select
        v-model="props.importFormData.company"
        label="配送会社"
        class="mr-2 ml-2"
        :items="fulfillmentCompanies"
        item-title="name"
        item-value="value"
      />
      <v-file-input class="mr-2" label="CSVを選択" />
      <v-divider />
      <v-card-actions>
        <v-spacer />
        <v-btn color="error" variant="text" @click="toggleImportDialog">
          キャンセル
        </v-btn>
        <v-btn color="primary" variant="outlined" @click="onSubmitImport">
          登録
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
  <v-dialog v-model="props.dialog.export" width="500">
    <v-card>
      <v-card-title class="text-h6 primaryLight">
        ファイルの出力
      </v-card-title>
      <v-divider />

      <v-select
        v-model="props.exportFormData.company"
        label="配送会社"
        class="mr-2 ml-2"
        :items="fulfillmentCompanies"
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
        <v-btn color="primary" variant="outlined" @click="onSubmitExport">
          登録
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
  <v-card class="mt-4" flat>
    <v-card-title class="d-flex flex-row">
      注文
      <v-spacer />
      <v-btn color="primary" variant="outlined" @click="toggleImportDialog">
        <v-icon start :icon="mdiImport" />
        Import
      </v-btn>
      <v-btn class="ml-4" color="secondary" variant="outlined" @click="toggleExportDialog">
        <v-icon start :icon="mdiExport" />
        Export
      </v-btn>
    </v-card-title>

    <v-card-text>
      <v-data-table-server
        :headers="headers"
        :items="props.orders"
        :items-per-page="props.tableItemsPerPage"
        :items-length="props.tableItemsTotal"
        no-data-text="表示する注文がありません"
        @update:page="onClickUpdatePage"
        @update:items-per-page="onClickUpdateItemsPerPage"
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
          {{ getOrderdAt(item.raw.orderedAt) }}
        </template>
        <template #[`item.actions`]="{ item }">
          <v-btn variant="outlined" color="primary" size="small" @click="onEdit(item.raw.id)">
            <v-icon size="small" :icon="mdiPencil" />
            詳細
          </v-btn>
        </template>
      </v-data-table-server>
    </v-card-text>
  </v-card>
</template>
