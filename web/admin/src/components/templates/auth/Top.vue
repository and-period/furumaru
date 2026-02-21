<script lang="ts" setup>
import {
  mdiAccountEdit,
  mdiEmail,
  mdiLock,
  mdiLink,
  mdiLogout,
  mdiChevronRight,
} from '@mdi/js'
import type { SettingMenu } from '~/types/props'

const props = defineProps({
  menus: {
    type: Array<SettingMenu>,
    default: () => [],
  },
})

const emit = defineEmits<{
  (e: 'click', action: () => void): void
}>()

const onClick = (action: () => void): void => {
  emit('click', action)
}

const getMenuIcon = (text: string): string => {
  if (text.includes('プロフィール')) return mdiAccountEdit
  if (text.includes('メールアドレス')) return mdiEmail
  if (text.includes('パスワード')) return mdiLock
  if (text.includes('SNS') || text.includes('アカウント連携')) return mdiLink
  if (text.includes('サインアウト')) return mdiLogout
  return mdiAccountEdit
}

const getMenuColor = (menu: SettingMenu): string => {
  if (menu.color === 'error') return 'error'
  if (menu.text.includes('パスワード')) return 'warning'
  return 'primary'
}

const getMenuVariant = (menu: SettingMenu): string => {
  return menu.color === 'error' ? 'outlined' : 'flat'
}

const getMenuDescription = (text: string): string => {
  if (text.includes('プロフィール')) return 'コーディネーター情報を編集'
  if (text.includes('メールアドレス')) return 'ログイン用メールアドレスを変更'
  if (text.includes('パスワード')) return '新しいパスワードに変更'
  if (text.includes('SNS') || text.includes('アカウント連携')) return 'Google・LINE等のSNS連携'
  if (text.includes('サインアウト')) return 'アカウントからログアウト'
  return '設定を変更'
}
</script>

<template>
  <v-container class="pa-6">
    <div class="mb-8">
      <h1 class="text-h4 font-weight-bold mb-2">
        <v-icon
          :icon="mdiAccountEdit"
          size="32"
          class="mr-3 text-primary"
        />
        アカウント設定
      </h1>
      <p class="text-body-1 text-grey-darken-1">
        プロフィール、セキュリティ設定を管理できます
      </p>
    </div>

    <v-row>
      <v-col
        v-for="(item, i) in props.menus"
        :key="i"
        cols="12"
        sm="6"
        md="4"
      >
        <v-card
          :class="[
            'setting-card',
            'h-100',
            item.color === 'error' ? 'danger-card' : '',
          ]"
          :color="item.color === 'error' ? 'error' : 'transparent'"
          :variant="item.color === 'error' ? 'outlined' : 'elevated'"
          elevation="2"
          hover
          @click="onClick(item.action)"
        >
          <v-card-text class="text-center pa-6">
            <div class="mb-4">
              <v-avatar
                :color="getMenuColor(item)"
                :variant="getMenuVariant(item)"
                size="56"
              >
                <v-icon
                  :icon="getMenuIcon(item.text)"
                  size="28"
                  :color="item.color === 'error' ? 'error' : 'white'"
                />
              </v-avatar>
            </div>
            <h3
              class="text-subtitle-1 font-weight-medium mb-2"
              :class="item.color === 'error' ? 'text-error' : ''"
            >
              {{ item.text }}
            </h3>
            <p
              class="text-body-2 text-grey-darken-1"
              :class="item.color === 'error' ? 'text-error-darken-1' : ''"
            >
              {{ getMenuDescription(item.text) }}
            </p>
          </v-card-text>
          <v-card-actions class="pa-4 pt-0">
            <v-spacer />
            <v-btn
              :color="getMenuColor(item)"
              :variant="item.color === 'error' ? 'outlined' : 'text'"
              size="small"
            >
              {{ item.color === 'error' ? 'サインアウト' : '変更する' }}
              <v-icon
                :icon="item.color === 'error' ? mdiLogout : mdiChevronRight"
                end
              />
            </v-btn>
          </v-card-actions>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<style scoped>
.setting-card {
  transition: all 0.3s ease;
  cursor: pointer;
  border-radius: 12px;
}

.setting-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 25px rgb(0 0 0 / 15%) !important;
}

.danger-card {
  border-color: rgb(244 67 54) !important;
}

.danger-card:hover {
  background-color: rgb(244 67 54 / 5%);
}

@media (width <= 600px) {
  .setting-card {
    margin-bottom: 16px;
  }
}
</style>
