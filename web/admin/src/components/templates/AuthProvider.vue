<script lang="ts" setup>
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
  <v-alert
    v-show="props.isAlert"
    :type="props.alertType"
    v-text="props.alertText"
  />

  <v-card class="mt-4 flat">
    <v-card-title>認証用の外部アカウント連携</v-card-title>

    <v-card-text>
      <div class="text-red">
        <p>※連携完了後、実際にログインで使えるようになるまでは少し時間が時間がかかる場合があります。</p>
        <p>※外部アカウント連携後、メールアドレス認証が利用できなくなるケースが存在します。</p>
      </div>

      <v-list>
        <v-list-item
          v-for="item in getItems"
          :key="item.type"
          :title="item.name"
        >
          <template #prepend>
            <v-avatar color="white">
              <v-img :src="item.image" />
            </v-avatar>
          </template>

          <template #append>
            <v-btn
              v-if="item.connected"
              color="unknown"
              disabled
            >
              連携済み
            </v-btn>
            <v-btn
              v-else
              color="primary"
              @click="item.action"
            >
              連携する
            </v-btn>
          </template>
        </v-list-item>
      </v-list>
    </v-card-text>
  </v-card>
</template>
