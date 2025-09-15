<script lang="ts" setup>
import {
  mdiPlayCircle,
  mdiVideo,
  mdiInformationOutline,
  mdiContentSave,
  mdiShopping,
  mdiPlus,
  mdiClose,
  mdiEarth,
} from '@mdi/js'
import useVuelidate from '@vuelidate/core'
import { getErrorMessage } from '~/lib/validations'
import { getProductThumbnailUrl, getExperienceThumbnailUrl } from '~/lib/formatter'
import { VideoStatus } from '~/types/api/v1'
import type { UpdateVideoRequest, Product, Experience } from '~/types/api/v1'
import type { DateTimeInput, UploadStatus } from '~/types/props'
import { TimeDataValidationRules, UpdateVideoValidationRules } from '~/types/validations'
import dayjs, { unix } from 'dayjs'

const props = defineProps({
  loading: {
    type: Boolean,
    default: false,
  },
  updatable: {
    type: Boolean,
    default: true,
  },
  formData: {
    type: Object as PropType<UpdateVideoRequest>,
    default: (): UpdateVideoRequest => ({
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
  selectedProducts: {
    type: Array<Product>,
    default: () => [],
  },
  selectedExperiences: {
    type: Array<Experience>,
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
})

const emit = defineEmits<{
  (e: 'update:form-data', formData: UpdateVideoRequest): void
  (e: 'update:video', file: File): void
  (e: 'update:thumbnail', file: File): void
  (e: 'click:link-product'): void
  (e: 'click:link-experience'): void
  (e: 'click:remove-product', productId: string): void
  (e: 'click:remove-experience', experienceId: string): void
  (e: 'submit'): void
}>()

const videoStatuses = [
  { value: VideoStatus.VideoStatusPrivate, title: '非公開' },
  { value: VideoStatus.VideoStatusLimited, title: '限定公開' },
  { value: VideoStatus.VideoStatusPublished, title: '公開中' },
  { value: VideoStatus.VideoStatusWaiting, title: '公開予定' },
]

const formDataValue = computed({
  get: (): UpdateVideoRequest => props.formData,
  set: (formData: UpdateVideoRequest): void =>
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
const videoStatus = computed<VideoStatus>(() => {
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

const formDataValidate = useVuelidate(
  UpdateVideoValidationRules,
  formDataValue,
)
const publishTimeDataValidate = useVuelidate(
  TimeDataValidationRules,
  publishTimeDataValue,
)

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
  emit('click:link-product')
}

const onClickLinkExperience = (): void => {
  emit('click:link-experience')
}

const onClickRemoveProduct = (productId: string): void => {
  emit('click:remove-product', productId)
}

const onClickRemoveExperience = (experienceId: string): void => {
  emit('click:remove-experience', experienceId)
}

const onSubmit = async (): Promise<void> => {
  const formDataValid = await formDataValidate.value.$validate()
  const publishTimeDataValid = await publishTimeDataValidate.value.$validate()
  if (!formDataValid || !publishTimeDataValid) {
    return
  }

  emit('submit')
}
</script>

<template>
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
            :readonly="!updatable"
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
            :readonly="!updatable"
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
                v-if="updatable"
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
                    v-if="updatable"
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
              :readonly="!updatable"
            />
          </div>

          <!-- 体験セクション -->
          <div>
            <div class="d-flex align-center justify-space-between mb-3">
              <div class="text-subtitle-1 font-weight-medium">
                紐づけ体験
              </div>
              <v-btn
                v-if="updatable"
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
                    v-if="updatable"
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
              :readonly="!updatable"
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
            :type="videoStatus === VideoStatus.VideoStatusPublished ? 'success' : 'info'"
            variant="tonal"
            density="compact"
            class="mb-4"
          >
            現在の状況: {{ videoStatuses.find(s => s.value === videoStatus)?.title || '不明' }}
          </v-alert>

          <v-checkbox
            v-model="formDataValue._public"
            label="動画を公開する"
            density="comfortable"
            :readonly="!updatable"
          />
          <v-checkbox
            v-model="formDataValue.limited"
            label="限定公開にする"
            density="comfortable"
            :readonly="!updatable"
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
                  :readonly="!updatable"
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
                  :readonly="!updatable"
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

      <!-- 注意事項 -->
      <v-card
        class="form-section-card"
        elevation="2"
      >
        <v-card-title class="d-flex align-center section-header">
          <v-icon
            :icon="mdiInformationOutline"
            size="24"
            class="mr-3 text-warning"
          />
          <span class="text-h6 font-weight-medium">注意事項</span>
        </v-card-title>
        <v-card-text class="pa-6">
          <v-list
            density="compact"
            class="bg-transparent"
          >
            <v-list-item class="px-0">
              <v-list-item-title class="text-body-2">
                • 動画ファイルは最大500MBまでアップロード可能です
              </v-list-item-title>
            </v-list-item>
            <v-list-item class="px-0">
              <v-list-item-title class="text-body-2">
                • サムネイル画像は16:9の比率を推奨します
              </v-list-item-title>
            </v-list-item>
            <v-list-item class="px-0">
              <v-list-item-title class="text-body-2">
                • 公開設定を変更すると即座に反映されます
              </v-list-item-title>
            </v-list-item>
          </v-list>
        </v-card-text>
      </v-card>
    </v-col>
  </v-row>

  <!-- 送信ボタン -->
  <div
    v-if="updatable"
    class="d-flex justify-end mt-6"
  >
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
      更新
    </v-btn>
  </div>
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
