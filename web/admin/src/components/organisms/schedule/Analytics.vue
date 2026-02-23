<script lang="ts" setup>
import {
  mdiChartLine,
  mdiEye,
  mdiClockOutline,
  mdiTrendingUp,
  mdiAccountMultiple,
} from '@mdi/js'
import dayjs from 'dayjs'
import type { BroadcastViewerLog } from '~/types/api/v1'

const AppChart = defineAsyncComponent(() => import('~/components/atoms/AppChart.vue'))

const props = defineProps({
  loading: {
    type: Boolean,
    default: false,
  },
  viewerLogs: {
    type: Array<BroadcastViewerLog>,
    default: () => [],
  },
  totalViewers: {
    type: Number,
    default: 0,
  },
})

const option = computed(() => {
  const labels: string[] = []
  const values: number[] = []

  props.viewerLogs.forEach((log) => {
    const startAt = dayjs.unix(log.startAt).format('YYYY-MM-DD HH:mm')
    labels.push(startAt)
    values.push(log.total)
  })

  return {
    title: {
      show: labels.length === 0,
      left: 'center',
      top: 'center',
      text: 'データがありません',
      textStyle: {
        color: '#9e9e9e',
        fontSize: 16,
      },
    },
    tooltip: {
      trigger: 'axis',
      backgroundColor: 'rgba(0, 0, 0, 0.8)',
      textStyle: {
        color: '#fff',
      },
      formatter: (params: any) => {
        if (!params || params.length === 0) return ''
        const data = params[0]
        return `${data.name}<br/>視聴者数: ${data.value}人`
      },
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '15%',
      containLabel: true,
    },
    xAxis: {
      type: 'category',
      data: labels,
      axisLabel: {
        rotate: 45,
        fontSize: 12,
        color: '#666',
      },
      axisLine: {
        lineStyle: {
          color: '#e0e0e0',
        },
      },
      name: '時刻',
      nameTextStyle: {
        color: '#666',
        fontSize: 12,
      },
    },
    yAxis: {
      minInterval: 1,
      type: 'value',
      name: '視聴者数',
      nameTextStyle: {
        color: '#666',
        fontSize: 12,
      },
      axisLabel: {
        color: '#666',
        fontSize: 12,
      },
      axisLine: {
        lineStyle: {
          color: '#e0e0e0',
        },
      },
      splitLine: {
        lineStyle: {
          color: '#f5f5f5',
          type: 'dashed',
        },
      },
    },
    series: [
      {
        type: 'line',
        data: values,
        smooth: true,
        lineStyle: {
          color: '#2196f3',
          width: 3,
        },
        itemStyle: {
          color: '#2196f3',
          borderColor: '#fff',
          borderWidth: 2,
        },
        areaStyle: {
          color: {
            type: 'linear',
            x: 0,
            y: 0,
            x2: 0,
            y2: 1,
            colorStops: [
              {
                offset: 0,
                color: 'rgba(33, 150, 243, 0.3)',
              },
              {
                offset: 1,
                color: 'rgba(33, 150, 243, 0.05)',
              },
            ],
          },
        },
      },
    ],
  }
})

const getTotalViewers = (): number => {
  return props.viewerLogs.reduce((sum, log) => sum + log.total, 0)
}

const getMaxViewers = (): number => {
  return props.viewerLogs.length > 0
    ? Math.max(...props.viewerLogs.map(log => log.total))
    : 0
}

const getAverageViewers = (): number => {
  if (props.viewerLogs.length === 0) return 0
  const total = getTotalViewers()
  return Math.round(total / props.viewerLogs.length)
}

const getBroadcastDuration = (): string => {
  if (props.viewerLogs.length === 0) return '--'

  const firstLog = props.viewerLogs[0]
  const lastLog = props.viewerLogs[props.viewerLogs.length - 1]
  const durationMinutes = Math.round((lastLog.startAt - firstLog.startAt) / 60)

  if (durationMinutes < 60) {
    return `${durationMinutes}分`
  }

  const hours = Math.floor(durationMinutes / 60)
  const minutes = durationMinutes % 60
  return `${hours}時間${minutes}分`
}
</script>

