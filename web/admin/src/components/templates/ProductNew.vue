<script lang="ts" setup>
import { mdiClose, mdiPlus } from '@mdi/js'
import useVuelidate from '@vuelidate/core'

import { AlertType } from '~/lib/hooks'
import { CreateProductRequest, ProducersResponseProducersInner, ProductTypesResponseProductTypesInner } from '~/types/api'
import {
  required,
  getErrorMessage,
  maxLength,
  minValue,
  maxValue,
  maxLengthArray
} from '~/lib/validations'
import { prefecturesList, cityList } from '~/constants'

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
    type: Object as PropType<CreateProductRequest>,
    default: (): CreateProductRequest => ({
      name: '',
      description: '',
      producerId: '',
      productTypeId: '',
      public: true,
      inventory: 0,
      weight: 0,
      itemUnit: '',
      itemDescription: '',
      media: [],
      price: 0,
      deliveryType: 1,
      box60Rate: 0,
      box80Rate: 0,
      box100Rate: 0,
      originPrefecture: '',
      originCity: ''
    })
  },
  producers: {
    type: Array<ProducersResponseProducersInner>,
    default: () => []
  },
  productTypes: {
    type: Array<ProductTypesResponseProductTypesInner>,
    default: () => []
  }
})

const emit = defineEmits<{
  (e: 'update:files', files: FileList): void
  (e: 'update:form-data', formData: CreateProductRequest): void
  (e: 'submit'): void
}>()

const breadcrumbsItem = [
  {
    title: '商品管理',
    href: '/products',
    disabled: false
  },
  {
    title: '商品登録',
    href: 'add',
    disabled: true
  }
]
const statusItems = [
  { text: '公開', value: true },
  { text: '非公開', value: false }
]
const deliveryTypeItems = [
  { text: '通常便', value: 1 },
  { text: '冷蔵便', value: 2 },
  { text: '冷凍便', value: 3 }
]
const itemUnits = ['個', '瓶']

