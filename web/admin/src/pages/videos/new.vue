<script setup lang="ts">
import { useVideoStore, useProductStore } from '~/store'
import type { CreateVideoRequest, Product, Experience } from '~/types/api'
import { ApiBaseError } from '~/types/exception'
import { getProductThumbnailUrl } from '~/lib/formatter'

const router = useRouter()
const videoStore = useVideoStore()

const productStore = useProductStore()
const { products } = storeToRefs(productStore)

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
  console.log('link product')
  isOpenLinkProductDialog.value = true
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

const handleClickLinkExperienceButton = () => {
  console.log('link experience')
  isOpenLinkExperienceDialog.value = true
}

const handleSubmit = async () => {
  try {
    await videoStore.createVideo(formData.value)
    router.push('/videos')
  }
  catch (error) {
    console.log('create error')
  }
}

const handleClickBackButton = () => {
  router.back()
}
</script>

<template>
  <div>
    <v-dialog v-model="isOpenLinkProductDialog">
      <v-card>
        <v-card-title>商品紐づけ</v-card-title>
        <v-card-text>
          <v-autocomplete
            v-model="linkTargetProductIds"
            :loading="productSearchStatus.isLoading"
            :items="products"
            label="商品名"
            messages="商品名を入力することで紐づける商品を検索できます。"
            item-title="name"
            item-value="id"
            multiple
            chips
            @update:search-text="searchProducts"
          >
            <template #chip="{ props: val, item }">
              <v-chip
                v-bind="val"
                :prepend-avatar="getProductThumbnailUrl(item.raw)"
                :text="item.raw.name"
                rounded
                class="px-4"
                variant="outlined"
              />
            </template>
            <template #item="{ props: val, item }">
              <v-list-item
                v-bind="val"
                :prepend-avatar="getProductThumbnailUrl(item.raw)"
                :title="item.raw.name"
              />
            </template>
          </v-autocomplete>
        </v-card-text>
        <v-card-actions>
          <v-btn
            color="primary"
            variant="outlined"
            @click="handleClickLinkProductAddButton"
          >
            追加
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <v-dialog v-model="isOpenLinkExperienceDialog">
      体験紐づけダイアログ
    </v-dialog>

    <v-card class="mb-16">
      <v-card-title>動画登録</v-card-title>
      <v-card-text>
        <organisms-video-form
          id="new-video-form"
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
          @click:link-product="handleClickLinkProductButton"
          @click:link-experience="handleClickLinkExperienceButton"
          @submit="handleSubmit"
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
          type="submit"
          form="new-video-form"
        >
          保存
        </v-btn>
      </div>
    </div>
  </div>
</template>
