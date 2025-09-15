<script lang="ts" setup>
import dayjs from 'dayjs'
import { storeToRefs } from 'pinia'
import { useAlert } from '~/lib/hooks'
import { useAuthStore, useCommonStore, useVideoStore, useProductStore, useExperienceStore } from '~/store'
import type { CreateVideoRequest, Product, Experience } from '~/types/api/v1'
import type { ImageUploadStatus } from '~/types/props'

const router = useRouter()
const commonStore = useCommonStore()
const authStore = useAuthStore()
const videoStore = useVideoStore()
const productStore = useProductStore()
const experienceStore = useExperienceStore()
const { alertType, isShow, alertText, show } = useAlert('error')

const { adminId } = storeToRefs(authStore)
const { products } = storeToRefs(productStore)
const { experiences } = storeToRefs(experienceStore)

const loading = ref<boolean>(false)
const formData = ref<CreateVideoRequest>({
  title: '',
  description: '',
  coordinatorId: '',
  productIds: [],
  experienceIds: [],
  thumbnailUrl: '',
  videoUrl: '',
  _public: false,
  limited: false,
  publishedAt: dayjs().unix(),
  displayProduct: false,
  displayExperience: false,
})

const selectedProducts = ref<Product[]>([])
const selectedExperiences = ref<Experience[]>([])

const videoUploadStatus = ref<ImageUploadStatus>({
  error: false,
  message: '',
})
const thumbnailUploadStatus = ref<ImageUploadStatus>({
  error: false,
  message: '',
})

const isLoading = (): boolean => {
  return loading.value
}

const handleUploadVideo = (file: File): void => {
  loading.value = true
  videoStore.uploadVideoFile(file)
    .then((url: string) => {
      formData.value.videoUrl = url
    })
    .catch(() => {
      videoUploadStatus.value.error = true
      videoUploadStatus.value.message = '動画のアップロードに失敗しました。'
    })
    .finally(() => {
      loading.value = false
    })
}

const handleUploadThumbnail = (file: File): void => {
  loading.value = true
  videoStore.uploadThumbnailFile(file)
    .then((url: string) => {
      formData.value.thumbnailUrl = url
    })
    .catch(() => {
      thumbnailUploadStatus.value.error = true
      thumbnailUploadStatus.value.message = 'サムネイルのアップロードに失敗しました。'
    })
    .finally(() => {
      loading.value = false
    })
}

const productDialog = ref<boolean>(false)
const experienceDialog = ref<boolean>(false)
const productFormData = ref<string[]>([])
const experienceFormData = ref<string[]>([])

const productSearchLoading = ref<boolean>(false)
const experienceSearchLoading = ref<boolean>(false)

const handleSearchProduct = async (name: string): Promise<void> => {
  try {
    productSearchLoading.value = true
    await productStore.searchProducts(name, undefined, productFormData.value)
  }
  catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
  finally {
    productSearchLoading.value = false
  }
}

const handleSearchExperience = async (title: string): Promise<void> => {
  try {
    experienceSearchLoading.value = true
    await experienceStore.searchExperiences(title)
  }
  catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
  finally {
    experienceSearchLoading.value = false
  }
}

const handleClickAddProducts = (): void => {
  const targetProducts = products.value.filter((product) => {
    return (
      !formData.value.productIds.includes(product.id)
      && productFormData.value.includes(product.id)
    )
  })

  selectedProducts.value.push(...targetProducts)
  formData.value.productIds.push(
    ...targetProducts.map(product => product.id),
  )
  productDialog.value = false
}

const handleClickAddExperiences = (): void => {
  const targetExperiences = experiences.value.filter((experience) => {
    return (
      !formData.value.experienceIds.includes(experience.id)
      && experienceFormData.value.includes(experience.id)
    )
  })

  selectedExperiences.value.push(...targetExperiences)
  formData.value.experienceIds.push(
    ...targetExperiences.map(experience => experience.id),
  )
  experienceDialog.value = false
}

const handleClickRemoveProduct = (productId: string): void => {
  formData.value.productIds = formData.value.productIds.filter(
    id => id !== productId,
  )
  productFormData.value = productFormData.value.filter(
    id => id !== productId,
  )
  selectedProducts.value = selectedProducts.value.filter(
    product => product.id !== productId,
  )
}

const handleClickRemoveExperience = (experienceId: string): void => {
  formData.value.experienceIds = formData.value.experienceIds.filter(
    id => id !== experienceId,
  )
  experienceFormData.value = experienceFormData.value.filter(
    id => id !== experienceId,
  )
  selectedExperiences.value = selectedExperiences.value.filter(
    experience => experience.id !== experienceId,
  )
}

const handleSubmit = async (): Promise<void> => {
  try {
    loading.value = true
    const req: CreateVideoRequest = {
      ...formData.value,
      coordinatorId: adminId.value,
    }
    await videoStore.createVideo(req)
    commonStore.addSnackbar({
      message: `${formData.value.title}を作成しました。`,
      color: 'info',
    })
    router.push('/videos')
  }
  catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    window.scrollTo({
      top: 0,
      behavior: 'smooth',
    })
    console.log(err)
  }
  finally {
    loading.value = false
  }
}

const handleClickBack = (): void => {
  router.back()
}

// 初期データ取得
handleSearchProduct('')
handleSearchExperience('')
</script>

<template>
  <templates-video-new
    v-model:form-data="formData"
    v-model:product-dialog="productDialog"
    v-model:experience-dialog="experienceDialog"
    v-model:product-form-data="productFormData"
    v-model:experience-form-data="experienceFormData"
    :loading="isLoading()"
    :is-alert="isShow"
    :alert-type="alertType"
    :alert-text="alertText"
    :video-upload-status="videoUploadStatus"
    :thumbnail-upload-status="thumbnailUploadStatus"
    :selected-products="selectedProducts"
    :selected-experiences="selectedExperiences"
    :searched-products="products"
    :searched-experiences="experiences"
    :product-search-loading="productSearchLoading"
    :experience-search-loading="experienceSearchLoading"
    @update:video="handleUploadVideo"
    @update:thumbnail="handleUploadThumbnail"
    @click:add-products="handleClickAddProducts"
    @click:add-experiences="handleClickAddExperiences"
    @click:remove-product="handleClickRemoveProduct"
    @click:remove-experience="handleClickRemoveExperience"
    @update:search-product="handleSearchProduct"
    @update:search-experience="handleSearchExperience"
    @submit="handleSubmit"
    @click:back="handleClickBack"
  />
</template>
