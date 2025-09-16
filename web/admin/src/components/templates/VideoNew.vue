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
  mdiClose,
  mdiEarth,
  mdiContentSave,
} from '@mdi/js'
import useVuelidate from '@vuelidate/core'
import dayjs, { unix } from 'dayjs'
import type { AlertType } from '~/lib/hooks'
import { getErrorMessage } from '~/lib/validations'
import { getProductThumbnailUrl, getExperienceThumbnailUrl } from '~/lib/formatter'
import { VideoStatus } from '~/types/api/v1'
import type { CreateVideoRequest, Product, Experience } from '~/types/api/v1'
import type { DateTimeInput, ImageUploadStatus } from '~/types/props'
import { CreateVideoValidationRules, TimeDataValidationRules } from '~/types/validations'
import { videoStatuses } from '~/constants'

const props = defineProps({
  loading: {
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
  formData: {
    type: Object as PropType<CreateVideoRequest>,
    default: (): CreateVideoRequest => ({
      title: '',
      description: '',
      coordinatorId: '',
      productIds: [],
      experienceIds: [],
      thumbnailUrl: '',
      videoUrl: '',
      _public: false,
      limited: false,
      publishedAt: 0,
      displayProduct: false,
      displayExperience: false,
    }),
  },
  videoUploadStatus: {
    type: Object,
    default: (): ImageUploadStatus => ({
      error: false,
      message: '',
    }),
  },
  thumbnailUploadStatus: {
    type: Object,
    default: (): ImageUploadStatus => ({
      error: false,
      message: '',
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
  productSearchLoading: {
    type: Boolean,
    default: false,
  },
  experienceSearchLoading: {
    type: Boolean,
    default: false,
  },
  productDialog: {
    type: Boolean,
    default: false,
  },
  experienceDialog: {
    type: Boolean,
    default: false,
  },
  productFormData: {
    type: Array<string>,
    default: () => [],
  },
  experienceFormData: {
    type: Array<string>,
    default: () => [],
  },
})

const emit = defineEmits<{
  (e: 'update:form-data', formData: CreateVideoRequest): void
  (e: 'update:product-dialog', value: boolean): void
  (e: 'update:experience-dialog', value: boolean): void
  (e: 'update:product-form-data', products: string[]): void
  (e: 'update:experience-form-data', experiences: string[]): void
  (e: 'update:video', file: File): void
  (e: 'update:thumbnail', file: File): void
  (e: 'click:add-products'): void
  (e: 'click:add-experiences'): void
  (e: 'click:remove-product', productId: string): void
  (e: 'click:remove-experience', experienceId: string): void
  (e: 'update:search-product', name: string): void
  (e: 'update:search-experience', title: string): void
  (e: 'submit'): void
  (e: 'click:back'): void
}>()

const formDataValue = computed({
  get: (): CreateVideoRequest => props.formData,
  set: (formData: CreateVideoRequest): void =>
    emit('update:form-data', formData),
})
const publishTimeDataValue = computed({
  get: (): DateTimeInput => ({
    date: unix(props.formData.publishedAt).format('YYYY-MM-DD'),
    time: unix(props.formData.publishedAt).format('HH:mm'),
  }),
  set: (timeData: DateTimeInput): void => {
    const publishedAt = dayjs(`${timeData.date} ${timeData.time}`)
    formDataValue.value.publishedAt = publishedAt.unix()
  },
})
const videoStatusValue = computed<VideoStatus>(() => {
  if (!formDataValue.value._public) {
    if (formDataValue.value.limited) {
      return VideoStatus.VideoStatusLimited // 限定公開
    }
    return VideoStatus.VideoStatusPrivate // 非公開
  }
  else {
    const now = Math.floor(Date.now() / 1000)
    if (formDataValue.value.publishedAt <= now) {
      return VideoStatus.VideoStatusPublished // 公開中
    }
    return VideoStatus.VideoStatusWaiting // 公開予定
  }
})
const productDialogValue = computed({
  get: (): boolean => props.productDialog,
  set: (value: boolean): void => emit('update:product-dialog', value),
})
const experienceDialogValue = computed({
  get: (): boolean => props.experienceDialog,
  set: (value: boolean): void => emit('update:experience-dialog', value),
})
const productFormDataValue = computed({
  get: (): string[] => props.productFormData,
  set: (item: string[]): void => emit('update:product-form-data', item),
})
const experienceFormDataValue = computed({
  get: (): string[] => props.experienceFormData,
  set: (item: string[]): void => emit('update:experience-form-data', item),
})

const formDataValidate = useVuelidate(
  CreateVideoValidationRules,
  formDataValue,
)
const publishTimeDataValidate = useVuelidate(
  TimeDataValidationRules,
  publishTimeDataValue,
)

const getStatus = (status: VideoStatus): string => {
  const value = videoStatuses.find(s => s.value === status)
  return value ? value.title : ''
}

const getStatusColor = (status: VideoStatus): string => {
  switch (status) {
    case VideoStatus.VideoStatusWaiting:
      return 'info'
    case VideoStatus.VideoStatusPublished:
      return 'primary'
    case VideoStatus.VideoStatusLimited:
      return 'secondary'
    case VideoStatus.VideoStatusPrivate:
      return 'warning'
    default:
      return ''
  }
}

const onClickBack = (): void => {
  emit('click:back')
}

const onSubmit = async (): Promise<void> => {
  const formDataValid = await formDataValidate.value.$validate()
  const publishTimeDataValid = await publishTimeDataValidate.value.$validate()
  if (!formDataValid || !publishTimeDataValid) {
    return
  }

  emit('submit')
}

const onChangePublishedAt = (): void => {
  const publishedAt = dayjs(
    `${publishTimeDataValue.value.date} ${publishTimeDataValue.value.time}`,
  )
  formDataValue.value.publishedAt = publishedAt.unix()
}

const onChangeVideoFile = (files?: FileList) => {
  if (!files || files.length === 0 || !files[0]) {
    return
  }
  emit('update:video', files[0])
}

const onChangeThumbnailFile = (files?: FileList) => {
  if (!files || files.length === 0 || !files[0]) {
    return
  }
  emit('update:thumbnail', files[0])
}

const onClickLinkProduct = (): void => {
  productDialogValue.value = true
}

const onClickLinkExperience = (): void => {
  experienceDialogValue.value = true
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

const onChangeSearchProduct = (name: string): void => {
  emit('update:search-product', name)
}

const onChangeSearchExperience = (title: string): void => {
  emit('update:search-experience', title)
}
</script>

<template>
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
          label="商品名で検索"
          messages="紐づける商品を検索して選択してください"
          item-title="name"
          item-value="id"
          multiple
          chips
          closable-chips
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
              :subtitle="`¥${item.raw.price.toLocaleString()}`"
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
          :loading="props.experienceSearchLoading"
          :items="props.searchedExperiences"
          label="体験名で検索"
          messages="紐づける体験を検索して選択してください"
          item-title="title"
          item-value="id"
          multiple
          chips
          closable-chips
          @update:search-text="onChangeSearchExperience"
        >
          <template #chip="{ props: val, item }">
            <v-chip
              v-bind="val"
              :prepend-avatar="getExperienceThumbnailUrl(item.raw)"
              :text="item.raw.title"
              rounded
              class="px-4"
              variant="outlined"
            />
          </template>
          <template #item="{ props: val, item }">
            <v-list-item
              v-bind="val"
              :prepend-avatar="getExperienceThumbnailUrl(item.raw)"
              :title="item.raw.title"
              :subtitle="`¥${item.raw.priceAdult.toLocaleString()} / 大人`"
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

  <v-container class="pa-6">
    <v-alert
      v-show="props.isAlert"
      :type="props.alertType"
      class="mb-6"
      v-text="props.alertText"
    />

    <div class="mb-6">
      <v-btn
        variant="text"
        :prepend-icon="mdiArrowLeft"
        @click="$router.back()"
      >
        戻る
      </v-btn>
      <h1 class="text-h4 font-weight-bold mt-2 mb-2">
        動画登録
      </h1>
      <p class="text-body-1 text-grey-darken-1">
        新しい動画コンテンツを登録します。各セクションを順番に入力してください。
      </p>
    </div>

    <!-- メインカード -->
    <v-card
      class="mt-6"
      elevation="0"
      rounded="lg"
    >
      <v-card-text class="pa-6">
        <v-row>
          <v-col
            cols="12"
            lg="8"
          >
            <!-- 基本情報セクション -->
            <v-card
              class="form-section-card mb-6"
              elevation="2"
            >
              <v-card-title class="d-flex align-center section-header">
                <v-icon
                  :icon="mdiPlayCircle"
                  size="24"
                  class="mr-3 text-primary"
                />
                <span class="text-h6 font-weight-medium">基本情報</span>
              </v-card-title>
              <v-card-text class="pa-6">
                <v-text-field
                  v-model="formDataValidate.title.$model"
                  :error-messages="getErrorMessage(formDataValidate.title.$errors)"
                  label="動画タイトル *"
                  variant="outlined"
                  density="comfortable"
                  class="mb-4"
                />
                <v-textarea
                  v-model="formDataValidate.description.$model"
                  :error-messages="
                    getErrorMessage(formDataValidate.description.$errors)
                  "
                  label="動画説明"
                  maxlength="2000"
                  variant="outlined"
                  density="comfortable"
                  rows="4"
                  counter
                />
              </v-card-text>
            </v-card>

            <!-- メディアファイル管理セクション -->
            <v-card
              class="form-section-card mb-6"
              elevation="2"
            >
              <v-card-title class="d-flex align-center section-header">
                <v-icon
                  :icon="mdiVideo"
                  size="24"
                  class="mr-3 text-primary"
                />
                <span class="text-h6 font-weight-medium">メディアファイル</span>
              </v-card-title>
              <v-card-text class="pa-6">
                <v-row>
                  <v-col
                    cols="12"
                    md="6"
                  >
                    <molecules-image-select-form
                      label="サムネイル画像"
                      :loading="props.loading"
                      :img-url="formDataValue.thumbnailUrl"
                      :error="props.thumbnailUploadStatus.hasError"
                      :error-messages="props.thumbnailUploadStatus.hasError ? props.thumbnailUploadStatus.errorMessage : ''"
                      @update:file="onChangeThumbnailFile"
                    />
                  </v-col>
                  <v-col
                    cols="12"
                    md="6"
                  >
                    <molecules-video-select-form
                      label="動画"
                      :loading="props.loading"
                      :video-url="formDataValue.videoUrl"
                      :error="props.videoUploadStatus.hasError"
                      :error-messages="props.videoUploadStatus.hasError ? props.videoUploadStatus.errorMessage : ''"
                      @update:file="onChangeVideoFile"
                    />
                  </v-col>
                </v-row>
              </v-card-text>
            </v-card>

            <!-- 関連コンテンツセクション -->
            <v-card
              class="form-section-card mb-6"
              elevation="2"
            >
              <v-card-title class="d-flex align-center section-header">
                <v-icon
                  :icon="mdiShopping"
                  size="24"
                  class="mr-3 text-primary"
                />
                <span class="text-h6 font-weight-medium">関連コンテンツ</span>
              </v-card-title>
              <v-card-text class="pa-6">
                <!-- 商品セクション -->
                <div class="mb-6">
                  <div class="d-flex align-center justify-space-between mb-3">
                    <div class="text-subtitle-1 font-weight-medium">
                      紐づけ商品
                    </div>
                    <v-btn
                      size="small"
                      variant="tonal"
                      color="primary"
                      @click="onClickLinkProduct"
                    >
                      <v-icon
                        start
                        :icon="mdiPlus"
                      />
                      商品を追加
                    </v-btn>
                  </div>
                  <div v-if="props.selectedProducts.length === 0">
                    <v-alert
                      type="info"
                      variant="tonal"
                      density="compact"
                    >
                      紐づけられた商品はありません
                    </v-alert>
                  </div>
                  <v-list
                    v-else
                    density="compact"
                  >
                    <v-list-item
                      v-for="product in props.selectedProducts"
                      :key="product.id"
                      :prepend-avatar="getProductThumbnailUrl(product)"
                    >
                      <v-list-item-title>{{ product.name }}</v-list-item-title>
                      <v-list-item-subtitle>¥{{ product.price.toLocaleString() }}</v-list-item-subtitle>
                      <template #append>
                        <v-btn
                          icon
                          size="small"
                          variant="text"
                          @click="onClickRemoveProduct(product.id)"
                        >
                          <v-icon :icon="mdiClose" />
                        </v-btn>
                      </template>
                    </v-list-item>
                  </v-list>
                  <v-checkbox
                    v-model="formDataValue.displayProduct"
                    label="動画再生時に商品を表示する"
                    density="compact"
                    class="mt-2"
                  />
                </div>

                <!-- 体験セクション -->
                <div>
                  <div class="d-flex align-center justify-space-between mb-3">
                    <div class="text-subtitle-1 font-weight-medium">
                      紐づけ体験
                    </div>
                    <v-btn
                      size="small"
                      variant="tonal"
                      color="primary"
                      @click="onClickLinkExperience"
                    >
                      <v-icon
                        start
                        :icon="mdiPlus"
                      />
                      体験を追加
                    </v-btn>
                  </div>
                  <div v-if="props.selectedExperiences.length === 0">
                    <v-alert
                      type="info"
                      variant="tonal"
                      density="compact"
                    >
                      紐づけられた体験はありません
                    </v-alert>
                  </div>
                  <v-list
                    v-else
                    density="compact"
                  >
                    <v-list-item
                      v-for="experience in props.selectedExperiences"
                      :key="experience.id"
                      :prepend-avatar="getExperienceThumbnailUrl(experience)"
                    >
                      <v-list-item-title>{{ experience.title }}</v-list-item-title>
                      <v-list-item-subtitle>¥{{ experience.priceAdult.toLocaleString() }} / 大人</v-list-item-subtitle>
                      <template #append>
                        <v-btn
                          icon
                          size="small"
                          variant="text"
                          @click="onClickRemoveExperience(experience.id)"
                        >
                          <v-icon :icon="mdiClose" />
                        </v-btn>
                      </template>
                    </v-list-item>
                  </v-list>
                  <v-checkbox
                    v-model="formDataValue.displayExperience"
                    label="動画再生時に体験を表示する"
                    density="compact"
                    class="mt-2"
                  />
                </div>
              </v-card-text>
            </v-card>
          </v-col>

          <!-- 右側のサイドバー -->
          <v-col
            cols="12"
            lg="4"
          >
            <!-- 公開設定セクション -->
            <v-card
              class="form-section-card mb-6"
              elevation="2"
            >
              <v-card-title class="d-flex align-center section-header">
                <v-icon
                  :icon="mdiEarth"
                  size="24"
                  class="mr-3 text-primary"
                />
                <span class="text-h6 font-weight-medium">公開設定</span>
              </v-card-title>
              <v-card-text class="pa-6">
                <v-alert
                  :color="getStatusColor(videoStatusValue)"
                  variant="tonal"
                  density="compact"
                  class="mb-4"
                >
                  現在の状況: {{ getStatus(videoStatusValue) }}
                </v-alert>

                <v-checkbox
                  v-model="formDataValue._public"
                  label="動画を公開する"
                  density="comfortable"
                />
                <v-checkbox
                  v-model="formDataValue.limited"
                  label="限定公開にする"
                  density="comfortable"
                />

                <div class="date-time-section">
                  <p class="text-subtitle-2 mb-3 text-grey-darken-1">
                    公開日時 *
                  </p>
                  <v-row>
                    <v-col
                      cols="12"
                      sm="6"
                    >
                      <v-text-field
                        v-model="publishTimeDataValidate.date.$model"
                        :error-messages="
                          getErrorMessage(publishTimeDataValidate.date.$errors)
                        "
                        label="日付"
                        type="date"
                        variant="outlined"
                        density="comfortable"
                        @update:model-value="onChangePublishedAt"
                      />
                    </v-col>
                    <v-col
                      cols="12"
                      sm="6"
                    >
                      <v-text-field
                        v-model="publishTimeDataValidate.time.$model"
                        :error-messages="
                          getErrorMessage(publishTimeDataValidate.time.$errors)
                        "
                        label="時刻"
                        type="time"
                        variant="outlined"
                        density="comfortable"
                        @update:model-value="onChangePublishedAt"
                      />
                    </v-col>
                  </v-row>
                </div>
              </v-card-text>
            </v-card>
          </v-col>
        </v-row>

        <!-- 送信ボタン -->
        <div class="d-flex justify-end mt-6">
          <v-btn
            :loading="loading"
            color="primary"
            variant="elevated"
            size="large"
            @click="onSubmit"
          >
            <v-icon
              :icon="mdiContentSave"
              start
            />
            登録
          </v-btn>
        </div>
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
  padding: 16px 24px;
}

.sticky-sidebar {
  position: sticky;
  top: 24px;
}

@media (width <= 1280px) {
  .sticky-sidebar {
    position: static;
  }
}

@media (width <= 600px) {
  .form-section-card {
    border-radius: 8px;
  }
}
</style>
