<script lang="ts" setup>
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
  loading: {
    type: Boolean,
    default: false
  },
  producersItems: {
    type: Array,
    default: () => {
      return []
    }
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

const formDataValue = computed({
  get: (): UpdateProductRequest => props.formData as UpdateProductRequest,
  set: (val: UpdateProductRequest) => emit('update:formData', val)
})

const handleImageUpload = (files?: FileList) => {
  emit('update:files', files)
}

const handleSubmit = () => {
  emit('submit')
}
</script>

<template>
  <div>
    <v-card-title>商品詳細</v-card-title>
    <v-skeleton-loader v-if="props.loading" loading type="article" />
    <div v-else>
      <organisms-product-update-form
        :form-data="formDataValue"
        :producers-items="props.producersItems"
        :product-types-items="props.productTypesItems"
        @update:files="handleImageUpload"
        @submit="handleSubmit"
      />
    </div>
  </div>
</template>
