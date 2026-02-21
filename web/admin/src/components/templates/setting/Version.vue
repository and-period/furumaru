<script lang="ts" setup>
import type { Dayjs } from 'dayjs'
import dayjs from 'dayjs'

const props = defineProps({
  version: {
    type: String,
    default: '',
  },
  builtAt: {
    type: Object as PropType<Dayjs>,
    default: () => dayjs(),
  },
})

const formattedBuiltAt = computed(() => {
  return props.builtAt.format('YYYY/MM/DD HH:mm')
})
</script>

<template>
  <v-card>
    <v-card-title>バージョン情報</v-card-title>

    <v-card-text>
      <v-row>
        <v-col cols="3">
          バージョン
        </v-col>
        <v-col cols="9">
          {{ props.version }}
        </v-col>
      </v-row>
      <v-row>
        <v-col cols="3">
          ビルド日時
        </v-col>
        <v-col cols="9">
          {{ formattedBuiltAt }}
        </v-col>
      </v-row>
    </v-card-text>
  </v-card>
</template>
