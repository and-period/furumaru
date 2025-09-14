<script lang="ts" setup>
import type { AlertType } from '~/lib/hooks'
import type { Product, Experience, UpdateVideoRequest, VideoComment, VideoViewerLog } from '~/types/api/v1'
import { getProductThumbnailUrl, getExperienceThumbnailUrl, dateTimeFormatter } from '~/lib/formatter'
import type { UploadStatus } from '~/types/props'
import type { VTabs } from 'vuetify/components'

const props = defineProps({
  loading: {
    type: Boolean,
    default: false,
  },
  currentTab: {
    type: Number,
    default: 0,
  },
  productSearchLoading: {
    type: Boolean,
    default: false,
  },
  experienceSearchLoading: {
    type: Boolean,
    default: false,
  },
  isAlert: {
    type: Boolean,
    default: false,
  },
  alertType: {
    type: String as PropType<AlertType>,
    default: undefined,
  },
  alertText: {
    type: String,
    default: '',
  },
  productDialog: {
    type: Boolean,
    default: false,
  },
  experienceDialog: {
    type: Boolean,
    default: false,
  },
  videoFormData: {
    type: Object as PropType<UpdateVideoRequest>,
    default: (): UpdateVideoRequest => ({
      videoUrl: '',
      coordinatorId: '',
      description: '',
      displayExperience: false,
      displayProduct: false,
      experienceIds: [],
      limited: false,
      productIds: [],
      _public: false,
      publishedAt: 0,
      thumbnailUrl: '',
      title: '',
    }),
  },
  productFormData: {
    type: Array<string>,
    default: () => [],
  },
  experienceFormData: {
    type: Array<string>,
    default: () => [],
  },
  viewerLogs: {
    type: Array<VideoViewerLog>,
    default: () => [],
  },
  videoUploadStatus: {
    type: Object as PropType<UploadStatus>,
    default: () => ({
      isUploading: false,
      hasError: false,
      errorMessage: '',
    }),
  },
  thumbnailUploadStatus: {
    type: Object as PropType<UploadStatus>,
    default: () => ({
      isUploading: false,
      hasError: false,
      errorMessage: '',
    }),
  },
  selectedProducts: {
    type: Array<Product>,
    default: () => [],
  },
  selectedExperiences: {
    type: Array<Experience>,
    default: () => [],
  },
  searchedProducts: {
    type: Array<Product>,
    default: () => [],
  },
  searchedExperiences: {
    type: Array<Experience>,
    default: () => [],
  },
  comments: {
    type: Array<VideoComment>,
    default: () => [],
  },
})

const emit = defineEmits<{
  (e: 'update:video-form-data', formData: UpdateVideoRequest): void
  (e: 'update:product-form-data', products: string[]): void
  (e: 'update:experience-form-data', experiences: string[]): void
  (e: 'update:product-dialog', value: boolean): void
  (e: 'update:experience-dialog', value: boolean): void
  (e: 'update:search-product', name: string): void
  (e: 'update:search-experience', title: string): void
  (e: 'submit'): void
  (e: 'submit:upload-video', file: File): void
  (e: 'submit:upload-thumbnail', file: File): void
  (e: 'click:add-products'): void
  (e: 'click:add-experiences'): void
  (e: 'click:remove-product', productId: string): void
  (e: 'click:remove-experience', experienceId: string): void
  (e: 'click:back'): void
  (e: 'update:current-tab', value: number): void
  (e: 'click:ban-comment', commentId: string): void
}>()

const tabs: VTabs[] = [
  { name: '動画情報', value: 'info' },
  { name: 'プレビュー', value: 'preview' },
  { name: '分析情報', value: 'analytics' },
]

const currentTabValue = computed({
  get: (): number => props.currentTab,
  set: (value: number): void => emit('update:current-tab', value),
})

const videoFormDataValue = computed({
  get: (): UpdateVideoRequest => props.videoFormData,
  set: (item: UpdateVideoRequest): void => emit('update:video-form-data', item),
})

const productFormDataValue = computed({
  get: (): string[] => props.productFormData,
  set: (item: string[]): void => emit('update:product-form-data', item),
})

const experienceFormDataValue = computed({
  get: (): string[] => props.experienceFormData,
  set: (item: string[]): void => emit('update:experience-form-data', item),
})

const productDialogValue = computed({
  get: (): boolean => props.productDialog,
  set: (value: boolean): void => emit('update:product-dialog', value),
})

const experienceDialogValue = computed({
  get: (): boolean => props.experienceDialog,
  set: (value: boolean): void => emit('update:experience-dialog', value),
})

const onChangeSearchProduct = (name: string): void => {
  emit('update:search-product', name)
}

const onChangeSearchExperience = (title: string): void => {
  emit('update:search-experience', title)
}

const onClickAddProducts = (): void => {
  emit('click:add-products')
}

const onClickAddExperiences = (): void => {
  emit('click:add-experiences')
}

