<script lang="ts" setup>
import dayjs, { unix } from 'dayjs'
import { use } from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'
import {
  TitleComponent,
  GridComponent,
  TooltipComponent,
} from 'echarts/components'
import { LineChart } from 'echarts/charts'
import VChart from 'vue-echarts'

import type { AlertType } from '~/lib/hooks'
import { TopOrderPeriodType } from '~/types'
import { PaymentMethodType } from '~/types/api/v1'
import type { TopOrderSalesTrend, TopOrdersResponse } from '~/types/api/v1'
import type { DateTimeInput } from '~/types/props'

use([
  GridComponent,
  CanvasRenderer,
  LineChart,
  TitleComponent,
  TooltipComponent,
])

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
  startAt: {
    type: Number,
    default: dayjs().add(-7, 'day').unix(),
  },
  endAt: {
    type: Number,
    default: dayjs().unix(),
  },
  periodType: {
    type: String as PropType<TopOrderPeriodType>,
    default: TopOrderPeriodType.DAY,
  },
  orders: {
    type: Object as PropType<TopOrdersResponse>,
    default: (): TopOrdersResponse => ({
      startAt: 0,
      endAt: 0,
      periodType: TopOrderPeriodType.DAY,
      orders: {
        value: 0,
        comparison: 0,
      },
      users: {
        value: 0,
        comparison: 0,
      },
      sales: {
        value: 0,
        comparison: 0,
      },
      payments: [],
      salesTrends: [],
    }),
  },
})

const emit = defineEmits<{
  (e: 'update:start-at', startAt: number): void
  (e: 'update:end-at', endAt: number): void
  (e: 'update:period-type', periodType: TopOrderPeriodType): void
}>()

const periodTypes = [
  { title: '日単位', value: TopOrderPeriodType.DAY },
  { title: '週単位', value: TopOrderPeriodType.WEEK },
  { title: '月単位', value: TopOrderPeriodType.MONTH },
]
const paymentHeaders = [
  { title: '支払い方法', key: 'paymentMethod' },
  { title: '利用率', value: 'orderRate' },
  { title: '注文数', key: 'orderCount' },
  { title: '利用者数', key: 'userCount' },
]

const startAtValue = computed<DateTimeInput>({
  get: (): DateTimeInput => ({
    date: unix(props.startAt).format('YYYY-MM-DD'),
    time: unix(props.startAt).format('HH:mm'),
  }),
  set: (val: DateTimeInput): void => {
    const startAt = dayjs(`${val.date} 00:00:00`)
    emit('update:start-at', startAt.unix())
  },
})
const endAtValue = computed<DateTimeInput>({
  get: (): DateTimeInput => ({
    date: unix(props.endAt).format('YYYY-MM-DD'),
    time: unix(props.endAt).format('HH:mm'),
  }),
  set: (val: DateTimeInput): void => {
    const endAt = dayjs(`${val.date} 00:00:00`).add(1, 'day')
    emit('update:end-at', endAt.unix())
  },
})
const periodTypeValue = computed<TopOrderPeriodType>({
  get: (): TopOrderPeriodType => props.periodType,
  set: (periodType: TopOrderPeriodType): void => emit('update:period-type', periodType),
})

const orderChartOption = computed(() => {
  const labels: string[] = []
  const values: number[] = []

  props.orders.salesTrends.forEach((trend: TopOrderSalesTrend) => {
    labels.push(trend.period)
    values.push(trend.salesTotal)
  })

  return {
    title: {
      show: labels.length === 0,
      left: 'center',
      top: 'center',
      text: 'データがありません',
      textStyle: {
        color: '#c0c0c0',
      },
    },
    tooltip: {
      trigger: 'axis',
    },
    xAxis: {
      type: 'category',
      data: labels,
      axisLabel: {
        rotate: 20,
      },
      name: '時刻',
    },
    yAxis: {
      minInterval: 1,
      type: 'value',
      name: '売上総額',
    },
    series: [
      {
        type: 'line',
        data: values,
        smooth: true,
      },
    ],
  }
})

const getComparison = (num: number): string => {
  let prefix = ''
  if (num > 0) {
    prefix = '+ '
  }
  return `${prefix}${num.toFixed()}`
}

const getComparisonColor = (num: number): string => {
  if (num === 0) {
    return 'text-grey'
  }
  return num > 0 ? 'text-primary' : 'text-error'
}

