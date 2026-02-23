<script lang="ts" setup>
import {
  mdiClose,
  mdiPlus,
  mdiPackageVariant,
  mdiImageMultiple,
  mdiStar,
  mdiCurrencyJpy,
  mdiPackage,
  mdiTruck,
  mdiCalendarClock,
  mdiTagMultiple,
  mdiArrowLeft,
  mdiContentSave,
  mdiFileDocument,
  mdiCommentText,
} from '@mdi/js'

import useVuelidate from '@vuelidate/core'
import dayjs, { unix } from 'dayjs'
import type { AlertType } from '~/lib/hooks'
import { Prefecture } from '~/types'
import {
  DeliveryType,
  ProductStatus,
  StorageMethodType,
  AdminType,
  ProductScope,
} from '~/types/api/v1'
import type { Category, CreateProductReviewRequest, Producer, Product, ProductTag, ProductType, UpdateProductRequest } from '~/types/api/v1'
import { getErrorMessage } from '~/lib/validations'
import {
  prefecturesList,
  cityList, productStatuses, storageMethodTypes, deliveryTypes, productItemUnits,
  productScopes,

} from '~/constants'
import type { PrefecturesListItem, CityListItem } from '~/constants'
import type { DateTimeInput } from '~/types/props'
import {
  NotSameTimeDataValidationRules,
  TimeDataValidationRules,
  UpdateProductValidationRules,
} from '~/types/validations'

