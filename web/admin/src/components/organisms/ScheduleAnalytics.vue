<script lang="ts" setup>
import dayjs from 'dayjs'
import type { ChartData, ChartOptions } from 'chart.js'
import { LineChart } from 'vue-chart-3'
import type { BroadcastViewerLog } from '~/types/api'

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

const chartData = computed<ChartData<'line'>>(() => {
  const labels: string[] = []
  const values: number[] = []

  props.viewerLogs.forEach((log) => {
    const startAt = dayjs.unix(log.startAt).format('YYYY-MM-DD HH:mm')
    labels.push(startAt)
    values.push(log.total)
  })

  return {
    labels,
    datasets: [
      {
        label: '視聴者数',
        data: values,
      },
    ],
  }
})

const chartOptions: ChartOptions<'line'> = {
  scales: {
    yAxes: {
      beginAtZero: true,
    },
  },
  responsive: true,
  maintainAspectRatio: false,
}
</script>

<template>
  <v-card :loading="loading">
    <v-card-text>
      <LineChart
        :chart-data="chartData"
        :options="chartOptions"
      />
    </v-card-text>
  </v-card>
</template>
