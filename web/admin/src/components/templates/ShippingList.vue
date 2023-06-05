<script lang="ts" setup>
import { mdiPlus } from '@mdi/js'
import { VDataTable } from 'vuetify/lib/labs/components'
import { dateTimeFormatter, moneyFormat } from '~/lib/formatter'
import { prefecturesList } from '~/constants'
import { AlertType } from '~/lib/hooks'
import { ShippingsResponseShippingsInner } from '~/types/api'

const props = defineProps({
  loading: {
    type: Boolean,
    default: false
  },
  deleteDialog: {
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
  },
  tableSortBy: {
    type: Array as PropType<VDataTable['sortBy']>,
    default: () => []
  }
})

const emit = defineEmits<{
  (e: 'click:update-page', page: number): void
  (e: 'click:update-items-per-page', page: number): void
  (e: 'click:row', notificationId: string): void
  (e: 'click:add'): void
  (e: 'click:delete', notificationId: string): void
  (e: 'update:delete-dialog', v: boolean): void
  (e: 'update:sort-by', sortBy: VDataTable['sortBy']): void
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

const onClickUpdateSortBy = (sortBy: VDataTable['sortBy']): void => {
  emit('update:sort-by', sortBy)
}

const onClickRow = (shippingId: string): void => {
  emit('click:row', shippingId)
}

const onClickAdd = (): void => {
  emit('click:add')
}
</script>

<template>
  <v-card-title class="d-flex flex-row">
    配送設定一覧
    <v-spacer />
    <v-btn variant="outlined" color="primary" @click="onClickAdd">
      <v-icon start :icon="mdiPlus" />
      配送情報登録
    </v-btn>
  </v-card-title>
  <v-card class="mt-4" flat :loading="loading">
    <v-card-text>
      <v-data-table-server
        :headers="headers"
        :items="shippings"
        :items-per-page="props.tableItemsPerPage"
        :items-length="props.tableItemsTotal"
        :sort-by="props.tableSortBy"
        :multi-sort="true"
        show-expand
        hover
        class="elevation-0"
        @update:page="onClickUpdatePage"
        @update:items-per-page="onClickUpdateItemsPerPage"
        @update:sort-by="onClickUpdateSortBy"
        @update:sort-desc="onClickUpdateSortBy"
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

        <template #expanded-item="{ item }">
          <td :colspan="headers.length" class="pa-4">
            <div v-for="n in [60, 80, 100]" :key="n">
              <div class="row my-2">
                サイズ{{ n }}詳細
              </div>
              <v-row
                v-for="(boxRate, i) in item.raw[`box${n}Rates`]"
                :key="i"
                class="align-center"
              >
                <v-col cols="1">
                  {{ boxRate.number }}
                </v-col>
                <v-col cols="1">
                  {{ boxRate.name }}
                </v-col>
                <v-col cols="1">
                  {{ moneyFormat(boxRate.price) }} 円
                </v-col>
                <v-col cols="9">
                  <v-select
                    v-model="boxRate.prefectures"
                    :items="prefecturesList"
                    :label="`${boxRate.prefectures.length}/${prefecturesList.length}`"
                    hide-details
                  >
                    <template #selection="{ item: selectItem, index }">
                      <v-chip v-if="index < 5" size="small">
                        <span>{{ selectItem.title }}</span>
                      </v-chip>
                      <span
                        v-if="index === 5"
                        class="grey--text text-caption"
                      >
                        (+{{ boxRate.prefectures.length - 5 }} others)
                      </span>
                    </template>
                  </v-select>
                </v-col>
              </v-row>
            </div>
          </td>
        </template>
      </v-data-table-server>
    </v-card-text>
  </v-card>
</template>
