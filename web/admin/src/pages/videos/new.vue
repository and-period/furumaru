<script setup lang="ts">
import { useVideoStore } from '~/store'
import type { CreateVideoRequest, Product, Experience } from '~/types/api'

const router = useRouter()
const videoStore = useVideoStore()

const formData = ref<CreateVideoRequest>({
  title: '',
  description: '',
  coordinatorId: '',
  productIds: [],
  experienceIds: [],
  thumbnailUrl: '',
  videoUrl: '',
  public: false,
  limited: false,
  publishedAt: 0,
  displayProduct: false,
  displayExperience: false,
})

const selectedProducts = ref<Product[]>([])
const selectedExperiences = ref<Experience[]>([])

const handleClickBackButton = () => {
  router.back()
}
</script>

<template>
  <div>
    <v-card class="mb-16">
      <v-card-title>動画登録</v-card-title>
      <v-card-text>
        <organisms-video-form
          id="new-video-form"
          v-model="formData"
          :selected-products="selectedProducts"
          :selected-experiences="selectedExperiences"
        />
      </v-card-text>
    </v-card>
    <div
      class="position-fixed bottom-0 left-0 w-100 bg-white pa-4 text-right elevation-3"
    >
      <div class="d-inline-flex ga-4">
        <v-btn
          color="secondary"
          variant="outlined"
          @click="handleClickBackButton"
        >
          戻る
        </v-btn>
        <v-btn
          color="primary"
          variant="outlined"
        >
          保存
        </v-btn>
      </div>
    </div>
  </div>
</template>
