<script lang="ts" setup>
import { mdiPlus } from '@mdi/js'
import { useVuelidate } from '@vuelidate/core'

import {
  required,
  getErrorMessage,
  maxLength,
  minValue,
  maxValue
} from '~/lib/validations'
import { UpdateProductRequest } from '~/types/api'

const props = defineProps({
  formData: {
    type: Object,
    default: (): UpdateProductRequest => ({
      name: '',
      description: '',
      producerId: '',
      productTypeId: '',
      public: false,
      inventory: 0,
      weight: 0,
      itemUnit: '',
      itemDescription: '',
      media: [],
      price: 0,
      deliveryType: 0,
      box60Rate: 0,
      box80Rate: 0,
      box100Rate: 0,
      originPrefecture: '',
      originCity: ''
    })
  },
  producersItems: {
    type: Array,
    default: () => []
  },
  productTypesItems: {
    type: Array,
    default: () => []
  }
})

const emit = defineEmits<{
  (e: 'update:formData', formData: UpdateProductRequest): void
  (e: 'update:files', files?: FileList): void
  (e: 'submit'): void
}>()

const statusItems = [
  { text: '公開', value: true },
  { text: '非公開', value: false }
]
const deliveryTypeItems = [
  { text: '通常便', value: 1 },
  { text: '冷蔵便', value: 2 },
  { text: '冷凍便', value: 3 }
]

const formDataValue = computed({
  get: (): UpdateProductRequest => props.formData as UpdateProductRequest,
  set: (val: UpdateProductRequest) => emit('update:formData', val)
})

const rules = computed(() => {
  return {
    name: { required, maxLength: maxLength(128) },
    inventory: { required, minValue: minValue(0) },
    price: { required, minValue: minValue(0) },
    weight: { required, minValue: minValue(0) },
    box60Rate: { minValue: minValue(0), maxValue: maxValue(100) },
    box80Rate: { minValue: minValue(0), maxValue: maxValue(100) },
    box100Rate: { minValue: minValue(0), maxValue: maxValue(100) }
  }
})

const v$ = useVuelidate<UpdateProductRequest>(rules, formDataValue)

const handleImageUpload = (files?: FileList) => {
  emit('update:files', files)
}

const handleSubmit = async () => {
  const result = await v$.value.$validate()
  if (!result) {
    window.scrollTo({
      top: 0,
      behavior: 'smooth'
    })
    return
  }
  emit('submit')
}
</script>

<template>
  <form @submit.prevent="handleSubmit">
    <div class="mb-4">
      <v-card class="mb-4">
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

      <v-card class="mb-4">
        <v-card-title>基本情報</v-card-title>
        <v-card-text>
          <molecules-product-name-and-description-form
            v-model:name="v$.name.$model"
            v-model:description="formDataValue.description"
            :name-error-message="getErrorMessage(v$.name.$errors)"
          />
        </v-card-text>
      </v-card>

      <v-card class="mb-4">
        <v-card-title>在庫</v-card-title>
        <v-card-text>
          <molecules-product-inventory-form
            v-model:inventory="v$.inventory.$model"
            v-model:item-unit="formDataValue.itemUnit"
            v-model:item-description="formDataValue.itemDescription"
            :inventory-error-message="getErrorMessage(v$.inventory.$errors)"
          />
        </v-card-text>
      </v-card>

      <v-card class="mb-4">
        <v-card-title>商品画像登録</v-card-title>
        <v-card-text>
          <molecules-product-media-form
            v-model:media="formDataValue.media"
            @update:files="handleImageUpload"
          />
        </v-card-text>
      </v-card>

      <v-card class="mb-4">
        <v-card-title>価格</v-card-title>
        <v-card-text>
          <v-text-field
            v-model.number="formDataValue.price"
            label="販売価格"
            type="number"
          />
        </v-card-text>
      </v-card>

      <v-card class="mb-4">
        <v-card-title>配送情報</v-card-title>
        <v-card-text>
          <molecules-product-delivery-form
            v-model:weight="formDataValue.weight"
            v-model:delivery-type="formDataValue.deliveryType"
            v-model:box60-rate="v$.box60Rate.$model"
            v-model:box80-rate="v$.box80Rate.$model"
            v-model:box100-rate="v$.box100Rate.$model"
            :box60-rate-error-message="getErrorMessage(v$.box60Rate.$errors)"
            :box80-rate-error-message="getErrorMessage(v$.box80Rate.$errors)"
            :box100-rate-error-message="getErrorMessage(v$.box100Rate.$errors)"
          />
        </v-card-text>
      </v-card>

      <v-card class="mb-4">
        <v-card-title>詳細情報</v-card-title>
        <v-card-text>
          <molecules-product-detail-form
            v-model:producer-id="formDataValue.producerId"
            v-model:product-type-id="formDataValue.productTypeId"
            v-model:origin-prefecture="formDataValue.originPrefecture"
            v-model:origin-city="formDataValue.originCity"
            :producers-items="producersItems"
            :product-types-items="productTypesItems"
          />
        </v-card-text>
      </v-card>

      <v-btn block variant="outlined" color="primary" type="submit">
        <v-icon start :icon="mdiPlus" />
        更新
      </v-btn>
    </div>
  </form>
</template>