const getPaymentMethod = (methodType: PaymentMethodType): string => {
  switch (methodType) {
    case PaymentMethodType.PaymentMethodTypeCash:
      return '代引支払い'
    case PaymentMethodType.PaymentMethodTypeCreditCard:
      return 'クレジットカード決済'
    case PaymentMethodType.PaymentMethodTypeKonbini:
      return 'コンビニ決済'
    case PaymentMethodType.PaymentMethodTypeBankTransfer:
      return '銀行振込決済'
    case PaymentMethodType.PaymentMethodTypePayPay:
      return 'QR決済（PayPay）'
    case PaymentMethodType.PaymentMethodTypeLinePay:
      return 'QR決済（LINE Pay）'
    case PaymentMethodType.PaymentMethodTypeMerpay:
      return 'QR決済（メルペイ）'
    case PaymentMethodType.PaymentMethodTypeRakutenPay:
      return 'QR決済（楽天ペイ）'
    case PaymentMethodType.PaymentMethodTypeAUPay:
      return 'QR決済（au PAY）'
    case PaymentMethodType.PaymentMethodTypePaidy:
      return 'ペイディ（Paidy）'
    case PaymentMethodType.PaymentMethodTypePayEasy:
      return 'ペイジー（Pay-easy）'
    default:
      return '不明'
  }
}

const getOrderRate = (rate: number): string => {
  return rate.toFixed()
}

const onChangeStartAt = (): void => {
  const startAt = dayjs(`${startAtValue.value.date} 00:00:00`)
  emit('update:start-at', startAt.unix())
}

const onChangeEndAt = (): void => {
  const endAt = dayjs(`${endAtValue.value.date} 00:00:00`)
  emit('update:end-at', endAt.unix())
}
</script>

<template>
  <v-container>
    <v-row>
      <v-col
        class="d-flex flex-column flex-md-row justify-center"
        cols="8"
      >
        <v-text-field
          v-model="startAtValue.date"
          type="date"
          variant="outlined"
          density="compact"
          class="mr-md-2"
          @change="onChangeStartAt"
        />
        <p class="text-subtitle-2 mx-4 pt-md-3 mb-4 pb-md-6">
          〜
        </p>
        <v-text-field
          v-model="endAtValue.date"
          type="date"
          variant="outlined"
          density="compact"
          @change="onChangeEndAt"
        />
      </v-col>
      <v-col cols="4">
        <v-select
          v-model="periodTypeValue"
          :items="periodTypes"
          variant="outlined"
        />
      </v-col>
    </v-row>

    <v-row>
      <v-col cols="4">
        <v-card class="py-2">
          <v-card-title class="mb-2">
            注文件数
          </v-card-title>
          <v-card-text>
            <div class="text-h5 mb-2">
              {{ orders.orders.value.toLocaleString() }} 件
            </div>
            <div :class="getComparisonColor(orders.orders.comparison)">
              {{ getComparison(orders.orders.comparison) }} &percnt;
            </div>
          </v-card-text>
        </v-card>
      </v-col>

      <v-col cols="4">
        <v-card class="py-2">
          <v-card-title class="mb-2">
            売上総額
          </v-card-title>
          <v-card-text>
            <div class="text-h5 mb-2">
              {{ orders.sales.value.toLocaleString() }} 円
            </div>
            <div :class="getComparisonColor(orders.sales.comparison)">
              {{ getComparison(orders.sales.comparison) }} &percnt;
            </div>
          </v-card-text>
        </v-card>
      </v-col>

      <v-col cols="4">
        <v-card class="py-2">
          <v-card-title class="mb-2">
            購入者数
          </v-card-title>
          <v-card-text>
            <div class="text-h5 mb-2">
              {{ orders.users.value.toLocaleString() }} 人
            </div>
            <div :class="getComparisonColor(orders.users.comparison)">
              {{ getComparison(orders.users.comparison) }} &percnt;
            </div>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>

    <v-row>
      <v-col cols="12">
        <v-card :loading="loading">
          <v-card-title>
            売上額推移
          </v-card-title>
          <v-card-text>
            <v-chart
              :option="orderChartOption"
              class="chart"
              autoresize
            />
          </v-card-text>
        </v-card>
      </v-col>

      <v-col cols="12">
        <v-card :loading="loading">
          <v-card-title>
            利用状況（支払い方法別）
          </v-card-title>
          <v-card-text>
            <v-data-table
              :loading="loading"
              :headers="paymentHeaders"
              :items="orders.payments"
              no-data-text="データがありません"
            >
              <template #[`item.paymentMethod`]="{ item }">
                {{ getPaymentMethod(item.paymentMethodType) }}
              </template>
              <template #[`item.orderCount`]="{ item }">
                {{ item.orderCount.toLocaleString() }} 件
              </template>
              <template #[`item.userCount`]="{ item }">
                {{ item.userCount.toLocaleString() }} 人
              </template>
              <template #[`item.orderRate`]="{ item }">
                {{ getOrderRate(item.rate) }} &percnt;
              </template>
            </v-data-table>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<style scoped lang="scss">
.chart {
  height: 400px;
}
</style>
