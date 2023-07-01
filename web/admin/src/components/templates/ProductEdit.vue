<script lang="ts" setup>
import { mdiClose, mdiPlus } from '@mdi/js'

import useVuelidate from '@vuelidate/core'
import { AlertType } from '~/lib/hooks'
import { CategoriesResponseCategoriesInner, DeliveryType, Prefecture, ProducersResponseProducersInner, ProductResponse, ProductTagsResponseProductTagsInner, ProductTypesResponseProductTypesInner, StorageMethodType, UpdateProductRequest } from '~/types/api'
import {
  required,
  getErrorMessage,
  maxLength,
  minValue,
  maxValue,
  maxLengthArray
} from '~/lib/validations'
import { prefecturesList, cityList, PrefecturesListItem, CityListItem } from '~/constants'

const props = defineProps({
  loading: {
    type: Boolean,
    default: false
  },
  isAlert: {
    type: Boolean,
    default: false
  },
  alertType: {
    type: String as PropType<AlertType>,
    default: undefined
  },
  alertText: {
    type: String,
    default: ''
  },
  formData: {
    type: Object as PropType<UpdateProductRequest>,
    default: (): UpdateProductRequest => ({
      name: '',
      description: '',
      public: false,
      producerId: '',
      productTypeId: '',
      productTagIds: [],
      media: [],
      price: 0,
      cost: 0,
      inventory: 0,
      weight: 0,
      itemUnit: '',
      itemDescription: '',
      deliveryType: DeliveryType.UNKNOWN,
      recommendedPoint1: '',
      recommendedPoint2: '',
      recommendedPoint3: '',
      expirationDate: 0,
      storageMethodType: StorageMethodType.UNKNOWN,
      box60Rate: 0,
      box80Rate: 0,
      box100Rate: 0,
      originPrefecture: Prefecture.HOKKAIDO,
      originCity: ''
    })
  },
  product: {
    type: Object as PropType<ProductResponse>,
    default: (): ProductResponse => ({
      id: '',
      name: '',
      description: '',
      public: false,
      producerId: '',
      producerName: '',
      categoryId: '',
      categoryName: '',
      productTypeId: '',
      productTypeName: '',
      productTypeIconUrl: '',
      productTypeIcons: [],
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
      originPrefecture: Prefecture.HOKKAIDO,
      originCity: '',
      createdAt: 0,
      updatedAt: 0
    })
  },
  selectedCategoryId: {
    type: String,
    default: null
  },
  producers: {
    type: Array<ProducersResponseProducersInner>,
    default: () => []
  },
  categories: {
    type: Array<CategoriesResponseCategoriesInner>,
    default: () => []
  },
  productTypes: {
    type: Array<ProductTypesResponseProductTypesInner>,
    default: () => []
  },
  productTags: {
    type: Array<ProductTagsResponseProductTagsInner>,
    default: () => []
  }
})

const emit = defineEmits<{
  (e: 'update:files', files: FileList): void
  (e: 'update:form-data', formData: UpdateProductRequest): void
  (e: 'update:selected-category-id', categoryId: string): void
  (e: 'submit'): void
}>()

const statuses = [
  { title: '公開', value: true },
  { title: '下書き', value: false }
]
const storageMethodTypes = [
  { title: '常温保存', value: StorageMethodType.NORMAL },
  { title: '冷暗所保存', value: StorageMethodType.COOL_DARK_PLACE },
  { title: '冷蔵保存', value: StorageMethodType.REFRIGERATED },
  { title: '冷凍保存', value: StorageMethodType.FROZEN }
]
const deliveryTypes = [
  { title: '通常便', value: DeliveryType.NORMAL },
  { title: '冷蔵便', value: DeliveryType.REFRIGERATED },
  { title: '冷凍便', value: DeliveryType.FROZEN }
]
const itemUnits = ['個', '瓶']

const rules = computed(() => ({
  name: { required, maxLength: maxLength(128) },
  description: { required },
  public: {},
  producerId: { required },
  productTypeId: { required },
  productTagIds: { maxLengthArray: maxLengthArray(8) },
  media: { maxLengthArray: maxLengthArray(8) },
  price: { required, minValue: minValue(0) },
  cost: { required, minValue: minValue(0) },
  inventory: { required, minValue: minValue(0) },
  weight: { required, minValue: minValue(0) },
  itemUnit: { required },
  itemDescription: { required },
  deliveryType: { required },
  recommendedPoint1: { maxValue: maxValue(128) },
  recommendedPoint2: { maxValue: maxValue(128) },
  recommendedPoint3: { maxValue: maxValue(128) },
  expirationDate: { required, minValue: minValue(0) },
  storageMethodType: { required },
  box60Rate: { required, minValue: minValue(0), maxValue: maxValue(100) },
  box80Rate: { required, minValue: minValue(0), maxValue: maxValue(100) },
  box100Rate: { required, minValue: minValue(0), maxValue: maxValue(100) },
  originPrefecture: {},
  originCity: {}
}))
const formDataValue = computed({
  get: (): UpdateProductRequest => props.formData,
  set: (v: UpdateProductRequest): void => emit('update:form-data', v)
})
const selectedCategoryIdValue = computed({
  get: (): string => props.selectedCategoryId || '',
  set: (categoryId: string): void => emit('update:selected-category-id', categoryId)
})
const cityListItems = computed(() => {
  const selectedPrefecture = prefecturesList.find((prefecture: PrefecturesListItem): boolean => {
    return props.formData.originPrefecture === prefecture.value
  })
  if (!selectedPrefecture) {
    return []
  }
  return cityList.filter((city: CityListItem): boolean => city.prefId === selectedPrefecture.id)
})
const thumbnailIndex = computed<number>({
  get: (): number => props.formData.media.findIndex(item => item.isThumbnail),
  set: (index: number): void => {
    if (formDataValue.value.media.length <= index) {
      return
    }
    formDataValue.value.media = formDataValue.value.media
      .map(item => ({
        ...item,
        isThumbnail: false
      }))
      .map((item, i) => {
        if (i !== index) {
          return item
        }
        return {
          ...item,
          isThumbnail: true
        }
      })
  }
})

