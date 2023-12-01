<script lang="ts" setup>
import { prefecturesList } from '~/constants'
import { Prefecture } from '~/types/api'

const props = defineProps({
  postalCode: {
    type: String,
    default: ''
  },
  prefectureCode: {
    type: Number as PropType<Prefecture>,
    default: Prefecture.UNKNOWN
  },
  city: {
    type: String,
    default: ''
  },
  addressLine1: {
    type: String,
    default: ''
  },
  addressLine2: {
    type: String,
    default: ''
  },
  loading: {
    type: Boolean,
    default: false
  },
  errorMessage: {
    type: String,
    default: ''
  }
})

const emit = defineEmits<{
  (e: 'update:postalCode', postalCode: any): void
  (e: 'update:prefecture', prefecture: any): void
  (e: 'update:city', city: any): void
  (e: 'update:addressLine1', address: any): void
  (e: 'update:addressLine2', address: any): void
  (e: 'click:search'): void
}>()

const postalCodeValue = computed({
  get: (): string => props.postalCode,
  set: (val: string) => emit('update:postalCode', val)
})
const prefectureValue = computed({
  get: (): number => props.prefectureCode,
  set: (val: number) => emit('update:prefecture', val)
})
const cityValue = computed({
  get: (): string => props.city,
  set: (val: string) => emit('update:city', val)
})
const addressLine1Value = computed({
  get: (): string => props.addressLine1,
  set: (val: string) => emit('update:addressLine1', val)
})
const addressLine2Value = computed({
  get: (): string => props.addressLine2,
  set: (val: string) => emit('update:addressLine2', val)
})

const handleSearch = () => {
  emit('click:search')
}
</script>

<template>
  <div>
    <div class="d-flex align-center">
      <v-text-field
        v-model="postalCodeValue"
        label="郵便番号"
        class="mr-4"
        :loading="props.loading"
        :messages="props.errorMessage"
        :error="props.errorMessage !== ''"
        @keydown.enter="handleSearch"
      />
      <v-btn color="primary" variant="outlined" size="small" @click="handleSearch">
        住所検索
      </v-btn>
      <v-spacer />
    </div>
    <v-select
      v-model="prefectureValue"
      label="都道府県"
      :items="prefecturesList"
      item-title="text"
      item-value="value"
      :loading="props.loading"
    />
    <v-text-field
      v-model="cityValue"
      label="市区町村"
      :loading="props.loading"
    />
    <v-text-field
      v-model="addressLine1Value"
      label="町名・番地"
      :loading="props.loading"
    />
    <v-text-field v-model="addressLine2Value" label="ビル名・号室（任意）" />
  </div>
</template>
