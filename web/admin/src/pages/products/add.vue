<script lang="ts" setup>
import { mdiPlus, mdiClose } from '@mdi/js'
import { useVuelidate } from '@vuelidate/core'
import { storeToRefs } from 'pinia'

import { useAlert } from '~/lib/hooks'
import { required, minValue, maxLength } from '~/lib/validations'
import {
  useCategoryStore,
  useProducerStore,
  useProductStore,
  useProductTypeStore
} from '~/store'
import { CreateProductRequest, UploadImageResponse } from '~/types/api'
import { prefecturesList, cityList } from '~/constants'

const productTypeStore = useProductTypeStore()
const categoryStore = useCategoryStore()
const producerStore = useProducerStore()

const { producers } = storeToRefs(producerStore)
const { productTypes } = storeToRefs(productTypeStore)

const fetchState = useAsyncData(async () => {
  await Promise.all([
    productTypeStore.fetchProductTypes(),
    categoryStore.fetchCategories(),
    producerStore.fetchProducers(20, 0, '')
  ])
})

const router = useRouter()

const { uploadProductImage, createProduct } = useProductStore()
const breadcrumbsItem = [
  {
    text: '商品管理',
    href: '/products',
    disabled: false
  },
  {
    text: '商品登録',
    href: 'add',
    disabled: true
  }
]

const statusItems = [
  { text: '公開', value: true },
  { text: '非公開', value: false }
]
const itemUnits = [
  { text: '個', value: '個' },
  { text: '瓶', value: '瓶' }
]
const deliveryTypeItems = [
  { text: '通常便', value: 1 },
  { text: '冷蔵便', value: 2 },
  { text: '冷凍便', value: 3 }
]

const formData = reactive<CreateProductRequest>({
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

const rules = computed(() => ({
  name: { required, maxLength: maxLength(128) },
  inventory: { required, minValue: minValue(0) },
  price: { required, minValue: minValue(0) },
  weight: { required, minValue: minValue(0) },
  itemUnit: { required },
  itemDescription: { required }
}))

const v$ = useVuelidate(rules, formData)

const cityListItems = computed(() => {
  const selectedPrefecture = prefecturesList.find(prefecture => formData.originPrefecture === prefecture.value)
  if (!selectedPrefecture) {
    return []
  } else {
    return cityList.filter(city => city.prefId === selectedPrefecture.id)
  }
})

const handleUpdateFormDataDescription = (htmlString: string) => {
  formData.description = htmlString
}

const handleImageUpload = async (files: FileList) => {
  for (const [index, file] of Array.from(files).entries()) {
    try {
      const uploadImage: UploadImageResponse = await uploadProductImage(file)
      formData.media.push({
        ...uploadImage,
        isThumbnail: index === 0
      })
    } catch (error) {
      console.log(error)
    }
  }
}

const handleDeleteThumbnailImageButton = (index: number) => {
  formData.media = formData.media.filter((_, i) => {
    return i !== index
  })
}

const { alertType, isShow, alertText, show } = useAlert('error')

const handleFormSubmit = async () => {
  const result = await v$.value.$validate()
  if (!result) {
    window.scrollTo({
      top: 0,
      behavior: 'smooth'
    })
    return
  }
  try {
    await createProduct(formData)
    router.push('/products')
  } catch (error) {
    show(error.message)
    window.scrollTo({
      top: 0,
      behavior: 'smooth'
    })
  }
}

const getErrorMessage = (key: string): string => {
  const error = v$.value.$errors.find((e) => {
    return e.$property === key
  })
  return error ? `${error.$message}` : ''
}

try {
  await fetchState.execute()
} catch (err) {
  console.log('failed to setup', err)
}
</script>

<template>
  <div>
    <v-card-title>商品登録</v-card-title>
    <v-breadcrumbs :items="breadcrumbsItem" large class="pa-0 mb-6" />

    <v-alert v-model="isShow" class="mb-4" :type="alertType" v-text="alertText" />

    <v-row>
      <v-col cols="8">
        <div class="mb-4">
          <v-card elevation="0" class="mb-4">
            <v-card-title>基本情報</v-card-title>
            <v-card-text>
              <v-select
                v-model="formData.producerId"
                label="販売店舗名"
                :items="producers"
                item-title="storeName"
                item-value="id"
              />

              <v-text-field
                v-model="v$.name.$model"
                label="商品名"
                outlined
                :error="v$.name.$error"
                :error-messages="getErrorMessage('name')"
              />
              <client-only>
                <tiptap-editor
                  label="商品詳細"
                  :value="formData.description"
                  @update:value="handleUpdateFormDataDescription"
                />
              </client-only>
            </v-card-text>
          </v-card>

          <v-card elevation="0" class="mb-4">
            <v-card-title>商品画像登録</v-card-title>
            <v-card-text>
              <v-radio-group>
                <v-row>
                  <v-col
                    v-for="(img, i) in formData.media"
                    :key="i"
                    cols="4"
                    class="d-flex flex-row align-center"
                  >
                    <v-sheet
                      border
                      rounded
                      variant="outlined"
                      width="100%"
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
                    </v-sheet>
                  </v-col>
                </v-row>
                <p v-show="formData.media.length > 0" class="mt-2">
                  ※ check された商品画像がサムネイルになります
                </p>
              </v-radio-group>
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
                v-model="v$.price.$model"
                label="販売価格"
                :error-messages="getErrorMessage('price')"
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
                    :error-messages="getErrorMessage('inventory')"
                    type="number"
                    label="在庫数"
                  />
                </v-col>
                <v-col cols="3">
                  <v-select
                    v-model="formData.itemUnit"
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
                  :error-messages="getErrorMessage('itemDescription')"
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
                  :error-messages="getErrorMessage('weight')"
                >
                  <template #append>
                    kg
                  </template>
                </v-text-field>
                <v-spacer />
              </div>
              <div class="d-flex">
                <v-select
                  v-model="formData.deliveryType"
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
                    v-model="formData[`box${size}Rate`]"
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
              v-model="formData.public"
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
                v-model="formData.productTypeId"
                label="品目"
                :items="productTypes"
                item-title="name"
                item-value="id"
              />
            </div>
            <v-select
              v-model="formData.originPrefecture"
              label="原産地（都道府県）"
              :items="prefecturesList"
              item-title="text"
              item-value="text"
            />
            <v-select
              v-model="formData.originCity"
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

    <v-btn block variant="outlined" @click="handleFormSubmit">
      <v-icon start :icon="mdiPlus" />
      登録
    </v-btn>
  </div>
</template>
