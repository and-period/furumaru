<script lang="ts" setup>
import { mdiPlus } from '@mdi/js'
import { VDataTable } from 'vuetify/lib/labs/components'
import { dateTimeFormatter, moneyFormat } from '~/lib/formatter'
import { AlertType } from '~/lib/hooks'
import { ShippingsResponseShippingsInner } from '~/types/api'

const props = defineProps({
  loading: {
    type: Boolean,
    default: false
  },
  isAlert: {
    type: Boolean,
    default: false
  },
  alertType: {
    type: String as PropType<AlertType>,
    default: undefined
  },
  alertText: {
    type: String,
    default: ''
  },
  shippings: {
    type: Array<ShippingsResponseShippingsInner>,
    default: () => []
  },
  tableItemsPerPage: {
    type: Number,
    default: 20
  },
  tableItemsTotal: {
    type: Number,
    default: 0
  }
})

const emit = defineEmits<{
  (e: 'click:update-page', page: number): void
  (e: 'click:update-items-per-page', page: number): void
  (e: 'click:row', notificationId: string): void
  (e: 'click:add'): void
}>()

const headers: VDataTable['headers'] = [
  {
    title: '名前',
    key: 'name'
  },
  {
    title: '配送無料オプション',
    key: 'hasFreeShipping'
  },
  {
    title: '更新日',
    key: 'updatedAt'
  }
]

const onClickUpdatePage = (page: number): void => {
  emit('click:update-page', page)
}

const onClickUpdateItemsPerPage = (page: number): void => {
  emit('click:update-items-per-page', page)
}

const onClickRow = (shippingId: string): void => {
  emit('click:row', shippingId)
}

const onClickAdd = (): void => {
  emit('click:add')
}
</script>

<template>
  <v-card class="mt-4" flat :loading="loading">
    <v-card-title class="d-flex flex-row">
      配送設定一覧
      <v-spacer />
      <v-btn variant="outlined" color="primary" @click="onClickAdd">
        <v-icon start :icon="mdiPlus" />
        配送情報登録
      </v-btn>
    </v-card-title>

    <v-card-text>
      <v-data-table-server
        :headers="headers"
        :items="shippings"
        :items-per-page="props.tableItemsPerPage"
        :items-length="props.tableItemsTotal"
        :multi-sort="true"
        hover
        class="elevation-0"
        @update:page="onClickUpdatePage"
        @update:items-per-page="onClickUpdateItemsPerPage"
        @click:row="(_: any, {item}: any) => onClickRow(item.raw.id)"
      >
        <template #[`item.hasFreeShipping`]="{ item }">
          <v-chip size="small">
            {{ item.raw.hasFreeShipping ? '有り' : '無し' }}
          </v-chip>
        </template>
        <template #[`item.updatedAt`]="{ item }">
          {{ dateTimeFormatter(item.raw.updatedAt) }}
        </template>
      </v-data-table-server>
    </v-card-text>
  </v-card>
</template>
