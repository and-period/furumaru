<script lang="ts" setup>
import { mdiClose, mdiRobotOutline } from '@mdi/js'
import type { UIMessage } from 'ai'
import type { FormFieldChange } from '~/types/ai'

const props = defineProps<{
  messages: UIMessage[]
  input: string
  loading: boolean
  error: Error | undefined
  hasPendingApproval: boolean
  pendingChanges: FormFieldChange[]
  pendingToolName: string
}>()

const emit = defineEmits<{
  'update:input': [value: string]
  'close': []
  'submit': []
  'approve': []
  'reject': []
}>()

const messagesContainer = ref<HTMLElement>()

function scrollToBottom() {
  nextTick(() => {
    if (messagesContainer.value) {
      messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
    }
  })
}

watch(() => props.messages.length, () => scrollToBottom())

onMounted(() => scrollToBottom())
</script>

<template>
  <v-navigation-drawer
    :model-value="true"
    location="right"
    width="400"
    permanent
    class="ai-chat-panel"
  >
    <!-- Header -->
    <template #prepend>
      <v-toolbar
        density="compact"
        color="primary"
      >
        <v-icon
          :icon="mdiRobotOutline"
          class="ml-4"
        />
        <v-toolbar-title class="text-body-1 font-weight-bold">
          AI アシスタント
        </v-toolbar-title>
        <v-btn
          :icon="mdiClose"
          variant="text"
          size="small"
          @click="emit('close')"
        />
      </v-toolbar>
    </template>

    <!-- Messages -->
    <div
      ref="messagesContainer"
      class="ai-chat-messages flex-grow-1 overflow-y-auto"
    >
      <!-- Welcome message -->
      <div
        v-if="messages.length === 0"
        class="pa-4 text-center"
      >
        <v-icon
          :icon="mdiRobotOutline"
          size="48"
          color="grey-lighten-1"
          class="mb-2"
        />
        <p class="text-body-2 text-medium-emphasis">
          商品情報を自然言語で伝えてください。<br>
          フォームへの入力をお手伝いします。
        </p>
        <div class="mt-3 d-flex flex-column ga-1">
          <v-chip
            size="small"
            variant="outlined"
            class="text-caption"
            @click="emit('update:input', '青森県産のふじりんごを登録したいです')"
          >
            例: 青森県産のふじりんごを登録したいです
          </v-chip>
          <v-chip
            size="small"
            variant="outlined"
            class="text-caption"
            @click="emit('update:input', '商品説明をもっと魅力的にしてください')"
          >
            例: 商品説明をもっと魅力的にしてください
          </v-chip>
          <v-chip
            size="small"
            variant="outlined"
            class="text-caption"
            @click="emit('update:input', 'おすすめポイントを考えてください')"
          >
            例: おすすめポイントを考えてください
          </v-chip>
        </div>
      </div>

      <!-- Chat messages -->
      <organisms-ai-assistant-ai-chat-message
        v-for="message in messages"
        :key="message.id"
        :message="message"
        :pending-changes="pendingChanges"
        :pending-tool-name="pendingToolName"
        :has-pending-approval="hasPendingApproval && message.id === messages[messages.length - 1]?.id"
        @approve="emit('approve')"
        @reject="emit('reject')"
      />

      <!-- Loading indicator -->
      <div
        v-if="loading && !hasPendingApproval"
        class="d-flex align-center ga-2 px-3 py-2"
      >
        <v-progress-circular
          indeterminate
          size="20"
          width="2"
          color="primary"
        />
        <span class="text-body-2 text-medium-emphasis">考え中...</span>
      </div>

      <!-- Error message -->
      <v-alert
        v-if="error"
        type="error"
        density="compact"
        variant="tonal"
        class="mx-3 my-2"
      >
        {{ error.message }}
      </v-alert>
    </div>

    <!-- Input -->
    <template #append>
      <organisms-ai-assistant-ai-chat-input
        :model-value="input"
        :loading="loading"
        :disabled="hasPendingApproval"
        @update:model-value="emit('update:input', $event)"
        @submit="emit('submit')"
      />
    </template>
  </v-navigation-drawer>
</template>

<style scoped>
.ai-chat-panel {
  border-left: 1px solid rgba(var(--v-border-color), var(--v-border-opacity));
}

.ai-chat-messages {
  min-height: 0;
}
</style>
