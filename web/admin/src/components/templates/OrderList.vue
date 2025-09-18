<script lang="ts" setup>
import { mdiImport, mdiExport, mdiContentCopy, mdiFileDocumentOutline } from '@mdi/js'
import type { VDataTable } from 'vuetify/lib/components/index.mjs'
import { unix } from 'dayjs'

import type { AlertType } from '~/lib/hooks'
import { CharacterEncodingType } from '~/types'
import { fulfillmentCompanies, characterEncodingTypes, orderStatuses, orderShippingTypes, orderTypes } from '~/constants'
import {
  OrderStatus,
  ShippingType,
  ShippingCarrier,
} from '~/types/api/v1'
import type { Coordinator, ExportOrdersRequest, Order, Promotion, User, OrderFulfillment, OrderType } from '~/types/api/v1'

// TODO: API設計が決まり次第型定義の厳格化
interface ImportFormData {
  company: boolean
}

const props = defineProps({
  selectedItemId: {
    type: String,
    default: '',
  },
  loading: {
    type: Boolean,
    default: false,
  },
  importDialog: {
    type: Boolean,
    default: false,
  },
  exportDialog: {
    type: Boolean,
    default: false,
  },
  isAlert: {
    type: Boolean,
    default: false,
  },
  alertType: {
    type: String as PropType<AlertType>,
    default: undefined,
  },
  alertText: {
    type: String,
    default: '',
  },
  tableSortBy: {
    type: Array as PropType<VDataTable['sortBy']>,
    default: () => [],
  },
  orders: {
    type: Array<Order>,
    default: () => [],
  },
  coordinators: {
    type: Array<Coordinator>,
    default: () => [],
  },
  customers: {
    type: Array<User>,
    default: () => [],
  },
  promotions: {
    type: Array<Promotion>,
    default: () => [],
  },
  tableItemsPerPage: {
    type: Number,
    default: 20,
  },
  tableItemsTotal: {
    type: Number,
    default: 0,
  },
  importFormData: {
    type: Object,
    default: () => ({
      company: false,
    }),
  },
  exportFormData: {
    type: Object as PropType<ExportOrdersRequest>,
    default: () => ({
      shippingCarrier: ShippingCarrier.ShippingCarrierUnknown,
      characterEncodingType: CharacterEncodingType.UTF8,
    }),
  },
})

const emit = defineEmits<{
  (e: 'click:row', orderId: string): void
  (e: 'click:update-page', page: number): void
  (e: 'click:update-items-per-page', page: number): void
  (e: 'update:import-dialog', dialog: boolean): void
  (e: 'update:export-dialog', dialog: boolean): void
  (e: 'update:import-form-data', formData: object): void
  (e: 'update:export-form-data', formData: ExportOrdersRequest): void
  (e: 'update:sort-by'): void
  (e: 'submit:import'): void
  (e: 'submit:export'): void
}>()

const headers: VDataTable['headers'] = [
  {
    title: '注文No.',
    key: 'managementId',
    sortable: false,
  },
  {
    title: '注文種別',
    key: 'orderType',
    sortable: false,
  },
  {
    title: '注文者',
    key: 'userId',
    sortable: false,
  },
  {
    title: 'ステータス',
    key: 'payment.status',
    sortable: false,
  },
  {
    title: '購入日時',
    key: 'payment.orderedAt',
    sortable: false,
  },
  {
    title: '購入金額',
    key: 'payment.total',
    sortable: false,
  },
  {
    title: '配送方法',
    key: 'shippingTypes',
    sortable: false,
  },
  {
    title: '伝票番号',
    key: 'trackingNumbers',
    sortable: false,
  },
]

const importDialogValue = computed({
  get: (): boolean => props.importDialog,
  set: (v: boolean): void => emit('update:import-dialog', v),
})
const exportDialogValue = computed({
  get: (): boolean => props.exportDialog,
  set: (v: boolean): void => emit('update:export-dialog', v),
})
const importFormDataValue = computed({
  get: (): ImportFormData => props.importFormData as ImportFormData,
  set: (v: ImportFormData): void => emit('update:import-form-data', v),
})
const exportFormDataValue = computed({
  get: (): ExportOrdersRequest => props.exportFormData,
  set: (v: ExportOrdersRequest): void => emit('update:export-form-data', v),
})

const getCustomerName = (userId: string): string => {
  const customer = props.customers.find((customer: User): boolean => {
    return customer.id === userId
  })
  return customer ? `${customer.lastname} ${customer.firstname}` : ''
}

const getStatus = (order: Order): string => {
  const value = orderStatuses.find(status => status.value === order.status)
  return value ? value.title : '不明'
}

