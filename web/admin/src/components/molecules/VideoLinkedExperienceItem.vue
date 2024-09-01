<script setup lang="ts">
import type { Experience } from '~/types/api'

interface Props {
  item: Experience
}

const props = defineProps<Props>()
</script>

<template>
  <div class="d-flex align-center ga-4">
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
    <div class="flex-grow-1 d-flex flex-column ga-3">
      <p class="text-subtitle-1">
        {{ item.title }}
      </p>
      <div class="text-gray">
        {{ item.description }}
      </div>
      <div>
        {{ item.hostCity }}
      </div>
    </div>
  </div>
</template>
