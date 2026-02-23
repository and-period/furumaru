<script lang="ts" setup>
import { use } from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'
import {
  TitleComponent,
  GridComponent,
  TooltipComponent,
  LegendComponent,
} from 'echarts/components'
import { LineChart, PieChart } from 'echarts/charts'
import VChart from 'vue-echarts'

use([
  GridComponent,
  CanvasRenderer,
  LineChart,
  PieChart,
  TitleComponent,
  TooltipComponent,
  LegendComponent,
])

defineProps<{
  option: Record<string, unknown>
  autoresize?: boolean
}>()
</script>

<template>
  <v-chart
    :option="option"
    :autoresize="autoresize"
    v-bind="$attrs"
  />
</template>
