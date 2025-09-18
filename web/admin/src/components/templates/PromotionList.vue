<script lang="ts" setup>
import { mdiDelete, mdiPlus, mdiTagOutline } from '@mdi/js'
import { unix } from 'dayjs'
import type { VDataTable } from 'vuetify/lib/components/index.mjs'

import type { AlertType } from '~/lib/hooks'
import {
  AdminType,
  DiscountType,
  PromotionStatus,
  PromotionTargetType,

} from '~/types/api/v1'
import type { Promotion, Shop } from '~/types/api/v1'

const props = defineProps({
  loading: {
    type: Boolean,
    default: false,
  },
  shopIds: {
    type: Array as PropType<string[]>,
    default: () => [],
  },
  adminType: {
    type: Number as PropType<AdminType>,
    default: AdminType.AdminTypeUnknown,
  },
  deleteDialog: {
    type: Boolean,
    default: false,
  },
  isAlert: {
    type: Boolean,
    default: false,
  },
  alertType: {
    type: String as PropType<AlertType>,
    default: undefined,
  },
  alertText: {
    type: String,
    default: '',
  },
  sortBy: {
    type: Array as PropType<VDataTable['sortBy']>,
    default: () => [],
  },
  promotions: {
    type: Array<Promotion>,
    default: () => [],
  },
  shops: {
    type: Array as PropType<Shop[]>,
    default: () => [],
  },
  tableItemsPerPage: {
    type: Number,
    default: 20,
  },
  tableItemsTotal: {
    type: Number,
    default: 0,
  },
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
    sortable: false,
  },
  {
    title: '対象マルシェ',
    key: 'target',
    sortable: false,
  },
  {
    title: 'ステータス',
    key: 'status',
    sortable: false,
  },
  {
    title: '割引コード',
    key: 'code',
    sortable: false,
  },
  {
    title: '割引額',
    key: 'discount',
    sortable: false,
  },
  {
    title: '使用期間',
    key: 'term',
    sortable: false,
  },
  {
    title: '使用回数',
    key: 'usedCount',
    sortable: false,
  },
  {
    title: '',
    key: 'actions',
    sortable: false,
  },
]

const selectedItem = ref<Promotion>()

const deleteDialogValue = computed({
  get: (): boolean => props.deleteDialog,
  set: (val: boolean): void => emit('update:delete-dialog', val),
})

const isRegisterable = (): boolean => {
  const registerable: AdminType[] = [AdminType.AdminTypeAdministrator, AdminType.AdminTypeCoordinator]
  return registerable.includes(props.adminType)
}

const isEditable = (promotion: Promotion): boolean => {
  switch (props.adminType) {
    case AdminType.AdminTypeAdministrator:
      return true
    case AdminType.AdminTypeCoordinator:
      return props.shopIds.includes(promotion.shopId)
    default:
      return false
  }
}

const getTarget = (promotion: Promotion): string => {
  switch (promotion.targetType) {
    case PromotionTargetType.PromotionTargetTypeAllShop:
      return '全て'
    case PromotionTargetType.PromotionTargetTypeSpecificShop: {
      const shop = props.shops.find((shop: Shop): boolean => shop.id === promotion.shopId)
      return shop?.name || ''
    }
    default:
      return ''
  }
}

const getDiscount = (
  discountType: DiscountType,
  discountRate: number,
): string => {
  switch (discountType) {
    case DiscountType.DiscountTypeAmount:
      return '￥' + discountRate.toLocaleString()
    case DiscountType.DiscountTypeRate:
      return discountRate + '％'
    case DiscountType.DiscountTypeFreeShipping:
      return '送料無料'
    default:
      return ''
  }
}

const getStatus = (status: PromotionStatus): string => {
  switch (status) {
    case PromotionStatus.PromotionStatusPrivate:
      return '非公開'
    case PromotionStatus.PromotionStatusWaiting:
      return '開始前'
    case PromotionStatus.PromotionStatusEnabled:
      return '有効'
    case PromotionStatus.PromotionStatusFinished:
      return '終了'
    default:
      return '無効'
  }
}

