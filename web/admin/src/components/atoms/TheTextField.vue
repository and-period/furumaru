<script lang="ts" setup>
const props = defineProps({
  autofocus: {
    type: Boolean,
    required: false,
    default: false
  },
  label: {
    type: String,
    required: false,
    default: ''
  },
  name: {
    type: String,
    required: false,
    default: ''
  },
  outlined: {
    type: Boolean,
    required: false,
    default: false
  },
  prependIcon: {
    type: String,
    required: false,
    default: undefined
  },
  appendIcon: {
    type: String,
    required: false,
    default: undefined
  },
  readonly: {
    type: Boolean,
    require: false,
    default: false
  },
  rules: {
    type: Object,
    required: false,
    default: () => ({})
  },
  type: {
    type: String,
    required: false,
    default: 'text'
  },
  value: {
    type: String,
    required: false,
    default: ''
  }
})

const emit = defineEmits<{
  (e: 'update:value', str: string): void
}>()

const formData = computed({
  get: () => props.value,
  set: (val: string) => emit('update:value', val)
})
</script>

<template>
  <validation-provider
    v-slot="{ errors, valid }"
    :name="props.label"
    :vid="props.name"
    :rules="props.rules"
  >
    <v-text-field
      v-model="formData"
      :type="props.type"
      :label="props.label"
      :error-messages="errors"
      :success="valid"
      :autofocus="props.autofocus"
      :outlined="props.outlined"
      :readonly="props.readonly"
      :prepend-icon="props.prependIcon"
      :append-icon="props.appendIcon"
    />
  </validation-provider>
</template>
