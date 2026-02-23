<script lang="ts" setup>
import dayjs from 'dayjs'
import type { VideoViewerLog } from '~/types/api/v1'

const AppChart = defineAsyncComponent(() => import('~/components/atoms/AppChart.vue'))

const props = defineProps({
  loading: {
    type: Boolean,
    default: false,
  },
  viewerLogs: {
    type: Array<VideoViewerLog>,
    default: () => [],
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
      <app-chart
        :option="option"
        class="chart"
        :autoresize="true"
      />
    </v-card-text>
  </v-card>
</template>

<style scoped lang="scss">
.chart {
  height: 400px;
}
</style>