const getStatusColor = (status: PromotionStatus): string => {
  switch (status) {
    case PromotionStatus.PromotionStatusPrivate:
      return 'warning'
    case PromotionStatus.PromotionStatusWaiting:
      return 'info'
    case PromotionStatus.PromotionStatusEnabled:
      return 'primary'
    case PromotionStatus.PromotionStatusFinished:
      return 'secondary'
    default:
      return 'error'
  }
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
  <v-alert
    v-show="props.isAlert"
    :type="props.alertType"
    v-text="props.alertText"
  />

  <v-dialog
    v-model="deleteDialogValue"
    width="500"
  >
    <v-card>
      <v-card-title class="text-h6 py-4">
        セール情報削除の確認
      </v-card-title>
      <v-card-text class="pb-4">
        <div class="text-body-1">
          「{{ selectedItem?.title || "" }}」を削除しますか？
        </div>
        <div class="text-body-2 text-medium-emphasis mt-2">
          この操作は取り消せません。
        </div>
      </v-card-text>
      <v-card-actions class="px-6 pb-4">
        <v-spacer />
        <v-btn
          color="medium-emphasis"
          variant="text"
          @click="onClickCloseDeleteDialog"
        >
          キャンセル
        </v-btn>
        <v-btn
          :loading="loading"
          color="error"
          variant="elevated"
          @click="onClickDelete"
        >
          削除
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>

  <v-card
    class="mt-6"
    elevation="0"
    rounded="lg"
  >
    <v-card-title class="d-flex flex-column flex-sm-row align-start align-sm-center justify-space-between pa-4 pa-sm-6 pb-4">
      <div class="d-flex align-center mb-3 mb-sm-0">
        <v-icon
          :icon="mdiTagOutline"
          size="28"
          class="mr-3 text-primary"
        />
        <div>
          <h1 class="text-h6 text-sm-h5 font-weight-bold text-primary">
            セール情報管理
          </h1>
          <p class="text-caption text-sm-body-2 text-medium-emphasis ma-0">
            セール・キャンペーン情報の登録・編集・削除を行います
          </p>
        </div>
      </div>
      <div class="d-flex ga-3 w-100 w-sm-auto">
        <v-btn
          v-show="isRegisterable()"
          variant="elevated"
          color="primary"
          :size="$vuetify.display.smAndDown ? 'default' : 'large'"
          class="w-100 w-sm-auto"
          @click="onClickAdd"
        >
          <v-icon
            start
            :icon="mdiPlus"
          />
          セール情報登録
        </v-btn>
      </div>
    </v-card-title>

    <v-card-text>
      <v-data-table-server
        :headers="headers"
        :loading="loading"
        :items="promotions"
        :items-per-page="props.tableItemsPerPage"
        :items-length="props.tableItemsTotal"
        :sort-by="props.sortBy"
        hover
        no-data-text="登録されているセール情報がありません。"
        @update:page="onClickUpdatePage"
        @update:items-per-page="onClickUpdateItemsPerPage"
        @update:sort-by="onClickUpdateSortBy"
        @click:row="(_: any, { item }: any) => onClickRow(item.id)"
      >
        <template #[`item.title`]="{ item }">
          {{ item.title }}
        </template>
        <template #[`item.target`]="{ item }">
          {{ getTarget(item) }}
        </template>
        <template #[`item.status`]="{ item }">
          <v-chip
            :color="getStatusColor(item.status)"
          >
            {{ getStatus(item.status) }}
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
            v-show="isEditable(item)"
            variant="outlined"
            color="error"
            size="small"
            :prepend-icon="mdiDelete"
            @click.stop="onClickOpenDeleteDialog(item)"
          >
            削除
          </v-btn>
        </template>
      </v-data-table-server>
    </v-card-text>
  </v-card>
</template>
