<script lang="ts" setup>
import { mdiDelete, mdiPlus } from '@mdi/js'
import { unix } from 'dayjs'
import { VDataTable } from 'vuetify/lib/labs/components.mjs'
import type { AlertType } from '~/lib/hooks'
import { AdminRole, DiscountType, type Promotion } from '~/types/api'

const props = defineProps({
  loading: {
    type: Boolean,
    default: false
  },
  role: {
    type: Number as PropType<AdminRole>,
    default: AdminRole.UNKNOWN
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
  sortBy: {
    type: Array as PropType<VDataTable['sortBy']>,
    default: () => []
  },
  promotions: {
    type: Array<Promotion>,
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
  (e: 'click:delete', notificationId: string): void
  (e: 'update:delete-dialog', v: boolean): void
  (e: 'update:sort-by', sortBy: VDataTable['sortBy']): void
}>()

const headers: VDataTable['headers'] = [
  {
    title: 'タイトル',
    key: 'title',
    sortable: false
  },
  {
    title: 'ステータス',
    key: 'public',
    sortable: false
  },
  {
    title: '割引コード',
    key: 'code',
    sortable: false
  },
  {
    title: '割引額',
    key: 'discount',
    sortable: false
  },
  {
    title: '使用期間',
    key: 'term',
    sortable: false
  },
  {
    title: '使用回数',
    key: 'total',
    sortable: false
  },
  {
    title: '',
    key: 'actions',
    sortable: false
  }
]

const selectedItem = ref<Promotion>()

const deleteDialogValue = computed({
  get: (): boolean => props.deleteDialog,
  set: (val: boolean): void => emit('update:delete-dialog', val)
})

const isRegisterable = (): boolean => {
  return props.role === AdminRole.ADMINISTRATOR
}

const isEditable = (): boolean => {
  return props.role === AdminRole.ADMINISTRATOR
}

const getDiscount = (discountType: number, discountRate: DiscountType): string => {
  switch (discountType) {
    case DiscountType.AMOUNT:
      return '￥' + discountRate.toLocaleString()
    case DiscountType.RATE:
      return discountRate + '％'
    case DiscountType.FREE_SHIPPING:
      return '送料無料'
    default:
      return ''
  }
}

const getStatus = (status: boolean): string => {
  return status ? '有効' : '無効'
}

const getStatusColor = (status: boolean): string => {
  return status ? 'primary' : 'error'
}

const getDay = (unixTime: number): string => {
  return unix(unixTime).format('YYYY/MM/DD HH:mm')
}

const getTerm = (promotion: Promotion): string => {
  return `${getDay(promotion.startAt)} ~ ${getDay(promotion.endAt)}`
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

const onClickOpenDeleteDialog = (promotion: Promotion): void => {
  selectedItem.value = promotion
  deleteDialogValue.value = true
}

const onClickCloseDeleteDialog = (): void => {
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

  <v-dialog v-model="deleteDialogValue" width="500">
    <v-card>
      <v-card-title class="text-h7">
        {{ selectedItem?.title || '' }}を本当に削除しますか？
      </v-card-title>
      <v-card-actions>
        <v-spacer />
        <v-btn color="error" variant="text" @click="onClickCloseDeleteDialog">
          キャンセル
        </v-btn>
        <v-btn :loading="loading" color="primary" variant="outlined" @click="onClickDelete">
          削除
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>

  <v-card class="mt-4" flat :loading="props.loading">
    <v-card-title class="d-flex flex-row">
      セール情報
      <v-spacer />
      <v-btn v-show="isRegisterable()" variant="outlined" color="primary" @click="onClickAdd">
        <v-icon start :icon="mdiPlus" />
        セール情報登録
      </v-btn>
    </v-card-title>

    <v-card-text>
      <v-data-table-server
        :headers="headers"
        :items="promotions"
        :items-per-page="props.tableItemsPerPage"
        :items-length="props.tableItemsTotal"
        :sort-by="props.sortBy"
        :multi-sort="true"
        hover
        no-data-text="登録されているセール情報がありません。"
        @update:page="onClickUpdatePage"
        @update:items-per-page="onClickUpdateItemsPerPage"
        @update:sort-by="onClickUpdateSortBy"
        @update:sort-desc="onClickUpdateSortBy"
        @click:row="(_: any, {item}: any) => onClickRow(item.id)"
      >
        <template #[`item.title`]="{ item }">
          {{ item.title }}
        </template>
        <template #[`item.public`]="{ item }">
          <v-chip size="small" :color="getStatusColor(item.public)">
            {{ getStatus(item.public) }}
          </v-chip>
        </template>
        <template #[`item.code`]="{ item }">
          {{ item.code }}
        </template>
        <template #[`item.discount`]="{ item }">
          {{ getDiscount(item.discountType, item.discountRate) }}
        </template>
        <template #[`item.term`]="{ item }">
          {{ getTerm(item) }}
        </template>
        <template #[`item.actions`]="{ item }">
          <v-btn
            v-show="isEditable()"
            color="primary"
            size="small"
            variant="outlined"
            @click.stop="onClickOpenDeleteDialog(item)"
          >
            <v-icon size="small" :icon="mdiDelete" />
            削除
          </v-btn>
        </template>
      </v-data-table-server>
    </v-card-text>
  </v-card>
</template>
