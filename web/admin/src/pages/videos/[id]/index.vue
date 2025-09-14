<script setup lang="ts">
import { useVideoStore, useProductStore, useExperienceStore, useVideoCommentStore } from '~/store'
import type { UpdateVideoRequest, Product, Experience } from '~/types/api/v1'
import { ApiBaseError } from '~/types/exception'
import { useAlert } from '~/lib/hooks'

const route = useRoute()
const router = useRouter()
const { alertType, isShow, alertText } = useAlert('error')

const videoStore = useVideoStore()
const { video, viewerLogs, products: sproducts, experiences: sexperiences } = storeToRefs(videoStore)

const videoCommentStore = useVideoCommentStore()
const { comments } = storeToRefs(videoCommentStore)

const productStore = useProductStore()
const { products } = storeToRefs(productStore)

const experienceStore = useExperienceStore()
const { experiences } = storeToRefs(experienceStore)

const videoId = route.params.id as string

const formData = ref<UpdateVideoRequest>({
  title: '',
  coordinatorId: '',
  description: '',
  productIds: [],
  experienceIds: [],
  thumbnailUrl: '',
  videoUrl: '',
  _public: false,
  limited: false,
  publishedAt: 0,
  displayProduct: false,
  displayExperience: false,
})

const selectedProducts = ref<Product[]>([])
const selectedExperiences = ref<Experience[]>([])

const { status } = useAsyncData(async () => {
  await Promise.all([
    videoStore.fetchVideo(videoId),
    videoStore.analyzeVideo(videoId),
    videoCommentStore.fetchAllComments(videoId),
  ])

  formData.value = { ...formData.value, ...video.value }
  selectedProducts.value = sproducts.value
  selectedExperiences.value = sexperiences.value
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

const isOpenLinkProductDialog = ref<boolean>(false)

const handleClickLinkProductButton = () => {
  isOpenLinkProductDialog.value = true
}

const handleClickLinkExperienceButton = () => {
  isOpenLinkExperienceDialog.value = true
}

const linkTargetProductIds = ref<string[]>([])

/**
 * 商品検索ステータス
 */
const productSearchStatus = ref<{
  isLoading: boolean
  hasError: boolean
  errorMessage: string
}>({
  isLoading: false,
  hasError: false,
  errorMessage: '',
})

/**
 * 商品検索関数
 * @param text
 */
const searchProducts = (text: string) => {
  try {
    productSearchStatus.value.isLoading = true
    productStore.searchProducts(text, undefined, linkTargetProductIds.value)
  }
  catch (error) {
    productSearchStatus.value.hasError = true
  }
  finally {
    productSearchStatus.value.isLoading = false
  }
}

/**
 * 紐づけ対象の商品をフォームに加える処理
 */
const handleClickLinkProductAddButton = () => {
  const targetProducts = products.value.filter((product) => {
    // 未選択の商品でかつ選択対象の商品IDリストに含まれる商品を対象とする
    return (
      !formData.value.productIds.includes(product.id)
      && linkTargetProductIds.value.includes(product.id)
    )
  })

  selectedProducts.value.push(...targetProducts)
  formData.value.productIds.push(
    ...targetProducts.map(product => product.id),
  )
  isOpenLinkProductDialog.value = false
}

searchProducts('')

const isOpenLinkExperienceDialog = ref<boolean>(false)
const linkTargetExperienceIds = ref<string[]>([])

/**
 * 体験検索ステータス
 */
const experienceSearchStatus = ref<{
  isLoading: boolean
  hasError: boolean
  errorMessage: string
}>({
  isLoading: false,
  hasError: false,
  errorMessage: '',
})

/**
 * 体験検索関数
 * @param text
 */
const searchExperiences = (text: string) => {
  try {
    experienceSearchStatus.value.isLoading = true
    experienceStore.searchExperiences(text)
  }
  catch (error) {
    experienceSearchStatus.value.hasError = true
  }
  finally {
    experienceSearchStatus.value.isLoading = false
  }
}

/**
 * 紐づけ対象の体験をフォームに加える処理
 */
const handleClickLinkExperienceAddButton = () => {
  // 未選択の体験でかつ選択対象の体験IDリストに含まれる体験を対象とする
  const targetExperiences = experiences.value
    .filter((experience) => {
      return linkTargetExperienceIds.value.includes(experience.id)
    })
    .filter((experience) => {
      return !formData.value.experienceIds.includes(experience.id)
    })

  selectedExperiences.value.push(...targetExperiences)
  formData.value.experienceIds.push(
    ...targetExperiences.map(experience => experience.id),
  )
  isOpenLinkExperienceDialog.value = false
}

searchExperiences('')

/**
 * 紐づけ済みの商品を削除する処理
 * @param productId
 */
const handleDeleteLinkedProduct = (productId: string) => {
  formData.value.productIds = formData.value.productIds.filter(
    id => id !== productId,
  )
  linkTargetProductIds.value = linkTargetProductIds.value.filter(
    id => id !== productId,
  )
  selectedProducts.value = selectedProducts.value.filter(
    product => product.id !== productId,
  )
}

/**
 * 紐づけ済みの体験を削除する処理
 */
const handleDeleteLinkedExperience = (experienceId: string) => {
  formData.value.experienceIds = formData.value.experienceIds.filter(
    id => id !== experienceId,
  )
  linkTargetExperienceIds.value = linkTargetExperienceIds.value.filter(
    id => id !== experienceId,
  )
  selectedExperiences.value = selectedExperiences.value.filter(
    experience => experience.id !== experienceId,
  )
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

const currentTab = ref<number>(0)

const handleBanComment = async (commentId: string) => {
  try {
    await videoCommentStore.disableComment(videoId, commentId, true)
  }
  catch (err) {
    if (err instanceof Error) {
      showError(err.message)
    }
    console.log(err)
  }
}
</script>

<template>
  <templates-video-edit
    v-model:video-form-data="formData"
    v-model:product-form-data="linkTargetProductIds"
    v-model:experience-form-data="linkTargetExperienceIds"
    v-model:product-dialog="isOpenLinkProductDialog"
    v-model:experience-dialog="isOpenLinkExperienceDialog"
    v-model:current-tab="currentTab"
    :is-alert="isShow"
    :alert-type="alertType"
    :alert-text="alertText"
    :loading="isInitLoading"
    :viewer-logs="viewerLogs"
    :video-upload-status="uploadVideoStatus"
    :thumbnail-upload-status="uploadThumbnailStatus"
    :product-search-loading="productSearchStatus.isLoading"
    :experience-search-loading="experienceSearchStatus.isLoading"
    :selected-products="selectedProducts"
    :selected-experiences="selectedExperiences"
    :searched-products="products"
    :searched-experiences="experiences"
    :comments="comments"
    @submit="handleSubmit"
    @submit:upload-video="handleUploadVideo"
    @submit:upload-thumbnail="handleUploadThumbnail"
    @click:add-products="handleClickLinkProductAddButton"
    @click:add-experiences="handleClickLinkExperienceAddButton"
    @click:remove-product="handleDeleteLinkedProduct"
    @click:remove-experience="handleDeleteLinkedExperience"
    @click:back="handleClickBackButton"
    @update:search-product="searchProducts"
    @update:search-experience="searchExperiences"
    @click:ban-comment="handleBanComment"
  />
</template>
