<script setup lang='ts'>
import dayjs from 'dayjs'
import VersionInfo from '../../version.json'

const formattedBuildData = computed(() => {
  return dayjs(VersionInfo.buildDate).format('YYYY/MM/DD HH:mm')
})

</script>

<template>
  <v-card>
    <v-card-title>
      バージョン情報
    </v-card-title>
    <v-card-text>
      <v-row>
        <v-col cols="3">
          バージョン
        </v-col>
        <v-col cols="9">
          {{ VersionInfo.version }}
        </v-col>
      </v-row>
      <v-row>
        <v-col cols="3">
          ビルド日時
        </v-col>
        <v-col cols="9">
          {{ formattedBuildData }}
        </v-col>
      </v-row>
    </v-card-text>
  </v-card>
</template>
