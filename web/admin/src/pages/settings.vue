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

<script lang="ts">
import { defineComponent, ref, useRouter } from '@nuxtjs/composition-api'

export default defineComponent({
  setup() {
    const router = useRouter()

    const menuList = ref<{ text: string; onClick: Function; class?: string }[]>(
      [
        {
          text: 'コーディネーター登録',
          onClick: () => {
            router.push('/coordinator/add')
          },
        },
        {
          text: 'ジャンル・品目登録',
          onClick: () => {
            router.push('/categories')
          },
        },
        {
          text: 'メールアドレス変更',
          onClick: () => {
            router.push('/accounts/email')
          },
        },
        {
          text: 'パスワード変更',
          onClick: () => {
            router.push('/accounts/password')
          },
        },
        {
          text: 'サインアウト',
          class: 'red--text ',
          onClick: () => {
            console.log('NOT IMPLEMENTED')
          },
        },
      ]
    )

    return {
      menuList,
    }
  },
})
</script>
