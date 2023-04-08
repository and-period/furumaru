<script lang="ts" setup>
import { prefecturesList, cityList } from '~/constants'

const props = defineProps({
  producerId: {
    type: String,
    default: '',
  },
  producersItems: {
    type: Array,
    default: () => {
      return []
    },
  },
  productTypeId: {
    type: String,
    default: '',
  },
  productTypesItems: {
    type: Array,
    default: () => {
      return []
    },
  },
  originPrefecture: {
    type: String,
    default: '',
  },
  originCity: {
    type: String,
    default: '',
  },
})

const emit = defineEmits<{
  (e: 'update:producerId', id: string): void
  (e: 'update:productTypeId', id: string): void
  (e: 'update:originPrefecture', prefecture: string): void
  (e: 'update:originCity', city: string): void
}>()

const producerIdValue = computed({
  get: () => props.producerId,
  set: (val: string) => emit('update:producerId', val),
})

const productTypeIdValue = computed({
  get: () => props.productTypeId,
  set: (val: string) => emit('update:productTypeId', val),
})

const originPrefectureValue = computed({
  get: () => props.originPrefecture,
  set: (val: string) => emit('update:originPrefecture', val),
})

const originCityValue = computed({
  get: () => props.originCity,
  set: (val: string) => emit('update:originCity', val),
})

const selectedPrefecture = computed(() => {
  return prefecturesList.find((item) => item.text === props.originPrefecture)
})

const filteredCityList = computed(() => {
  if (selectedPrefecture.value) {
    return cityList.filter(
      (item) => item.prefId === selectedPrefecture.value?.id
    )
  } else {
    return []
  }
})
</script>

<template>
  <div>
    <div class="d-flex">
      <v-select
        v-model="productTypeIdValue"
        label="品目"
        :items="props.productTypesItems"
        item-text="name"
        item-value="id"
      />
    </div>
    <div class="d-flex">
      <v-select
        v-model="originPrefectureValue"
        :items="prefecturesList"
        item-text="text"
        item-value="text"
        class="mr-4"
        label="原産地（都道府県）"
      />
      <v-select
        v-model="originCityValue"
        :items="filteredCityList"
        label="原産地（市町村）"
        messages="先に原産地を選択してください。"
      />
    </div>
    <v-select
      v-model="producerIdValue"
      label="店舗名"
      :items="props.producersItems"
      item-text="storeName"
      item-value="id"
    />
  </div>
</template>