<template>
  <div class="analytics-container">
    <!-- サマリー統計セクション -->
    <v-row class="mb-6">
      <v-col
        cols="12"
        sm="6"
        lg="3"
      >
        <v-card
          class="stats-card"
          elevation="2"
        >
          <v-card-text class="d-flex align-center pa-4">
            <v-avatar
              color="primary"
              class="mr-4"
              size="48"
            >
              <v-icon
                :icon="mdiEye"
                size="24"
                color="white"
              />
            </v-avatar>
            <div>
              <p class="text-body-2 text-grey-darken-1 mb-1">
                総視聴者数
              </p>
              <p class="text-h6 font-weight-bold mb-0">
                {{ props.totalViewers }}人
              </p>
            </div>
          </v-card-text>
        </v-card>
      </v-col>
      <v-col
        cols="12"
        sm="6"
        lg="3"
      >
        <v-card
          class="stats-card"
          elevation="2"
        >
          <v-card-text class="d-flex align-center pa-4">
            <v-avatar
              color="success"
              class="mr-4"
              size="48"
            >
              <v-icon
                :icon="mdiTrendingUp"
                size="24"
                color="white"
              />
            </v-avatar>
            <div>
              <p class="text-body-2 text-grey-darken-1 mb-1">
                最大同時視聴者数
              </p>
              <p class="text-h6 font-weight-bold mb-0">
                {{ getMaxViewers() }}人
              </p>
            </div>
          </v-card-text>
        </v-card>
      </v-col>
      <v-col
        cols="12"
        sm="6"
        lg="3"
      >
        <v-card
          class="stats-card"
          elevation="2"
        >
          <v-card-text class="d-flex align-center pa-4">
            <v-avatar
              color="warning"
              class="mr-4"
              size="48"
            >
              <v-icon
                :icon="mdiAccountMultiple"
                size="24"
                color="white"
              />
            </v-avatar>
            <div>
              <p class="text-body-2 text-grey-darken-1 mb-1">
                平均視聴者数
              </p>
              <p class="text-h6 font-weight-bold mb-0">
                {{ getAverageViewers() }}人
              </p>
            </div>
          </v-card-text>
        </v-card>
      </v-col>
      <v-col
        cols="12"
        sm="6"
        lg="3"
      >
        <v-card
          class="stats-card"
          elevation="2"
        >
          <v-card-text class="d-flex align-center pa-4">
            <v-avatar
              color="info"
              class="mr-4"
              size="48"
            >
              <v-icon
                :icon="mdiClockOutline"
                size="24"
                color="white"
              />
            </v-avatar>
            <div>
              <p class="text-body-2 text-grey-darken-1 mb-1">
                配信時間
              </p>
              <p class="text-h6 font-weight-bold mb-0">
                {{ getBroadcastDuration() }}
              </p>
            </div>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>

    <!-- 視聴者数推移チャートセクション -->
    <v-card
      class="form-section-card"
      elevation="2"
      :loading="loading"
    >
      <v-card-title class="d-flex align-center section-header">
        <v-icon
          :icon="mdiChartLine"
          size="24"
          class="mr-3 text-primary"
        />
        <span class="text-h6 font-weight-medium">視聴者数推移</span>
      </v-card-title>
      <v-card-text class="pa-6">
        <div v-if="props.viewerLogs.length > 0">
          <app-chart
            :option="option"
            class="chart"
            :autoresize="true"
          />
        </div>
        <div
          v-else
          class="chart-placeholder"
        >
          <v-icon
            :icon="mdiChartLine"
            size="64"
            class="text-grey-lighten-1 mb-4"
          />
          <p class="text-body-1 text-grey-darken-1 mb-2">
            分析データがありません
          </p>
          <p class="text-body-2 text-grey-darken-1">
            配信が開始されると視聴者数の推移が表示されます
          </p>
        </div>
      </v-card-text>
    </v-card>
  </div>
</template>

<style scoped lang="scss">
.form-section-card {
  border-radius: 12px;
  max-width: none;
}

.section-header {
  background: linear-gradient(90deg, rgb(33 150 243 / 5%) 0%, rgb(33 150 243 / 0%) 100%);
  border-bottom: 1px solid rgb(0 0 0 / 5%);
  padding: 20px 24px;
}

.analytics-container {
  min-height: 400px;
}

.stats-card {
  border-radius: 12px;
  height: 100%;
}

.chart {
  height: 400px;
}

.chart-placeholder {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 400px;
  background: rgb(250 250 250);
  border-radius: 8px;
  border: 2px dashed rgb(224 224 224);
}

@media (width <= 600px) {
  .form-section-card, .stats-card {
    border-radius: 8px;
  }

  .section-header {
    padding: 16px 20px;
  }

  .chart {
    height: 300px;
  }

  .chart-placeholder {
    height: 300px;
  }
}
</style>
