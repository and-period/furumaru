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
}

defineProps<Props>()

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
</script>

<template>
  <v-form
    :id="id"
    class="d-flex flex-column ga-4"
  >
    <v-text-field
      v-model="formData.title"
      label="タイトル"
      required
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
    />

    <molecules-video-select-form
      label="動画"
      :video-url="formData.videoUrl"
    />
    <molecules-icon-select-form
      label="サムネイル"
      :img-url="formData.thumbnailUrl"
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
      <template
        v-for="product in selectedProducts"
        :key="product.id"
      >
        <molecules-video-linked-product-item :item="product" />
      </template>
      <v-btn
        class="w-100 mt-3"
        variant="outlined"
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
      >
        <v-icon :icon="mdiPlus" />
        体験を紐づける
      </v-btn>
    </div>
  </v-form>
</template>
