<script lang="ts" setup>
import { mdiDelete, mdiPlus, mdiContentCopy, mdiPackageVariant, mdiCoffee } from '@mdi/js'
import type { VDataTable } from 'vuetify/components'
import { productStatuses } from '~/constants'

import { getResizedImages } from '~/lib/helpers'
import type { AlertType } from '~/lib/hooks'
import { ProductStatus, AdminType } from '~/types/api/v1'
import type { Product, ProductMedia, Category, ProductTag, ProductType, Producer } from '~/types/api/v1'

const props = defineProps({
  selectedItemId: {
    type: String,
    default: '',
  },
  loading: {
    type: Boolean,
    default: false,
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
    title: '価格',
    key: 'price',
    sortable: false,
  },
  {
    title: '在庫',
    key: 'inventory',
    sortable: false,
  },
  {
    title: 'ジャンル',
    key: 'categoryName',
    sortable: false,
  },
  {
    title: '品目',
    key: 'productTypeName',
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
  return props.adminType === AdminType.AdminTypeCoordinator
}

const isDeletable = (): boolean => {
  const targets: AdminType[] = [AdminType.AdminTypeAdministrator, AdminType.AdminTypeCoordinator]
  return targets.includes(props.adminType)
}

const getCategoryName = (categoryId: string): string => {
  const category = props.categories.find((category: Category): boolean => {
    return category.id === categoryId
  })
  return category ? category.name : ''
}

const getProductTypeName = (productTypeId: string): string => {
  const productType = props.productTypes.find(
    (productType: ProductType): boolean => {
      return productType.id === productTypeId
    },
  )
  return productType ? productType.name : ''
}

const getProducerName = (producerId: string): string => {
  const producer = props.producers.find((producer: Producer): boolean => {
    return producer.id === producerId
  })
  return producer ? producer.username : ''
}

const getThumbnail = (media: ProductMedia[]): string => {
  const thumbnail = media.find((media: ProductMedia) => {
    return media.isThumbnail
  })
  return thumbnail?.url || ''
}

const getResizedThumbnails = (media: ProductMedia[]): string => {
  const thumbnail = media.find((media: ProductMedia) => {
    return media.isThumbnail
  })
  if (!thumbnail) {
    return ''
  }
  return getResizedImages(thumbnail.url)
}

const getStatus = (status: ProductStatus): string => {
  const value = productStatuses.find(s => s.value === status)
  return value ? value.title : '不明'
}

const getStatusColor = (status: ProductStatus): string => {
  switch (status) {
    case ProductStatus.ProductStatusPresale:
      return 'info'
    case ProductStatus.ProductStatusForSale:
      return 'primary'
    case ProductStatus.ProductStatusOutOfSale:
      return 'secondary'
    case ProductStatus.ProductStatusPrivate:
      return 'warning'
    case ProductStatus.ProductStatusArchived:
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
      <v-card-title class="text-h6 py-4">
        商品削除の確認
      </v-card-title>
      <v-card-text class="pb-4">
        <div class="text-body-1">
          「{{ selectedItem?.name }}」を削除しますか？
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
          @click="toggleDeleteDialog"
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
          :icon="mdiPackageVariant"
          size="28"
          class="mr-3 text-primary"
        />
        <div>
          <h1 class="text-h6 text-sm-h5 font-weight-bold text-primary">
            商品管理
          </h1>
          <p class="text-caption text-sm-body-2 text-medium-emphasis ma-0">
            商品の登録・編集・削除を行います
          </p>
        </div>
      </div>
      <div class="d-flex flex-column flex-sm-row ga-2 ga-sm-3 w-100 w-sm-auto">
        <v-btn
          v-show="isRegisterable()"
          variant="outlined"
          color="secondary"
          :size="$vuetify.display.smAndDown ? 'default' : 'large'"
          class="w-100 w-sm-auto"
          :disabled="selectedItemId === ''"
          @click="onClickCopyItem"
        >
          <v-icon
            start
            :icon="mdiContentCopy"
          />
          商品複製
          <v-tooltip
            v-if="selectedItemId === ''"
            activator="parent"
            location="bottom"
          >
            コピー元の商品をチェックする必要があります。
          </v-tooltip>
        </v-btn>
        <v-btn
          v-show="isRegisterable()"
          variant="elevated"
          color="primary"
          :size="$vuetify.display.smAndDown ? 'default' : 'large'"
          class="w-100 w-sm-auto"
          @click="onClickNew"
        >
          <v-icon
            start
            :icon="mdiPlus"
          />
          商品登録
        </v-btn>
      </div>
    </v-card-title>

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
          <v-avatar size="40">
            <v-img
              v-if="getThumbnail(item.media) !== ''"
              cover
              :src="getThumbnail(item.media)"
              :srcset="getResizedThumbnails(item.media)"
            />
            <v-icon
              v-else
              :icon="mdiCoffee"
              color="grey"
            />
          </v-avatar>
        </template>
        <template #[`item.name`]="{ item }">
          {{ item.name.length > 24 ? item.name.slice(0, 24) + '...' : item.name }}
        </template>
        <template #[`item.status`]="{ item }">
          <v-chip :color="getStatusColor(item.status)">
            {{ getStatus(item.status) }}
          </v-chip>
        </template>
        <template #[`item.price`]="{ item }">
          &yen; {{ item.price.toLocaleString() }}
        </template>
        <template #[`item.inventory`]="{ item }">
          <div :class="getInventoryColor(item.inventory)">
            {{ item.inventory.toLocaleString() }}
          </div>
        </template>
        <template #[`item.categoryName`]="{ item }">
          {{ getCategoryName(item.categoryId) }}
        </template>
        <template #[`item.productTypeName`]="{ item }">
          {{ getProductTypeName(item.productTypeId) }}
        </template>
        <template #[`item.producerName`]="{ item }">
          {{ getProducerName(item.producerId) }}
        </template>
        <template #[`item.actions`]="{ item }">
          <v-btn
            v-show="isDeletable()"
            variant="outlined"
            color="error"
            size="small"
            :prepend-icon="mdiDelete"
            @click.stop="toggleDeleteDialog(item)"
          >
            削除
          </v-btn>
        </template>
      </v-data-table-server>
    </v-card-text>
  </v-card>
</template>
