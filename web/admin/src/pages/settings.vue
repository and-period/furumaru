<script lang="ts" setup>
import { useAuthStore } from '~/store'
import { SettingMenu } from '~/types/props'

const router = useRouter()
const authStore = useAuthStore()

const menus: SettingMenu[] = [
  {
    text: 'カテゴリー・品目登録',
    action: () => router.push('/categories')
  },
  {
    text: '配送設定',
    action: () => router.push('/shippings')
  },
  {
    text: 'メールアドレス変更',
    action: () => router.push('/accounts/email')
  },
  {
    text: 'パスワード変更',
    action: () => router.push('/accounts/password')
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
  <templates-setting-top-page :menus="menus" @click="handleClick" />
</template>
