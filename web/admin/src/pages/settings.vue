<script lang="ts">
import { useAuthStore } from '~/store'

const router = useRouter()
const { logout } = useAuthStore()

const menuList = ref<{ text: string; onClick: Function; class?: string }[]>([
  {
    text: 'カテゴリー・品目登録',
    onClick: () => {
      router.push('/categories')
    }
  },
  {
    text: '配送設定',
    onClick: () => {
      router.push('/shippings')
    }
  },
  {
    text: 'メールアドレス変更',
    onClick: () => {
      router.push('/accounts/email')
    }
  },
  {
    text: 'パスワード変更',
    onClick: () => {
      router.push('/accounts/password')
    }
  },
  {
    text: 'サインアウト',
    class: 'red--text ',
    onClick: () => {
      logout()
      router.push('/')
    }
  }
])
</script>

<template>
  <div>
    <v-card-title>システム設定</v-card-title>
    <v-list>
      <div v-for="(item, i) in menuList" :key="i">
        <v-list-item :class="item.class" @click="item.onClick">
          <v-list-item-content>
            <v-list-item-title>{{ item.text }}</v-list-item-title>
          </v-list-item-content>
        </v-list-item>
        <v-divider />
      </div>
    </v-list>
  </div>
</template>
