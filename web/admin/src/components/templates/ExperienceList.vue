<script lang="ts" setup>
import { mdiDelete, mdiPlus, mdiContentCopy } from '@mdi/js'
import type { VDataTable } from 'vuetify/lib/components/index.mjs'

import { getResizedImages } from '~/lib/helpers'
import type { AlertType } from '~/lib/hooks'
import {
  type Product,
  type ProductMediaInner,
  ProductStatus,
  type Category,
  type ProductTag,
  type ProductType,
  type Producer,
  AdminRole,
} from '~/types/api'

const props = defineProps({
  selectedItemId: {
    type: String,
    default: '',
  },
  loading: {
    type: Boolean,
    default: false,
  },
  role: {
    type: Number as PropType<AdminRole>,
    default: AdminRole.UNKNOWN,
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
  categories: {
    type: Array<Category>,
    default: () => [],
  },
  producers: {
    type: Array<Producer>,
    default: () => [],
  },
  products: {
    type: Array<Product>,
    default: () => [],
  },
  productTags: {
    type: Array<ProductTag>,
    default: () => [],
  },
  productTypes: {
    type: Array<ProductType>,
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
  (e: 'click:show', productId: string): void
  (e: 'click:new'): void
  (e: 'click:copyItem'): void
  (e: 'click:delete', productId: string): void
  (e: 'update:delete-dialog', v: boolean): void
  (e: 'update:selectedItemId', v: string): void
}>()

const headers: VDataTable['headers'] = [
  {
    title: '',
    key: 'media',
    width: 80,
    sortable: false,
  },
  {
    title: '商品名',
    key: 'name',
    sortable: false,
  },
  {
    title: 'ステータス',
    key: 'status',
    sortable: false,
  },
  {
    title: 'カテゴリ',
    key: 'experienceType',
    sortable: false,
  },
  {
    title: '場所',
    key: 'place',
    sortable: false,
  },
  {
    title: '生産者名',
    key: 'producerName',
    sortable: false,
  },
  {
    title: '',
    key: 'actions',
    sortable: false,
  },
]

const handleUpdateSelectItemId = (itemIds: string[]): void => {
  if (itemIds.length === 0) {
    emit('update:selectedItemId', '')
  }
  else {
    emit('update:selectedItemId', itemIds[0])
  }
}

const selectedItem = ref<Product>()

const deleteDialogValue = computed({
  get: (): boolean => props.deleteDialog,
  set: (val: boolean): void => emit('update:delete-dialog', val),
})

const isRegisterable = (): boolean => {
  return props.role === AdminRole.COORDINATOR
}

const isDeletable = (): boolean => {
  const targets: AdminRole[] = [AdminRole.ADMINISTRATOR, AdminRole.COORDINATOR]
  return targets.includes(props.role)
}

const getStatus = (status: ProductStatus): string => {
  switch (status) {
    case ProductStatus.PRESALE:
      return '予約販売'
    case ProductStatus.FOR_SALE:
      return '販売中'
    case ProductStatus.OUT_OF_SALES:
      return '販売終了'
    case ProductStatus.PRIVATE:
      return '非公開'
    case ProductStatus.ARCHIVED:
      return 'アーカイブ済み'
    default:
      return ''
  }
}

const getStatusColor = (status: ProductStatus): string => {
  switch (status) {
    case ProductStatus.PRESALE:
      return 'info'
    case ProductStatus.FOR_SALE:
      return 'primary'
    case ProductStatus.OUT_OF_SALES:
      return 'secondary'
    case ProductStatus.PRIVATE:
      return 'warning'
    case ProductStatus.ARCHIVED:
      return 'error'
    default:
      return ''
  }
}

const getInventoryColor = (inventory: number): string => {
  return inventory > 0 ? '' : 'text-error'
}

const toggleDeleteDialog = (product?: Product): void => {
  if (product) {
    selectedItem.value = product
  }
  deleteDialogValue.value = !deleteDialogValue.value
}

const onUpdatePage = (page: number): void => {
  emit('click:update-page', page)
}

const onUpdateItemsPerPage = (page: number): void => {
  emit('click:update-items-per-page', page)
}

const onClickShow = (productId: string): void => {
  emit('click:show', productId)
}

const onClickNew = (): void => {
  emit('click:new')
}

const onClickDelete = (): void => {
  emit('click:delete', selectedItem?.value?.id || '')
}

const onClickCopyItem = (): void => {
  emit('click:copyItem')
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
      <v-card-text class="text-h7">
        {{ selectedItem?.name }}を本当に削除しますか？
      </v-card-text>
      <v-card-actions>
        <v-spacer />
        <v-btn
          color="error"
          variant="text"
          @click="toggleDeleteDialog"
        >
          キャンセル
        </v-btn>
        <v-btn
          :loading="loading"
          color="primary"
          variant="outlined"
          @click="onClickDelete"
        >
          削除
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>

  <v-card
    class="mt-4"
    flat
  >
    <v-card-title class="d-flex flex-row">
      体験管理
    </v-card-title>

    <div class="d-flex w-100 px-6 ga-2">
      <v-spacer />
      <v-btn
        v-show="isRegisterable()"
        variant="outlined"
        color="primary"
        @click="onClickNew"
      >
        <v-icon
          start
          :icon="mdiPlus"
        />
        体験登録
      </v-btn>
    </div>

    <v-card-text>
      <v-data-table-server
        :headers="headers"
        :loading="loading"
        :items="props.products"
        :items-per-page="props.tableItemsPerPage"
        :items-length="props.tableItemsTotal"
        hover
        select-strategy="single"
        show-select
        no-data-text="登録されている商品がありません。"
        @update:model-value="handleUpdateSelectItemId"
        @update:page="onUpdatePage"
        @update:items-per-page="onUpdateItemsPerPage"
        @click:row="(_: any, { item }: any) => onClickShow(item.id)"
      >
        <template #[`item.media`]="{ item }">
          <v-img
            aspect-ratio="1/1"
            :max-height="56"
            :max-width="80"
          />
        </template>
        <template #[`item.actions`]="{ item }">
          <v-btn
            v-show="isDeletable()"
            variant="outlined"
            color="primary"
            size="small"
            @click.stop="toggleDeleteDialog(item)"
          >
            <v-icon
              size="small"
              :icon="mdiDelete"
            />
            削除
          </v-btn>
        </template>
      </v-data-table-server>
    </v-card-text>
  </v-card>
</template>
