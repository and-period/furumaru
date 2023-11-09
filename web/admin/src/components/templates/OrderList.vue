<script lang="ts" setup>
import { mdiImport, mdiExport } from '@mdi/js'
import { VDataTable } from 'vuetify/lib/labs/components.mjs'
import { unix } from 'dayjs'

import type { AlertType } from '~/lib/hooks'
import {
  DeliveryType,
  FulfillmentStatus,
  PaymentStatus,
  type Coordinator,
  type Order,
  type Promotion,
  type User
} from '~/types/api'

// TODO: API設計が決まり次第型定義の厳格化
interface ImportFormData {
  company: boolean
}
interface ExportFormData {
  company: boolean
}

const props = defineProps({
  loading: {
    type: Boolean,
    default: false
  },
  importDialog: {
    type: Boolean,
    default: false
  },
  exportDialog: {
    type: Boolean,
    default: false
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
    type: Array<Order>,
    default: () => []
  },
  coordinators: {
    type: Array<Coordinator>,
    default: () => []
  },
  customers: {
    type: Array<User>,
    default: () => []
  },
  promotions: {
    type: Array<Promotion>,
    default: () => []
  },
  tableItemsPerPage: {
    type: Number,
    default: 20
  },
  tableItemsTotal: {
    type: Number,
    default: 0
  },
  importFormData: {
    type: Object,
    default: () => ({
      company: false
    })
  },
  exportFormData: {
    type: Object,
    default: () => ({
      company: false
    })
  }
})

const emit = defineEmits<{
  (e: 'click:row', orderId: string): void
  (e: 'click:edit', orderId: string): void
  (e: 'click:update-page', page: number): void
  (e: 'click:update-items-per-page', page: number): void
  (e: 'update:import-dialog', dialog: boolean): void
  (e: 'update:export-dialog', dialog: boolean): void
  (e: 'update:import-form-data', formData: Object): void
  (e: 'update:export-form-data', formData: Object): void
  (e: 'submit:import'): void
  (e: 'submit:export'): void
}>()

const headers: VDataTable['headers'] = [
  {
    title: '注文者',
    key: 'userId'
  },
  {
    title: '支払いステータス',
    key: 'payment.status'
  },
  {
    title: '配送ステータス',
    key: 'fulfillment.status'
  },
  {
    title: '購入金額',
    key: 'payment.total'
  },
  {
    title: '配送方法',
    key: 'fulfillment.shippingMethod'
  },
  {
    title: '伝票番号',
    key: 'fulfillment.trackingNumber'
  },
  {
    title: '購入日時',
    key: 'createdAt'
  }
]
const fulfillmentCompanies = [
  { title: '佐川急便', value: '佐川急便' },
  { title: 'ヤマト運輸', value: 'ヤマト運輸' }
]

const importDialogValue = computed({
  get: (): boolean => props.importDialog,
  set: (v: boolean): void => emit('update:import-dialog', v)
})
const exportDialogValue = computed({
  get: (): boolean => props.exportDialog,
  set: (v: boolean): void => emit('update:export-dialog', v)
})
const importFormDataValue = computed({
  get: (): ImportFormData => props.importFormData as ImportFormData,
  set: (v: ImportFormData): void => emit('update:import-form-data', v)
})
const exportFormDataValue = computed({
  get: (): ExportFormData => props.importFormData as ExportFormData,
  set: (v: ExportFormData): void => emit('update:export-form-data', v)
})

const getCustomerName = (userId: string): string => {
  const customer = props.customers.find((customer: User): boolean => {
    return customer.id === userId
  })
  return customer ? `${customer.lastname} ${customer.firstname}` : ''
}

const getPaymentStatus = (status: PaymentStatus): string => {
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

const getPaymentStatusColor = (status: PaymentStatus): string => {
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

const getFulfillmentStatus = (status: FulfillmentStatus): string => {
  switch (status) {
    case FulfillmentStatus.FULFILLED:
      return '発送済み'
    case FulfillmentStatus.UNFULFILLED:
      return '未発送'
    default:
      return '不明'
  }
}

const getFulfillmentStatusColor = (status: FulfillmentStatus): string => {
  switch (status) {
    case FulfillmentStatus.FULFILLED:
      return 'primary'
    case FulfillmentStatus.UNFULFILLED:
      return 'secondary'
    default:
      return 'unknown'
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

const getCreatedAt = (createdAt: number): string => {
  return unix(createdAt).format('YYYY/MM/DD HH:mm')
}

const toggleImportDialog = (): void => {
  importDialogValue.value = !importDialogValue.value
}

const toggleExportDialog = (): void => {
  exportDialogValue.value = !exportDialogValue.value
}

const onClickUpdatePage = (page: number): void => {
  emit('click:update-page', page)
}

const onClickUpdateItemsPerPage = (page: number): void => {
  emit('click:update-items-per-page', page)
}

const onClickRow = (orderId: string): void => {
  emit('click:row', orderId)
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

  <v-dialog v-model="importDialogValue" width="500">
    <v-card>
      <v-card-title class="text-h6 primaryLight">
        ファイルの取り込み
      </v-card-title>

      <v-card-text>
        <v-select
          v-model="importFormDataValue.company"
          label="配送会社"
          class="mr-2 ml-2"
          :items="fulfillmentCompanies"
          item-title="title"
          item-value="value"
        />
        <v-file-input class="mr-2" label="CSVを選択" />
      </v-card-text>
      <v-card-actions>
        <v-spacer />
        <v-btn color="error" variant="text" @click="toggleImportDialog">
          キャンセル
        </v-btn>
        <v-btn color="primary" variant="outlined" :loading="loading" @click="onSubmitImport">
          登録
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>

  <v-dialog v-model="exportDialogValue" width="500">
    <v-card>
      <v-card-title class="text-h6 primaryLight">
        ファイルの出力
      </v-card-title>

      <v-card-text>
        <v-select
          v-model="exportFormDataValue.company"
          label="配送会社"
          class="mr-2 ml-2"
          :items="fulfillmentCompanies"
          item-title="title"
          item-value="value"
        />
        <v-file-input class="mr-2" label="CSVを選択" />
      </v-card-text>
      <v-card-actions>
        <v-spacer />
        <v-btn color="error" variant="text" @click="toggleExportDialog">
          キャンセル
        </v-btn>
        <v-btn color="primary" variant="outlined" :loading="loading" @click="onSubmitExport">
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
        hover
        no-data-text="表示する注文がありません"
        @update:page="onClickUpdatePage"
        @update:items-per-page="onClickUpdateItemsPerPage"
        @click:row="(_: any, {item}: any) => onClickRow(item.id)"
      >
        <template #[`item.userId`]="{ item }">
          {{ getCustomerName(item.userId) }}
        </template>
        <template #[`item.payment.status`]="{ item }">
          <v-chip size="small" :color="getPaymentStatusColor(item.payment.status)">
            {{ getPaymentStatus(item.payment.status) }}
          </v-chip>
        </template>
        <template #[`item.fulfillment.status`]="{ item }">
          <v-chip size="small" :color="getFulfillmentStatusColor(item.fulfillment.status)">
            {{ getFulfillmentStatus(item.fulfillment.status) }}
          </v-chip>
        </template>
        <template #[`item.fulfillment.shippingMethod`]="{ item }">
          {{ getShippingMethod(item.fulfillment.shippingMethod) }}
        </template>
        <template #[`item.createdAt`]="{ item }">
          {{ getCreatedAt(item.createdAt) }}
        </template>
      </v-data-table-server>
    </v-card-text>
  </v-card>
</template>
