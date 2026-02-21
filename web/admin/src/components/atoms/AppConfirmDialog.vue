<script lang="ts" setup>
interface Props {
  modelValue: boolean
  title?: string
  message?: string
  loading?: boolean
  confirmColor?: string
  confirmText?: string
}

const props = withDefaults(defineProps<Props>(), {
  title: '削除の確認',
  message: '本当に削除してもよろしいですか？',
  loading: false,
  confirmColor: 'error',
  confirmText: '削除',
})

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  'confirm': []
  'cancel': []
}>()

const dialogValue = computed({
  get: (): boolean => props.modelValue,
  set: (val: boolean): void => emit('update:modelValue', val),
})

function onCancel() {
  emit('cancel')
  dialogValue.value = false
}

function onConfirm() {
  emit('confirm')
}
</script>

<template>
  <v-dialog
    v-model="dialogValue"
    width="500"
  >
    <v-card>
      <v-card-title class="text-h6 py-4">
        {{ title }}
      </v-card-title>
      <v-card-text class="pb-4">
        <div class="text-body-1">
          {{ message }}
        </div>
        <div class="text-body-2 text-medium-emphasis mt-2">
          この操作は取り消せません。
        </div>
      </v-card-text>
      <v-card-actions class="px-6 pb-4">
        <v-spacer />
        <v-btn
          color="medium-emphasis"
          variant="text"
          @click="onCancel"
        >
          キャンセル
        </v-btn>
        <v-btn
          :loading="loading"
          :color="confirmColor"
          variant="elevated"
          @click="onConfirm"
        >
          {{ confirmText }}
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>
