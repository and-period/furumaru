<script lang="ts" setup>
import { mdiClose, mdiPlus } from '@mdi/js'

import useVuelidate from '@vuelidate/core'
import dayjs, { unix } from 'dayjs'
import type { AlertType } from '~/lib/hooks'
import {

  DeliveryType,
  Prefecture,

  ProductStatus,

  StorageMethodType,

  AdminType,
} from '~/types/api'
import type { Category, Producer, Product, ProductTag, ProductType, UpdateProductRequest } from '~/types/api'
import { getErrorMessage } from '~/lib/validations'
import {
  prefecturesList,
  cityList,

} from '~/constants'
import type { PrefecturesListItem, CityListItem } from '~/constants'
import type { DateTimeInput } from '~/types/props'
import {
  TimeDataValidationRules,
  UpdateProductValidationRules,
} from '~/types/validations'

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
  adminType: {
    type: Number as PropType<AdminType>,
    default: AdminType.UNKNOWN,
  },
  formData: {
    type: Object as PropType<UpdateProductRequest>,
    default: (): UpdateProductRequest => ({
      name: '',
      description: '',
      public: false,
      productTypeId: '',
      productTagIds: [],
      media: [],
      price: 0,
      cost: 0,
      inventory: 0,
      weight: 0,
      itemUnit: '',
      itemDescription: '',
      deliveryType: DeliveryType.NORMAL,
      recommendedPoint1: '',
      recommendedPoint2: '',
      recommendedPoint3: '',
      expirationDate: 0,
      storageMethodType: StorageMethodType.NORMAL,
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
      public: false,
      status: ProductStatus.UNKNOWN,
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
      storageMethodType: StorageMethodType.UNKNOWN,
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
  (e: 'update:selected-category-id', categoryId: string): void
  (e: 'update:search-category', name: string): void
  (e: 'update:search-product-type', name: string): void
  (e: 'update:search-product-tag', name: string): void
  (e: 'submit'): void
}>()

const statuses = [
  { title: '公開', value: true },
  { title: '下書き', value: false },
]
const productStatuses = [
  { title: '予約販売', value: ProductStatus.PRESALE },
  { title: '販売中', value: ProductStatus.FOR_SALE },
  { title: '販売期間外', value: ProductStatus.OUT_OF_SALES },
  { title: '非公開', value: ProductStatus.PRIVATE },
  { title: 'アーカイブ済み', value: ProductStatus.ARCHIVED },
  { title: '不明', value: ProductStatus.UNKNOWN },
]
const storageMethodTypes = [
  { title: '常温保存', value: StorageMethodType.NORMAL },
  { title: '冷暗所保存', value: StorageMethodType.COOL_DARK_PLACE },
  { title: '冷蔵保存', value: StorageMethodType.REFRIGERATED },
  { title: '冷凍保存', value: StorageMethodType.FROZEN },
]
const deliveryTypes = [
  { title: '通常便', value: DeliveryType.NORMAL },
  { title: '冷蔵便', value: DeliveryType.REFRIGERATED },
  { title: '冷凍便', value: DeliveryType.FROZEN },
]
const itemUnits = ['個', '瓶']

const formDataValue = computed({
  get: (): UpdateProductRequest => props.formData,
  set: (v: UpdateProductRequest): void => emit('update:form-data', v),
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
  if (!formDataValue.value.public) {
    return ProductStatus.PRIVATE
  }
  if (!formDataValue.value.startAt || !formDataValue.value.endAt) {
    return ProductStatus.UNKNOWN
  }
  const now = dayjs()
  const startAt = unix(formDataValue.value.startAt)
  const endAt = unix(formDataValue.value.endAt)
  if (now.isBefore(startAt)) {
    return ProductStatus.PRESALE
  }
  if (now.isAfter(endAt)) {
    return ProductStatus.OUT_OF_SALES
  }
  return ProductStatus.FOR_SALE
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

const isUpdatable = (): boolean => {
  if (props.product.status === ProductStatus.ARCHIVED) {
    return false
  }
  const targets: AdminType[] = [AdminType.ADMINISTRATOR, AdminType.COORDINATOR]
  return targets.includes(props.adminType)
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

// Expose the onSubmit method to parent components
defineExpose({
  onSubmit,
})
</script>

<template>
  <v-alert
    v-show="props.isAlert"
    :type="props.alertType"
    v-text="props.alertText"
  />

  <v-row>
    <v-col
      cols="12"
      lg="8"
    >
      <div class="mb-4">
        <v-card
          elevation="0"
          class="mb-4"
        >
          <v-card-title>基本情報</v-card-title>
          <v-card-text>
            <v-select
              v-model="producerIdValue"
              label="生産者名"
              :items="producers"
              item-title="username"
              item-value="id"
              readonly
            />
            <v-text-field
              v-model="formDataValidate.name.$model"
              :error-messages="getErrorMessage(formDataValidate.name.$errors)"
              label="商品名"
              outlined
            />
            <v-textarea
              v-model="formDataValidate.description.$model"
              :error-messages="
                getErrorMessage(formDataValidate.description.$errors)
              "
              label="商品説明"
              maxlength="2000"
            />
          </v-card-text>

          <v-card-subtitle>商品画像登録</v-card-subtitle>
          <v-card-text>
            <v-radio-group
              v-model="thumbnailIndex"
              :error-messages="getErrorMessage(formDataValidate.media.$errors)"
            >
              <molecules-sortable-product-thumbnail
                v-model="formDataValue.media"
                @click="onClickThumbnail"
                @delete="onDeleteThumbnail"
              />
            </v-radio-group>
            <p
              v-show="formDataValue.media.length > 0"
              class="mt-2"
            >
              ※ check された商品画像がサムネイルになります
            </p>
            <div class="mb-2">
              <atoms-file-upload-filed
                text="商品画像"
                @update:files="onClickImageUpload"
              />
            </div>
          </v-card-text>

          <v-card-text>
            <v-text-field
              v-model="formDataValidate.recommendedPoint1.$model"
              :error-messages="
                getErrorMessage(formDataValidate.recommendedPoint1.$errors)
              "
              label="おすすめポイント1"
            />
            <v-text-field
              v-model="formDataValidate.recommendedPoint2.$model"
              :error-messages="
                getErrorMessage(formDataValidate.recommendedPoint2.$errors)
              "
              label="おすすめポイント2"
            />
            <v-text-field
              v-model="formDataValidate.recommendedPoint3.$model"
              :error-messages="
                getErrorMessage(formDataValidate.recommendedPoint3.$errors)
              "
              label="おすすめポイント3"
            />
            <v-text-field
              v-model.number="formDataValidate.expirationDate.$model"
              :error-messages="
                getErrorMessage(formDataValidate.expirationDate.$errors)
              "
              label="賞味期限"
              type="number"
              min="0"
              suffix="日"
            />
            <v-select
              v-model="formDataValidate.storageMethodType.$model"
              :error-messages="
                getErrorMessage(formDataValidate.storageMethodType.$errors)
              "
              label="保存方法"
              :items="storageMethodTypes"
            />
          </v-card-text>
        </v-card>

        <v-card
          elevation="0"
          class="mb-4"
        >
          <v-card-title>価格設定</v-card-title>
          <v-card-text>
            <v-text-field
              v-model.number="formDataValidate.price.$model"
              :error-messages="getErrorMessage(formDataValidate.price.$errors)"
              label="販売価格(税込)"
              type="number"
              min="0"
              suffix="円"
            />
            <v-text-field
              v-model.number="formDataValidate.cost.$model"
              :error-messages="getErrorMessage(formDataValidate.cost.$errors)"
              label="原価(税込)"
              type="number"
              min="0"
              suffix="円"
            />

            <v-table>
              <thead>
                <tr>
                  <th class="text-left">
                    コーディネーター様への支払い金額
                  </th>
                  <th class="text-left">
                    {{ getBenefits().toLocaleString() }} 円
                  </th>
                </tr>
              </thead>
              <tbody>
                <tr>
                  <td>販売価格</td>
                  <td>{{ formDataValue.price.toLocaleString() }} 円</td>
                </tr>
                <tr>
                  <td>原価</td>
                  <td>{{ formDataValue.cost.toLocaleString() }} 円</td>
                </tr>
                <tr>
                  <td>手数料(10%)</td>
                  <td>{{ getCommission().toLocaleString() }} 円</td>
                </tr>
              </tbody>
            </v-table>
          </v-card-text>
        </v-card>

        <v-card
          elevation="0"
          class="mb-4"
        >
          <v-card-title>在庫設定</v-card-title>
          <v-card-text>
            <v-row>
              <v-col cols="9">
                <v-text-field
                  v-model.number="formDataValidate.inventory.$model"
                  :error-messages="
                    getErrorMessage(formDataValidate.inventory.$errors)
                  "
                  label="在庫数"
                  type="number"
                  min="0"
                />
              </v-col>
              <v-col cols="3">
                <v-combobox
                  v-model="formDataValidate.itemUnit.$model"
                  :error-messages="
                    getErrorMessage(formDataValidate.itemUnit.$errors)
                  "
                  label="単位"
                  :items="itemUnits"
                />
              </v-col>
            </v-row>

            <div class="d-flex align-center">
              <v-text-field
                v-model="formDataValidate.itemDescription.$model"
                :error-messages="
                  getErrorMessage(formDataValidate.itemDescription.$errors)
                "
                label="内容説明(発送時に使用)"
                placeholder="1個あたり、3kg程のみかんが入っています。(40~50個)"
              />
            </div>
          </v-card-text>
        </v-card>

        <v-card
          elevation="0"
          class="mb-4"
        >
          <v-card-title>配送設定</v-card-title>
          <v-card-text>
            <v-select
              v-model="formDataValidate.deliveryType.$model"
              :error-messages="
                getErrorMessage(formDataValidate.deliveryType.$errors)
              "
              label="配送種別"
              :items="deliveryTypes"
            />
            <v-text-field
              v-model.number="formDataValidate.weight.$model"
              :error-messages="getErrorMessage(formDataValidate.weight.$errors)"
              label="重さ"
              suffix="kg"
            />
            <v-row>
              <v-col cols="3">
                箱のサイズ
              </v-col>
              <v-col cols="9">
                占有率
              </v-col>
            </v-row>
            <v-row
              v-for="(size, i) in [60, 80, 100]"
              :key="i"
            >
              <v-col
                cols="3"
                align-self="center"
              >
                <p class="mb-0 mx-6 text-h6">
                  {{ size }}
                </p>
              </v-col>
              <v-col cols="9">
                <v-text-field
                  v-model.number="formDataValidate[`box${size}Rate`].$model"
                  :error-messages="
                    getErrorMessage(formDataValidate[`box${size}Rate`].$errors)
                  "
                  label="占有率"
                  type="number"
                  min="0"
                  max="100"
                  suffix="％"
                />
              </v-col>
            </v-row>
          </v-card-text>
        </v-card>
      </div>
    </v-col>

    <v-col
      cols="12"
      lg="4"
    >
      <v-card
        elevation="0"
        class="mb-4"
      >
        <v-card-title>販売設定</v-card-title>
        <v-card-text>
          <v-select
            v-model="productStatus"
            label="販売状況"
            :items="productStatuses"
            item-title="title"
            item-value="value"
            variant="plain"
            readonly
          />
          <v-select
            v-model="formDataValidate.public.$model"
            :error-messages="getErrorMessage(formDataValidate.public.$errors)"
            label="公開状況"
            :items="statuses"
          />
          <p class="text-subtitle-2 text-grey py-2">
            販売開始日時
          </p>
          <div class="d-flex flex-column flex-md-row justify-center">
            <v-text-field
              v-model="startTimeDataValidate.date.$model"
              :error-messages="
                getErrorMessage(startTimeDataValidate.date.$errors)
              "
              type="date"
              variant="outlined"
              density="compact"
              class="mr-md-2"
              @update:model-value="onChangeStartAt"
            />
            <v-text-field
              v-model="startTimeDataValidate.time.$model"
              :error-messages="
                getErrorMessage(startTimeDataValidate.time.$errors)
              "
              type="time"
              variant="outlined"
              density="compact"
              @update:model-value="onChangeStartAt"
            />
          </div>
          <p class="text-subtitle-2 text-grey py-2">
            販売終了日時
          </p>
          <div class="d-flex flex-column flex-md-row justify-center">
            <v-text-field
              v-model="endTimeDataValidate.date.$model"
              :error-messages="
                getErrorMessage(endTimeDataValidate.date.$errors)
              "
              type="date"
              variant="outlined"
              density="compact"
              class="mr-md-2"
              @update:model-value="onChangeEndAt"
            />
            <v-text-field
              v-model="endTimeDataValidate.time.$model"
              :error-messages="
                getErrorMessage(endTimeDataValidate.time.$errors)
              "
              type="time"
              variant="outlined"
              density="compact"
              @update:model-value="onChangeEndAt"
            />
          </div>
        </v-card-text>
      </v-card>

      <v-card
        elevation="0"
        class="mb-4"
      >
        <v-card-title>詳細情報</v-card-title>
        <v-card-text>
          <v-autocomplete
            v-model="selectedCategoryIdValue"
            label="カテゴリ"
            :items="categories"
            item-title="name"
            item-value="id"
            clearable
            @update:search="onChangeSearchCategory"
          />
          <v-autocomplete
            v-model="formDataValidate.productTypeId.$model"
            :error-messages="
              getErrorMessage(formDataValidate.productTypeId.$errors)
            "
            label="品目"
            :items="productTypes"
            item-title="name"
            item-value="id"
            no-data-text="カテゴリを先に選択してください。"
            clearable
            @update:search="onChangeSearchProductType"
          />
          <v-select
            v-model="formDataValidate.originPrefectureCode.$model"
            :error-messages="
              getErrorMessage(formDataValidate.originPrefectureCode.$errors)
            "
            label="原産地（都道府県）"
            :items="prefecturesList"
            item-title="text"
            item-value="value"
          />
          <v-select
            v-model="formDataValidate.originCity.$model"
            :error-messages="
              getErrorMessage(formDataValidate.originCity.$errors)
            "
            :items="cityListItems"
            item-title="text"
            item-value="text"
            label="原産地（市町村）"
            no-data-text="原産地（都道府県）を先に選択してください。"
          />
          <v-autocomplete
            v-model="formDataValidate.productTagIds.$model"
            label="商品タグ"
            :error-messages="
              getErrorMessage(formDataValidate.productTagIds.$errors)
            "
            :items="productTags"
            item-title="name"
            item-value="id"
            chips
            closable-chips
            multiple
            density="comfortable"
            @update:search="onChangeSearchProductTag"
          />
        </v-card-text>
      </v-card>
    </v-col>
  </v-row>
</template>
