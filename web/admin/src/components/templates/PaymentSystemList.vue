<script lang="ts" setup>
import { mdiPencil } from '@mdi/js'
import type { VDataTable } from 'vuetify/lib/components/index.mjs'

import { paymentsList } from '~/constants'
import type { PaymentListItem } from '~/constants'
import type { AlertType } from '~/lib/hooks'
import { PaymentSystemStatus } from '~/types/api/v1'
import type { PaymentMethodType, PaymentSystem } from '~/types/api/v1'

const props = defineProps({
  loading: {
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
  systems: {
    type: Array<PaymentSystem>,
    default: () => [],
  },
})

const emit = defineEmits<{
  (e: 'submit', methodType: PaymentMethodType): void
}>()

const headers: VDataTable['headers'] = [
  {
    title: '決済システム',
    key: 'methodType',
    sortable: false,
  },
  {
    title: '状態',
    key: 'status',
    sortable: false,
  },
  {
    title: '',
    key: 'actions',
    sortable: false,
  },
]

const getPaymentSystemName = (methodType: PaymentMethodType): string => {
  const payment = paymentsList.find((paymentItem: PaymentListItem): boolean => {
    return paymentItem.value === methodType
  })
  return payment?.name || ''
}

const getPaymentSystemStatus = (status: PaymentSystemStatus): string => {
  switch (status) {
    case PaymentSystemStatus.PaymentSystemStatusInUse:
      return '利用可能'
    case PaymentSystemStatus.PaymentSystemStatusOutage:
      return '停止中'
    default:
      return '不明'
  }
}

const getUpdateButtonLabel = (status: PaymentSystemStatus): string => {
  switch (status) {
    case PaymentSystemStatus.PaymentSystemStatusInUse:
      return '停止する'
    case PaymentSystemStatus.PaymentSystemStatusOutage:
      return '利用開始'
    default:
      return ''
  }
}

const getUpdateButtonColor = (status: PaymentSystemStatus): string => {
  switch (status) {
    case PaymentSystemStatus.PaymentSystemStatusInUse:
      return 'error'
    case PaymentSystemStatus.PaymentSystemStatusOutage:
      return 'primary'
    default:
      return ''
  }
}

const onSubmit = (methodType: PaymentMethodType): void => {
  emit('submit', methodType)
}
</script>

<template>
  <v-alert
    v-show="isAlert"
    :type="alertType"
    v-text="alertText"
  />

  <v-card
    class="mt-4"
    flat
  >
    <v-card-title>決済システム状態管理</v-card-title>

    <v-card-text>
      <v-data-table-server
        :headers="headers"
        :loading="loading"
        :items="systems"
        :items-length="systems.length"
        no-data-text="登録されている決済システムがありません。"
      >
        <template #[`item.methodType`]="{ item }">
          {{ getPaymentSystemName(item.methodType) }}
        </template>
        <template #[`item.status`]="{ item }">
          {{ getPaymentSystemStatus(item.status) }}
        </template>
        <template #[`item.actions`]="{ item }">
          <v-btn
            variant="outlined"
            :color="getUpdateButtonColor(item.status)"
            @click.stop="onSubmit(item.methodType)"
          >
            <v-icon
              size="small"
              :icon="mdiPencil"
            />
            {{ getUpdateButtonLabel(item.status) }}
          </v-btn>
        </template>
      </v-data-table-server>
    </v-card-text>
  </v-card>
</template>
