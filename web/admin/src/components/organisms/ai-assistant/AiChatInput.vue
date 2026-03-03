<script lang="ts" setup>
import { mdiSend } from '@mdi/js'

const model = defineModel<string>({ required: true })

defineProps<{
  loading: boolean
  disabled: boolean
}>()

const emit = defineEmits<{
  submit: [event: Event]
}>()

function handleKeydown(event: KeyboardEvent) {
  // Ctrl+Enter or Cmd+Enter to send
  if (event.key === 'Enter' && (event.ctrlKey || event.metaKey)) {
    event.preventDefault()
    emit('submit', event)
  }
}
</script>

<template>
  <form
    class="ai-chat-input pa-3"
    @submit.prevent="emit('submit', $event)"
  >
    <v-textarea
      v-model="model"
      variant="outlined"
      density="compact"
      placeholder="商品情報を入力してください..."
      rows="2"
      auto-grow
      max-rows="5"
      hide-details
      :disabled="disabled"
      @keydown="handleKeydown"
    >
      <template #append-inner>
        <v-btn
          :icon="mdiSend"
          size="small"
          variant="text"
          color="primary"
          :loading="loading"
          :disabled="!model.trim() || disabled"
          type="submit"
        />
      </template>
    </v-textarea>
    <div class="text-caption text-medium-emphasis mt-1">
      Ctrl+Enter で送信
    </div>
  </form>
</template>

<style scoped>
.ai-chat-input {
  border-top: 1px solid rgba(var(--v-border-color), var(--v-border-opacity));
}
</style>
