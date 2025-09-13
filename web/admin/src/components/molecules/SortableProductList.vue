<script setup lang="ts">
import { useSortable } from '@vueuse/integrations/useSortable'
import { mdiDrag } from '@mdi/js'
import { getResizedImages } from '~/lib/helpers'
import type { Product, ProductMedia } from '~/types/api/v1'

const model = defineModel<any>()

interface Props {
  products: Product[]
}

defineProps<Props>()

const sortableRef = ref<HTMLElement | null>(null)
useSortable(sortableRef, model)

const getProductThumbnailUrl = (product: Product): string => {
  const thumbnail = product.media?.find((media: ProductMedia) => {
    return media.isThumbnail
  })
  return thumbnail ? thumbnail.url : ''
}

const getProductThumbnails = (product: Product): string => {
  const thumbnail = product.media?.find((media: ProductMedia) => {
    return media.isThumbnail
  })
  return thumbnail ? getResizedImages(thumbnail.url) : ''
}

const getProductInventoryColor = (product: Product): string => {
  return product.inventory > 0 ? '' : 'text-error'
}
</script>

<template>
  <v-table>
    <thead>
      <tr>
        <th width="5" />
        <th />
        <th>商品名</th>
        <th>価格</th>
        <th>在庫</th>
      </tr>
    </thead>

    <tbody ref="sortableRef">
      <tr
        v-for="product in products"
        :key="product.id"
      >
        <td>
          <v-icon :icon="mdiDrag" />
        </td>
        <td>
          <v-img
            aspect-ratio="1/1"
            :max-height="56"
            :max-width="80"
            :src="getProductThumbnailUrl(product)"
            :srcset="getProductThumbnails(product)"
          />
        </td>
        <td>{{ product.name }}</td>
        <td>{{ product.price }}</td>
        <td :class="getProductInventoryColor(product)">
          {{ product.inventory }}
        </td>
      </tr>
    </tbody>
  </v-table>
</template>
