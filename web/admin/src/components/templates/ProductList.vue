<script lang="ts" setup>
import { mdiDelete, mdiPlus } from '@mdi/js'
import { VDataTable } from 'vuetify/lib/labs/components.mjs'
import { PrefecturesListItem, prefecturesList } from '~/constants'

import { getResizedImages } from '~/lib/helpers'
import { AlertType } from '~/lib/hooks'
import { ProductsResponseProductsInner, ProductsResponseProductsInnerMediaInner, ImageSize, ProductsResponseProductsInnerMediaInnerImagesInner, Prefecture, ProductStatus } from '~/types/api'

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
  products: {
    type: Array<ProductsResponseProductsInner>,
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
  (e: 'click:show', productId: string): void
  (e: 'click:new'): void
  (e: 'click:delete', productId: string): void
  (e: 'update:delete-dialog', v: boolean): void
}>()

const headers: VDataTable['headers'] = [
  {
    title: '',
    key: 'media',
    width: 80,
    sortable: false
  },
  {
    title: '商品名',
    key: 'name',
    sortable: false
  },
  {
    title: 'ステータス',
    key: 'status',
    sortable: false
  },
  {
    title: '価格',
    key: 'price',
    sortable: false
  },
  {
    title: '在庫',
    key: 'inventory',
    sortable: false
  },
  {
    title: 'ジャンル',
    key: 'categoryName',
    sortable: false
  },
  {
    title: '品目',
    key: 'productTypeName',
    sortable: false
  },
  {
    title: '原産地',
    key: 'originPrefecture',
    sortable: false
  },
  {
    title: '生産者名',
    key: 'producerName',
    sortable: false
  },
  {
    title: '',
    key: 'actions',
    sortable: false
  }
]

const selectedItem = ref<ProductsResponseProductsInner>()

const deleteDialogValue = computed({
  get: (): boolean => props.deleteDialog,
  set: (val: boolean): void => emit('update:delete-dialog', val)
})

const getThumbnail = (media: ProductsResponseProductsInnerMediaInner[]): string => {
  const thumbnail = media.find((media: ProductsResponseProductsInnerMediaInner) => {
    return media.isThumbnail
  })
  return thumbnail?.url || ''
}

const getResizedThumbnails = (media: ProductsResponseProductsInnerMediaInner[]): string => {
  const thumbnail = media.find((media: ProductsResponseProductsInnerMediaInner) => {
    return media.isThumbnail
  })
  if (!thumbnail) {
    return ''
  }
  return getResizedImages(thumbnail.images)
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
    default:
      return ''
  }
}

const getInventoryColor = (inventory: number): string => {
  return inventory > 0 ? '' : 'text-error'
}

const getPrefecture = (prefecture: Prefecture): string => {
  const pref = prefecturesList.find((val: PrefecturesListItem): boolean => prefecture === val.value)
  return pref?.text || ''
}

const toggleDeleteDialog = (product?: ProductsResponseProductsInner): void => {
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
</script>

<template>
  <v-alert v-show="props.isAlert" :type="props.alertType" v-text="props.alertText" />

  <v-dialog v-model="deleteDialogValue" width="500">
    <v-card>
      <v-card-text class="text-h7">
        {{ selectedItem?.name }}を本当に削除しますか？
      </v-card-text>
      <v-card-actions>
        <v-spacer />
        <v-btn color="error" variant="text" @click="toggleDeleteDialog">
          キャンセル
        </v-btn>
        <v-btn color="primary" variant="outlined" @click="onClickDelete">
          削除
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>

  <v-card class="mt-4" flat :loading="props.loading">
    <v-card-title class="d-flex flex-row">
      商品管理
      <v-spacer />
      <v-btn variant="outlined" color="primary" @click="onClickNew">
        <v-icon start :icon="mdiPlus" />
        商品登録
      </v-btn>
    </v-card-title>

    <v-card-text>
      <v-data-table-server
        :headers="headers"
        :items="props.products"
        :items-per-page="props.tableItemsPerPage"
        :items-length="props.tableItemsTotal"
        hover
        no-data-text="登録されている商品がありません。"
        @update:page="onUpdatePage"
        @update:items-per-page="onUpdateItemsPerPage"
        @click:row="(_: any, { item }:any) => onClickShow(item.raw.id)"
      >
        <template #[`item.media`]="{ item }">
          <v-img aspect-ratio="1/1" :max-height="56" :max-width="80" :src="getThumbnail(item.raw.media)" :srcset="getResizedThumbnails(item.raw.media)" />
        </template>
        <template #[`item.status`]="{ item }">
          <v-chip :color="getStatusColor(item.raw.status)">
            {{ getStatus(item.raw.status) }}
          </v-chip>
        </template>
        <template #[`item.inventory`]="{ item }">
          <div :class="getInventoryColor(item.raw.inventory)">
            {{ item.raw.inventory }}
          </div>
        </template>
        <template #[`item.originPrefecture`]="{ item }">
          {{ getPrefecture(item.raw.originPrefecture) }}
        </template>
        <template #[`item.actions`]="{ item }">
          <v-btn variant="outlined" color="primary" size="small" @click.stop="toggleDeleteDialog(item.raw)">
            <v-icon size="small" :icon="mdiDelete" />
            削除
          </v-btn>
        </template>
      </v-data-table-server>
    </v-card-text>
  </v-card>
</template>
