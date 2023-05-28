<script lang="ts" setup>
import { mdiPlus } from '@mdi/js'
import { AlertType } from '~/lib/hooks'
import { CreateProductRequest, ProducersResponseProducersInner, ProductTypesResponseProductTypesInner } from '~/types/api'

const props = defineProps({
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
    default: () => ({})
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

const formDataValue = computed({
  get: (): CreateProductRequest => props.formData,
  set: (v: CreateProductRequest): void => emit('update:form-data', v)
})

const onClickImageUpload = (files: FileList): void => {
  emit('update:files', files)
}
const onSubmit = (): void => {
  emit('submit')
}
</script>

<template>
  <v-card-title>商品登録</v-card-title>
  <v-breadcrumbs :items="breadcrumbsItem" large class="pa-0 mb-6" />

  <v-alert v-show="props.isAlert" :type="props.alertType" v-text="props.alertText" />

  <organisms-product-form
    v-model:form-data="formDataValue"
    :producers-items="producers"
    :product-types-items="productTypes"
    @update:files="onClickImageUpload"
  />

  <v-btn block variant="outlined" @click="onSubmit">
    <v-icon start :icon="mdiPlus" />
    登録
  </v-btn>
</template>
