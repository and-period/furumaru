<script setup lang="ts">
import type { Experience } from '~/types/api/v1'

interface Props {
  item: Experience
}

const props = defineProps<Props>()

interface Emits {
  (e: 'click:delete', id: string): void
}

const emits = defineEmits<Emits>()

const handleClickDeleteButton = () => {
  emits('click:delete')
}
</script>

<template>
  <v-row>
    <v-col cols="4">
      <v-responsive
        :aspect-ratio="16 / 9"
        class="border"
        :width="250"
      >
        <video
          :src="item.promotionVideoUrl"
          controls
          class="w-100"
        />
      </v-responsive>
    </v-col>
    <v-col
      cols="7"
      class="flex-grow-1 d-flex flex-column ga-3"
    >
      <p class="text-subtitle-1">
        {{ item.title }}
      </p>
      <div class="text-gray">
        {{ item.description }}
      </div>
      <div>
        {{ item.hostCity }}
      </div>
    </v-col>
    <v-col
      cols="1"
      class="d-flex align-center"
    >
      <v-btn
        v-btn
        color="error"
        variant="outlined"
        size="small"
        @click="handleClickDeleteButton"
      >
        削除
      </v-btn>
    </v-col>
  </v-row>
</template>
