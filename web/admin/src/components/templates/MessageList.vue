<script lang="ts" setup>
import { mdiEmailOpenOutline, mdiEmailOutline, mdiMagnify, mdiInbox, mdiClockOutline } from '@mdi/js'
import dayjs from 'dayjs'
import relativeTime from 'dayjs/plugin/relativeTime'
import 'dayjs/locale/ja'
import type { AlertType } from '~/lib/hooks'
import type { Message } from '~/types/api/v1'

dayjs.extend(relativeTime)
dayjs.locale('ja')

interface Props {
  loading: boolean
  isAlert: boolean
  alertType: AlertType
  alertText: string
  message: Message
  messages: Message[]
}

const props = defineProps<Props>()

const emit = defineEmits<{
  (e: 'click:message', messageId: string): void
}>()

const onClickMessage = (messageId: string): void => {
  emit('click:message', messageId)
}

const selectedMessageId = ref<string>('')
const searchQuery = ref<string>('')
const isMobile = ref<boolean>(false)
const showDetail = ref<boolean>(false)

const parsedMessage = computed<string>(() => {
  if (props.message && props.message.body) {
    return props.message.body.replaceAll('\\n', '\n')
  }
  else {
    return ''
  }
})

const filteredMessages = computed(() => {
  if (!searchQuery.value) return props.messages

  const query = searchQuery.value.toLowerCase()
  return props.messages.filter(msg =>
    msg.title.toLowerCase().includes(query)
    || msg.body?.toLowerCase().includes(query),
  )
})

const formatDate = (timestamp: number) => {
  const date = dayjs.unix(timestamp)
  const now = dayjs()

  if (date.isSame(now, 'day')) {
    return date.format('HH:mm')
  }
  else if (date.isSame(now.subtract(1, 'day'), 'day')) {
    return '昨日'
  }
  else if (date.isAfter(now.subtract(7, 'day'))) {
    return date.fromNow()
  }
  else {
    return date.format('MM/DD')
  }
}

const getMessagePreview = (body?: string) => {
  if (!body) return ''
  const cleanBody = body.replaceAll('\\n', ' ').trim()
  return cleanBody.length > 50 ? cleanBody.substring(0, 50) + '...' : cleanBody
}

const handleMessageClick = (messageId: string) => {
  selectedMessageId.value = messageId
  if (isMobile.value) {
    showDetail.value = true
  }
  onClickMessage(messageId)
}

const handleBackToList = () => {
  showDetail.value = false
}

const handleCloseMessage = () => {
  selectedMessageId.value = ''
  if (isMobile.value) {
    showDetail.value = false
  }
}

onMounted(() => {
  isMobile.value = window.innerWidth < 768
  window.addEventListener('resize', () => {
    isMobile.value = window.innerWidth < 768
  })
})
</script>

