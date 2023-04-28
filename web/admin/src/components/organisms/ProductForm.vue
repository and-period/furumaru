<script lang="ts" setup>
import { mdiClose } from '@mdi/js'
import { useVuelidate } from '@vuelidate/core'

import {
  required,
  getErrorMessage,
  maxLength,
  minValue,
  maxValue,
  maxLengthArray
} from '~/lib/validations'
import { UpdateProductRequest, CreateProductRequest, ProducersResponse, ProductTypesResponse } from '~/types/api'
import { prefecturesList, cityList } from '~/constants'

type FormData = CreateProductRequest | UpdateProductRequest

interface Props {
  formData: FormData
  producersItems: ProducersResponse['producers'],
  productTypesItems: ProductTypesResponse['productTypes'],

}

const props = defineProps<Props>()

interface Emits {
  (e: 'update:formData', formData: FormData): void
  (e: 'update:files', files: FileList): void
}

const emits = defineEmits<Emits>()

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

const formDataValue = computed({
  get: (): FormData => props.formData,
  set: (val: FormData) => emits('update:formData', val)
})

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
  itemDescription: { required }
}))

const v$ = useVuelidate<FormData>(rules, formDataValue)

const cityListItems = computed(() => {
  const selectedPrefecture = prefecturesList.find(prefecture => props.formData.originPrefecture === prefecture.text)
  if (!selectedPrefecture) {
    return []
  } else {
    return cityList.filter(city => city.prefId === selectedPrefecture.id)
  }
})

const thumbnailIndex = computed<number>({
  get: () => props.formData.media.findIndex(item => item.isThumbnail),
  set: (index: number) => {
    if (index < formDataValue.value.media.length) {
      emits('update:formData', {
        ...formDataValue.value,
        media: formDataValue.value.media.map((item) => {
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
      })
    }
  }
})

const handleImageUpload = (files: FileList) => {
  emits('update:files', files)
}

const handleDeleteThumbnailImageButton = (i: number) => {
  const targetItem = props.formData.media.find((_, index) => index === i)
  if (!targetItem) {
    return
  }

  const newMedia = targetItem.isThumbnail
    ? props.formData.media.filter((_, index) => index !== i).map((item, i) => {
      return i === 0 ? { ...item, isThumbnail: true } : item
    })
    : props.formData.media.filter((_, index) => index !== i)
  emits('update:formData', {
    ...props.formData,
    media: newMedia
  })
}
</script>

<template>
  <v-row>
    <v-col cols="8">
      <div class="mb-4">
        <v-card elevation="0" class="mb-4">
          <v-card-title>基本情報</v-card-title>
          <v-card-text>
            <v-select
              v-model="v$.producerId.$model"
              :error-messages="getErrorMessage(v$.producerId.$errors)"
              label="販売店舗名"
              :items="producersItems"
              item-title="storeName"
              item-value="id"
            />

            <v-text-field
              v-model="v$.name.$model"
              label="商品名"
              outlined
              :error="v$.name.$error"
              :error-messages="getErrorMessage(v$.name.$errors)"
            />
            <client-only>
              <tiptap-editor
                v-model="v$.description.$model"
                :error-message="getErrorMessage(v$.description.$errors)"
                label="商品説明"
                class="mt-4"
              />
            </client-only>
          </v-card-text>
        </v-card>

        <v-card elevation="0" class="mb-4">
          <v-card-title>商品画像登録</v-card-title>
          <v-card-text>
            <v-radio-group v-model="thumbnailIndex" :error-messages="getErrorMessage(v$.media.$errors)">
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
                        <v-btn :icon="mdiClose" color="error" variant="text" size="small" @click="handleDeleteThumbnailImageButton(i)" />
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
                @update:files="handleImageUpload"
              />
            </div>
          </v-card-text>
        </v-card>

        <v-card elevation="0" class="mb-4">
          <v-card-title>価格</v-card-title>
          <v-card-text>
            <v-text-field
              v-model.number="v$.price.$model"
              label="販売価格"
              type="number"
              :error-messages="getErrorMessage(v$.price.$errors)"
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
                  v-model="v$.inventory.$model"
                  :error-messages="getErrorMessage(v$.inventory.$errors)"
                  type="number"
                  label="在庫数"
                />
              </v-col>
              <v-col cols="3">
                <v-combobox
                  v-model="v$.itemUnit.$model"
                  :error-messages="getErrorMessage(v$.itemUnit.$errors)"
                  label="単位"
                  :items="itemUnits"
                  item-title="text"
                  item-value="value"
                />
              </v-col>
            </v-row>

            <div class="d-flex align-center">
              <v-text-field
                v-model="v$.itemDescription.$model"
                label="単位説明"
                :error-messages="getErrorMessage(v$.itemDescription.$errors)"
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
                v-model.number="v$.weight.$model"
                label="重さ"
                :error-messages="getErrorMessage(v$.weight.$errors)"
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
                  v-model.number="v$[`box${size}Rate`].$model"
                  :error-messages="getErrorMessage(v$[`box${size}Rate`].$errors)"
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
              v-model="v$.productTypeId.$model"
              :error-messages="getErrorMessage(v$.productTypeId.$errors)"
              label="品目"
              :items="productTypesItems"
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
</template>

<style lang="scss">
.thumbnail-border {
  border: 2px;
  border-style: solid;
  border-color: rgb(var(--v-theme-secondary));;
}

</style>
