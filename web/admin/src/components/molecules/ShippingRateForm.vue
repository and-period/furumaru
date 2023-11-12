<script setup lang='ts'>
import useVuelidate from '@vuelidate/core'
import { type PrefecturesListSelectItems } from '~/lib/prefectures'
import { required, getErrorMessage, minValue, minLengthArray } from '~/lib/validations'
import type { UpdateDefaultShippingRate, UpsertShippingRate } from '~/types/api'

interface Props {
  modelValue: UpdateDefaultShippingRate | UpsertShippingRate
  selectablePrefectureList: PrefecturesListSelectItems[]
}

const props = defineProps<Props>()

interface Emits {
  (e: 'update:modelValue', val: UpdateDefaultShippingRate | UpsertShippingRate): void
  (e: 'click:selectAll'): void
}

const emits = defineEmits<Emits>()

const formDataValue = computed({
  get: () => props.modelValue,
  set: (val: UpdateDefaultShippingRate | UpsertShippingRate) => emits('update:modelValue', val)
})

const rules = computed(() => ({
  name: { required },
  price: { required, minValue: minValue(1) },
  prefectureCodes: { minLengthArray: minLengthArray(1) }
}))

const v$ = useVuelidate<UpdateDefaultShippingRate | UpsertShippingRate>(rules, formDataValue)

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
      <v-list-item ripple @click="handleClickSelectAll" @mousedown.prevent>
        <v-list-item-title>すべて選択</v-list-item-title>
      </v-list-item>
    </template>
  </v-select>
</template>
