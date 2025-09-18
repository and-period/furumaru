<script lang="ts" setup>
import dayjs, { unix } from 'dayjs'
import { use } from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'
import {
  TitleComponent,
  GridComponent,
  TooltipComponent, LegendComponent,
} from 'echarts/components'
import { LineChart, PieChart } from 'echarts/charts'
import VChart from 'vue-echarts'

import type { AlertType } from '~/lib/hooks'
import { TopOrderPeriodType } from '~/types'
import { PaymentMethodType } from '~/types/api/v1'
import type { TopOrderSalesTrend, TopOrdersResponse } from '~/types/api/v1'
import type { DateTimeInput } from '~/types/props'
import { mdiAccountMultipleOutline, mdiCartOutline, mdiCashMultiple, mdiChartLine, mdiCreditCardOutline, mdiCurrencyJpy, mdiTable } from '@mdi/js'

use([
  GridComponent,
  CanvasRenderer,
  LineChart,
  PieChart,
  TitleComponent,
  TooltipComponent,
  LegendComponent,
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
    grid: {
      left: '3%',
      right: '4%',
      bottom: '10%',
      containLabel: true,
    },
    xAxis: {
      type: 'category',
      data: labels,
      boundaryGap: false, // 線グラフの点を目盛りの位置に配置
      axisLabel: {
        rotate: 30,
        interval: 0, // すべてのラベルを表示
        align: 'right', // ラベルを右寄せで回転
      },
      name: '期間',
      nameLocation: 'middle',
      nameGap: 30,
    },
    yAxis: {
      minInterval: 1,
      type: 'value',
      name: '売上総額（円）',
      axisLabel: {
        formatter: (value: number) => `¥${value.toLocaleString()}`,
      },
    },
    series: [
      {
        type: 'line',
        data: values,
        smooth: true,
        symbolSize: 6, // データポイントのサイズ
        itemStyle: {
          color: '#1976D2', // プライマリカラー
        },
        lineStyle: {
          width: 2,
        },
        emphasis: {
          focus: 'series',
          itemStyle: {
            shadowBlur: 10,
            shadowColor: 'rgba(0,0,0,0.3)',
          },
        },
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

const averageOrderValue = computed(() => {
  if (props.orders.orders.value === 0) {
    return 0
  }
  return Math.floor(props.orders.sales.value / props.orders.orders.value)
})

const sortedPayments = computed(() => {
  return [...props.orders.payments].sort((a, b) => b.orderCount - a.orderCount)
})

const paymentChartOption = computed(() => {
  const data = props.orders.payments
    .map(p => ({
      value: p.orderCount,
      name: getPaymentMethod(p.paymentMethodType),
    }))
    .sort((a, b) => b.value - a.value) // 降順（大きい順）にソート

  return {
    title: {
      show: data.length === 0,
      left: 'center',
      top: 'center',
      text: 'データがありません',
      textStyle: {
        color: '#c0c0c0',
      },
    },
    tooltip: {
      trigger: 'item',
      formatter: '{b}: {c}件 ({d}%)',
    },
    legend: {
      type: 'scroll',
      orient: 'vertical',
      right: 10,
      top: 20,
      bottom: 20,
    },
    series: [{
      type: 'pie',
      radius: ['40%', '70%'],
      center: ['40%', '50%'],
      avoidLabelOverlap: false,
      itemStyle: {
        borderRadius: 10,
        borderColor: '#fff',
        borderWidth: 2,
      },
      label: {
        show: false,
        position: 'center',
      },
      emphasis: {
        label: {
          show: true,
          fontSize: 16,
          fontWeight: 'bold',
          formatter: '{b}\n{c}件 ({d}%)',
        },
      },
      labelLine: {
        show: false,
      },
      data,
      // 色の設定（大きい順に濃い色から薄い色へ）
      color: [
        '#1976D2', // 濃い青
        '#42A5F5', // 青
        '#66BB6A', // 緑
        '#FFA726', // オレンジ
        '#EF5350', // 赤
        '#AB47BC', // 紫
        '#78909C', // グレー
        '#8D6E63', // 茶色
      ],
    }],
  }
})

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
    <!-- Date Range & Period Selection -->
    <v-card
      class="mb-4"
      elevation="2"
    >
      <v-card-text>
        <v-row align="center">
          <v-col
            cols="12"
            md="9"
          >
            <div class="d-flex align-center flex-wrap ga-3">
              <v-text-field
                v-model="startAtValue.date"
                type="date"
                label="開始日"
                variant="outlined"
                density="compact"
                hide-details
                class="flex-grow-1"
                style="max-width: 200px"
                @change="onChangeStartAt"
              />
              <div class="text-subtitle-2 px-2">
                〜
              </div>
              <v-text-field
                v-model="endAtValue.date"
                type="date"
                label="終了日"
                variant="outlined"
                density="compact"
                hide-details
                class="flex-grow-1"
                style="max-width: 200px"
                @change="onChangeEndAt"
              />
            </div>
          </v-col>
          <v-col
            cols="12"
            md="3"
          >
            <v-select
              v-model="periodTypeValue"
              :items="periodTypes"
              label="集計単位"
              variant="outlined"
              density="compact"
              hide-details
            />
          </v-col>
        </v-row>
      </v-card-text>
    </v-card>

    <!-- KPI Cards -->
    <v-row class="mb-4">
      <v-col
        cols="12"
        md="3"
      >
        <v-card
          class="kpi-card"
          elevation="2"
        >
          <v-card-text>
            <div class="d-flex align-center mb-2">
              <v-icon
                color="primary"
                size="small"
                class="mr-2"
                :icon="mdiCurrencyJpy"
              />
              <span class="text-caption text-grey-darken-1">売上総額</span>
            </div>
            <div class="text-h5 font-weight-bold mb-1">
              ¥{{ orders.sales.value.toLocaleString() }}
            </div>
            <v-chip
              :color="orders.sales.comparison >= 0 ? 'success' : 'error'"
              size="x-small"
              label
            >
              <v-icon
                size="x-small"
                class="mr-1"
              >
                {{ orders.sales.comparison >= 0 ? 'mdi-trending-up' : 'mdi-trending-down' }}
              </v-icon>
              {{ getComparison(orders.sales.comparison) }}%
            </v-chip>
          </v-card-text>
        </v-card>
      </v-col>

      <v-col
        cols="12"
        md="3"
      >
        <v-card
          class="kpi-card"
          elevation="2"
        >
          <v-card-text>
            <div class="d-flex align-center mb-2">
              <v-icon
                color="info"
                size="small"
                class="mr-2"
                :icon="mdiCartOutline"
              />
              <span class="text-caption text-grey-darken-1">注文件数</span>
            </div>
            <div class="text-h5 font-weight-bold mb-1">
              {{ orders.orders.value.toLocaleString() }}
              <span class="text-body-2">件</span>
            </div>
            <v-chip
              :color="orders.orders.comparison >= 0 ? 'success' : 'error'"
              size="x-small"
              label
            >
              <v-icon
                size="x-small"
                class="mr-1"
              >
                {{ orders.orders.comparison >= 0 ? 'mdi-trending-up' : 'mdi-trending-down' }}
              </v-icon>
              {{ getComparison(orders.orders.comparison) }}%
            </v-chip>
          </v-card-text>
        </v-card>
      </v-col>

      <v-col
        cols="12"
        md="3"
      >
        <v-card
          class="kpi-card"
          elevation="2"
        >
          <v-card-text>
            <div class="d-flex align-center mb-2">
              <v-icon
                color="success"
                size="small"
                class="mr-2"
                :icon="mdiAccountMultipleOutline"
              />
              <span class="text-caption text-grey-darken-1">購入者数</span>
            </div>
            <div class="text-h5 font-weight-bold mb-1">
              {{ orders.users.value.toLocaleString() }}
              <span class="text-body-2">人</span>
            </div>
            <v-chip
              :color="orders.users.comparison >= 0 ? 'success' : 'error'"
              size="x-small"
              label
            >
              <v-icon
                size="x-small"
                class="mr-1"
              >
                {{ orders.users.comparison >= 0 ? 'mdi-trending-up' : 'mdi-trending-down' }}
              </v-icon>
              {{ getComparison(orders.users.comparison) }}%
            </v-chip>
          </v-card-text>
        </v-card>
      </v-col>

      <v-col
        cols="12"
        md="3"
      >
        <v-card
          class="kpi-card"
          elevation="2"
        >
          <v-card-text>
            <div class="d-flex align-center mb-2">
              <v-icon
                color="warning"
                size="small"
                class="mr-2"
                :icon="mdiCashMultiple"
              />
              <span class="text-caption text-grey-darken-1">平均単価</span>
            </div>
            <div class="text-h5 font-weight-bold mb-1">
              ¥{{ averageOrderValue.toLocaleString() }}
            </div>
            <div class="text-caption text-grey">
              売上 ÷ 注文数
            </div>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>

    <!-- Charts Row -->
    <v-row>
      <v-col
        cols="12"
        lg="8"
      >
        <v-card
          :loading="loading"
          elevation="2"
        >
          <v-card-title class="d-flex align-center">
            <v-icon
              color="primary"
              class="mr-2"
              :icon="mdiChartLine"
            />
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

      <v-col
        cols="12"
        lg="4"
      >
        <v-card
          :loading="loading"
          elevation="2"
        >
          <v-card-title class="d-flex align-center">
            <v-icon
              color="info"
              class="mr-2"
              :icon="mdiCreditCardOutline"
            />
            支払い方法別
          </v-card-title>
          <v-card-text>
            <v-chart
              :option="paymentChartOption"
              class="payment-chart"
              autoresize
            />
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>

    <!-- Payment Details Table -->
    <v-row>
      <v-col cols="12">
        <v-card
          :loading="loading"
          elevation="2"
        >
          <v-card-title class="d-flex align-center">
            <v-icon
              color="secondary"
              class="mr-2"
              :icon="mdiTable"
            />
          </v-card-title>
          <v-card-text>
            <v-data-table
              :loading="loading"
              :headers="paymentHeaders"
              :items="sortedPayments"
              no-data-text="データがありません"
              disable-pagination
              hide-default-footer
            >
              <template #[`item.paymentMethod`]="{ item, index }">
                <div class="d-flex align-center">
                  <v-chip
                    v-if="index < 3"
                    size="x-small"
                    :color="index === 0 ? 'gold' : index === 1 ? 'grey' : 'brown'"
                    class="mr-2"
                  >
                    {{ index + 1 }}
                  </v-chip>
                  {{ getPaymentMethod(item.paymentMethodType) }}
                </div>
              </template>
              <template #[`item.orderCount`]="{ item }">
                {{ item.orderCount.toLocaleString() }} 件
              </template>
              <template #[`item.userCount`]="{ item }">
                {{ item.userCount.toLocaleString() }} 人
              </template>
              <template #[`item.orderRate`]="{ item }">
                <v-chip
                  size="small"
                  variant="tonal"
                  color="primary"
                >
                  {{ getOrderRate(item.rate) }}%
                </v-chip>
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

.payment-chart {
  height: 380px;
}

.kpi-card {
  transition: all 0.3s;

  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
  }
}

:deep(.v-chip) {
  font-weight: 600;
}
</style>
