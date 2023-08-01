<script setup lang='ts'>
import useVuelidate from '@vuelidate/core'
import { PrefecturesListSelectItems } from '~/lib/prefectures'
import { required, getErrorMessage, minValue, minLengthArray } from '~/lib/validations'
import { CreateShippingRate } from '~/types/api'

interface Props {
  modelValue: CreateShippingRate
  selectablePrefectureList: PrefecturesListSelectItems[]
}

const props = defineProps<Props>()

interface Emits {
  (e: 'update:modelValue', val: CreateShippingRate): void
  (e: 'click:selectAll'): void
}

const emits = defineEmits<Emits>()

const formDataValue = computed({
  get: () => props.modelValue,
  set: (val: CreateShippingRate) => emits('update:modelValue', val)
})

const rules = computed(() => {
  return {
    name: { required },
    price: { required, minValue: minValue(1) },
    prefectures: { minLengthArray: minLengthArray(1) }
  }
})

const v$ = useVuelidate<CreateShippingRate>(rules, formDataValue)

const handleClickSelectAll = () => {
  emits('click:selectAll')
}
</script>

<template>
  <v-text-field
    v-model="v$.name.$model"
    :error-messages="getErrorMessage(v$.name.$errors)"
    label="名前"
  />
  <v-text-field
    v-model.number="formDataValue.price"
    :error-messages="getErrorMessage(v$.price.$errors)"
    label="配送価格"
    type="number"
    suffix="円"
  />
  <v-select
    v-model="v$.prefectures.$model"
    :error-messages="getErrorMessage(v$.prefectures.$errors)"
    label="都道府県"
    chips
    multiple
    :items="selectablePrefectureList"
    item-title="text"
    item-value="value"
  >
    <template #prepend-item>
      <v-list-item ripple @click="handleClickSelectAll" @mousedown.prevent>
        <v-list-item-title>すべて選択</v-list-item-title>
      </v-list-item>
    </template>
  </v-select>
</template>
