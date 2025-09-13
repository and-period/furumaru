<script lang="ts" setup>
import { useAuthStore } from '~/store'
import { AdminType } from '~/types/api/v1'
import type { SettingMenu } from '~/types/props'

const router = useRouter()
const authStore = useAuthStore()

const { adminType } = storeToRefs(authStore)

const menus: SettingMenu[] = [
  {
    text: 'プロフィール変更',
    action: () => router.push('/accounts/coordinator'),
    adminTypes: [AdminType.AdminTypeCoordinator],
  },
  {
    text: 'メールアドレス変更',
    action: () => router.push('/accounts/email'),
    adminTypes: [AdminType.AdminTypeAdministrator, AdminType.AdminTypeCoordinator, AdminType.AdminTypeProducer],
  },
  {
    text: 'パスワード変更',
    action: () => router.push('/accounts/password'),
    adminTypes: [AdminType.AdminTypeAdministrator, AdminType.AdminTypeCoordinator, AdminType.AdminTypeProducer],
  },
  {
    text: 'SNSアカウント連携',
    action: () => router.push('/accounts/providers'),
    adminTypes: [AdminType.AdminTypeAdministrator, AdminType.AdminTypeCoordinator, AdminType.AdminTypeProducer],
  },
  {
    text: 'サインアウト',
    color: 'error',
    action: () => {
      authStore.logout()
      router.push('/signin')
    },
    adminTypes: [AdminType.AdminTypeAdministrator, AdminType.AdminTypeCoordinator, AdminType.AdminTypeProducer],
  },
]

const getMenus = (): SettingMenu[] => {
  return menus.filter((drawer: SettingMenu): boolean => {
    return drawer.adminTypes?.includes(adminType.value) || false
  })
}

const handleClick = (action: () => void): void => {
  action()
}
</script>

<template>
  <templates-auth-top
    :menus="getMenus()"
    @click="handleClick"
  />
</template>
