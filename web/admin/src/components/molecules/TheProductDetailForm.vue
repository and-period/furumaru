<template>
  <div>
    <div class="d-flex">
      <v-select
        v-model="productTypeIdValue"
        label="品目"
        :items="productTypesItems"
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
      :items="producersItems"
      item-text="storeName"
      item-value="id"
    />
  </div>
</template>

<script lang="ts">
import { computed, defineComponent } from '@vue/composition-api'

import { prefecturesList, cityList } from '~/constants'

export default defineComponent({
  props: {
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
  },

  setup(props, { emit }) {
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
      return prefecturesList.find(
        (item) => item.text === props.originPrefecture
      )
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

    return {
      // 定数
      prefecturesList,
      // リアクティブ変数
      filteredCityList,
      producerIdValue,
      productTypeIdValue,
      originPrefectureValue,
      originCityValue,
    }
  },
})
</script>
