<script lang="ts" setup>
import { mdiDelete, mdiPlus } from '@mdi/js'
import { VDataTable } from 'vuetify/lib/labs/components'

import { getResizedImages } from '~/lib/helpers'
import { AlertType } from '~/lib/hooks'
import { ProductsResponseProductsInner, ProductsResponseProductsInnerMediaInner, ImageSize, ProductsResponseProductsInnerMediaInnerImagesInner } from '~/types/api'

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
  products: {
    type: Array<ProductsResponseProductsInner>,
    default: () => []
  },
  tableItemsTotal: {
    type: Number,
    default: 0
  },
  tableItemsPerPage: {
    type: Number,
    default: 20
  }
})

const emit = defineEmits<{
  (e: 'click:update-page', page: number): void
  (e: 'click:update-items-per-page', page: number): void
  (e: 'click:show', productId: string): void
  (e: 'click:new'): void
  (e: 'click:delete', productId: string): void
}>()

const headers: VDataTable['headers'] = [
  {
    title: '',
    key: 'media'
  },
  {
    title: '商品名',
    key: 'name'
  },
  {
    title: 'ステータス',
    key: 'public'
  },
  {
    title: '価格',
    key: 'price'
  },
  {
    title: '在庫',
    key: 'inventory'
  },
  {
    title: 'ジャンル',
    key: 'categoryName'
  },
  {
    title: '品目',
    key: 'productTypeName'
  },
  {
    title: '農園名',
    key: 'storeName'
  },
  {
    title: 'Action',
    key: 'actions',
    sortable: false
  }
]

const deleteDialog = ref<boolean>(false)
const deleteProduct = ref<ProductsResponseProductsInner>()

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

const getPublished = (published: boolean): string => {
  return published ? '公開' : '非公開'
}

const getPublishedColor = (published: boolean): string => {
  return published ? 'primary' : 'warning'
}

const toggleDeleteDialog = (product?: ProductsResponseProductsInner): void => {
  if (product) {
    deleteProduct.value = product
  }
  deleteDialog.value = !deleteDialog.value
}

const onUpdatePage = (page: number): void => {
  emit('click:update-page', page)
}

const onUpdateItemsPerPage = (page: number): void => {
  emit('click:update-items-per-page', page)
}

const onClickShow = (_: Event, { item }: any): void => {
  const product = item.raw as ProductsResponseProductsInner
  emit('click:show', product.id)
}

const onClickNew = (): void => {
  emit('click:new')
}

const onClickDelete = (): void => {
  emit('click:delete', deleteProduct?.value?.id || '')
  deleteDialog.value = false
}
</script>

<template>
  <v-alert v-show="props.isAlert" :type="props.alertType" v-text="props.alertText" />
  <v-dialog v-model="deleteDialog" width="500">
    <v-card>
      <v-card-text class="text-h7">
        {{ deleteProduct?.name }}を本当に削除しますか？
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
        :items-length="props.tableItemsTotal"
        :items-per-page="props.tableItemsPerPage"
        no-data-text="登録されている商品がありません。"
        hover
        @update:page="onUpdatePage"
        @update:items-per-page="onUpdateItemsPerPage"
        @click:row="onClickShow"
      >
        <template #[`item.media`]="{ item }">
          <v-img aspect-ratio="1/1" :src="getThumbnail(item.raw.media)" :srcset="getResizedThumbnails(item.raw.media)" />
        </template>
        <template #[`item.public`]="{ item }">
          <v-chip :color="getPublishedColor(item.raw.public)">
            {{ getPublished(item.raw.public) }}
          </v-chip>
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
