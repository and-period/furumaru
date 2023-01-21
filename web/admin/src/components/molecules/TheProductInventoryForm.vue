<template>
  <div>
    <div class="d-flex">
      <v-text-field
        v-model.number="inventoryValue"
        :error-messages="inventoryErrorMessage"
        type="number"
        label="在庫数"
      />
      <v-spacer />
    </div>

    <div class="d-flex">
      <v-select v-model="itemUnitValue" label="単位" :items="itemUnitItems" />
      <v-spacer />
    </div>

    <div class="d-flex align-center">
      <v-text-field
        v-model="itemDescriptionValue"
        label="単位説明"
        :error-messages="itemDescriptionErrorMessage"
      />
      <p class="ml-12 mb-0">ex) 1kg → 5個入り</p>
      <v-spacer />
    </div>
  </div>
</template>

<script lang="ts">
import { computed, defineComponent } from '@vue/composition-api'

export default defineComponent({
  props: {
    inventory: {
      type: Number,
      default: 0,
    },
    inventoryErrorMessage: {
      type: String,
      default: '',
    },
    itemUnit: {
      type: String,
      default: '個',
    },
    itemDescription: {
      type: String,
      default: '',
    },
    itemDescriptionErrorMessage: {
      type: String,
      default: '',
    },
  },

  setup(props, { emit }) {
    const itemUnitItems = ['個', '瓶']

    const inventoryValue = computed({
      get: () => props.inventory,
      set: (val: number) => emit('update:inventory', Number(val)),
    })

    const itemUnitValue = computed({
      get: () => props.itemUnit,
      set: (val: string) => emit('update:itemUnit', val),
    })

    const itemDescriptionValue = computed({
      get: () => props.itemDescription,
      set: (val: string) => emit('update:itemDescription', val),
    })

    return {
      // 定数
      itemUnitItems,
      // リアクティブ変数
      inventoryValue,
      itemUnitValue,
      itemDescriptionValue,
    }
  },
})
</script>
