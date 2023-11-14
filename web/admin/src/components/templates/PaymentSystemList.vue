<script lang="ts" setup>
import type { AlertType } from '~/lib/hooks'
import { PaymentMethodType, type PaymentSystem } from '~/types/api'

const props = defineProps({
  loading: {
    type: Boolean,
    default: false
  },
  isAlert: {
    type: Boolean,
    default: false
  },
  alertType: {
    type: String as PropType<AlertType>,
    default: undefined
  },
  alertText: {
    type: String,
    default: ''
  },
  systems: {
    type: Array<PaymentSystem>,
    default: () => []
  }
})

const emit = defineEmits<{
  (e: 'submit', methodType: PaymentMethodType): void
}>()
</script>

<template>
  <v-alert v-show="isAlert" :type="alertType" v-text="alertText" />

  <v-card class="mt-4" flat :loading="loading">
    <v-card-title>決済システム状態管理</v-card-title>

    <v-card-text>
      <div>TODO: ここにもろもろの実装をする</div>
      <div>systems: {{ systems }}</div>
    </v-card-text>
  </v-card>
</template>
