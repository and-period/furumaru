<template>
  <div>
    <v-text-field
      v-model="nameValue"
      label="商品名"
      :error-messages="nameErrorMessage"
    />

    <client-only>
      <tiptap-editor
        label="商品詳細"
        :value="description"
        @update:value="handleUpdateFormDataDescription"
      />
    </client-only>
  </div>
</template>

<script lang="ts">
import { computed, defineComponent } from '@vue/composition-api'

export default defineComponent({
  props: {
    name: {
      type: String,
      default: '',
    },
    description: {
      type: String,
      default: '',
    },
    nameErrorMessage: {
      type: String,
      default: '',
    },
  },

  setup(props, { emit }) {
    const nameValue = computed({
      get: () => props.name,
      set: (val: string) => emit('update:name', val),
    })

    const handleUpdateFormDataDescription = (htmlString: string) => {
      emit('update:description', htmlString)
    }

    return {
      nameValue,
      handleUpdateFormDataDescription,
    }
  },
})
</script>
