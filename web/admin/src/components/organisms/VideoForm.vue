<script setup lang="ts">
import dayjs, { unix } from 'dayjs'
import { mdiPlus } from '@mdi/js'
import type {
  CreateVideoRequest,
  UpdateVideoRequest,
  Product,
  Experience,
} from '~/types/api'

interface Props {
  id: string
  selectedProducts: Product[]
  selectedExperiences: Experience[]
  videoIsUploading: boolean
  videoHasError: boolean
  videoErrorMessage: string
  thumbnailIsUploading: boolean
  thumbnailHasError: boolean
  thumbnailErrorMessage: string
}

defineProps<Props>()

interface Emits {
  (e: 'click:link-product'): void
  (e: 'click:link-experience'): void
  (e: 'update:video', files: File): void
  (e: 'update:thumbnail', files: File): void
  (e: 'submit'): void
}

const emits = defineEmits<Emits>()

const formData = defineModel<CreateVideoRequest | UpdateVideoRequest>({
  required: true,
})

/**
 * 公開範囲の選択肢
 */
const publicSelectItems = [
  {
    title: '公開',
    value: true,
  },
  {
    title: '非公開',
    value: false,
  },
]

/**
 * publishedAtの情報から日付フォームの初期値を設定する関数
 * @param publishedAt unixtimeの時刻情報
 */
const initializePublishedAtDate = (publishedAt: number): string => {
  if (publishedAt === 0) {
    return dayjs().format('YYYY-MM-DD')
  }
  else {
    return unix(publishedAt).format('YYYY-MM-DD')
  }
}

/**
 * publishedAtの情報から時刻フォームの初期値を設定する関数
 * @param publishedAt unixtimeの時刻情報
 */
const initializePublishedAtTime = (publishedAt: number): string => {
  if (publishedAt === 0) {
    return '12:00'
  }
  else {
    return unix(publishedAt).format('HH:mm')
  }
}

/**
 * 公開時間の日付フォーム用のデータ
 * 初期値はpublishedAtの情報が0なら今日の日付
 */
const publishedAtDate = ref<string>(
  initializePublishedAtDate(formData.value.publishedAt),
)

/**
 * 公開時間の時刻フォーム用のデータ
 * 初期値はpublishedAtの情報が0なら00:00
 */
const publishedAtTime = ref<string>(
  initializePublishedAtTime(formData.value.publishedAt),
)

/**
 * 公開時間のフォームの値をセットする関数
 */
watch(
  () => {
    return [publishedAtDate.value, publishedAtTime.value]
  },
  () => {
    const publishedAt: number = dayjs(
      `${publishedAtDate.value} ${publishedAtTime.value}`,
    ).unix()
    formData.value.publishedAt = publishedAt
  },
)

/**
 * タイトルのバリデーションルール
 */
const titleRules = [
  (v: string) => !!v || 'タイトルは必須です',
  (v: string) =>
    (v && v.length <= 128) || 'タイトルは128文字以内で入力してください',
]

/**
 * 詳細のバリデーションルール
 */
const descriptionRules = [
  (v: string) => v.length <= 2000 || '詳細は2000文字以内で入力してください',
]

/**
 * 動画ファイルの変更時の処理
 */
const handleChangeVideoFile = (files?: FileList): void => {
  if (files) {
    emits('update:video', files[0])
  }
}

/**
 * サムネイル画像の変更時の処理
 */
const handleChangeThumbnailFile = (files?: FileList): void => {
  if (files) {
    emits('update:thumbnail', files[0])
  }
}

/**
 * 商品を紐づけるボタンクリック時の処理
 */
const handleClickLinkProductButton = (): void => {
  emits('click:link-product')
}

/**
 * 体験を紐づけるボタンクリック時の処理
 */
const handleClickLinkExperienceButton = (): void => {
  emits('click:link-experience')
}

/**
 * フォームの送信時の処理
 */
const handleSubmit = () => {
  emits('submit')
}
</script>

<template>
  <v-form
    :id="id"
    class="d-flex flex-column ga-4"
    @submit.prevent="handleSubmit"
  >
    <v-text-field
      v-model="formData.title"
      label="タイトル"
      required
      :rules="titleRules"
    />
    <v-select
      v-model="formData.public"
      label="公開状況"
      :items="publicSelectItems"
      items-title="title"
      item-value="value"
    />
    <div class="w-100">
      <p class="text-subtitle-2 text-gray mb-2">
        公開時間
      </p>
      <div class="d-flex ga-4 w-50">
        <v-text-field
          v-model="publishedAtDate"
          type="date"
          density="compact"
          variant="outlined"
          :disabled="!formData.public"
        />
        <v-text-field
          v-model="publishedAtTime"
          type="time"
          density="compact"
          variant="outlined"
          :disabled="!formData.public"
        />
        <div class="pt-3">
          ～
        </div>
      </div>
    </div>
    <v-textarea
      v-model="formData.description"
      label="詳細"
      :rules="descriptionRules"
    />

    <molecules-video-select-form
      label="動画"
      :video-url="formData.videoUrl"
      :loading="videoIsUploading"
      :error="videoHasError"
      :message="videoErrorMessage"
      @update:file="handleChangeVideoFile"
    />
    <molecules-image-select-form
      label="サムネイル"
      :img-url="formData.thumbnailUrl"
      :loading="thumbnailIsUploading"
      :error="thumbnailHasError"
      :message="thumbnailErrorMessage"
      @update:file="handleChangeThumbnailFile"
    />

    <div>
      <p class="text-subtitle-2 text-gray mb-2">
        カテゴリ
      </p>
      <v-checkbox
        v-model="formData.displayProduct"
        label="商品"
        density="compact"
        :hide-details="true"
      />
      <v-checkbox
        v-model="formData.displayExperience"
        label="体験"
        density="compact"
        :hide-details="true"
      />
    </div>

    <div>
      <p class="text-subtitle-2 text-gray mb-2">
        商品を紐づける
      </p>
      <template v-if="selectedProducts.length === 0">
        <v-alert
          dense
          variant="outlined"
          type="info"
          class="my-2"
        >
          商品が紐づけられていません
        </v-alert>
      </template>
      <template v-else>
        <v-table>
          <thead>
            <tr>
              <th />

              <th>商品名</th>
              <th>価格</th>
              <th>在庫</th>
            </tr>
          </thead>
          <tbody>
            <tr
              v-for="product in selectedProducts"
              :key="product.id"
            >
              <molecules-video-linked-product-item :item="product" />
            </tr>
          </tbody>
        </v-table>
      </template>
      <v-btn
        class="w-100 mt-3"
        variant="outlined"
        @click="handleClickLinkProductButton"
      >
        <v-icon :icon="mdiPlus" />
        商品を紐づける
      </v-btn>
    </div>

    <div>
      <p class="text-subtitle-2 text-gray mb-2">
        体験を紐づける
      </p>
      <template
        v-for="experience in selectedExperiences"
        :key="experience.id"
      >
        <molecules-video-linked-experience-item :item="experience" />
      </template>

      <v-btn
        class="w-100 mt-3"
        variant="outlined"
        @click="handleClickLinkExperienceButton"
      >
        <v-icon :icon="mdiPlus" />
        体験を紐づける
      </v-btn>
    </div>
  </v-form>
</template>