<template>
  <v-container
    fluid
    class="pa-0"
  >
    <v-alert
      v-show="isAlert"
      class="mb-4 mx-4"
      :type="alertType"
      v-text="alertText"
    />

    <!-- Header Bar -->
    <v-card
      flat
      class="mb-2"
    >
      <v-card-text class="pb-2">
        <v-row align="center">
          <v-col
            cols="12"
            md="6"
          >
            <h2 class="text-h5 font-weight-bold">
              <v-icon
                color="primary"
                class="mr-2"
              >
                {{ mdiInbox }}
              </v-icon>
              メッセージ
            </h2>
          </v-col>
          <v-col
            cols="12"
            md="6"
          >
            <v-text-field
              v-model="searchQuery"
              :prepend-inner-icon="mdiMagnify"
              placeholder="メッセージを検索"
              variant="outlined"
              density="compact"
              hide-details
              clearable
            />
          </v-col>
        </v-row>
      </v-card-text>
    </v-card>

    <v-row
      no-gutters
      class="message-container"
    >
      <!-- Message List -->
      <v-col
        v-show="!isMobile || !showDetail"
        :cols="isMobile ? 12 : 12"
        :md="5"
        :lg="4"
        class="message-list-column"
      >
        <v-card
          elevation="2"
          class="h-100 d-flex flex-column"
        >
          <v-card-title class="py-3 px-4 bg-grey-lighten-4">
            <div class="d-flex align-center justify-space-between">
              <span class="text-subtitle-1 font-weight-medium">
                受信ボックス
                <v-chip
                  v-if="messages.length > 0"
                  size="x-small"
                  class="ml-2"
                >
                  {{ messages.length }}
                </v-chip>
              </span>
            </div>
          </v-card-title>
          <v-divider />

          <v-list
            v-if="filteredMessages.length > 0"
            lines="two"
            class="flex-grow-1 overflow-y-auto"
          >
            <v-list-item
              v-for="item in filteredMessages"
              :key="item.id"
              :class="{
                'bg-blue-lighten-5': selectedMessageId === item.id,
                'unread-message': !item.read,
              }"
              @click="handleMessageClick(item.id)"
            >
              <template #prepend>
                <v-icon
                  :color="item.read ? 'grey' : 'primary'"
                  class="mr-3"
                >
                  {{ item.read ? mdiEmailOpenOutline : mdiEmailOutline }}
                </v-icon>
              </template>

              <v-list-item-title class="d-flex align-center justify-space-between mb-1">
                <span :class="{ 'font-weight-bold': !item.read }">
                  {{ item.title }}
                </span>
                <v-chip
                  size="x-small"
                  variant="text"
                  class="ml-2"
                >
                  <v-icon
                    size="x-small"
                    class="mr-1"
                  >
                    {{ mdiClockOutline }}
                  </v-icon>
                  {{ formatDate(item.receivedAt) }}
                </v-chip>
              </v-list-item-title>

              <v-list-item-subtitle>
                {{ getMessagePreview(item.body) }}
              </v-list-item-subtitle>

              <template
                v-if="!item.read"
                #append
              >
                <v-badge
                  color="primary"
                  dot
                  inline
                />
              </template>
            </v-list-item>
          </v-list>

          <!-- Empty State -->
          <v-card-text
            v-else
            class="flex-grow-1 d-flex flex-column align-center justify-center text-grey"
          >
            <v-icon
              size="64"
              class="mb-4"
            >
              {{ mdiInbox }}
            </v-icon>
            <p class="text-h6">
              メッセージがありません
            </p>
            <p class="text-body-2">
              新しいメッセージが届くとここに表示されます
            </p>
          </v-card-text>
        </v-card>
      </v-col>

      <!-- Message Detail -->
      <v-col
        v-show="!isMobile || showDetail"
        :cols="isMobile ? 12 : 12"
        :md="7"
        :lg="8"
      >
        <v-card
          elevation="2"
          class="h-100 d-flex flex-column"
        >
          <div v-if="message.title">
            <!-- Message Header -->
            <v-card-title class="bg-white">
              <v-row align="center">
                <v-col cols="12">
                  <div class="d-flex align-center">
                    <v-btn
                      v-if="isMobile"
                      icon
                      variant="text"
                      size="small"
                      class="mr-2"
                      @click="handleBackToList"
                    >
                      <v-icon>mdi-arrow-left</v-icon>
                    </v-btn>
                    <div class="flex-grow-1">
                      <h3 class="text-h6 font-weight-medium mb-1">
                        {{ message.title }}
                      </h3>
                      <div class="text-caption text-grey">
                        <v-icon
                          size="small"
                          class="mr-1"
                        >
                          {{ mdiClockOutline }}
                        </v-icon>
                        受信日時: {{ dayjs.unix(message.receivedAt).format('YYYY年MM月DD日 HH:mm') }}
                      </div>
                    </div>
                    <div class="d-flex align-center ga-2">
                      <v-chip
                        v-if="message.read"
                        color="success"
                        size="small"
                        variant="tonal"
                      >
                        既読
                      </v-chip>
                      <v-chip
                        v-else
                        color="primary"
                        size="small"
                      >
                        未読
                      </v-chip>
                      <v-btn
                        icon
                        variant="text"
                        size="small"
                        class="ml-2"
                        @click="handleCloseMessage"
                      >
                        <v-icon>mdi-close</v-icon>
                        <v-tooltip
                          activator="parent"
                          location="bottom"
                        >
                          メッセージを閉じる
                        </v-tooltip>
                      </v-btn>
                    </div>
                  </div>
                </v-col>
              </v-row>
            </v-card-title>

            <v-divider />

            <!-- Message Body -->
            <v-card-text class="message-content flex-grow-1 overflow-y-auto pa-6">
              <div
                class="message-body"
                v-text="parsedMessage"
              />
            </v-card-text>
          </div>

          <!-- Empty State -->
          <v-card-text
            v-else
            class="flex-grow-1 d-flex flex-column align-center justify-center text-grey"
          >
            <v-icon
              size="64"
              class="mb-4"
            >
              {{ mdiEmailOutline }}
            </v-icon>
            <p class="text-h6">
              メッセージを選択してください
            </p>
            <p class="text-body-2">
              左の一覧からメッセージを選択すると内容が表示されます
            </p>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<style scoped>
.message-container {
  height: calc(100vh - 200px);
  min-height: 600px;
}

.message-list-column {
  height: 100%;
}

.message-list-column .v-card {
  height: 100%;
}

.message-list-column .v-list {
  max-height: calc(100vh - 320px);
}

.unread-message {
  background-color: rgb(33 150 243 / 4%);
}

.unread-message:hover {
  background-color: rgb(33 150 243 / 8%) !important;
}

.v-list-item {
  cursor: pointer;
  transition: all 0.2s;
}

.v-list-item:hover {
  background-color: rgb(0 0 0 / 4%);
}

.message-content {
  background-color: #fafafa;
}

.message-body {
  white-space: pre-wrap;
  line-height: 1.6;
  font-size: 14px;
  color: #424242;
  max-width: 800px;
}

@media (width <= 768px) {
  .message-container {
    height: calc(100vh - 150px);
  }

  .message-list-column .v-list {
    max-height: calc(100vh - 250px);
  }
}
</style>
