<script setup lang="ts">
import dayjs, { unix } from 'dayjs'
import { mdiPencil } from '@mdi/js'
import useVuelidate from '@vuelidate/core'
import type {
  Live,
  Producer,
  Product,
  ProductMediaInner,
  UpdateLiveRequest,
} from '~/types/api'
import {
  TimeDataValidationRules,
  UpdateLiveValidationRules,
} from '~/types/validations'
import { getErrorMessage } from '~/lib/validations'
import type { DateTimeInput } from '~/types/props'

interface Props {
  /**
   * ライブで生産者が扱う商品などの情報
   */
  item: Live
  /**
   * 生産者のサムネイル画像のURL
   */
  producerThumbnailUrl: string
  /**
   * 生産者のサムネイル画像のsrcset属性
   */
  producerThumbnailsSrcset: string
  /**
   * 生産者の名前
   */
  producerName: string
  /**
   * 生産者が扱う商品のリスト
   */
  products: Product[]
  /**
   * ライブに紐づいている商品のリスト
   */
  liveProducts: Product[]
  /**
   * 生産者のリスト
   */
  producers: Producer[]
  /**
   * フォーム送信中かどうかのフラグ
   */
  loading: boolean
}

const props = defineProps<Props>()

interface Emits {
  (e: 'search:product', producerId: string, keyword: string): void
  (e: 'submit:update', id: string, formData: UpdateLiveRequest): void
  (e: 'submit:delete', id: string): void
}

const emits = defineEmits<Emits>()

/**
 * 商品のサムネイル情報を取得する関数
 */
const getProductThumbnailUrl = (product: Product): string => {
  const thumbnail = product.media?.find((media: ProductMediaInner) => {
    return media.isThumbnail
  })
  return thumbnail ? thumbnail.url : ''
}

/**
 * unixTimeを日付文字列に変換する関数
 * @param unixTime
 */
const getDay = (unixTime: number): string => {
  return unix(unixTime).format('YYYY/MM/DD HH:mm')
}

/**
 * 生産者のライブ配信期間
 */
const liveTerm = computed<string>(() => {
  return `${getDay(props.item.startAt)} ~ ${getDay(props.item.endAt)}`
})

/**
 * 編集ダイアログの表示フラグ
 */
const updateDialogValue = ref<boolean>(false)

/**
 * 更新フォームのデータ
 */
const formData = ref<UpdateLiveRequest>({
  productIds: props.item.productIds,
  comment: props.item.comment,
  startAt: props.item.startAt,
  endAt: props.item.endAt,
})

/**
 * 開始時間のフォームデータ
 */
const updateStartTimeDataValue = computed({
  get: (): DateTimeInput => ({
    date: unix(formData.value.startAt).format('YYYY-MM-DD'),
    time: unix(formData.value.startAt).format('HH:mm'),
  }),
  set: (timeData: DateTimeInput): void => {
    const startAt = dayjs(`${timeData.date} ${timeData.time}`)
    formData.value.startAt = startAt.unix()
  },
})

/**
 * 終了時間のフォームデータ
 */
const updateEndTimeDataValue = computed({
  get: (): DateTimeInput => ({
    date: unix(formData.value.endAt).format('YYYY-MM-DD'),
    time: unix(formData.value.endAt).format('HH:mm'),
  }),
  set: (timeData: DateTimeInput): void => {
    const endAt = dayjs(`${timeData.date} ${timeData.time}`)
    formData.value.endAt = endAt.unix()
  },
})

/**
 * 更新フォームのバリデーションルール
 */
const updateFormDataValidate = useVuelidate(
  UpdateLiveValidationRules,
  formData,
)

/**
 * 更新フォームの開始日のバリデーションルール
 */
const updateStartTimeDataValidate = useVuelidate(
  TimeDataValidationRules,
  updateStartTimeDataValue,
)

/**
 * 更新フォームの終了日のバリデーションルール
 */
const updateEndTimeDataValidate = useVuelidate(
  TimeDataValidationRules,
  updateEndTimeDataValue,
)

/**
 * 開始時刻の更新処理
 */
const onChangeUpdateStartAt = (): void => {
  const startAt = dayjs(
    `${updateStartTimeDataValue.value.date} ${updateStartTimeDataValue.value.time}`,
  )
  formData.value.startAt = startAt.unix()
}

/**
 * 終了時刻の更新処理
 */
const onChangeUpdateEndAt = (): void => {
  const endAt = dayjs(
    `${updateEndTimeDataValue.value.date} ${updateEndTimeDataValue.value.time}`,
  )
  formData.value.endAt = endAt.unix()
}

/**
 * 編集ボタンがクリックされたときの処理
 */
const handleClickEditButton = () => {
  updateDialogValue.value = true
}

/**
 * キャンセルボタンがクリックされたときの処理
 */
const handleClickCloseDialogButton = () => {
  updateDialogValue.value = false
}

/**
 * 生産者が扱う商品を検索する処理
 * @param keyword
 */
const onSearchProductFromUpdate = (keyword: string) => {
  emits('search:product', props.item.producerId, keyword)
}

/**
 * 更新処理
 */
