<script lang="ts" setup>
import { mdiPlus, mdiMagnify } from '@mdi/js'
import { VDataTable } from 'vuetify/lib/labs/components'

import { usePagination } from '~/lib/hooks/'
import { useProductStore } from '~/store'
import { ImageSize, ProductsResponseProductsInner, ProductsResponseProductsInnerMediaInner } from '~/types/api'

const router = useRouter()
const productStore = useProductStore()
const products = computed(() => productStore.products)
const totalItems = computed(() => productStore.totalItems)

const {
  updateCurrentPage,
  itemsPerPage,
  handleUpdateItemsPerPage,
  options,
  offset
} = usePagination()

watch(itemsPerPage, () => {
  productStore.fetchProducts(itemsPerPage.value, 0)
})

const handleUpdatePage = async (page: number) => {
  updateCurrentPage(page)
  await productStore.fetchProducts(itemsPerPage.value, offset.value)
}

const fetchState = useAsyncData(async () => {
  try {
    await productStore.fetchProducts(itemsPerPage.value, offset.value)
  } catch (error) {
    console.log(error)
  }
})

const searchWord = ref<string>('')

const handleRowClick = (
  _: any,
  { item }: { item: VDataTable['itemValue'] }
): void => {
  router.push(`/products/${item.raw.id}`)
}

const handleClickAddBtn = () => {
  router.push('/products/add')
}

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
    title: '種類',
    key: 'type'
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
  }
]

const isLoading = (): boolean => {
  return fetchState?.pending?.value || false
}

const getThumbnail = (product: ProductsResponseProductsInner): string => {
  const thumbnail = product.media?.find((media: ProductsResponseProductsInnerMediaInner) => {
    return media.isThumbnail
  })
  return thumbnail?.url || ''
}

const getResizedThumbnails = (product: ProductsResponseProductsInner): string => {
  const thumbnail = product.media?.find((media: ProductsResponseProductsInnerMediaInner) => {
    return media.isThumbnail
  })
  if (!thumbnail) {
    return ''
  }
  const images: string[] = thumbnail.images.map((image): string => {
    switch (image.size) {
      case ImageSize.SMALL:
        return `${thumbnail.url} 1x`
      case ImageSize.MEDIUM:
        return `${thumbnail.url} 2x`
      case ImageSize.LARGE:
        return `${thumbnail.url} 3x`
      default:
        return thumbnail.url
    }
  })
  return images.join(', ')
}

try {
  await fetchState.execute()
} catch (err) {
  console.log('failed to setup', err)
}
</script>

<template>
  <div>
    <v-card-title class="d-flex flex-row">
      商品管理
      <v-spacer />
      <v-btn variant="outlined" color="primary" @click="handleClickAddBtn">
        <v-icon start :icon="mdiPlus" />
        商品登録
      </v-btn>
    </v-card-title>

    <v-card class="mt-4" flat :loading="isLoading()">
      <v-card-text>
        <div class="d-flex align-center mb-4">
          <v-spacer />
          <v-text-field
            v-model="searchWord"
            :append-icon="mdiMagnify"
            label="商品名"
            hide-details
            single-line
            variant="underlined"
          />
        </div>

        <v-data-table-server
          v-model:items-per-page="itemsPerPage"
          :headers="headers"
          :items="products"
          no-data-text="登録されている商品がありません。"
          :items-length="totalItems"
          :footer-props="options"
          @update:items-per-page="handleUpdateItemsPerPage"
          @update:page="handleUpdatePage"
          @click:row="handleRowClick"
        >
          <template #[`item.media`]="{ item }">
            <v-avatar>
              <v-img :src="getThumbnail(item.raw)" :srcset="getResizedThumbnails(item.raw)" />
            </v-avatar>
          </template>
          <template #[`item.public`]="{ item }">
            <v-chip :color="item.raw.public ? 'primary' : 'warning'">
              {{ item.public ? '公開' : '非公開' }}
            </v-chip>
          </template>
        </v-data-table-server>
      </v-card-text>
    </v-card>
  </div>
</template>
