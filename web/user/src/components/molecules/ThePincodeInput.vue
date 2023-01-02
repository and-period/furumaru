<template>
  <client-only>
    <pincode-input v-model="formData" :length="props.length" placeholder="・" />
  </client-only>
</template>

<script lang="ts" setup>
const props = defineProps({
  value: {
    type: String,
    required: true,
  },
  length: {
    type: Number,
    default: 4,
  },
})

const emits = defineEmits<{
  (name: 'update:value', val: string): void
}>()

const formData = computed({
  get: () => props.value,
  set: (val: string) => emits('update:value', val),
})
</script>

<style lang="scss" scoped>
div.vue-pincode-input-wrapper {
  display: block;
}

::v-deep input.vue-pincode-input {
  background-color: #f7f8f9 !important;
  border: 1px solid rgba(0, 0, 0, 0.5);
  box-sizing: border-box;
  border-radius: 6px;
  box-shadow: none;
}
</style>
