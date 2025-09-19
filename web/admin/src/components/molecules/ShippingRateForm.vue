<script setup lang='ts'>
import useVuelidate from '@vuelidate/core'
import type { PrefecturesListSelectItems } from '~/lib/prefectures'
import { getErrorMessage } from '~/lib/validations'
import type { CreateShippingRate, UpdateDefaultShippingRate, UpdateShippingRate } from '~/types/api/v1'
import { ShippingRateValidationRules } from '~/types/validations'

interface Props {
  modelValue: CreateShippingRate | UpdateShippingRate | UpdateDefaultShippingRate
  selectablePrefectureList: PrefecturesListSelectItems[]
}

const props = defineProps<Props>()

interface Emits {
  (e: 'update:modelValue', val: CreateShippingRate | UpdateShippingRate | UpdateDefaultShippingRate): void
  (e: 'click:selectAll'): void
}

const emits = defineEmits<Emits>()

const formDataValue = computed({
  get: () => props.modelValue,
  set: (val: CreateShippingRate | UpdateShippingRate | UpdateDefaultShippingRate) => emits('update:modelValue', val),
})

const v$ = useVuelidate<CreateShippingRate | UpdateShippingRate | UpdateDefaultShippingRate>(ShippingRateValidationRules, formDataValue)

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
    v-model="v$.prefectureCodes.$model"
    :error-messages="getErrorMessage(v$.prefectureCodes.$errors)"
    label="都道府県"
    chips
    multiple
    :items="selectablePrefectureList"
    item-title="text"
    item-value="value"
  >
    <template #prepend-item>
      <v-list-item
        ripple
        @click="handleClickSelectAll"
        @mousedown.prevent
      >
        <v-list-item-title>すべて選択</v-list-item-title>
      </v-list-item>
    </template>
  </v-select>
</template>
