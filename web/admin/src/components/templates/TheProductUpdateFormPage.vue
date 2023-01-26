<template>
  <div>
    <v-card-title>商品詳細</v-card-title>
    <v-skeleton-loader v-if="loading" :loading="loading" type="article" />
    <div v-else>
      <the-product-update-form
        :form-data="formDataValue"
        :producers-items="producersItems"
        :product-types-items="productTypesItems"
        @update:files="handleImageUpload"
        @submit="handleSubmit"
      />
    </div>
  </div>
</template>

<script lang="ts">
import { computed, defineComponent, PropType } from '@vue/composition-api'

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
    loading: {
      type: Boolean,
      default: false,
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
    const formDataValue = computed({
      get: (): UpdateProductRequest => props.formData,
      set: (val: UpdateProductRequest) => emit('update:formData', val),
    })

    const handleImageUpload = (files: FileList) => {
      emit('update:files', files)
    }

    const handleSubmit = () => {
      emit('submit')
    }

    return {
      formDataValue,
      // 関数
      handleImageUpload,
      handleSubmit,
    }
  },
})
</script>