const validate = useVuelidate(rules, formDataValue)

const getCommission = (): number => {
  return Math.trunc(formDataValue.value.price * 0.1)
}

const getBenefits = (): number => {
  return formDataValue.value.price - (formDataValue.value.cost + getCommission())
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
    ? props.formData.media.filter((_, index) => index !== i).map((item, i) => {
      return i === 0 ? { ...item, isThumbnail: true } : item
    })
    : props.formData.media.filter((_, index) => index !== i)
  formDataValue.value.media = media
}

const onSubmit = (): void => {
  const valid = validate.value.$validate()
  if (!valid) {
    return
  }

  emit('submit')
}
</script>

<template>
  <v-alert v-show="props.isAlert" :type="props.alertType" v-text="props.alertText" />

  <v-row>
    <v-col sm="12" md="12" lg="8">
      <div class="mb-4">
        <v-card elevation="0" class="mb-4">
          <v-card-title>基本情報</v-card-title>
          <v-card-text>
            <v-select
              v-model="validate.producerId.$model"
              :error-messages="getErrorMessage(validate.producerId.$errors)"
              label="生産者名"
              :items="producers"
              item-title="username"
              item-value="id"
            />
            <v-text-field
              v-model="validate.name.$model"
              :error-messages="getErrorMessage(validate.name.$errors)"
              label="商品名"
              outlined
            />
          </v-card-text>

          <v-card-subtitle>商品説明</v-card-subtitle>
          <v-card-text>
            <client-only>
              <tiptap-editor
                v-model="validate.description.$model"
                :error-message="getErrorMessage(validate.description.$errors)"
                class="mb-4"
              />
            </client-only>
          </v-card-text>

          <v-card-subtitle>商品画像登録</v-card-subtitle>
          <v-card-text>
            <v-radio-group v-model="thumbnailIndex" :error-messages="getErrorMessage(validate.media.$errors)">
              <v-row>
                <v-col
                  v-for="(img, i) in formDataValue.media"
                  :key="i"
                  cols="4"
                  class="d-flex flex-row align-center"
                >
                  <v-card
                    rounded
                    variant="outlined"
                    width="100%"
                    :class="{'thumbnail-border': img.isThumbnail }"
                    @click="onClickThumbnail(i)"
                  >
                    <v-img
                      :src="img.url"
                      aspect-ratio="1"
                    >
                      <div class="d-flex col">
                        <v-radio :value="i" color="primary" />
                        <v-btn :icon="mdiClose" color="error" variant="text" size="small" @click="onDeleteThumbnail(i)" />
                      </div>
                    </v-img>
                  </v-card>
                </v-col>
              </v-row>
            </v-radio-group>
            <p v-show="formDataValue.media.length > 0" class="mt-2">
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
              v-model="validate.recommendedPoint1.$model"
              :error-messages="getErrorMessage(validate.recommendedPoint1.$errors)"
              label="おすすめポイント1"
            />
            <v-text-field
              v-model="validate.recommendedPoint2.$model"
              :error-messages="getErrorMessage(validate.recommendedPoint2.$errors)"
              label="おすすめポイント2"
            />
            <v-text-field
              v-model="validate.recommendedPoint3.$model"
              :error-messages="getErrorMessage(validate.recommendedPoint3.$errors)"
              label="おすすめポイント3"
            />
            <v-text-field
              v-model.number="validate.expirationDate.$model"
              :error-messages="getErrorMessage(validate.expirationDate.$errors)"
              label="賞味期限"
              type="number"
              min="0"
              suffix="日"
            />
            <v-select
              v-model="validate.storageMethodType.$model"
              :error-messages="getErrorMessage(validate.storageMethodType.$errors)"
              label="保存方法"
              :items="storageMethodTypes"
            />
          </v-card-text>
        </v-card>

        <v-card elevation="0" class="mb-4">
          <v-card-title>価格設定</v-card-title>
          <v-card-text>
            <v-text-field
              v-model.number="validate.price.$model"
              :error-messages="getErrorMessage(validate.price.$errors)"
              label="販売価格(税込)"
              type="number"
              min="0"
              suffix="円"
            />
            <v-text-field
              v-model.number="validate.cost.$model"
              :error-messages="getErrorMessage(validate.cost.$errors)"
              label="原価"
              type="number"
              min="0"
              suffix="円"
            />

            <v-table>
              <thead>
                <tr>
                  <th class="text-left">
                    コーディネータ様への支払い金額
                  </th>
                  <th class="text-left">
                    {{ getBenefits() }} 円
                  </th>
                </tr>
              </thead>
              <tbody>
                <tr>
                  <td>販売価格</td>
                  <td>{{ formDataValue.price }} 円</td>
                </tr>
                <tr>
                  <td>原価</td>
                  <td>{{ formDataValue.cost }} 円</td>
                </tr>
                <tr>
                  <td>手数料(10%)</td>
                  <td>{{ getCommission() }} 円</td>
                </tr>
              </tbody>
            </v-table>
          </v-card-text>
        </v-card>

        <v-card elevation="0" class="mb-4">
          <v-card-title>在庫設定</v-card-title>
          <v-card-text>
            <v-row>
              <v-col cols="9">
                <v-text-field
                  v-model.number="validate.inventory.$model"
                  :error-messages="getErrorMessage(validate.inventory.$errors)"
                  label="在庫数"
                  type="number"
                  min="0"
                />
              </v-col>
              <v-col cols="3">
                <v-combobox
                  v-model="validate.itemUnit.$model"
                  :error-messages="getErrorMessage(validate.itemUnit.$errors)"
                  label="単位"
                  :items="itemUnits"
                />
              </v-col>
            </v-row>

            <div class="d-flex align-center">
              <v-text-field
                v-model="validate.itemDescription.$model"
                :error-messages="getErrorMessage(validate.itemDescription.$errors)"
                label="内容説明(発送時に使用)"
                placeholder="1個あたり、3kg程のみかんが入っています。(40~50個)"
              />
            </div>
          </v-card-text>
        </v-card>

        <v-card elevation="0" class="mb-4">
          <v-card-title>配送設定</v-card-title>
          <v-card-text>
            <v-text-field
              v-model.number="validate.weight.$model"
              :error-messages="getErrorMessage(validate.weight.$errors)"
              label="重さ"
              suffix="kg"
            />
            <v-select
              v-model="validate.deliveryType.$model"
              :error-messages="getErrorMessage(validate.deliveryType.$errors)"
              label="配送種別"
              :items="deliveryTypes"
            />

            <v-row>
              <v-col cols="3">
                箱のサイズ
              </v-col>
              <v-col cols="9">
                占有率
              </v-col>
            </v-row>
            <v-row v-for="(size, i) in [60, 80, 100]" :key="i">
              <v-col cols="3" align-self="center">
                <p class="mb-0 mx-6 text-h6">
                  {{ size }}
                </p>
              </v-col>
              <v-col cols="9">
                <v-text-field
                  v-model.number="validate[`box${size}Rate`].$model"
                  :error-messages="getErrorMessage(validate[`box${size}Rate`].$errors)"
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

    <v-col sm="12" md="12" lg="4">
      <v-card elevation="0" class="mb-4">
        <v-card-title>商品ステータス</v-card-title>
        <v-card-text>
          <v-select
            v-model="validate.public.$model"
            :error-messages="getErrorMessage(validate.public.$errors)"
            label="ステータス"
            :items="statuses"
          />
        </v-card-text>
      </v-card>

      <v-card elevation="0" class="mb-4">
        <v-card-title>詳細情報</v-card-title>
        <v-card-text>
          <v-autocomplete
            v-model="selectedCategoryIdValue"
            label="カテゴリ"
            :items="categories"
            item-title="name"
            item-value="id"
          />
          <v-autocomplete
            v-model="validate.productTypeId.$model"
            :error-messages="getErrorMessage(validate.productTypeId.$errors)"
            label="品目"
            :items="productTypes"
            item-title="name"
            item-value="id"
            no-data-text="カテゴリを先に選択してください。"
          />
          <v-autocomplete
            v-model="validate.originPrefecture.$model"
            :error-messages="getErrorMessage(validate.originPrefecture.$errors)"
            label="原産地（都道府県）"
            :items="prefecturesList"
            item-title="text"
            item-value="value"
          />
          <v-autocomplete
            v-model="validate.originCity.$model"
            :error-messages="getErrorMessage(validate.originCity.$errors)"
            :items="cityListItems"
            item-title="text"
            item-value="text"
            label="原産地（市町村）"
            no-data-text="原産地（都道府県）を先に選択してください。"
          />
          <v-autocomplete
            v-model="validate.productTagIds.$model"
            label="商品タグ"
            :error-messages="getErrorMessage(validate.productTagIds.$errors)"
            :items="productTags"
            item-title="name"
            item-value="id"
            closable-chips
            multiple
            density="comfortable"
          />
        </v-card-text>
      </v-card>
    </v-col>
  </v-row>

  <v-btn block variant="outlined" @click="onSubmit">
    <v-icon start :icon="mdiPlus" />
    更新
  </v-btn>
</template>
