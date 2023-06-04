<script lang="ts" setup>
import { useAuthStore } from '~/store'
import { SettingMenu } from '~/types/props'

const router = useRouter()
const authStore = useAuthStore()

const menus: SettingMenu[] = [
  {
    text: 'メールアドレス変更',
    action: () => router.push('/accounts/email')
  },
  {
    text: 'パスワード変更',
    action: () => router.push('/accounts/password')
  },
  {
    text: 'バージョン情報',
    action: () => router.push('/version')
  },
  {
    text: 'サインアウト',
    color: 'error',
    action: () => {
      authStore.logout()
      router.push('/')
    }
  }
]

const handleClick = (action: () => void): void => {
  action()
}
</script>

<template>
  <templates-setting-top :menus="menus" @click="handleClick" />
</template>
