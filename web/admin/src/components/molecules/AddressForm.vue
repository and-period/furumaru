<script lang="ts" setup>
import { prefecturesList } from '~/constants'

interface Props {
  postalCode: string;
  prefectureCode: number;
  city: string;
  addressLine1: string;
  addressLine2: string;
  loading: boolean;
  errorMessage: string;
}

const props = defineProps<Props>()

const emit = defineEmits<{
  (e: 'update:postalCode', postalCode: string): void;
  (e: 'update:prefectureCode', prefecture: number): void;
  (e: 'update:city', city: string): void;
  (e: 'update:addressLine1', address: string): void;
  (e: 'update:addressLine2', address: string): void;
  (e: 'click:search'): void;
}>()

const postalCodeValue = computed({
  get: (): string => props.postalCode,
  set: (val: string) => emit('update:postalCode', val)
})
const prefectureValue = computed({
  get: (): number => props.prefectureCode,
  set: (val: number) => emit('update:prefectureCode', val)
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
      <v-btn
        color="primary"
        variant="outlined"
        size="small"
        @click="handleSearch"
      >
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
