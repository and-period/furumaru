<script lang="ts" setup>
import dayjs from 'dayjs'
import { use } from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'
import {
  TitleComponent,
  GridComponent,
  TooltipComponent,
} from 'echarts/components'
import { LineChart } from 'echarts/charts'
import VChart from 'vue-echarts'
import type { BroadcastViewerLog } from '~/types/api/v1'

const props = defineProps({
  loading: {
    type: Boolean,
    default: false,
  },
  viewerLogs: {
    type: Array<BroadcastViewerLog>,
    default: () => [],
  },
})

use([
  GridComponent,
  CanvasRenderer,
  LineChart,
  TitleComponent,
  TooltipComponent,
])

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
      name: '視聴者数',
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
</script>

<template>
  <v-card :loading="loading">
    <v-card-text>
      <v-chart
        :option="option"
        class="chart"
        autoresize
      />
    </v-card-text>
  </v-card>
</template>

<style scoped lang="scss">
.chart {
  height: 400px;
}
</style>
