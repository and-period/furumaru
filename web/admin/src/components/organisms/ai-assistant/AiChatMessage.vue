<script lang="ts" setup>
import { mdiRobotOutline, mdiAccountOutline } from '@mdi/js'
import type { UIMessage } from 'ai'
import type { FormFieldChange } from '~/types/ai'

defineProps<{
  message: UIMessage
  pendingChanges: FormFieldChange[]
  pendingToolName: string
  hasPendingApproval: boolean
}>()

defineEmits<{
  approve: []
  reject: []
}>()

interface MessagePart {
  type: string
  text?: string
  toolInvocationId?: string
  toolName?: string
  state?: string
}

function hasToolInvocation(message: UIMessage): boolean {
  return message.parts?.some((p: MessagePart) => p.type === 'tool-invocation') ?? false
}

function getTextContent(message: UIMessage): string {
  if (message.parts && message.parts.length > 0) {
    return message.parts
      .filter((p: MessagePart) => p.type === 'text')
      .map((p: MessagePart) => p.text || '')
      .join('')
  }
  return ''
}
</script>

<template>
  <div
    class="ai-chat-message d-flex ga-2 px-3 py-2"
    :class="{ 'flex-row-reverse': message.role === 'user' }"
  >
    <v-avatar
      :color="message.role === 'user' ? 'primary' : 'grey-lighten-2'"
      size="28"
    >
      <v-icon
        :icon="message.role === 'user' ? mdiAccountOutline : mdiRobotOutline"
        size="16"
        :color="message.role === 'user' ? 'white' : 'grey-darken-2'"
      />
    </v-avatar>
    <div
      class="ai-chat-bubble pa-2 rounded-lg text-body-2"
      :class="message.role === 'user' ? 'ai-chat-bubble-user' : 'ai-chat-bubble-assistant'"
      style="max-width: 85%"
    >
      <div
        v-if="getTextContent(message)"
        class="ai-chat-text"
        style="white-space: pre-wrap; overflow-wrap: break-word;"
      >
        {{ getTextContent(message) }}
      </div>

      <!-- Tool invocation with pending approval -->
      <organisms-ai-assistant-ai-form-preview
        v-if="message.role === 'assistant' && hasToolInvocation(message) && hasPendingApproval"
        :changes="pendingChanges"
        :tool-name="pendingToolName"
        class="mt-2"
        @approve="$emit('approve')"
        @reject="$emit('reject')"
      />

      <!-- Tool invocation already processed -->
      <v-chip
        v-else-if="message.role === 'assistant' && hasToolInvocation(message)"
        size="x-small"
        color="success"
        variant="tonal"
        class="mt-1"
      >
        フォームに反映済み
      </v-chip>
    </div>
  </div>
</template>

<style scoped>
.ai-chat-bubble-user {
  background: rgb(var(--v-theme-primary));
  color: white;
}

.ai-chat-bubble-assistant {
  background: rgb(var(--v-theme-surface-variant));
}
</style>
