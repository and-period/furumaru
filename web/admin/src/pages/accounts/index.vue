<script lang="ts" setup>
import { useAuthStore } from '~/store'
import { AdminRole } from '~/types/api';
import type { SettingMenu } from '~/types/props'

const router = useRouter()
const authStore = useAuthStore()

const { role } = storeToRefs(authStore)

const menus: SettingMenu[] = [
  {
    text: 'プロフィール変更',
    action: () => router.push('/accounts/coordinator'),
    roles: [AdminRole.COORDINATOR]
  },
  {
    text: 'メールアドレス変更',
    action: () => router.push('/accounts/email'),
    roles: [AdminRole.ADMINISTRATOR, AdminRole.COORDINATOR, AdminRole.PRODUCER]
  },
  {
    text: 'パスワード変更',
    action: () => router.push('/accounts/password'),
    roles: [AdminRole.ADMINISTRATOR, AdminRole.COORDINATOR, AdminRole.PRODUCER]
  },
  {
    text: 'サインアウト',
    color: 'error',
    action: () => {
      authStore.logout()
      router.push('/signin')
    },
    roles: [AdminRole.ADMINISTRATOR, AdminRole.COORDINATOR, AdminRole.PRODUCER]
  }
]

const getMenus = (): SettingMenu[] => {
  return menus.filter((drawer: SettingMenu): boolean => {
    return drawer.roles?.includes(role.value) || false
  })
}

const handleClick = (action: () => void): void => {
  action()
}
</script>

<template>
  <templates-auth-top :menus="getMenus()" @click="handleClick" />
</template>
