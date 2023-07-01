<script lang="ts" setup>
import { useAuthStore } from '~/store'
import { SettingMenu } from '~/types/props'

const router = useRouter()

const menus: SettingMenu[] = [
  {
    text: '管理者管理',
    action: () => router.push('/administrators')
  },
  {
    text: 'コーディネータ管理',
    action: () => router.push('/coordinators')
  },
  {
    text: 'カテゴリー・品目管理',
    action: () => router.push('/categories')
  },
  {
    text: '商品タグ管理',
    action: () => router.push('/product-tags')
  },
  {
    text: '配送設定管理',
    action: () => router.push('/shippings')
  }
]

const handleClick = (action: () => void): void => {
  action()
}
</script>

<template>
  <templates-system-top :menus="menus" @click="handleClick" />
</template>
