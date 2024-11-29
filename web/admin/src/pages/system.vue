<script lang="ts" setup>
import { storeToRefs } from 'pinia'
import { useAuthStore } from '~/store'
import { AdminRole } from '~/types/api'
import type { SettingMenu } from '~/types/props'

const router = useRouter()
const authStore = useAuthStore()

const { role } = storeToRefs(authStore)

const menus: SettingMenu[] = [
  {
    text: '管理者管理',
    action: () => router.push('/administrators'),
    roles: [AdminRole.ADMINISTRATOR],
  },
  {
    text: 'コーディネーター管理',
    action: () => router.push('/coordinators'),
    roles: [AdminRole.ADMINISTRATOR],
  },
  {
    text: 'カテゴリー・品目管理',
    action: () => router.push('/categories'),
    roles: [AdminRole.ADMINISTRATOR, AdminRole.COORDINATOR],
  },
  {
    text: '商品タグ管理',
    action: () => router.push('/product-tags'),
    roles: [AdminRole.ADMINISTRATOR, AdminRole.COORDINATOR],
  },
  {
    text: '体験種別管理',
    action: () => router.push('/experience-types'),
    roles: [AdminRole.ADMINISTRATOR, AdminRole.COORDINATOR],
  },
  {
    text: '配送設定管理',
    action: () => router.push('/shippings'),
    roles: [AdminRole.COORDINATOR],
  },
  {
    text: 'デフォルト配送設定管理',
    action: () => router.push('/shippings/default'),
    roles: [AdminRole.ADMINISTRATOR],
  },
  {
    text: '決済システム管理',
    action: () => router.push('/payment-systems'),
    roles: [AdminRole.ADMINISTRATOR],
  },
]

const getMenus = (): SettingMenu[] => {
  return menus.filter((menu: SettingMenu): boolean => {
    return menu.roles.includes(role.value)
  })
}

const handleClick = (action: () => void): void => {
  action()
}
</script>

<template>
  <templates-system-top
    :menus="getMenus()"
    @click="handleClick"
  />
</template>
