<template>
  <form>
    <div class="mb-4">
      <v-card class="mb-4">
        <v-card-title>商品ステータス</v-card-title>
        <v-card-text>
          <v-select
            v-model="formData.public"
            label="ステータス"
            :items="statusItems"
          />
        </v-card-text>
      </v-card>

      <v-card>
        <v-card-title>基本情報</v-card-title>
        <v-card-text>
          <v-text-field
            v-model="v$.name.$model"
            label="商品名"
            :error-messages="getErrorMessage(v$.name.$errors)"
          />
        </v-card-text>
      </v-card>

      {{ formDataValue }}
    </div>
  </form>
</template>

<script lang="ts">
import { computed, defineComponent, PropType } from '@vue/composition-api'
import { useVuelidate } from '@vuelidate/core'

import { required, getErrorMessage, maxLength } from '~/lib/validations'
import { UpdateProductRequest } from '~/types/api'

export default defineComponent({
  props: {
    formData: {
      type: Object as PropType<UpdateProductRequest>,
      default: () => {
        return {
          name: '',
          description: '',
          producerId: '',
          productTypeId: '',
          public: false,
          inventory: 0,
          weight: 0,
          itemUnit: '',
          itemDescription: '',
          media: [],
          price: 0,
          deliveryType: 0,
          box60Rate: 0,
          box80Rate: 0,
          box100Rate: 0,
          originPrefecture: '',
          originCity: '',
        }
      },
    },
  },

  setup(props, { emit }) {
    const statusItems = [
      { text: '公開', value: true },
      { text: '非公開', value: false },
    ]
    const deliveryTypeItems = [
      { text: '通常便', value: 1 },
      { text: '冷蔵便', value: 2 },
      { text: '冷凍便', value: 3 },
    ]

    const formDataValue = computed({
      get: (): UpdateProductRequest => props.formData,
      set: (val: UpdateProductRequest) => emit('update:formData', val),
    })

    const rules = computed(() => {
      return {
        name: { required, maxLength: maxLength(128) },
        hasFreeShipping: { required },
      }
    })

    const v$ = useVuelidate<UpdateProductRequest>(rules, formDataValue)

    return {
      // 定数
      statusItems,
      deliveryTypeItems,
      // リアクティブ変数
      v$,
      formDataValue,
      // 関数
      getErrorMessage,
    }
  },
})
</script>
