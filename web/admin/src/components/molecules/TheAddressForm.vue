<template>
  <div>
    <div class="d-flex align-center">
      <v-text-field
        v-model="postalCodeValue"
        label="郵便番号"
        class="mr-4"
        :loading="loading"
        :messages="errorMessage"
        :error="errorMessage !== ''"
        @keydown.enter="handleSearch"
      />
      <v-btn color="primary" outlined small @click="handleSearch">
        住所検索
      </v-btn>
      <v-spacer />
    </div>
    <v-text-field
      v-model="prefectureValue"
      label="都道府県"
      :loading="loading"
    />
    <v-text-field v-model="cityValue" label="市区町村" :loading="loading" />
    <v-text-field
      v-model="addressLine1Value"
      label="町名・番地"
      :loading="loading"
    />
    <v-text-field v-model="addressLine2Value" label="ビル名・号室（任意）" />
  </div>
</template>

<script lang="ts">
import { computed, defineComponent } from '@vue/composition-api'

export default defineComponent({
  props: {
    postalCode: {
      type: String,
      default: '',
    },
    prefecture: {
      type: String,
      default: '',
    },
    city: {
      type: String,
      default: '',
    },
    addressLine1: {
      type: String,
      default: '',
    },
    addressLine2: {
      type: String,
      default: '',
    },
    loading: {
      type: Boolean,
      default: false,
    },
    errorMessage: {
      type: String,
      default: '',
    },
  },
  setup(props, { emit }) {
    const postalCodeValue = computed({
      get: (): any => props.postalCode,
      set: (val: any) => emit('update:postalCode', val),
    })

    const prefectureValue = computed({
      get: (): any => props.prefecture,
      set: (val: any) => emit('update:prefecture', val),
    })

    const cityValue = computed({
      get: (): any => props.city,
      set: (val: any) => emit('update:city', val),
    })

    const addressLine1Value = computed({
      get: (): any => props.addressLine1,
      set: (val: any) => emit('update:addressLine1', val),
    })

    const addressLine2Value = computed({
      get: (): any => props.addressLine2,
      set: (val: any) => emit('update:addressLine2', val),
    })

    const handleSearch = () => {
      emit('click:search')
    }

    return {
      handleSearch,
      postalCodeValue,
      prefectureValue,
      cityValue,
      addressLine1Value,
      addressLine2Value,
    }
  },
})
</script>
