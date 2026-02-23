<script lang="ts" setup>
import { mdiLink, mdiArrowLeft, mdiCheckCircle } from '@mdi/js'
import type { AlertType } from '~/lib/hooks'
import type { AuthProvider } from '~/types/api/v1'
import { AuthProviderType } from '~/types/api/v1'

interface listItem {
  name: string
  type: AuthProviderType
  image: string
  connected: boolean
  action: () => void
}

const props = defineProps({
  loading: {
    type: Boolean,
    default: false,
  },
  isAlert: {
    type: Boolean,
    default: false,
  },
  alertType: {
    type: String as PropType<AlertType>,
    default: undefined,
  },
  alertText: {
    type: String,
    default: '',
  },
  providers: {
    type: Array<AuthProvider>,
    default: () => [],
  },
})

const emit = defineEmits<{
  (e: 'click:google'): void
  (e: 'click:line'): void
}>()

const onClickGoogle = (): void => {
  emit('click:google')
}

const onClickLINE = (): void => {
  emit('click:line')
}

const items: listItem[] = [
  {
    name: 'Google',
    type: AuthProviderType.AuthProviderTypeGoogle,
    image: '/sns/google.png',
    connected: false,
    action: onClickGoogle,
  },
  {
    name: 'LINE',
    type: AuthProviderType.AuthProviderTypeLINE,
    image: '/sns/line.png',
    connected: false,
    action: onClickLINE,
  },
]

const getItems = computed(() => {
  return items.map((item): listItem => {
    const provider = props.providers.find((provider: AuthProvider) => provider.type === item.type)
    item.connected = provider ? true : false
    return item
  })
})
</script>

<template>
  <v-container class="pa-6">
    <atoms-app-alert
      :show="props.isAlert"
      :type="props.alertType"
      :text="props.alertText"
      class="mb-6"
    />

    <div class="mb-6">
      <v-btn
        variant="text"
        :icon="mdiArrowLeft"
        class="mb-4"
        @click="$router.back()"
      >
        戻る
      </v-btn>
      <h1 class="text-h4 font-weight-bold mb-2">
        <v-icon
          :icon="mdiLink"
          size="32"
          class="mr-3 text-primary"
        />
        SNSアカウント連携
      </h1>
      <p class="text-body-1 text-grey-darken-1">
        外部SNSアカウントと連携して、簡単ログインを有効にします。
      </p>
    </div>

    <v-card
      elevation="2"
      class="provider-card"
    >
      <v-card-text class="pa-6">
        <v-alert
          type="warning"
          variant="outlined"
          class="mb-6"
        >
          <div class="text-body-2">
            <strong>重要な注意事項</strong><br><br>
            • 連携完了後、実際にログインで使えるようになるまで数分かかる場合があります<br>
            • 外部アカウント連携後、メールアドレス認証が利用できなくなるケースがあります
          </div>
        </v-alert>

        <div class="mb-4">
          <h3 class="text-subtitle-1 font-weight-medium mb-4 text-grey-darken-1">
            利用可能SNSサービス
          </h3>
        </div>

        <v-row>
          <v-col
            v-for="item in getItems"
            :key="item.type"
            cols="12"
            sm="6"
          >
            <v-card
              :class="[
                'provider-item',
                item.connected ? 'connected' : 'disconnected',
              ]"
              :variant="item.connected ? 'outlined' : 'elevated'"
              :color="item.connected ? 'success' : 'transparent'"
              elevation="1"
            >
              <v-card-text class="text-center pa-6">
                <div class="mb-4">
                  <v-avatar
                    size="64"
                    :class="item.connected ? 'connected-avatar' : ''"
                  >
                    <v-img
                      :src="item.image"
                      :alt="item.name || '認証プロバイダ'"
                    />
                    <v-icon
                      v-if="item.connected"
                      :icon="mdiCheckCircle"
                      size="20"
                      class="connected-badge"
                      color="success"
                    />
                  </v-avatar>
                </div>
                <h4 class="text-h6 font-weight-medium mb-2">
                  {{ item.name }}
                </h4>
                <p class="text-body-2 text-grey-darken-1 mb-4">
                  {{ item.connected ? '連携済み' : '未連携' }}
                </p>
                <v-btn
                  v-if="item.connected"
                  color="success"
                  variant="outlined"
                  disabled
                  block
                >
                  <v-icon
                    :icon="mdiCheckCircle"
                    start
                  />
                  連携済み
                </v-btn>
                <v-btn
                  v-else
                  color="primary"
                  variant="elevated"
                  block
                  @click="item.action"
                >
                  <v-icon
                    :icon="mdiLink"
                    start
                  />
                  {{ item.name }}で連携
                </v-btn>
              </v-card-text>
            </v-card>
          </v-col>
        </v-row>
      </v-card-text>
    </v-card>
  </v-container>
</template>

<style scoped>
.provider-card {
  border-radius: 12px;
  max-width: 800px;
  margin: 0 auto;
}

.provider-item {
  border-radius: 12px;
  transition: all 0.3s ease;
  cursor: default;
  height: 100%;
}

.provider-item.disconnected:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 15px rgb(0 0 0 / 10%) !important;
}

.provider-item.connected {
  border-color: rgb(76 175 80) !important;
  background: rgb(76 175 80 / 5%);
}

.connected-avatar {
  position: relative;
}

.connected-badge {
  position: absolute;
  bottom: -4px;
  right: -4px;
  background: white;
  border-radius: 50%;
  padding: 2px;
}

@media (width <= 600px) {
  .provider-card {
    margin: 0;
  }
}
</style>
