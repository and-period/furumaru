<script setup lang="ts">
import { useVideoStore } from '~/store'
import type { UpdateVideoRequest, Product, Experience } from '~/types/api'

const route = useRoute()
const router = useRouter()
const videoStore = useVideoStore()

const videoId = route.params.id as string

const formData = ref<UpdateVideoRequest>({
  title: '',
  description: '',
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

const { status } = useAsyncData(async () => {
  const res = await videoStore.fetchVideo(videoId)
  formData.value.title = res.video.title
  formData.value.public = res.video.public
  formData.value.publishedAt = res.video.publishedAt
  formData.value.description = res.video.description
  formData.value.videoUrl = res.video.videoUrl
  formData.value.thumbnailUrl = res.video.thumbnailUrl
  formData.value.productIds = res.video.productIds
  formData.value.experienceIds = res.video.experienceIds

  selectedProducts.value = res.products
  selectedExperiences.value = res.experiences
})

const isInitLoading = computed<boolean>(() => status.value === 'pending')

const handleClickBackButton = () => {
  router.back()
}
</script>

<template>
  <div>
    <v-card
      class="mb-16"
      :loading="isInitLoading"
    >
      <v-card-title>動画編集</v-card-title>
      <v-card-text>
        <template v-if="!isInitLoading">
          <organisms-video-form
            id="update-video-form"
            v-model="formData"
            :selected-products="selectedProducts"
            :selected-experiences="selectedExperiences"
          />
        </template>
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
