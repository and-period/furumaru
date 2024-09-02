<script setup lang="ts">
import { useVideoStore } from '~/store'
import type { UpdateVideoRequest, Product, Experience } from '~/types/api'
import { ApiBaseError } from '~/types/exception'

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

/**
 * 動画ファイルアップロードステータス
 */
const uploadVideoStatus = ref<{
  isUploading: boolean
  hasError: boolean
  errorMessage: string
}>({
  isUploading: false,
  hasError: false,
  errorMessage: '',
})

/**
 * 動画ファイルアップロード関数
 */
const handleUploadVideo = async (file: File) => {
  try {
    uploadVideoStatus.value.isUploading = true
    const newUrl = await videoStore.uploadVideoFile(file)
    formData.value.videoUrl = newUrl
  }
  catch (error) {
    uploadVideoStatus.value.hasError = true
    if (error instanceof ApiBaseError) {
      uploadVideoStatus.value.errorMessage = error.message
    }
    else {
      uploadVideoStatus.value.errorMessage
        = '動画ファイルのアップロードに失敗しました。不明なエラーが発生しました。'
    }
  }
  finally {
    uploadVideoStatus.value.isUploading = false
  }
}

/**
 * サムネイル画像アップロードステータス
 */
const uploadThumbnailStatus = ref<{
  isUploading: boolean
  hasError: boolean
  errorMessage: string
}>({
  isUploading: false,
  hasError: false,
  errorMessage: '',
})

/**
 * サムネイル画像ファイルアップロード関数
 * @param file
 */
const handleUploadThumbnail = async (file: File) => {
  try {
    uploadThumbnailStatus.value.isUploading = true
    const newUrl = await videoStore.uploadThumbnailFile(file)
    formData.value.thumbnailUrl = newUrl
  }
  catch (error) {
    uploadThumbnailStatus.value.hasError = true
    if (error instanceof ApiBaseError) {
      uploadThumbnailStatus.value.errorMessage = error.message
    }
    else {
      uploadThumbnailStatus.value.errorMessage
        = 'サムネイル画像のアップロードに失敗しました。不明なエラーが発生しました。'
    }
  }
  finally {
    uploadThumbnailStatus.value.isUploading = false
  }
}

const handleSubmit = async () => {
  try {
    await videoStore.updateVideo(videoId, formData.value)
    router.push('/videos')
  }
  catch (error) {
    console.log('update error')
  }
}

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
            :video-is-uploading="uploadVideoStatus.isUploading"
            :video-has-error="uploadVideoStatus.hasError"
            :video-error-message="uploadVideoStatus.errorMessage"
            :thumbnail-is-uploading="uploadThumbnailStatus.isUploading"
            :thumbnail-has-error="uploadThumbnailStatus.hasError"
            :thumbnail-error-message="uploadThumbnailStatus.errorMessage"
            @update:video="handleUploadVideo"
            @update:thumbnail="handleUploadThumbnail"
            @submit="handleSubmit"
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
          type="submit"
          form="update-video-form"
        >
          保存
        </v-btn>
      </div>
    </div>
  </div>
</template>
