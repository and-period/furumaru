<script lang="ts" setup>
import {
  mdiPlus,
  mdiArrowLeft,
  mdiPlayCircle,
  mdiImageMultiple,
  mdiVideo,
  mdiFileDocument,
  mdiShopping,
  mdiMapMarker,
  mdiChartLine,
  mdiEye,
  mdiComment,
  mdiAccount,
  mdiClose,
} from '@mdi/js'
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
  <v-container class="pa-6">
    <!-- 商品紐づけダイアログ -->
    <v-dialog
      v-model="productDialogValue"
      max-width="600"
    >
      <v-card>
        <v-card-title class="d-flex align-center">
          <v-icon
            :icon="mdiShopping"
            class="mr-2"
          />
          商品を紐づける
        </v-card-title>
        <v-card-text>
          <v-autocomplete
            v-model="productFormDataValue"
            :loading="props.productSearchLoading"
            :items="props.searchedProducts"
            label="商品名"
            messages="紐づける商品を検索して選択してください"
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
          <v-spacer />
          <v-btn
            variant="text"
            @click="productDialogValue = false"
          >
            キャンセル
          </v-btn>
          <v-btn
            color="primary"
            variant="elevated"
            @click="onClickAddProducts"
          >
            追加
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- 体験紐づけダイアログ -->
    <v-dialog
      v-model="experienceDialogValue"
      max-width="600"
    >
      <v-card>
        <v-card-title class="d-flex align-center">
          <v-icon
            :icon="mdiMapMarker"
            class="mr-2"
          />
          体験を紐づける
        </v-card-title>
        <v-card-text>
          <v-autocomplete
            v-model="experienceFormDataValue"
            :items="props.searchedExperiences"
            item-title="title"
            item-value="id"
            label="体験名"
            messages="紐づける体験を検索して選択してください"
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
          <v-spacer />
          <v-btn
            variant="text"
            @click="experienceDialogValue = false"
          >
            キャンセル
          </v-btn>
          <v-btn
            color="primary"
            variant="elevated"
            @click="onClickAddExperiences"
          >
            追加
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <atoms-app-alert
      :show="props.isAlert"
      :type="props.alertType"
      :text="props.alertText"
      class="mb-6"
    />

    <div class="mb-6">
      <v-btn
        variant="text"
        :prepend-icon="mdiArrowLeft"
        @click="onClickBack"
      >
        戻る
      </v-btn>
      <h1 class="text-h4 font-weight-bold mt-2 mb-2">
        動画編集
      </h1>
      <p class="text-body-1 text-grey-darken-1">
        動画コンテンツの詳細情報を編集・管理します。タブを切り替えて各種設定を行ってください。
      </p>
    </div>

    <v-card
      class="form-section-card mb-6"
      elevation="2"
      :loading="props.loading"
    >
      <v-card-title class="section-header pa-0">
        <v-tabs
          v-model="currentTabValue"
          class="w-100"
          density="comfortable"
        >
          <v-tab
            value="info"
            class="tab-item"
          >
            <v-icon
              :icon="mdiFileDocument"
              size="20"
              class="mr-2"
            />
            基本情報
          </v-tab>
          <v-tab
            value="preview"
            class="tab-item"
          >
            <v-icon
              :icon="mdiEye"
              size="20"
              class="mr-2"
            />
            プレビュー
          </v-tab>
          <v-tab
            value="analytics"
            class="tab-item"
          >
            <v-icon
              :icon="mdiChartLine"
              size="20"
              class="mr-2"
            />
            分析情報
          </v-tab>
          <v-tab
            value="comments"
            class="tab-item"
          >
            <v-icon
              :icon="mdiComment"
              size="20"
              class="mr-2"
            />
            コメント管理
          </v-tab>
        </v-tabs>
      </v-card-title>
      <v-card-text class="pa-0">
        <v-window
          v-model="currentTabValue"
          class="tab-content"
        >
          <v-window-item value="info">
            <div class="pa-6">
              <template v-if="!props.loading">
                <organisms-video-show
                  v-model:form-data="videoFormDataValue"
                  :loading="props.loading"
                  :updatable="true"
                  :selected-products="props.selectedProducts"
                  :selected-experiences="props.selectedExperiences"
                  :video-upload-status="props.videoUploadStatus"
                  :thumbnail-upload-status="props.thumbnailUploadStatus"
                  @update:video="onSubmitUploadVideo"
                  @update:thumbnail="onSubmitUploadThumbnail"
                  @click:link-product="onClickOpenProductDialog"
                  @click:link-experience="onClickOpenExperienceDialog"
                  @click:remove-product="onClickRemoveProduct"
                  @click:remove-experience="onClickRemoveExperience"
                  @submit="onSubmit"
                />
              </template>
            </div>
          </v-window-item>
          <v-window-item value="preview">
            <div class="pa-6">
              <organisms-video-preview
                :form-data="videoFormDataValue"
                :comments="props.comments"
                @click:ban-comment="onClickBanComment"
              />
            </div>
          </v-window-item>
          <v-window-item value="analytics">
            <div class="pa-6">
              <organisms-video-analytics
                :loading="loading"
                :viewer-logs="viewerLogs"
              />
            </div>
          </v-window-item>
          <v-window-item value="comments">
            <div class="pa-6">
              <div class="mb-4">
                <h3 class="text-h6 mb-2">
                  コメント管理
                </h3>
                <p class="text-body-2 text-medium-emphasis">
                  視聴者からのコメントを管理します
                </p>
              </div>
              <v-list v-if="props.comments.length > 0">
                <v-list-item
                  v-for="comment in props.comments"
                  :key="comment.id"
                  class="mb-2"
                >
                  <template #prepend>
                    <v-avatar color="grey-lighten-1">
                      <v-icon :icon="mdiAccount" />
                    </v-avatar>
                  </template>
                  <v-list-item-title>{{ comment.username }}</v-list-item-title>
                  <v-list-item-subtitle>{{ comment.comment }}</v-list-item-subtitle>
                  <template #append>
                    <v-btn
                      v-if="!comment.disabled"
                      size="small"
                      variant="text"
                      color="error"
                      @click="onClickBanComment(comment.id)"
                    >
                      非表示にする
                    </v-btn>
                    <v-chip
                      v-else
                      size="small"
                      color="error"
                    >
                      非表示中
                    </v-chip>
                  </template>
                </v-list-item>
              </v-list>
              <v-alert
                v-else
                type="info"
                variant="tonal"
              >
                コメントはまだありません
              </v-alert>
            </div>
          </v-window-item>
        </v-window>
      </v-card-text>
    </v-card>
  </v-container>
</template>

<style scoped>
.form-section-card {
  border-radius: 12px;
  max-width: none;
}

.section-header {
  background: linear-gradient(90deg, rgb(33 150 243 / 5%) 0%, rgb(33 150 243 / 0%) 100%);
  border-bottom: 1px solid rgb(0 0 0 / 5%);
  padding: 0;
}

.tab-item {
  text-transform: none;
  font-weight: 500;
}

.tab-content {
  min-height: 400px;
}

@media (width <= 600px) {
  .form-section-card {
    border-radius: 8px;
  }

  .tab-item {
    min-width: auto;
    font-size: 0.875rem;
  }
}
</style>
