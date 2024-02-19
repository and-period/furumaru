<script lang="ts" setup>
import { mdiImport, mdiExport } from '@mdi/js'
import type { VDataTable } from 'vuetify/lib/components/index.mjs'
import { unix } from 'dayjs'

import type { AlertType } from '~/lib/hooks'
import {
  CharacterEncodingType,
  OrderStatus,
  ShippingType,
  ShippingCarrier,
  type Coordinator,
  type ExportOrdersRequest,
  type Order,
  type Promotion,
  type User,
  type OrderFulfillment
} from '~/types/api'

// TODO: API設計が決まり次第型定義の厳格化
interface ImportFormData {
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
    type: Object as PropType<ExportOrdersRequest>,
    default: () => ({
      shippingCarrier: ShippingCarrier.UNKNOWN,
      characterEncodingType: CharacterEncodingType.UTF8
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
    title: '注文No.',
    key: 'managementId',
    sortable: false
  },
  {
    title: '注文者',
    key: 'userId',
    sortable: false
  },
  {
    title: 'ステータス',
    key: 'payment.status',
    sortable: false
  },
  {
    title: '購入日時',
    key: 'payment.orderedAt',
    sortable: false
  },
  {
    title: '購入金額',
    key: 'payment.total',
    sortable: false
  },
  {
    title: '配送方法',
    key: 'shippingTypes',
    sortable: false
  },
  {
    title: '伝票番号',
    key: 'trackingNumbers',
    sortable: false
  }
]
const fulfillmentCompanies = [
  { title: '指定なし', value: ShippingCarrier.UNKNOWN },
  { title: '佐川急便', value: ShippingCarrier.SAGAWA },
  { title: 'ヤマト運輸', value: ShippingCarrier.YAMATO }
]
const characterEncodingTypes = [
  { title: 'UTF-8', value: CharacterEncodingType.UTF8 },
  { title: 'Shift-JIS', value: CharacterEncodingType.ShiftJIS }
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
  get: (): ExportOrdersRequest => props.exportFormData,
  set: (v: ExportOrdersRequest): void => emit('update:export-form-data', v)
})

const getCustomerName = (userId: string): string => {
  const customer = props.customers.find((customer: User): boolean => {
    return customer.id === userId
  })
  return customer ? `${customer.lastname} ${customer.firstname}` : ''
}

const getStatus = (order: Order): string => {
  switch (order.status) {
    case OrderStatus.UNPAID:
      return '支払い待ち'
    case OrderStatus.WAITING:
      return '受注待ち'
    case OrderStatus.PREPARING:
      return '発送準備中'
    case OrderStatus.SHIPPED:
      return '発送完了'
    case OrderStatus.COMPLETED:
      return '完了'
    case OrderStatus.CANCELED:
      return 'キャンセル'
    case OrderStatus.REFUNDED:
      return '返金'
    case OrderStatus.FAILED:
      return '失敗'
    default:
      return '不明'
  }
}

const getStatusColor = (order: Order): string => {
  switch (order.status) {
    case OrderStatus.UNPAID:
      return 'secondary'
    case OrderStatus.WAITING:
      return 'secondary'
    case OrderStatus.PREPARING:
      return 'info'
    case OrderStatus.SHIPPED:
      return 'info'
    case OrderStatus.COMPLETED:
      return 'primary'
    case OrderStatus.CANCELED:
      return 'warning'
    case OrderStatus.REFUNDED:
      return 'warning'
    case OrderStatus.FAILED:
      return 'error'
    default:
      return 'unknown'
  }
}

const getShippingType = (shippingType: ShippingType): string => {
  switch (shippingType) {
    case ShippingType.NORMAL:
      return '常温・冷蔵便'
    case ShippingType.FROZEN:
      return '冷凍便'
    default:
      return '不明'
  }
}

const getShippingTypes = (fulfillments: OrderFulfillment[]): string => {
  let types: ShippingType[] = []
  fulfillments.forEach((fulfillment: OrderFulfillment): void => {
    if (!types.includes(fulfillment.shippingType)) {
      types.push(fulfillment.shippingType)
    }
  })
  types = types.sort((a, b) => a - b)
  const res: string[] = types.map((type: ShippingType): string => {
    return getShippingType(type)
  })
  return res.join('\n')
}

const getTrackingNumbers = (fulfillments: OrderFulfillment[]): string => {
  const numbers: string[] = []
  fulfillments.forEach((fulfillment: OrderFulfillment): void => {
    if (fulfillment.trackingNumber !== '') {
      numbers.push(fulfillment.trackingNumber)
    }
  })
  return numbers.join('\n')
}

const getOrderedAt = (orderedAt: number): string => {
  if (orderedAt === 0) {
    return '-'
  }
  return unix(orderedAt).format('YYYY/MM/DD HH:mm')
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
          v-model="exportFormDataValue.shippingCarrier"
          label="配送会社"
          :items="fulfillmentCompanies"
        />
        <v-select
          v-model="exportFormDataValue.characterEncodingType"
          label="文字エンコード種別"
          :items="characterEncodingTypes"
        />
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
      <!-- <v-btn color="primary" variant="outlined" @click="toggleImportDialog">
        <v-icon start :icon="mdiImport" />
        Import
      </v-btn> -->
      <v-btn class="ml-4" color="secondary" variant="outlined" @click="toggleExportDialog">
        <v-icon start :icon="mdiExport" />
        エクスポート
      </v-btn>
    </v-card-title>

    <v-card-text>
      <v-data-table-server
        :headers="headers"
        :loading="loading"
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
          <v-chip size="small" :color="getStatusColor(item)">
            {{ getStatus(item) }}
          </v-chip>
        </template>
        <template #[`item.payment.orderedAt`]="{ item }">
          {{ getOrderedAt(item.payment.orderedAt) }}
        </template>
        <template #[`item.payment.total`]="{ item }">
          &yen; {{ item.payment.total.toLocaleString() }}
        </template>
        <template #[`item.shippingTypes`]="{ item }">
          {{ getShippingTypes(item.fulfillments) }}
        </template>
        <template #[`item.trackingNumbers`]="{ item }">
          {{ getTrackingNumbers(item.fulfillments) }}
        </template>
      </v-data-table-server>
    </v-card-text>
  </v-card>
</template>
