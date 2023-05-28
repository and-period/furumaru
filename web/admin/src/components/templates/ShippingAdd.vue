<script lang="ts" setup>
import { AlertType } from '~/lib/hooks';
import { CreateShippingRequest } from '~/types/api';

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
    type: Object as PropType<CreateShippingRequest>,
    default: () => ({})
  },
})

const emit = defineEmits<{
  (e: 'update:form-data', v: CreateShippingRequest): void
  (e: 'submit'): void
}>()

const formDataValue = computed({
  get: (): CreateShippingRequest => props.formData,
  set: (v: CreateShippingRequest): void => emit('update:form-data', v)
})

const addBox60RateItem = () => {
  formDataValue.value.box60Rates.push({
    name: '',
    price: 0,
    prefectures: []
  })
}

const addBox80RateItem = () => {
  formDataValue.value.box80Rates.push({
    name: '',
    price: 0,
    prefectures: []
  })
}

const addBox100RateItem = () => {
  formDataValue.value.box100Rates.push({
    name: '',
    price: 0,
    prefectures: []
  })
}

const onClickRemoveItem =(rate: '60' | '80' | '100', index: number) => {
  switch (rate) {
    case '60':
      formDataValue.value.box60Rates.splice(index, 1)
      break
    case '80':
      formDataValue.value.box80Rates.splice(index, 1)
      break
    case '100':
      formDataValue.value.box100Rates.splice(index, 1)
      break
  }
}

const onSubmit = (): void => {
  emit('submit')
}
</script>

<template>
  <v-card-title>配送情報登録</v-card-title>

  <v-alert v-show="props.isAlert" :type="props.alertType" v-text="props.alertText" />

  <organisms-shipping-form
    v-model="formData"
    @click:add-box60-rate-item="addBox60RateItem"
    @click:add-box80-rate-item="addBox80RateItem"
    @click:add-box100-rate-item="addBox100RateItem"
    @click:remove-item-button="onClickRemoveItem"
    @submit="onSubmit"
  />
</template>