const onSubmitUpdate = async () => {
  // フォームのバリデーションチェック
  const formDataValid = await updateFormDataValidate.value.$validate()
  const startTimeDataValid
    = await updateStartTimeDataValidate.value.$validate()
  const endTimeDataValid = await updateEndTimeDataValidate.value.$validate()

  if (!formDataValid || !startTimeDataValid || !endTimeDataValid) {
    return
  }

  emits('submit:update', props.item.id, formData.value)
}

/**
 * 削除処理
 */
const onSubmitDelete = () => {
  emits('submit:delete', props.item.id)
}

watch(
  () => props.loading,
  (value) => {
    // ローディングが終了したらダイアログを閉じる
    if (!value) {
      updateDialogValue.value = false
    }
  },
)

watch(
  () => formData.value.productIds,
  () => {
    console.log('並び替えが発生', formData.value.productIds)
    // ダイアログが閉じられた状態でフォーム要素のproductIdsが更新されたら更新処理を走らせる
    if (!updateDialogValue.value) {
      onSubmitUpdate()
    }
  },
)
</script>

<template>
  <div>
    <v-dialog
      v-model="updateDialogValue"
      width="500"
    >
      <v-card>
        <v-card-title class="text-h6 primaryLight">
          スケジュール更新
        </v-card-title>
        <v-card-text>
          <p class="text-subtitle-2 text-grey pb-2">
            ライブ配開始日時
          </p>
          <div class="d-flex flex-column flex-md-row justify-center">
            <v-text-field
              v-model="updateStartTimeDataValidate.date.$model"
              :error-messages="
                getErrorMessage(updateStartTimeDataValidate.date.$errors)
              "
              type="date"
              variant="outlined"
              density="compact"
              class="mr-md-2"
              @update:model-value="onChangeUpdateStartAt"
            />
            <v-text-field
              v-model="updateStartTimeDataValidate.time.$model"
              :error-messages="
                getErrorMessage(updateStartTimeDataValidate.time.$errors)
              "
              type="time"
              variant="outlined"
              density="compact"
              @update:model-value="onChangeUpdateStartAt"
            />
          </div>
          <p class="text-subtitle-2 text-grey pb-2">
            ライブ配終了日時
          </p>
          <div class="d-flex flex-column flex-md-row justify-center">
            <v-text-field
              v-model="updateEndTimeDataValidate.date.$model"
              :error-messages="
                getErrorMessage(updateEndTimeDataValidate.date.$errors)
              "
              type="date"
              variant="outlined"
              density="compact"
              class="mr-md-2"
              @update:model-value="onChangeUpdateEndAt"
            />
            <v-text-field
              v-model="updateEndTimeDataValidate.time.$model"
              :error-messages="
                getErrorMessage(updateEndTimeDataValidate.time.$errors)
              "
              type="time"
              variant="outlined"
              density="compact"
              @update:model-value="onChangeUpdateEndAt"
            />
          </div>
          <v-autocomplete
            :model-value="item.producerId"
            label="生産者"
            :items="producers"
            item-title="username"
            item-value="id"
            readonly
          />
          <v-autocomplete
            v-model="updateFormDataValidate.productIds.$model"
            :error-messages="
              getErrorMessage(updateFormDataValidate.productIds.$errors)
            "
            label="関連する商品"
            :items="products"
            item-title="name"
            item-value="id"
            chips
            closable-chips
            clearable
            multiple
            density="comfortable"
            @update:search="onSearchProductFromUpdate"
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
          <v-textarea
            v-model="updateFormDataValidate.comment.$model"
            :error-messages="
              getErrorMessage(updateFormDataValidate.comment.$errors)
            "
            label="概要"
            maxlength="2000"
          />
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn
            color=""
            variant="text"
            @click="handleClickCloseDialogButton"
          >
            キャンセル
          </v-btn>
          <v-btn
            :loading="loading"
            color="error"
            variant="text"
            @click="onSubmitDelete"
          >
            削除
          </v-btn>
          <v-btn
            :loading="loading"
            color="primary"
            variant="text"
            @click="onSubmitUpdate"
          >
            更新
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <v-card>
      <v-card-title>
        <v-list-item>
          <template #prepend>
            <v-avatar>
              <v-img
                cover
                :src="producerThumbnailUrl"
                :srcset="producerThumbnailsSrcset"
              />
            </v-avatar>
          </template>
          <v-list-item-title>{{ producerName }}</v-list-item-title>
          <v-list-item-subtitle>{{ liveTerm }}</v-list-item-subtitle>
          <template #append>
            <v-btn
              variant="outlined"
              color="primary"
              size="small"
              @click.stop="handleClickEditButton"
            >
              <v-icon
                size="small"
                :icon="mdiPencil"
              />
            </v-btn>
          </template>
        </v-list-item>
      </v-card-title>

      <v-card-text>
        <v-row>
          <v-col sm="12">
            <p class="text-subtitle-2 text-grey pb-2">
              概要
            </p>
            <p
              class="text-subtitle-2"
              v-html="item.comment"
            />
          </v-col>
          <v-col sm="12">
            <p class="text-subtitle-2 text-grey pb-2">
              関連商品
            </p>

            <molecules-sortable-product-list
              v-model="formData.productIds"
              :products="liveProducts"
            />
          </v-col>
        </v-row>
      </v-card-text>
    </v-card>
  </div>
</template>
