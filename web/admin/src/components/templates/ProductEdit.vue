<script lang="ts" setup>
import { mdiPlus } from '@mdi/js'
import { AlertType } from '~/lib/hooks'
import { ProducersResponseProducersInner, ProductTypesResponseProductTypesInner, UpdateProductRequest } from '~/types/api'

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
    type: Object as PropType<UpdateProductRequest>,
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
  (e: 'update:form-data', formData: UpdateProductRequest): void
  (e: 'submit'): void
}>()

const formDataValue = computed({
  get: (): UpdateProductRequest => props.formData,
  set: (v: UpdateProductRequest): void => emit('update:form-data', v)
})

const onClickImageUpload = (files: FileList): void => {
  emit('update:files', files)
}
const onSubmit = (): void => {
  emit('submit')
}
</script>

<template>
  <v-alert v-show="props.isAlert" :type="props.alertType" v-text="props.alertText" />

  <organisms-product-form
    v-model:form-data="formDataValue"
    :producers-items="producers"
    :product-types-items="productTypes"
    @update:files="onClickImageUpload"
  />

  <v-btn block variant="outlined" @click="onSubmit">
    <v-icon start :icon="mdiPlus" />
    更新
  </v-btn>
</template>
