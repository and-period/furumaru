<script lang="ts" setup>
import { mdiDelete, mdiPlus } from '@mdi/js'
import { unix } from 'dayjs'
import { VDataTable } from 'vuetify/lib/labs/components'
import { AlertType } from '~/lib/hooks'
import { PromotionsResponsePromotionsInner } from '~/types/api'

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
  promotions: {
    type: Array<PromotionsResponsePromotionsInner>,
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
    title: 'タイトル',
    key: 'title'
  },
  {
    title: 'ステータス',
    key: 'public'
  },
  {
    title: '割引コード',
    key: 'code'
  },
  {
    title: '割引方法',
    key: 'discount'
  },
  {
    title: '使用開始',
    key: 'startAt'
  },
  {
    title: '使用終了',
    key: 'endAt'
  },
  {
    title: 'Actions',
    key: 'actions',
    sortable: false
  }
]

const selectedItem = ref<PromotionsResponsePromotionsInner>()

const deleteDialogValue = computed({
  get: () => props.deleteDialog,
  set: (val: boolean) => emit('update:delete-dialog', val)
})

const getDiscount = (discountType: number, discountRate: number): string => {
  switch (discountType) {
    case 1:
      return '-' + discountRate + '円'
    case 2:
      return '-' + discountRate + '%'
    case 3:
      return '送料無料'
    default:
      return ''
  }
}

const getStatus = (status: boolean): string => {
  if (status) {
    return '有効'
  } else {
    return '無効'
  }
}

const getStatusColor = (status: boolean): string => {
  if (status) {
    return 'primary'
  } else {
    return 'error'
  }
}

const getDay = (unixTime: number): string => {
  return unix(unixTime).format('YYYY/MM/DD HH:mm')
}

const onClickUpdatePage = (page: number): void => {
  emit('click:update-page', page)
}

const onClickUpdateItemsPerPage = (page: number): void => {
  emit('click:update-items-per-page', page)
}

const onClickUpdateSortBy = (sortBy: VDataTable['sortBy']): void => {
  emit('update:sort-by', sortBy)
}

const onClickRow = (promotionId: string): void => {
  emit('click:row', promotionId)
}

const onClickOpen = (promotion: PromotionsResponsePromotionsInner): void => {
  selectedItem.value = promotion
  deleteDialogValue.value = true
}

const onClickClose = (): void => {
  deleteDialogValue.value = false
}

const onClickAdd = (): void => {
  emit('click:add')
}

const onClickDelete = (): void => {
  emit('click:delete', selectedItem?.value?.id || '')
}
</script>

<template>
  <v-alert v-show="props.isAlert" :type="props.alertType" v-text="props.alertText" />
  <v-card-title class="d-flex flex-row">
    セール情報
    <v-spacer />
    <v-btn variant="outlined" color="primary" @click="onClickAdd">
      <v-icon start :icon="mdiPlus" />
      セール情報登録
    </v-btn>
  </v-card-title>

  <v-dialog v-model="deleteDialogValue" width="500">
    <v-card>
      <v-card-title class="text-h7">
        {{ selectedItem?.title || '' }}を本当に削除しますか？
      </v-card-title>
      <v-card-actions>
        <v-spacer />
        <v-btn color="error" variant="text" @click="onClickClose">
          キャンセル
        </v-btn>
        <v-btn color="primary" variant="outlined" @click="onClickDelete">
          削除
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>

  <v-card class="mt-4" flat :loading="props.loading">
    <v-card-text>
      <v-data-table-server
        :headers="headers"
        :items="promotions"
        :items-per-page="props.tableItemsPerPage"
        :items-length="props.tableItemsTotal"
        :sort-by="props.tableSortBy"
        :multi-sort="true"
        hover
        no-data-text="登録されているセール情報がありません。"
        @update:page="onClickUpdatePage"
        @update:items-per-page="onClickUpdateItemsPerPage"
        @update:sort-by="onClickUpdateSortBy"
        @update:sort-desc="onClickUpdateSortBy"
        @click:row="(_: any, {item}: any) => onClickRow(item.raw.id)"
      >
        <template #[`item.title`]="{ item }">
          {{ item.raw.title }}
        </template>
        <template #[`item.public`]="{ item }">
          <v-chip size="small" :color="getStatusColor(item.raw.public)">
            {{ getStatus(item.raw.public) }}
          </v-chip>
        </template>
        <template #[`item.code`]="{ item }">
          {{ item.raw.code }}
        </template>
        <template #[`item.discount`]="{ item }">
          {{ getDiscount(item.raw.discountType, item.raw.discountRate) }}
        </template>
        <template #[`item.startAt`]="{ item }">
          {{ getDay(item.raw.startAt) }}
        </template>
        <template #[`item.endAt`]="{ item }">
          {{ getDay(item.raw.endAt) }}
        </template>
        <template #[`item.actions`]="{ item }">
          <v-btn
            color="primary"
            size="small"
            variant="outlined"
            @click.stop="onClickOpen(item.raw)"
          >
            <v-icon size="small" :icon="mdiDelete" />
            削除
          </v-btn>
        </template>
      </v-data-table-server>
    </v-card-text>
  </v-card>
</template>
