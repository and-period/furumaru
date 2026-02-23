<script setup lang="ts">
import type { Product } from '~/types/api/v1'
import { getProductThumbnailUrl } from '~/lib/formatter'

interface Props {
  item: Product
}

defineProps<Props>()

interface Emits {
  (e: 'click:delete'): void
}

const emits = defineEmits<Emits>()

const handleClickDeleteButton = () => {
  emits('click:delete')
}
</script>

<template>
  <td>
    <v-img
      aspect-ratio="1/1"
      :max-height="56"
      :max-width="80"
      :src="getProductThumbnailUrl(item)"
      :alt="item.name || '商品画像'"
    />
  </td>
  <td>{{ item.name }}</td>
  <td>{{ item.price }}</td>
  <td>
    {{ item.inventory }}
  </td>
  <td width="36">
    <v-btn
      variant="outlined"
      color="error"
      size="small"
      @click="handleClickDeleteButton"
    >
      削除
    </v-btn>
  </td>
</template>
