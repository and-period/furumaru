<script lang="ts" setup>
const props = defineProps({
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
})

const emit = defineEmits<{
  (e: 'update:inventory', inventoty: number): void
  (e: 'update:itemUnit', unit: string): void
  (e: 'update:itemDescription', description: string): void
}>()

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
</script>

<template>
  <div>
    <div class="d-flex">
      <v-text-field
        v-model.number="inventoryValue"
        :error-messages="props.inventoryErrorMessage"
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
        :error-messages="props.itemDescriptionErrorMessage"
      />
      <p class="ml-12 mb-0">ex) 1kg → 5個入り</p>
      <v-spacer />
    </div>
  </div>
</template>
