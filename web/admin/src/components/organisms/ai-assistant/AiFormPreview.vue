<script lang="ts" setup>
import { mdiCheck, mdiClose } from '@mdi/js'
import type { FormFieldChange } from '~/types/ai'

const DELIVERY_TYPE_LABELS: Record<number, string> = {
  1: '常温便',
  2: '冷蔵便',
  3: '冷凍便',
}

const STORAGE_METHOD_LABELS: Record<number, string> = {
  1: '常温保存',
  2: '冷暗所保存',
  3: '冷蔵保存',
  4: '冷凍保存',
}

const SCOPE_LABELS: Record<number, string> = {
  1: '全体公開',
  2: 'LINE限定',
  3: '下書き',
}

defineProps<{
  changes: FormFieldChange[]
  toolName: string
}>()

const emit = defineEmits<{
  approve: []
  reject: []
}>()

function formatValue(field: string, value: unknown): string {
  if (value === '' || value === null || value === undefined) {
    return '（未設定）'
  }
  if (field === 'deliveryType') {
    return DELIVERY_TYPE_LABELS[value as number] || String(value)
  }
  if (field === 'storageMethodType') {
    return STORAGE_METHOD_LABELS[value as number] || String(value)
  }
  if (field === 'scope') {
    return SCOPE_LABELS[value as number] || String(value)
  }
  if (field === 'price' || field === 'cost') {
    return `¥${Number(value).toLocaleString()}`
  }
  if (field === 'weight') {
    return `${value} kg`
  }
  if (field === 'expirationDate') {
    return `${value} 日`
  }
  if (field.includes('Rate')) {
    return `${value}%`
  }
  if (typeof value === 'string' && value.length > 100) {
    return `${value.substring(0, 100)}...`
  }
  return String(value)
}
</script>

<template>
  <v-card
    variant="outlined"
    class="ai-form-preview mx-2 mb-2"
    density="compact"
  >
    <v-card-title class="text-body-2 font-weight-bold pa-3 pb-1">
      フォーム更新プレビュー
    </v-card-title>
    <v-card-text class="pa-3 pt-1">
      <v-list
        density="compact"
        class="pa-0"
      >
        <v-list-item
          v-for="change in changes"
          :key="change.field"
          class="pa-0 mb-1"
        >
          <div class="text-caption text-medium-emphasis">
            {{ change.label }}
          </div>
          <div class="text-body-2">
            {{ formatValue(change.field, change.newValue) }}
          </div>
        </v-list-item>
      </v-list>
    </v-card-text>
    <v-card-actions class="pa-3 pt-0">
      <v-btn
        color="primary"
        size="small"
        variant="flat"
        :prepend-icon="mdiCheck"
        @click="emit('approve')"
      >
        反映する
      </v-btn>
      <v-btn
        size="small"
        variant="text"
        :prepend-icon="mdiClose"
        @click="emit('reject')"
      >
        やめる
      </v-btn>
    </v-card-actions>
  </v-card>
</template>

<style scoped>
.ai-form-preview {
  border-color: rgb(var(--v-theme-primary));
  background: rgb(var(--v-theme-primary), 0.04);
}
</style>
