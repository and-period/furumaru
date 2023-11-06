<script lang="ts" setup>
import { useAuthStore } from '~/store'
import type { SettingMenu } from '~/types/props'

const router = useRouter()
const authStore = useAuthStore()

const menus: SettingMenu[] = [
  {
    text: 'メールアドレス変更',
    action: () => router.push('/accounts/email'),
    roles: []
  },
  {
    text: 'パスワード変更',
    action: () => router.push('/accounts/password'),
    roles: []
  },
  {
    text: 'サインアウト',
    color: 'error',
    action: () => {
      authStore.logout()
      router.push('/')
    },
    roles: []
  }
]

const handleClick = (action: () => void): void => {
  action()
}
</script>

<template>
  <templates-auth-top :menus="menus" @click="handleClick" />
</template>
