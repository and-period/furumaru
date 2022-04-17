<template>
  <client-only>
    <pincode-input v-model="formData" :length="length" placeholder="*" />
  </client-only>
</template>

<script lang="ts">
import {
  computed,
  defineComponent,
  SetupContext,
} from '@nuxtjs/composition-api'

export default defineComponent({
  props: {
    value: {
      type: String,
      required: true,
    },
    length: {
      type: Number,
      default: 4,
    },
  },
  setup(props, { emit }: SetupContext) {
    const formData = computed({
      get: () => props.value,
      set: (val: string) => emit('update:value', val),
    })

    return { formData }
  },
})
</script>
