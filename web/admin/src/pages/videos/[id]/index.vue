<script setup lang="ts">
import { useVideoStore, useProductStore, useExperienceStore, useVideoCommentStore, useCommonStore } from '~/store'
import { useUnsavedChangesGuard } from '~/composables/useUnsavedChangesGuard'
import type { UpdateVideoRequest, Product, Experience } from '~/types/api/v1'
import { ApiBaseError } from '~/types/exception'
import { useAlert } from '~/lib/hooks'
import dayjs from 'dayjs'

const route = useRoute()
const router = useRouter()
const commonStore = useCommonStore()
const { alertType, isShow, alertText, show: showError } = useAlert('error')

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

const { captureSnapshot, markAsSaved, showLeaveDialog, confirmLeave, cancelLeave }
  = useUnsavedChangesGuard(formData)

const fetchState = useAsyncData('video-detail', async (): Promise<void> => {
  try {
    await Promise.all([
      videoStore.fetchVideo(videoId),
      videoCommentStore.fetchAllComments(videoId),
    ])

    const endAt = dayjs().unix()
    let startAt = dayjs().subtract(3, 'month').unix()
    if (video.value?.publishedAt && video.value.publishedAt > startAt) {
      startAt = video.value.publishedAt
    }

    await videoStore.analyzeVideo(videoId, startAt, endAt)

    formData.value = { ...formData.value, ...video.value }
    selectedProducts.value = sproducts.value
    selectedExperiences.value = sexperiences.value
    captureSnapshot()
  }
  catch (err) {
    if (err instanceof Error) {
      showError(err.message)
    }
    console.log(err)
  }
})

const isLoading = (): boolean => {
  return fetchState?.pending?.value || false
}

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
const linkTargetExperienceIds = ref<string[]>([])

const productSearchLoading = ref<boolean>(false)
const experienceSearchLoading = ref<boolean>(false)

const handleSearchProduct = async (name: string): Promise<void> => {
  try {
    productSearchLoading.value = true
    await productStore.searchProducts(name, undefined, linkTargetProductIds.value)
  }
  catch (err) {
    if (err instanceof Error) {
      showError(err.message)
    }
    console.log(err)
  }
  finally {
    productSearchLoading.value = false
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

const handleSearchExperience = async (title: string): Promise<void> => {
  try {
    experienceSearchLoading.value = true
    await experienceStore.searchExperiences(title)
  }
  catch (err) {
    if (err instanceof Error) {
      showError(err.message)
    }
    console.log(err)
  }
  finally {
    experienceSearchLoading.value = false
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

const isOpenLinkExperienceDialog = ref<boolean>(false)

// 初期データ取得
handleSearchProduct('')
handleSearchExperience('')

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

const handleSubmit = async (): Promise<void> => {
  try {
    await videoStore.updateVideo(videoId, formData.value)
    commonStore.addSnackbar({
      message: `${formData.value.title}を更新しました。`,
      color: 'info',
    })
    markAsSaved()
    router.push('/videos')
  }
  catch (err) {
    if (err instanceof Error) {
      showError(err.message)
    }
    window.scrollTo({
      top: 0,
      behavior: 'smooth',
    })
    console.log(err)
  }
}

const handleClickBack = (): void => {
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

try {
  await fetchState.execute()
}
catch (err) {
  console.log('failed to setup', err)
}
</script>

<template>
  <div>
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
      :loading="isLoading()"
      :viewer-logs="viewerLogs"
      :video-upload-status="uploadVideoStatus"
      :thumbnail-upload-status="uploadThumbnailStatus"
      :product-search-loading="productSearchLoading"
      :experience-search-loading="experienceSearchLoading"
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
      @click:back="handleClickBack"
      @update:search-product="handleSearchProduct"
      @update:search-experience="handleSearchExperience"
      @click:ban-comment="handleBanComment"
    />

    <atoms-app-confirm-dialog
      v-model="showLeaveDialog"
      title="未保存の変更があります"
      message="ページを離れると入力内容が失われます。よろしいですか？"
      confirm-text="破棄して離れる"
      confirm-color="warning"
      @confirm="confirmLeave"
      @cancel="cancelLeave"
    />
  </div>
</template>