const props = defineProps({
  loading: {
    type: Boolean,
    default: false,
  },
  reviewLoading: {
    type: Boolean,
    default: false,
  },
  currentTab: {
    type: Number,
    default: 0,
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
  adminType: {
    type: Number as PropType<AdminType>,
    default: AdminType.AdminTypeUnknown,
  },
  reviewFormData: {
    type: Object as PropType<CreateProductReviewRequest>,
    default: (): CreateProductReviewRequest => ({
      title: '',
      comment: '',
      rate: 5,
    }),
  },
  formData: {
    type: Object as PropType<UpdateProductRequest>,
    default: (): UpdateProductRequest => ({
      name: '',
      description: '',
      scope: ProductScope.ProductScopePublic,
      productTypeId: '',
      productTagIds: [],
      media: [],
      price: 0,
      cost: 0,
      inventory: 0,
      weight: 0,
      itemUnit: '',
      itemDescription: '',
      deliveryType: DeliveryType.DeliveryTypeNormal,
      recommendedPoint1: '',
      recommendedPoint2: '',
      recommendedPoint3: '',
      expirationDate: 0,
      storageMethodType: StorageMethodType.StorageMethodTypeNormal,
      box60Rate: 0,
      box80Rate: 0,
      box100Rate: 0,
      originPrefectureCode: Prefecture.HOKKAIDO,
      originCity: '',
      startAt: dayjs().unix(),
      endAt: dayjs().unix(),
    }),
  },
  product: {
    type: Object as PropType<Product>,
    default: (): Product => ({
      id: '',
      name: '',
      description: '',
      scope: ProductScope.ProductScopePublic,
      status: ProductStatus.ProductStatusUnknown,
      coordinatorId: '',
      producerId: '',
      categoryId: '',
      productTypeId: '',
      productTagIds: [],
      media: [],
      price: 0,
      cost: 0,
      inventory: 0,
      weight: 0,
      itemUnit: '',
      itemDescription: '',
      deliveryType: 0,
      recommendedPoint1: '',
      recommendedPoint2: '',
      recommendedPoint3: '',
      expirationDate: 0,
      storageMethodType: StorageMethodType.StorageMethodTypeUnknown,
      box60Rate: 0,
      box80Rate: 0,
      box100Rate: 0,
      originPrefectureCode: Prefecture.HOKKAIDO,
      originCity: '',
      startAt: dayjs().unix(),
      endAt: dayjs().unix(),
      createdAt: 0,
      updatedAt: 0,
    }),
  },
  selectedCategoryId: {
    type: String,
    default: null,
  },
  producers: {
    type: Array<Producer>,
    default: () => [],
  },
  categories: {
    type: Array<Category>,
    default: () => [],
  },
  productTypes: {
    type: Array<ProductType>,
    default: () => [],
  },
  productTags: {
    type: Array<ProductTag>,
    default: () => [],
  },
})

const emit = defineEmits<{
  (e: 'update:files', files: FileList): void
  (e: 'update:form-data', formData: UpdateProductRequest): void
  (e: 'update:review-form-data', reviewFormData: CreateProductReviewRequest): void
  (e: 'update:selected-category-id', categoryId: string): void
  (e: 'update:current-tab', value: number): void
  (e: 'update:search-category', name: string): void
  (e: 'update:search-product-type', name: string): void
  (e: 'update:search-product-tag', name: string): void
  (e: 'submit'): void
  (e: 'submit:review'): void
}>()

const formDataValue = computed({
  get: (): UpdateProductRequest => props.formData,
  set: (v: UpdateProductRequest): void => emit('update:form-data', v),
})
const reviewFormDataValue = computed({
  get: (): CreateProductReviewRequest => props.reviewFormData,
  set: (v: CreateProductReviewRequest): void => emit('update:review-form-data', v),
})
const currentTabValue = computed({
  get: (): number => props.currentTab,
  set: (v: number): void => emit('update:current-tab', v),
})
const startTimeDataValue = computed({
  get: (): DateTimeInput => ({
    date: unix(props.formData.startAt).format('YYYY-MM-DD'),
    time: unix(props.formData.startAt).format('HH:mm'),
  }),
  set: (timeData: DateTimeInput): void => {
    const startAt = dayjs(`${timeData.date} ${timeData.time}`)
    formDataValue.value.startAt = startAt.unix()
  },
})
const endTimeDataValue = computed({
  get: (): DateTimeInput => ({
    date: unix(props.formData.endAt).format('YYYY-MM-DD'),
    time: unix(props.formData.endAt).format('HH:mm'),
  }),
  set: (timeData: DateTimeInput): void => {
    const endAt = dayjs(`${timeData.date} ${timeData.time}`)
    formDataValue.value.endAt = endAt.unix()
  },
})
const productStatus = computed<ProductStatus>(() => {
  if (formDataValue.value.scope === ProductScope.ProductScopePrivate) {
    return ProductStatus.ProductStatusPrivate
  }
  if (!formDataValue.value.startAt || !formDataValue.value.endAt) {
    return ProductStatus.ProductStatusUnknown
  }
  const now = dayjs()
  const startAt = unix(formDataValue.value.startAt)
  const endAt = unix(formDataValue.value.endAt)
  if (now.isBefore(startAt)) {
    return ProductStatus.ProductStatusPresale
  }
  if (now.isAfter(endAt)) {
    return ProductStatus.ProductStatusOutOfSale
  }
  return ProductStatus.ProductStatusForSale
})
const selectedCategoryIdValue = computed({
  get: (): string => props.selectedCategoryId || '',
  set: (categoryId: string): void =>
    emit('update:selected-category-id', categoryId),
})
const cityListItems = computed(() => {
  const selectedPrefecture = prefecturesList.find(
    (prefecture: PrefecturesListItem): boolean => {
      return props.formData.originPrefectureCode === prefecture.value
    },
  )
  if (!selectedPrefecture) {
    return []
  }
  return cityList.filter(
    (city: CityListItem): boolean => city.prefId === selectedPrefecture.id,
  )
})
const thumbnailIndex = computed<number>({
  get: (): number => props.formData.media.findIndex(item => item.isThumbnail),
  set: (index: number): void => {
    if (formDataValue.value.media.length <= index) {
      return
    }
    formDataValue.value.media = formDataValue.value.media.map((item, i) => {
      if (i === index) {
        return {
          ...item,
          isThumbnail: true,
        }
      }
      else {
        return {
          ...item,
          isThumbnail: false,
        }
      }
    })
  },
})
const producerIdValue = computed(() => props.product.producerId)
const isAdministrator = computed(() => props.adminType === AdminType.AdminTypeAdministrator)

const formDataValidate = useVuelidate(
  UpdateProductValidationRules,
  formDataValue,
)
const startTimeDataValidate = useVuelidate(
  TimeDataValidationRules,
  startTimeDataValue,
)
const endTimeDataValidate = useVuelidate(
  TimeDataValidationRules,
  endTimeDataValue,
)
const notSameTimeValidate = useVuelidate(
  () => NotSameTimeDataValidationRules(props.formData.startAt, '販売開始日時'),
  formDataValue,
)

const getStatus = (status: ProductStatus): string => {
  const value = productStatuses.find(s => s.value === status)
  return value ? value.title : '不明'
}

const getStatusColor = (status: ProductStatus): string => {
  switch (status) {
    case ProductStatus.ProductStatusPresale:
      return 'info'
    case ProductStatus.ProductStatusForSale:
      return 'primary'
    case ProductStatus.ProductStatusOutOfSale:
      return 'secondary'
    case ProductStatus.ProductStatusPrivate:
      return 'warning'
    case ProductStatus.ProductStatusArchived:
      return 'error'
    default:
      return ''
  }
}

const onChangeStartAt = (): void => {
  const startAt = dayjs(
    `${startTimeDataValue.value.date} ${startTimeDataValue.value.time}`,
  )
  formDataValue.value.startAt = startAt.unix()
}

const onChangeEndAt = (): void => {
  const endAt = dayjs(
    `${endTimeDataValue.value.date} ${endTimeDataValue.value.time}`,
  )
  formDataValue.value.endAt = endAt.unix()
}

const getCommission = (): number => {
  return Math.trunc(formDataValue.value.price * 0.1)
}

const getBenefits = (): number => {
  return (
    formDataValue.value.price - (formDataValue.value.cost + getCommission())
  )
}

const onChangeSearchCategory = (name: string): void => {
  emit('update:search-category', name)
}

const onChangeSearchProductType = (name: string): void => {
  emit('update:search-product-type', name)
}

const onChangeSearchProductTag = (name: string): void => {
  emit('update:search-product-tag', name)
}

const onClickImageUpload = (files?: FileList): void => {
  if (!files) {
    return
  }

  emit('update:files', files)
}

const onClickThumbnail = (i: number): void => {
  thumbnailIndex.value = i
}

const onDeleteThumbnail = (i: number): void => {
  const targetItem = props.formData.media.find((_, index) => index === i)
  if (!targetItem) {
    return
  }

  const media = targetItem.isThumbnail
    ? props.formData.media
        .filter((_, index) => index !== i)
        .map((item, i) => {
          return i === 0 ? { ...item, isThumbnail: true } : item
        })
    : props.formData.media.filter((_, index) => index !== i)
  formDataValue.value.media = media
}

const onSubmit = async (): Promise<void> => {
  const formDataValid = await formDataValidate.value.$validate()
  const startTimeDataValid = await startTimeDataValidate.value.$validate()
  const endTimeDataValid = await endTimeDataValidate.value.$validate()
  if (!formDataValid || !startTimeDataValid || !endTimeDataValid) {
    return
  }

  emit('submit')
}

const onSubmitReview = (): void => {
  emit('submit:review')
}
</script>

<template>
  <v-container class="pa-6">
    <atoms-app-alert
      :show="props.isAlert"
      :type="props.alertType"
      :text="props.alertText"
      class="mb-6"
    />

    <!-- ヘッダー -->
    <div class="mb-6">
      <v-btn
        variant="text"
        :prepend-icon="mdiArrowLeft"
        @click="$router.back()"
      >
        戻る
      </v-btn>
      <h1 class="text-h4 font-weight-bold mt-2 mb-2">
        商品編集
      </h1>
      <p class="text-body-1 text-grey-darken-1">
        商品情報を編集・管理します。タブを切り替えて各種設定を行ってください。
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
            :value="0"
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
            v-if="isAdministrator"
            :value="1"
            class="tab-item"
          >
            <v-icon
              :icon="mdiCommentText"
              size="20"
              class="mr-2"
            />
            レビュー投稿
          </v-tab>
        </v-tabs>
      </v-card-title>
      <v-card-text class="pa-0">
        <v-window
          v-model="currentTabValue"
          class="tab-content"
        >
          <v-window-item :value="0">
            <div class="pa-6">
              <v-row>
                <v-col
                  cols="12"
                  lg="8"
                >
                  <!-- 基本情報セクション -->
                  <v-card
                    class="mb-6"
                    variant="outlined"
                  >
                    <v-card-title class="d-flex align-center section-header">
                      <v-icon
                        :icon="mdiPackageVariant"
                        size="24"
                        class="mr-3 text-primary"
                      />
                      <span class="text-h6 font-weight-medium">基本情報</span>
                    </v-card-title>
                    <v-card-text class="pa-6">
                      <v-select
                        v-model="producerIdValue"
                        label="生産者名"
                        :items="producers"
                        item-title="username"
                        item-value="id"
                        variant="outlined"
                        density="comfortable"
                        readonly
                        class="mb-4"
                      />
                      <v-text-field
                        v-model="formDataValidate.name.$model"
                        :error-messages="getErrorMessage(formDataValidate.name.$errors)"
                        label="商品名 *"
                        variant="outlined"
                        density="comfortable"
                        class="mb-4"
                      />
                      <v-textarea
                        v-model="formDataValidate.description.$model"
                        :error-messages="
                          getErrorMessage(formDataValidate.description.$errors)
                        "
                        label="商品説明 *"
                        variant="outlined"
                        density="comfortable"
                        rows="4"
                        maxlength="2000"
                        counter
                      />
                    </v-card-text>
                  </v-card>

                  <!-- 商品画像管理セクション -->
                  <v-card
                    :loading="props.loading"
                    class="form-section-card mb-6"
                    elevation="2"
                  >
                    <v-card-title class="d-flex align-center section-header">
                      <v-icon
                        :icon="mdiImageMultiple"
                        size="24"
                        class="mr-3 text-primary"
                      />
                      <span class="text-h6 font-weight-medium">商品画像管理</span>
                    </v-card-title>
                    <v-card-text class="pa-6">
                      <div class="mb-4">
                        <atoms-file-upload-filed
                          text="商品画像をアップロード"
                          @update:files="onClickImageUpload"
                        />
                      </div>

                      <v-radio-group
                        v-if="formDataValue.media.length > 0"
                        v-model="thumbnailIndex"
                        :error-messages="getErrorMessage(formDataValidate.media.$errors)"
                        class="image-gallery"
                      >
                        <div class="mb-3">
                          <v-chip
                            color="primary"
                            variant="outlined"
                            size="small"
                          >
                            サムネイルを選択してください
                          </v-chip>
                        </div>
                        <v-row>
                          <v-col
                            v-for="(img, i) in formDataValue.media"
                            :key="i"
                            cols="6"
                            sm="4"
                            md="3"
                          >
                            <v-card
                              class="image-card"
                              :class="{ 'thumbnail-selected': img.isThumbnail }"
                              @click="onClickThumbnail(i)"
                            >
                              <v-img
                                :src="img.url"
                                aspect-ratio="1"
                                class="image-preview"
                                alt="商品画像"
                              >
                                <div class="image-overlay">
                                  <v-radio
                                    :value="i"
                                    color="primary"
                                    class="thumbnail-radio"
                                  />
                                  <v-btn
                                    :icon="mdiClose"
                                    color="error"
                                    variant="text"
                                    size="small"
                                    class="delete-btn"
                                    aria-label="画像を削除"
                                    @click.stop="onDeleteThumbnail(i)"
                                  />
                                </div>
                              </v-img>
                              <v-card-text class="pa-2 text-center">
                                <v-chip
                                  v-if="img.isThumbnail"
                                  color="primary"
                                  size="x-small"
                                  variant="elevated"
                                >
                                  サムネイル
                                </v-chip>
                                <span
                                  v-else
                                  class="text-caption text-grey"
                                >
                                  画像 {{ i + 1 }}
                                </span>
                              </v-card-text>
                            </v-card>
                          </v-col>
                        </v-row>
                      </v-radio-group>
                    </v-card-text>
                  </v-card>

                  <!-- おすすめポイント・商品詳細セクション -->
                  <v-card
                    class="form-section-card mb-6"
                    elevation="2"
                  >
                    <v-card-title class="d-flex align-center section-header">
                      <v-icon
                        :icon="mdiStar"
                        size="24"
                        class="mr-3 text-primary"
                      />
                      <span class="text-h6 font-weight-medium">おすすめポイント・商品詳細</span>
                    </v-card-title>
                    <v-card-text class="pa-6">
                      <div class="mb-4">
                        <p class="text-subtitle-2 mb-3 text-grey-darken-1">
                          おすすめポイント
                        </p>
                        <v-text-field
                          v-model="formDataValidate.recommendedPoint1.$model"
                          :error-messages="
                            getErrorMessage(formDataValidate.recommendedPoint1.$errors)
                          "
                          label="ポイント 1"
                          variant="outlined"
                          density="comfortable"
                          class="mb-3"
                        />
                        <v-text-field
                          v-model="formDataValidate.recommendedPoint2.$model"
                          :error-messages="
                            getErrorMessage(formDataValidate.recommendedPoint2.$errors)
                          "
                          label="ポイント 2"
                          variant="outlined"
                          density="comfortable"
                          class="mb-3"
                        />
                        <v-text-field
                          v-model="formDataValidate.recommendedPoint3.$model"
                          :error-messages="
                            getErrorMessage(formDataValidate.recommendedPoint3.$errors)
                          "
                          label="ポイント 3"
                          variant="outlined"
                          density="comfortable"
                        />
                      </div>

                      <v-row>
                        <v-col
                          cols="12"
                          sm="6"
                        >
                          <v-text-field
                            v-model.number="formDataValidate.expirationDate.$model"
                            :error-messages="
                              getErrorMessage(formDataValidate.expirationDate.$errors)
                            "
                            label="賞味期限"
                            type="number"
                            min="0"
                            suffix="日"
                            variant="outlined"
                            density="comfortable"
                          />
                        </v-col>
                        <v-col
                          cols="12"
                          sm="6"
                        >
                          <v-select
                            v-model="formDataValidate.storageMethodType.$model"
                            :error-messages="
                              getErrorMessage(formDataValidate.storageMethodType.$errors)
                            "
                            label="保存方法 *"
                            :items="storageMethodTypes"
                            variant="outlined"
                            density="comfortable"
                          />
                        </v-col>
                      </v-row>
                    </v-card-text>
                  </v-card>

                  <!-- 価格設定セクション -->
                  <v-card
                    class="form-section-card mb-6"
                    elevation="2"
                  >
                    <v-card-title class="d-flex align-center section-header">
                      <v-icon
                        :icon="mdiCurrencyJpy"
                        size="24"
                        class="mr-3 text-primary"
                      />
                      <span class="text-h6 font-weight-medium">価格設定</span>
                    </v-card-title>
                    <v-card-text class="pa-6">
                      <v-row>
                        <v-col
                          cols="12"
                          sm="6"
                        >
                          <v-text-field
                            v-model.number="formDataValidate.price.$model"
                            :error-messages="getErrorMessage(formDataValidate.price.$errors)"
                            label="販売価格(税込) *"
                            type="number"
                            min="0"
                            suffix="円"
                            variant="outlined"
                            density="comfortable"
                          />
                        </v-col>
                        <v-col
                          cols="12"
                          sm="6"
                        >
                          <v-text-field
                            v-model.number="formDataValidate.cost.$model"
                            :error-messages="getErrorMessage(formDataValidate.cost.$errors)"
                            label="原価(税込) *"
                            type="number"
                            min="0"
                            suffix="円"
                            variant="outlined"
                            density="comfortable"
                          />
                        </v-col>
                      </v-row>

                      <v-card
                        variant="outlined"
                        class="profit-calculation mt-4"
                      >
                        <v-card-title class="text-h6 d-flex align-center">
                          <v-icon
                            icon="mdi-calculator"
                            size="20"
                            class="mr-2"
                          />
                          利益計算
                        </v-card-title>
                        <v-card-text>
                          <v-row>
                            <v-col cols="8">
                              コーディネーター様への支払い金額
                            </v-col>
                            <v-col
                              cols="4"
                              class="text-right font-weight-bold text-h6 text-primary"
                            >
                              {{ getBenefits().toLocaleString() }} 円
                            </v-col>
                          </v-row>
                          <v-divider class="my-2" />
                          <v-row dense>
                            <v-col
                              cols="8"
                              class="text-body-2"
                            >
                              販売価格
                            </v-col>
                            <v-col
                              cols="4"
                              class="text-right"
                            >
                              {{ formDataValue.price.toLocaleString() }} 円
                            </v-col>
                          </v-row>
                          <v-row dense>
                            <v-col
                              cols="8"
                              class="text-body-2"
                            >
                              原価
                            </v-col>
                            <v-col
                              cols="4"
                              class="text-right"
                            >
                              -{{ formDataValue.cost.toLocaleString() }} 円
                            </v-col>
                          </v-row>
                          <v-row dense>
                            <v-col
                              cols="8"
                              class="text-body-2"
                            >
                              手数料(10%)
                            </v-col>
                            <v-col
                              cols="4"
                              class="text-right"
                            >
                              -{{ getCommission().toLocaleString() }} 円
                            </v-col>
                          </v-row>
                        </v-card-text>
                      </v-card>
                    </v-card-text>
                  </v-card>

                  <!-- 在庫設定セクション -->
                  <v-card
                    class="form-section-card mb-6"
                    elevation="2"
                  >
                    <v-card-title class="d-flex align-center section-header">
                      <v-icon
                        :icon="mdiPackage"
                        size="24"
                        class="mr-3 text-primary"
                      />
                      <span class="text-h6 font-weight-medium">在庫設定</span>
                    </v-card-title>
                    <v-card-text class="pa-6">
                      <v-row>
                        <v-col cols="9">
                          <v-text-field
                            v-model.number="formDataValidate.inventory.$model"
                            :error-messages="
                              getErrorMessage(formDataValidate.inventory.$errors)
                            "
                            label="在庫数 *"
                            type="number"
                            min="0"
                            variant="outlined"
                            density="comfortable"
                          />
                        </v-col>
                        <v-col cols="3">
                          <v-combobox
                            v-model="formDataValidate.itemUnit.$model"
                            :error-messages="
                              getErrorMessage(formDataValidate.itemUnit.$errors)
                            "
                            label="単位"
                            :items="productItemUnits"
                            variant="outlined"
                            density="comfortable"
                          />
                        </v-col>
                      </v-row>

                      <v-text-field
                        v-model="formDataValidate.itemDescription.$model"
                        :error-messages="
                          getErrorMessage(formDataValidate.itemDescription.$errors)
                        "
                        label="内容説明(発送時に使用)"
                        placeholder="1個あたり、3kg程のみかんが入っています。(40~50個)"
                        variant="outlined"
                        density="comfortable"
                      />
                    </v-card-text>
                  </v-card>

                  <!-- 配送設定セクション -->
                  <v-card
                    class="form-section-card mb-6"
                    elevation="2"
                  >
                    <v-card-title class="d-flex align-center section-header">
                      <v-icon
                        :icon="mdiTruck"
                        size="24"
                        class="mr-3 text-primary"
                      />
                      <span class="text-h6 font-weight-medium">配送設定</span>
                    </v-card-title>
                    <v-card-text class="pa-6">
                      <v-row>
                        <v-col
                          cols="12"
                          sm="6"
                        >
                          <v-select
                            v-model="formDataValidate.deliveryType.$model"
                            :error-messages="
                              getErrorMessage(formDataValidate.deliveryType.$errors)
                            "
                            label="配送種別 *"
                            :items="deliveryTypes"
                            variant="outlined"
                            density="comfortable"
                          />
                        </v-col>
                        <v-col
                          cols="12"
                          sm="6"
                        >
                          <v-text-field
                            v-model.number="formDataValidate.weight.$model"
                            :error-messages="getErrorMessage(formDataValidate.weight.$errors)"
                            label="重さ *"
                            type="number"
                            min="0"
                            step="0.1"
                            suffix="kg"
                            variant="outlined"
                            density="comfortable"
                          />
                        </v-col>
                      </v-row>
                      <div class="box-size-section mt-4">
                        <p class="text-subtitle-2 mb-3 text-grey-darken-1">
                          箱サイズ占有率
                        </p>
                        <v-row
                          v-for="(size, i) in [60, 80, 100]"
                          :key="i"
                        >
                          <v-col
                            cols="4"
                            class="d-flex align-center"
                          >
                            <v-chip
                              color="primary"
                              variant="outlined"
                              size="large"
                            >
                              {{ size }}cm
                            </v-chip>
                          </v-col>
                          <v-col cols="8">
                            <v-text-field
                              v-model.number="formDataValidate[`box${size}Rate`].$model"
                              :error-messages="
                                getErrorMessage(formDataValidate[`box${size}Rate`].$errors)
                              "
                              label="占有率"
                              type="number"
                              min="0"
                              max="100"
                              suffix="%"
                              variant="outlined"
                              density="comfortable"
                            />
                          </v-col>
                        </v-row>
                      </div>
                    </v-card-text>
                  </v-card>
                </v-col>

                <v-col
                  cols="12"
                  lg="4"
                >
                  <!-- 販売設定セクション -->
                  <v-card
                    class="form-section-card mb-6"
                    elevation="2"
                  >
                    <v-card-title class="d-flex align-center section-header">
                      <v-icon
                        :icon="mdiCalendarClock"
                        size="24"
                        class="mr-3 text-primary"
                      />
                      <span class="text-h6 font-weight-medium">販売設定</span>
                    </v-card-title>
                    <v-card-text class="pa-6">
                      <v-alert
                        :color="getStatusColor(productStatus)"
                        variant="tonal"
                        density="compact"
                        class="mb-4"
                      >
                        現在の状況: {{ getStatus(productStatus) }}
                      </v-alert>

                      <v-select
                        v-model="formDataValidate.scope.$model"
                        :error-messages="getErrorMessage(formDataValidate.scope.$errors)"
                        label="公開状況 *"
                        :items="productScopes"
                        variant="outlined"
                        density="comfortable"
                        class="mb-4"
                      />
                      <div class="date-time-section">
                        <p class="text-subtitle-2 mb-3 text-grey-darken-1">
                          販売開始日時 *
                        </p>
                        <v-row>
                          <v-col
                            cols="12"
                            sm="6"
                          >
                            <v-text-field
                              v-model="startTimeDataValidate.date.$model"
                              :error-messages="
                                getErrorMessage(startTimeDataValidate.date.$errors)
                              "
                              label="日付"
                              type="date"
                              variant="outlined"
                              density="comfortable"
                              @update:model-value="onChangeStartAt"
                            />
                          </v-col>
                          <v-col
                            cols="12"
                            sm="6"
                          >
                            <v-text-field
                              v-model="startTimeDataValidate.time.$model"
                              :error-messages="
                                getErrorMessage(startTimeDataValidate.time.$errors)
                              "
                              label="時刻"
                              type="time"
                              variant="outlined"
                              density="comfortable"
                              @update:model-value="onChangeStartAt"
                            />
                          </v-col>
                        </v-row>

                        <p class="text-subtitle-2 mb-3 mt-4 text-grey-darken-1">
                          販売終了日時 *
                        </p>
                        <v-row>
                          <v-col
                            cols="12"
                            sm="6"
                          >
                            <v-text-field
                              v-model="endTimeDataValidate.date.$model"
                              :error-messages="
                                getErrorMessage(endTimeDataValidate.date.$errors)
                              "
                              label="日付"
                              type="date"
                              variant="outlined"
                              density="comfortable"
                              @update:model-value="onChangeEndAt"
                            />
                          </v-col>
                          <v-col
                            cols="12"
                            sm="6"
                          >
                            <v-text-field
                              v-model="endTimeDataValidate.time.$model"
                              :error-messages="
                                getErrorMessage(notSameTimeValidate.endAt.$errors)
                              "
                              label="時刻"
                              type="time"
                              variant="outlined"
                              density="comfortable"
                              @update:model-value="onChangeEndAt"
                            />
                          </v-col>
                        </v-row>
                      </div>
                    </v-card-text>
                  </v-card>

                  <!-- 詳細分類セクション -->
                  <v-card
                    class="form-section-card mb-6"
                    elevation="2"
                  >
                    <v-card-title class="d-flex align-center section-header">
                      <v-icon
                        :icon="mdiTagMultiple"
                        size="24"
                        class="mr-3 text-primary"
                      />
                      <span class="text-h6 font-weight-medium">詳細分類</span>
                    </v-card-title>
                    <v-card-text class="pa-6">
                      <v-autocomplete
                        v-model="selectedCategoryIdValue"
                        label="カテゴリ *"
                        :items="categories"
                        item-title="name"
                        item-value="id"
                        variant="outlined"
                        density="comfortable"
                        class="mb-4"
                        clearable
                        @update:search="onChangeSearchCategory"
                      />
                      <v-autocomplete
                        v-model="formDataValidate.productTypeId.$model"
                        :error-messages="
                          getErrorMessage(formDataValidate.productTypeId.$errors)
                        "
                        label="品目 *"
                        :items="productTypes"
                        item-title="name"
                        item-value="id"
                        variant="outlined"
                        density="comfortable"
                        class="mb-4"
                        no-data-text="カテゴリを先に選択してください。"
                        clearable
                        @update:search="onChangeSearchProductType"
                      />
                      <v-row>
                        <v-col
                          cols="12"
                          sm="6"
                        >
                          <v-autocomplete
                            v-model="formDataValidate.originPrefectureCode.$model"
                            :error-messages="
                              getErrorMessage(formDataValidate.originPrefectureCode.$errors)
                            "
                            label="原産地（都道府県） *"
                            :items="prefecturesList"
                            item-title="text"
                            item-value="value"
                            variant="outlined"
                            density="comfortable"
                          />
                        </v-col>
                        <v-col
                          cols="12"
                          sm="6"
                        >
                          <v-autocomplete
                            v-model="formDataValidate.originCity.$model"
                            :error-messages="
                              getErrorMessage(formDataValidate.originCity.$errors)
                            "
                            :items="cityListItems"
                            item-title="text"
                            item-value="text"
                            label="原産地（市町村）"
                            variant="outlined"
                            density="comfortable"
                            no-data-text="原産地（都道府県）を先に選択してください。"
                          />
                        </v-col>
                      </v-row>
                      <v-autocomplete
                        v-model="formDataValidate.productTagIds.$model"
                        label="商品タグ"
                        :error-messages="
                          getErrorMessage(formDataValidate.productTagIds.$errors)
                        "
                        :items="productTags"
                        item-title="name"
                        item-value="id"
                        variant="outlined"
                        density="comfortable"
                        chips
                        closable-chips
                        multiple
                        @update:search="onChangeSearchProductTag"
                      />
                    </v-card-text>
                  </v-card>
                </v-col>
              </v-row>

              <!-- 送信ボタン -->
              <div class="d-flex justify-end gap-3 mt-8">
                <v-btn
                  variant="text"
                  size="large"
                  @click="$router.back()"
                >
                  キャンセル
                </v-btn>
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
                  変更を保存
                </v-btn>
              </div>
            </div>
          </v-window-item>

          <!-- レビュー投稿タブ -->
          <v-window-item
            v-if="isAdministrator"
            :value="1"
          >
            <div class="pa-6">
              <div class="mb-4">
                <h3 class="text-h6 mb-2">
                  ダミーレビュー投稿
                </h3>
                <p class="text-body-2 text-medium-emphasis">
                  商品のダミーレビューを投稿します。入力したレビューは商品詳細ページに表示されます。
                </p>
              </div>

              <v-card
                class="mb-6"
                variant="outlined"
              >
                <v-card-text class="pa-6">
                  <v-text-field
                    v-model="reviewFormDataValue.title"
                    label="レビュータイトル *"
                    variant="outlined"
                    density="comfortable"
                    class="mb-4"
                    placeholder="例: とても美味しかったです！"
                  />

                  <div class="mb-4">
                    <p class="text-subtitle-2 mb-2 text-grey-darken-1">
                      評価 *
                    </p>
                    <v-rating
                      v-model="reviewFormDataValue.rate"
                      color="amber"
                      active-color="amber"
                      hover
                      length="5"
                      size="large"
                    />
                    <p class="text-body-2 text-medium-emphasis mt-1">
                      {{ reviewFormDataValue.rate }} / 5
                    </p>
                  </div>

                  <v-textarea
                    v-model="reviewFormDataValue.comment"
                    label="レビューコメント *"
                    variant="outlined"
                    density="comfortable"
                    rows="4"
                    maxlength="2000"
                    counter
                    placeholder="例: 新鮮で甘みがあり、家族みんなで美味しくいただきました。また購入したいです。"
                  />
                </v-card-text>
              </v-card>

              <div class="d-flex justify-end gap-3">
                <v-btn
                  :loading="reviewLoading"
                  color="primary"
                  variant="elevated"
                  size="large"
                  @click="onSubmitReview"
                >
                  <v-icon
                    :icon="mdiCommentText"
                    start
                  />
                  レビューを投稿
                </v-btn>
              </div>
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

.image-gallery {
  margin-top: 16px;
}

.image-card {
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s ease;
  border: 2px solid transparent;
}

.image-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 15px rgb(0 0 0 / 10%);
}

.thumbnail-selected {
  border-color: rgb(33 150 243);
  background: rgb(33 150 243 / 5%);
}

.image-preview {
  position: relative;
}

.image-overlay {
  position: absolute;
  top: 4px;
  right: 4px;
  display: flex;
  gap: 4px;
}

.thumbnail-radio {
  background: rgb(255 255 255 / 90%);
  border-radius: 50%;
}

.delete-btn {
  background: rgb(255 255 255 / 90%) !important;
}

.profit-calculation {
  border-radius: 8px;
}

.box-size-section {
  border-top: 1px solid rgb(0 0 0 / 10%);
  padding-top: 16px;
}

.date-time-section {
  border-top: 1px solid rgb(0 0 0 / 10%);
  padding-top: 16px;
}

@media (width <= 600px) {
  .form-section-card {
    border-radius: 8px;
  }

  .section-header {
    padding: 16px 20px;
  }

  .image-card {
    margin-bottom: 16px;
  }
}
</style>