const onClickRemoveProduct = (productId: string): void => {
  emit('click:remove-product', productId)
}

const onClickRemoveExperience = (experienceId: string): void => {
  emit('click:remove-experience', experienceId)
}

const onClickOpenProductDialog = (): void => {
  emit('update:product-dialog', true)
}

const onClickOpenExperienceDialog = (): void => {
  emit('update:experience-dialog', true)
}

const onClickBack = (): void => {
  emit('click:back')
}

const onSubmit = (): void => {
  emit('submit')
}

const onSubmitUploadVideo = (file: File): void => {
  emit('submit:upload-video', file)
}

const onSubmitUploadThumbnail = (file: File): void => {
  emit('submit:upload-thumbnail', file)
}

const onClickBanComment = (commentId: string): void => {
  emit('click:ban-comment', commentId)
}
</script>

<template>
  <div>
    <v-dialog v-model="productDialogValue">
      <v-card>
        <v-card-title>商品紐づけ</v-card-title>
        <v-card-text>
          <v-autocomplete
            v-model="productFormDataValue"
            :loading="props.productSearchLoading"
            :items="props.searchedProducts"
            label="商品名"
            messages="商品名を入力することで紐づける商品を検索できます。"
            item-title="name"
            item-value="id"
            multiple
            chips
            @update:search-text="onChangeSearchProduct"
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
            @click="onClickAddProducts"
          >
            追加
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <v-dialog v-model="experienceDialogValue">
      <v-card>
        <v-card-title>体験紐づけ</v-card-title>
        <v-card-text>
          <v-autocomplete
            v-model="experienceFormDataValue"
            :items="props.searchedExperiences"
            item-title="title"
            item-value="id"
            label="体験名"
            messages="体験名を入力することで紐づける体験を検索できます。"
            :loading="props.experienceSearchLoading"
            multiple
            chips
            @update:search-text="onChangeSearchExperience"
          >
            <template #chip="{ props: val, item }">
              <v-chip
                v-bind="val"
                :text="item.raw.title"
                rounded
                class="px-4"
                variant="outlined"
                :prepend-avatar="getExperienceThumbnailUrl(item.raw)"
              />
            </template>
            <template #item="{ props: val, item }">
              <v-list-item
                v-bind="val"
                :title="item.raw.title"
                :prepend-avatar="getExperienceThumbnailUrl(item.raw)"
              />
            </template>
          </v-autocomplete>
        </v-card-text>
        <v-card-actions>
          <v-btn
            color="primary"
            variant="outlined"
            @click="onClickAddExperiences"
          >
            追加
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <v-card
      class="mb-16"
      :loading="props.loading"
    >
      <v-card-title>動画編集</v-card-title>
      <v-tabs
        v-model="currentTabValue"
        grow
      >
        <v-tab
          v-for="item in tabs"
          :key="item.value"
          :value="item.value"
        >
          {{ item.name }}
        </v-tab>
      </v-tabs>
      <v-tabs-window v-model="currentTabValue">
        <v-tabs-window-item value="info">
          <v-card-text>
            <template v-if="!props.loading">
              <organisms-video-form
                id="update-video-form"
                v-model="videoFormDataValue"
                :selected-products="props.selectedProducts"
                :selected-experiences="props.selectedExperiences"
                :video-is-uploading="props.videoUploadStatus.isUploading"
                :video-has-error="props.videoUploadStatus.hasError"
                :video-error-message="props.videoUploadStatus.errorMessage"
                :thumbnail-is-uploading="props.thumbnailUploadStatus.isUploading"
                :thumbnail-has-error="props.thumbnailUploadStatus.hasError"
                :thumbnail-error-message="props.thumbnailUploadStatus.errorMessage"
                @update:video="onSubmitUploadVideo"
                @update:thumbnail="onSubmitUploadThumbnail"
                @click:link-product="onClickOpenProductDialog"
                @click:link-experience="onClickOpenExperienceDialog"
                @click:delete-linked-product="onClickRemoveProduct"
                @click:delete-linked-experience="onClickRemoveExperience"
                @submit="onSubmit"
              />
            </template>
          </v-card-text>
        </v-tabs-window-item>
        <v-tabs-window-item value="preview">
          <organisms-video-preview
            :form-data="videoFormDataValue"
            :comments="props.comments"
            @click:ban-comment="onClickBanComment"
          />
        </v-tabs-window-item>
        <v-tabs-window-item value="analytics">
          <organisms-video-analytics
            :loading="loading"
            :viewer-logs="viewerLogs"
          />
        </v-tabs-window-item>
      </v-tabs-window>
    </v-card>
    <div
      class="position-fixed bottom-0 left-0 w-100 bg-white pa-4 text-right elevation-3"
    >
      <div class="d-inline-flex ga-4">
        <v-btn
          color="secondary"
          variant="outlined"
          @click="onClickBack"
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
