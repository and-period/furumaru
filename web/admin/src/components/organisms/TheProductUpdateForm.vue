<template>
  <form @submit.prevent="handleSubmit">
    <div class="mb-4">
      <v-card class="mb-4">
        <v-card-title>商品ステータス</v-card-title>
        <v-card-text>
          <v-select
            v-model="formData.public"
            label="ステータス"
            :items="statusItems"
          />
        </v-card-text>
      </v-card>

      <v-card class="mb-4">
        <v-card-title>基本情報</v-card-title>
        <v-card-text>
          <the-product-name-and-description-form
            :name.sync="v$.name.$model"
            :description.sync="formData.description"
            :name-error-message="getErrorMessage(v$.name.$errors)"
          />
        </v-card-text>
      </v-card>

      <v-card class="mb-4">
        <v-card-title>在庫</v-card-title>
        <v-card-text>
          <the-product-inventory-form
            :inventory.sync="v$.inventory.$model"
            :item-unit.sync="formData.itemUnit"
            :item-description.sync="formData.itemDescription"
            :inventory-error-message="getErrorMessage(v$.inventory.$errors)"
          />
        </v-card-text>
      </v-card>

      <v-card class="mb-4">
        <v-card-title>商品画像登録</v-card-title>
        <v-card-text>
          <the-product-media-form
            :media.sync="formData.media"
            @update:files="handleImageUpload"
          />
        </v-card-text>
      </v-card>

      <v-card class="mb-4">
        <v-card-title>詳細情報</v-card-title>
        <v-card-text>
          <the-product-detail-form
            :producer-id.sync="formData.producerId"
            :product-type-id.sync="formData.productTypeId"
            :origin-prefecture.sync="formData.originPrefecture"
            :origin-city.sync="formData.originCity"
            :producers-items="producersItems"
            :product-types-items="productTypesItems"
          />
        </v-card-text>
      </v-card>

      <v-card class="mb-4">
        <v-card-title>価格</v-card-title>
        <v-card-text>
          <v-text-field
            v-model.number="formData.price"
            label="販売価格"
            type="number"
          />
        </v-card-text>
      </v-card>

      <v-card class="mb-4">
        <v-card-title>配送情報</v-card-title>
        <v-card-text>
          <the-product-delivery-form
            :weight.sync="formData.weight"
            :delivery-type.sync="formData.deliveryType"
            :box60-rate.sync="v$.box60Rate.$model"
            :box80-rate.sync="v$.box80Rate.$model"
            :box100-rate.sync="v$.box100Rate.$model"
            :box60-rate-error-message="getErrorMessage(v$.box60Rate.$errors)"
            :box80-rate-error-message="getErrorMessage(v$.box80Rate.$errors)"
            :box100-rate-error-message="getErrorMessage(v$.box100Rate.$errors)"
          />
        </v-card-text>
      </v-card>

      <v-btn block outlined color="primary" type="submit">
        <v-icon left>mdi-plus</v-icon>
        更新
      </v-btn>
    </div>
  </form>
</template>

<script lang="ts">
import { computed, defineComponent, PropType } from '@vue/composition-api'
import { useVuelidate } from '@vuelidate/core'

import {
  required,
  getErrorMessage,
  maxLength,
  minValue,
  maxValue,
} from '~/lib/validations'
import { UpdateProductRequest } from '~/types/api'

export default defineComponent({
  props: {
    formData: {
      type: Object as PropType<UpdateProductRequest>,
      default: () => {
        return {
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
          originCity: '',
        }
      },
    },
    producersItems: {
      type: Array,
      default: () => {
        return []
      },
    },
    productTypesItems: {
      type: Array,
      default: () => {
        return []
      },
    },
  },

  setup(props, { emit }) {
    const statusItems = [
      { text: '公開', value: true },
      { text: '非公開', value: false },
    ]
    const deliveryTypeItems = [
      { text: '通常便', value: 1 },
      { text: '冷蔵便', value: 2 },
      { text: '冷凍便', value: 3 },
    ]

    const formDataValue = computed({
      get: (): UpdateProductRequest => props.formData,
      set: (val: UpdateProductRequest) => emit('update:formData', val),
    })

    const rules = computed(() => {
      return {
        name: { required, maxLength: maxLength(128) },
        inventory: { required, minValue: minValue(0) },
        price: { required, minValue: minValue(0) },
        weight: { required, minValue: minValue(0) },
        box60Rate: { minValue: minValue(0), maxValue: maxValue(100) },
        box80Rate: { minValue: minValue(0), maxValue: maxValue(100) },
        box100Rate: { minValue: minValue(0), maxValue: maxValue(100) },
      }
    })

    const v$ = useVuelidate<UpdateProductRequest>(rules, formDataValue)

    const handleImageUpload = (files: FileList) => {
      emit('update:files', files)
    }

    const handleSubmit = async () => {
      const result = await v$.value.$validate()
      if (!result) {
        window.scrollTo({
          top: 0,
          behavior: 'smooth',
        })
        return
      }
      emit('submit')
    }

    return {
      // 定数
      statusItems,
      deliveryTypeItems,
      // リアクティブ変数
      v$,
      formDataValue,
      // 関数
      getErrorMessage,
      handleImageUpload,
      handleSubmit,
    }
  },
})
</script>