const rules = computed(() => ({
  name: { required, maxLength: maxLength(128) },
  description: { required },
  media: { maxLengthArray: maxLengthArray(8) },
  producerId: { required },
  productTypeId: { required },
  inventory: { required, minValue: minValue(0) },
  price: { required, minValue: minValue(0) },
  weight: { required, minValue: minValue(0) },
  box60Rate: { required, minValue: minValue(0), maxValue: maxValue(100) },
  box80Rate: { required, minValue: minValue(0), maxValue: maxValue(100) },
  box100Rate: { required, minValue: minValue(0), maxValue: maxValue(100) },
  itemUnit: { required },
  itemDescription: { required },
  deliveryType: {},
  originPrefecture: {},
  originCity: {}
}))
const formDataValue = computed({
  get: (): CreateProductRequest => props.formData,
  set: (formData: CreateProductRequest): void => emit('update:form-data', formData)
})
const cityListItems = computed(() => {
  const selectedPrefecture = prefecturesList.find(prefecture => props.formData.originPrefecture === prefecture.text)
  if (!selectedPrefecture) {
    return []
  } else {
    return cityList.filter(city => city.prefId === selectedPrefecture.id)
  }
})
const thumbnailIndex = computed<number>({
  get: (): number => props.formData.media.findIndex(item => item.isThumbnail),
  set: (index: number): void => {
    if (formDataValue.value.media.length <= index) {
      return
    }
    formDataValue.value.media = formDataValue.value.media.map((item) => {
      return {
        ...item,
        isThumbnail: false
      }
    }).map((item, i) => {
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

const onClickImageUpload = (files?: FileList): void => {
  if (!files) {
    return
  }

  emit('update:files', files)
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

  <v-card-title>商品登録</v-card-title>
  <v-breadcrumbs :items="breadcrumbsItem" large class="pa-0 mb-6" />

  <v-row>
    <v-col cols="8">
      <div class="mb-4">
        <v-card elevation="0" class="mb-4">
          <v-card-title>基本情報</v-card-title>
          <v-card-text>
            <v-select
              v-model="validate.producerId.$model"
              :error-messages="getErrorMessage(validate.producerId.$errors)"
              label="販売店舗名"
              :items="producers"
              item-title="storeName"
              item-value="id"
            />

            <v-text-field
              v-model="validate.name.$model"
              label="商品名"
              outlined
              :error="validate.name.$error"
              :error-messages="getErrorMessage(validate.name.$errors)"
            />
            <client-only>
              <tiptap-editor
                v-model="validate.description.$model"
                :error-message="getErrorMessage(validate.description.$errors)"
                label="商品説明"
                class="mt-4"
              />
            </client-only>
          </v-card-text>
        </v-card>

        <v-card elevation="0" class="mb-4">
          <v-card-title>商品画像登録</v-card-title>
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
                  >
                    <v-img
                      :src="img.url"
                      aspect-ratio="1"
                    >
                      <div class="d-flex col">
                        <v-radio :value="i" />
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
        </v-card>

        <v-card elevation="0" class="mb-4">
          <v-card-title>価格</v-card-title>
          <v-card-text>
            <v-text-field
              v-model.number="validate.price.$model"
              label="販売価格"
              type="number"
              :error-messages="getErrorMessage(validate.price.$errors)"
            >
              <template #prepend>
                &yen;
              </template>
            </v-text-field>
          </v-card-text>
        </v-card>

        <v-card elevation="0" class="mb-4">
          <v-card-title>在庫</v-card-title>
          <v-card-text>
            <v-row>
              <v-col cols="9">
                <v-text-field
                  v-model="validate.inventory.$model"
                  :error-messages="getErrorMessage(validate.inventory.$errors)"
                  type="number"
                  label="在庫数"
                />
              </v-col>
              <v-col cols="3">
                <v-combobox
                  v-model="validate.itemUnit.$model"
                  :error-messages="getErrorMessage(validate.itemUnit.$errors)"
                  label="単位"
                  :items="itemUnits"
                  item-title="text"
                  item-value="value"
                />
              </v-col>
            </v-row>

            <div class="d-flex align-center">
              <v-text-field
                v-model="validate.itemDescription.$model"
                label="単位説明"
                :error-messages="getErrorMessage(validate.itemDescription.$errors)"
              />
              <p class="ml-12 mb-0">
                ex) 1kg → 5個入り
              </p>
              <v-spacer />
            </div>
          </v-card-text>
        </v-card>

        <v-card elevation="0" class="mb-4">
          <v-card-title>配送情報</v-card-title>
          <v-card-text>
            <div class="d-flex">
              <v-text-field
                v-model.number="validate.weight.$model"
                label="重さ"
                :error-messages="getErrorMessage(validate.weight.$errors)"
              >
                <template #append>
                  kg
                </template>
              </v-text-field>
              <v-spacer />
            </div>
            <div class="d-flex">
              <v-select
                v-model="formDataValue.deliveryType"
                :items="deliveryTypeItems"
                item-title="text"
                item-value="value"
                label="配送種別"
              />
              <v-spacer />
            </div>

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
                  type="number"
                  min="0"
                  max="100"
                  label="占有率"
                >
                  <template #append>
                    %
                  </template>
                </v-text-field>
              </v-col>
            </v-row>
          </v-card-text>
        </v-card>
      </div>
    </v-col>

    <v-col cols="4">
      <v-card elevation="0" class="mb-4">
        <v-card-title>商品ステータス</v-card-title>
        <v-card-text>
          <v-select
            v-model="formDataValue.public"
            label="ステータス"
            :items="statusItems"
            item-title="text"
            item-value="value"
          />
        </v-card-text>
      </v-card>

      <v-card elevation="0" class="mb-4">
        <v-card-title>詳細情報</v-card-title>
        <v-card-text>
          <div class="d-flex">
            <v-select
              v-model="validate.productTypeId.$model"
              :error-messages="getErrorMessage(validate.productTypeId.$errors)"
              label="品目"
              :items="productTypes"
              item-title="name"
              item-value="id"
            />
          </div>
          <v-select
            v-model="formDataValue.originPrefecture"
            label="原産地（都道府県）"
            :items="prefecturesList"
            item-title="text"
            item-value="text"
          />
          <v-select
            v-model="formDataValue.originCity"
            :items="cityListItems"
            item-title="text"
            item-value="text"
            label="原産地（市町村）"
            no-data-text="原産地（都道府県）を先に選択してください。"
          />
        </v-card-text>
      </v-card>
    </v-col>
  </v-row>

  <v-btn block variant="outlined" @click="onSubmit">
    <v-icon start :icon="mdiPlus" />
    登録
  </v-btn>
</template>

<style lang="scss">
.thumbnail-border {
  border: 2px;
  border-style: solid;
  border-color: rgb(var(--v-theme-secondary));;
}
</style>