const getStatusColor = (order: Order): string => {
  switch (order.status) {
    case OrderStatus.OrderStatusUnpaid:
      return 'secondary'
    case OrderStatus.OrderStatusWaiting:
      return 'secondary'
    case OrderStatus.OrderStatusPreparing:
      return 'info'
    case OrderStatus.OrderStatusShipped:
      return 'info'
    case OrderStatus.OrderStatusCompleted:
      return 'primary'
    case OrderStatus.OrderStatusCanceled:
      return 'warning'
    case OrderStatus.OrderStatusRefunded:
      return 'warning'
    case OrderStatus.OrderStatusFailed:
      return 'error'
    default:
      return 'unknown'
  }
}

const getOrderType = (orderType: OrderType): string => {
  const value = orderTypes.find(type => type.value === orderType)
  return value ? value.title : '不明'
}

const getShippingType = (shippingType: ShippingType): string => {
  const value = orderShippingTypes.find(type => type.value === shippingType)
  return value ? value.title : '不明'
}

const getShippingTypes = (fulfillments: OrderFulfillment[]): string => {
  if (fulfillments.length === 0) {
    return '-'
  }
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
  if (fulfillments.length === 0) {
    return '-'
  }
  if (fulfillments[0]?.shippingType === ShippingType.ShippingTypePickup) {
    return '-'
  }
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
  <v-alert
    v-show="props.isAlert"
    :type="props.alertType"
    v-text="props.alertText"
  />

  <v-dialog
    v-model="importDialogValue"
    width="500"
  >
    <v-card>
      <v-card-title class="text-h6 py-4">
        ファイルの取り込み
      </v-card-title>

      <v-card-text class="pb-4">
        <v-select
          v-model="importFormDataValue.company"
          label="配送会社"
          class="mb-4"
          :items="fulfillmentCompanies"
          item-title="title"
          item-value="value"
        />
        <v-file-input
          label="CSVを選択"
        />
      </v-card-text>
      <v-card-actions class="px-6 pb-4">
        <v-spacer />
        <v-btn
          color="medium-emphasis"
          variant="text"
          @click="toggleImportDialog"
        >
          キャンセル
        </v-btn>
        <v-btn
          color="primary"
          variant="elevated"
          :loading="loading"
          @click="onSubmitImport"
        >
          取り込み
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>

  <v-dialog
    v-model="exportDialogValue"
    width="500"
  >
    <v-card>
      <v-card-title class="text-h6 py-4">
        ファイルの出力
      </v-card-title>

      <v-card-text class="pb-4">
        <v-select
          v-model="exportFormDataValue.shippingCarrier"
          label="配送会社"
          class="mb-4"
          :items="fulfillmentCompanies"
        />
        <v-select
          v-model="exportFormDataValue.characterEncodingType"
          label="文字エンコード種別"
          :items="characterEncodingTypes"
        />
      </v-card-text>
      <v-card-actions class="px-6 pb-4">
        <v-spacer />
        <v-btn
          color="medium-emphasis"
          variant="text"
          @click="toggleExportDialog"
        >
          キャンセル
        </v-btn>
        <v-btn
          color="primary"
          variant="elevated"
          :loading="loading"
          @click="onSubmitExport"
        >
          出力
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>

  <v-card
    class="mt-6"
    elevation="0"
    rounded="lg"
  >
    <v-card-title class="d-flex flex-column flex-sm-row align-start align-sm-center justify-space-between pa-4 pa-sm-6 pb-4">
      <div class="d-flex align-center mb-3 mb-sm-0">
        <v-icon
          :icon="mdiFileDocumentOutline"
          size="28"
          class="mr-3 text-primary"
        />
        <div>
          <h1 class="text-h6 text-sm-h5 font-weight-bold text-primary">
            注文管理
          </h1>
          <p class="text-caption text-sm-body-2 text-medium-emphasis ma-0">
            注文の確認・管理・エクスポートを行います
          </p>
        </div>
      </div>
      <div class="d-flex ga-3 w-100 w-sm-auto">
        <v-btn
          variant="outlined"
          color="info"
          :size="$vuetify.display.smAndDown ? 'default' : 'large'"
          class="w-100 w-sm-auto"
          @click="toggleExportDialog"
        >
          <v-icon
            start
            :icon="mdiExport"
          />
          エクスポート
        </v-btn>
      </div>
    </v-card-title>

    <v-card-text>
      <v-data-table-server
        :headers="headers"
        :loading="loading"
        :items="props.orders"
        :items-per-page="props.tableItemsPerPage"
        :items-length="props.tableItemsTotal"
        :sort-by="props.tableSortBy"
        hover
        select-strategy="single"
        no-data-text="登録されている注文がありません。"
        @update:page="onClickUpdatePage"
        @update:items-per-page="onClickUpdateItemsPerPage"
        @update:sort-by="emit('update:sort-by')"
        @click:row="(_: any, { item }: any) => onClickRow(item.id)"
      >
        <template #[`item.userId`]="{ item }">
          {{ getCustomerName(item.userId) }}
        </template>
        <template #[`item.payment.status`]="{ item }">
          <v-chip :color="getStatusColor(item)">
            {{ getStatus(item) }}
          </v-chip>
        </template>
        <template #[`item.orderType`]="{ item }">
          {{ getOrderType(item.type) }}
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
