<script lang="ts" setup>
import { storeToRefs } from 'pinia'
import { useAuthStore } from '~/store'
import { AdminType } from '~/types/api/v1'
import type { SettingMenu } from '~/types/props'

const router = useRouter()
const authStore = useAuthStore()

const { adminType } = storeToRefs(authStore)

const menus: SettingMenu[] = [
  {
    text: '管理者管理',
    action: () => router.push('/administrators'),
    adminTypes: [AdminType.AdminTypeAdministrator],
  },
  {
    text: 'コーディネーター管理',
    action: () => router.push('/coordinators'),
    adminTypes: [AdminType.AdminTypeCoordinator],
  },
  {
    text: 'カテゴリー・品目管理',
    action: () => router.push('/categories'),
    adminTypes: [AdminType.AdminTypeAdministrator, AdminType.AdminTypeCoordinator],
  },
  {
    text: '商品タグ管理',
    action: () => router.push('/product-tags'),
    adminTypes: [AdminType.AdminTypeAdministrator, AdminType.AdminTypeCoordinator],
  },
  {
    text: '体験種別管理',
    action: () => router.push('/experience-types'),
    adminTypes: [AdminType.AdminTypeAdministrator, AdminType.AdminTypeCoordinator],
  },
  {
    text: '配送設定管理',
    action: () => router.push('/shippings'),
    adminTypes: [AdminType.AdminTypeCoordinator],
  },
  {
    text: 'デフォルト配送設定管理',
    action: () => router.push('/shippings/default'),
    adminTypes: [AdminType.AdminTypeAdministrator],
  },
  {
    text: 'スポット種別管理',
    action: () => router.push('/spot-types'),
    adminTypes: [AdminType.AdminTypeAdministrator, AdminType.AdminTypeCoordinator],
  },
  {
    text: '決済システム管理',
    action: () => router.push('/payment-systems'),
    adminTypes: [AdminType.AdminTypeAdministrator],
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
