<script lang="ts" setup>
import { useAuthStore } from '~/store'
import { AdminRole } from '~/types/api'
import { SettingMenu } from '~/types/props'

const router = useRouter()

const menus: SettingMenu[] = [
  {
    text: '管理者管理',
    action: () => router.push('/administrators'),
    roles: [AdminRole.ADMINISTRATOR]
  },
  {
    text: 'コーディネータ管理',
    action: () => router.push('/coordinators'),
    roles: [AdminRole.ADMINISTRATOR]
  },
  {
    text: 'カテゴリー・品目管理',
    action: () => router.push('/categories'),
    roles: [AdminRole.ADMINISTRATOR, AdminRole.COORDINATOR]
  },
  {
    text: '商品タグ管理',
    action: () => router.push('/product-tags'),
    roles: [AdminRole.ADMINISTRATOR, AdminRole.COORDINATOR]
  },
  {
    text: '配送設定管理',
    action: () => router.push('/shippings'),
    roles: [AdminRole.ADMINISTRATOR, AdminRole.COORDINATOR]
  }
]

const handleClick = (action: () => void): void => {
  action()
}
</script>

<template>
  <templates-system-top :menus="menus" @click="handleClick" />
</template>
