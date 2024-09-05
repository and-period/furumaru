<script setup lang="ts">
import type { Product } from '~/types/api'
import { getProductThumbnailUrl } from '~/lib/formatter'

interface Props {
  item: Product
}

defineProps<Props>()
</script>

<template>
  <td>
    <v-img
      aspect-ratio="1/1"
      :max-height="56"
      :max-width="80"
      :src="getProductThumbnailUrl(item)"
    />
  </td>
  <td>{{ item.name }}</td>
  <td>{{ item.price }}</td>
  <td>
    {{ item.inventory }}
  </td>
</template>
