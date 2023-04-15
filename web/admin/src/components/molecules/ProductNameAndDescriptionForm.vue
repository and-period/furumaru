<script lang="ts" setup>
const props = defineProps({
  name: {
    type: String,
    default: ''
  },
  description: {
    type: String,
    default: ''
  },
  nameErrorMessage: {
    type: String,
    default: ''
  }
})

const emit = defineEmits<{
  (e: 'update:name', name: string): void
  (e: 'update:description', description: string): void
}>()

const nameValue = computed({
  get: () => props.name,
  set: (val: string) => emit('update:name', val)
})

const handleUpdateFormDataDescription = (htmlString: string) => {
  emit('update:description', htmlString)
}
</script>

<template>
  <div>
    <v-text-field
      v-model="nameValue"
      label="商品名"
      :error-messages="props.nameErrorMessage"
    />

    <client-only>
      <tiptap-editor
        label="商品詳細"
        :value="props.description"
        @update:value="handleUpdateFormDataDescription"
      />
    </client-only>
  </div>
</template>
