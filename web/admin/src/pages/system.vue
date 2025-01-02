<script lang="ts" setup>
import { storeToRefs } from 'pinia'
import { useAuthStore } from '~/store'
import { AdminType } from '~/types/api'
import type { SettingMenu } from '~/types/props'

const router = useRouter()
const authStore = useAuthStore()

const { adminType } = storeToRefs(authStore)

const menus: SettingMenu[] = [
  {
    text: '管理者管理',
    action: () => router.push('/administrators'),
    adminTypes: [AdminType.ADMINISTRATOR],
  },
  {
    text: 'コーディネーター管理',
    action: () => router.push('/coordinators'),
    adminTypes: [AdminType.ADMINISTRATOR],
  },
  {
    text: 'カテゴリー・品目管理',
    action: () => router.push('/categories'),
    adminTypes: [AdminType.ADMINISTRATOR, AdminType.COORDINATOR],
  },
  {
    text: '商品タグ管理',
    action: () => router.push('/product-tags'),
    adminTypes: [AdminType.ADMINISTRATOR, AdminType.COORDINATOR],
  },
  {
    text: '体験種別管理',
    action: () => router.push('/experience-types'),
    adminTypes: [AdminType.ADMINISTRATOR, AdminType.COORDINATOR],
  },
  {
    text: '配送設定管理',
    action: () => router.push('/shippings'),
    adminTypes: [AdminType.COORDINATOR],
  },
  {
    text: 'デフォルト配送設定管理',
    action: () => router.push('/shippings/default'),
    adminTypes: [AdminType.ADMINISTRATOR],
  },
  {
    text: 'スポット種別管理',
    action: () => router.push('/spot-types'),
    adminTypes: [AdminType.ADMINISTRATOR, AdminType.COORDINATOR],
  },
  {
    text: '決済システム管理',
    action: () => router.push('/payment-systems'),
    adminTypes: [AdminType.ADMINISTRATOR],
  },
]

const getMenus = (): SettingMenu[] => {
  return menus.filter((menu: SettingMenu): boolean => {
    return menu.adminTypes.includes(adminType.value)
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
