<template>
  <div>
    <div class="d-flex">
      <v-text-field
        v-model.number="weightValue"
        label="重さ"
        :error-messages="weightErrorMessage"
      >
        <template #append>kg</template>
      </v-text-field>
      <v-spacer />
    </div>
    <div class="d-flex">
      <v-select
        v-model.number="deliveryTypeValue"
        :items="deliveryTypeItems"
        label="配送種別"
      />
      <v-spacer />
    </div>

    <v-list>
      <v-list-item>
        <v-list-item-action>箱のサイズ</v-list-item-action>
        <v-list-item-content> 占有率 </v-list-item-content>
      </v-list-item>

      <v-list-item>
        <v-list-item-action>
          <p class="mb-0 mx-6 text-h6">60</p>
        </v-list-item-action>
        <v-list-item-content>
          <v-text-field
            v-model.number="box60RateValue"
            type="number"
            min="0"
            max="100"
            label="占有率"
            :error-messages="box60RateErrorMessage"
          >
            <template #append>%</template>
          </v-text-field>
        </v-list-item-content>
      </v-list-item>

      <v-list-item>
        <v-list-item-action>
          <p class="mb-0 mx-6 text-h6">80</p>
        </v-list-item-action>
        <v-list-item-content>
          <v-text-field
            v-model.number="box80RateValue"
            type="number"
            min="0"
            max="100"
            label="占有率"
            :error-messages="box80RateErrorMessage"
          >
            <template #append>%</template>
          </v-text-field>
        </v-list-item-content>
      </v-list-item>

      <v-list-item>
        <v-list-item-action>
          <p class="mb-0 mx-6 text-h6">100</p>
        </v-list-item-action>
        <v-list-item-content>
          <v-text-field
            v-model.number="box100RateValue"
            type="number"
            min="0"
            max="100"
            label="占有率"
            :error-messages="box100RateErrorMessage"
          >
            <template #append>%</template>
          </v-text-field>
        </v-list-item-content>
      </v-list-item>
    </v-list>
  </div>
</template>

<script lang="ts">
import { computed, defineComponent } from '@vue/composition-api'

export default defineComponent({
  props: {
    weight: {
      type: Number,
      default: 0,
    },
    weightErrorMessage: {
      type: String,
      default: '',
    },
    deliveryType: {
      type: Number,
      default: 0,
    },
    box60Rate: {
      type: Number,
      default: 0,
    },
    box60RateErrorMessage: {
      type: String,
      default: '',
    },
    box80Rate: {
      type: Number,
      default: 0,
    },
    box80RateErrorMessage: {
      type: String,
      default: '',
    },
    box100Rate: {
      type: Number,
      default: 0,
    },
    box100RateErrorMessage: {
      type: String,
      default: '',
    },
  },

  setup(props, { emit }) {
    const deliveryTypeItems = [
      { text: '通常便', value: 1 },
      { text: '冷蔵便', value: 2 },
      { text: '冷凍便', value: 3 },
    ]

    const weightValue = computed({
      get: () => props.weight,
      set: (val: number) => emit('update:weight', Number(val)),
    })

    const deliveryTypeValue = computed({
      get: () => props.deliveryType,
      set: (val: number) => emit('update:deliveryType', val),
    })

    const box60RateValue = computed({
      get: () => props.box60Rate,
      set: (val: number) => emit('update:box60Rate', Number(val)),
    })

    const box80RateValue = computed({
      get: () => props.box80Rate,
      set: (val: number) => emit('update:box80Rate', Number(val)),
    })

    const box100RateValue = computed({
      get: () => props.box100Rate,
      set: (val: number) => emit('update:box100Rate', Number(val)),
    })

    return {
      // 定数
      deliveryTypeItems,
      // リアクティブ変数
      weightValue,
      deliveryTypeValue,
      box60RateValue,
      box80RateValue,
      box100RateValue,
    }
  },
})
</script>
