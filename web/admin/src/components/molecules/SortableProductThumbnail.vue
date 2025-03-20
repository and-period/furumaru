<script setup lang="ts">
import { useSortable } from '@vueuse/integrations/useSortable'
import { mdiClose } from '@mdi/js'
import type { ProductMediaInner } from '~/types/api'

const model = defineModel<ProductMediaInner[]>({ required: true })

interface Emits {
  (e: 'click', i: number): void
  (e: 'delete', i: number): void
}

const emits = defineEmits<Emits>()

const element = ref<HTMLDivElement | null>()

useSortable(element, model)

const handleClick = (i: number) => {
  emits('click', i)
}

const handleDelete = (i: number) => {
  emits('delete', i)
}
</script>

<template>
  <v-row ref="element">
    <v-col
      v-for="(img, i) in model"
      :key="img.url"
      cols="4"
    >
      <v-card
        rounded
        variant="outlined"
        width="100%"
        :class="{ 'thumbnail-border': img.isThumbnail }"
        @click="handleClick(i)"
      >
        <v-img
          :src="img.url"
          aspect-ratio="1"
        >
          <div class="d-flex col">
            <v-radio
              :value="i"
              color="primary"
            />
            <v-btn
              :icon="mdiClose"
              color="error"
              variant="text"
              size="small"
              @click="handleDelete(i)"
            />
          </div>
        </v-img>
      </v-card>
    </v-col>
  </v-row>
</template>
